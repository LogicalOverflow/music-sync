// Package ssh contains the ssh control interface
package ssh

import (
	"encoding/json"
	"fmt"
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/LogicalOverflow/music-sync/ssh/parser"
	"github.com/LogicalOverflow/music-sync/util"
	"github.com/chzyer/readline"
	"github.com/gliderlabs/ssh"
	"io"
	"os"
	"strings"
)

// HostKeyFile is the path to the host key to use for the ssh control interface
var HostKeyFile string

// UserAuth contains the data to authenticate a user connection to the ssh control interface
type UserAuth struct {
	Password string `json:"password"`
	PubKey   []byte `json:"pubKey"`
}

var logger = log.GetLogger("ssh")

// ReadUsersFile reads a json file containing all users and passwords allowed to access the ssh control interface
// It returns a dictionary user->password
func ReadUsersFile(filename string) (map[string]UserAuth, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open users file (%s): %v", filename, err)
	}
	users := make(map[string]UserAuth)
	if err := json.NewDecoder(f).Decode(&users); err != nil {
		return nil, fmt.Errorf("failed to decode users file (%s): %v", filename, err)
	}
	return users, nil
}

// StartSSH starts the ssh control interface on listening on address and accepting all users with the respective
// passwords in the users dict (user->password)
func StartSSH(address string, users map[string]UserAuth) {
	ssh.Handle(func(s ssh.Session) {
		defer s.Close()

		logger.Infof("new ssh connection from %s as %s", s.RemoteAddr(), s.User())
		pty, resize, ok := s.Pty()
		if !ok {

			logger.Infof("no pty req from %s as %s", s.RemoteAddr(), s.User())
			return
		}
		width, height := pty.Window.Width, pty.Window.Height
		go func() {
			for w := range resize {
				width, height = w.Width, w.Height
			}
		}()

		tcfg := &readline.Config{
			Prompt: "\033[31mmusic-sync»\033[0m ",

			AutoComplete: &sshAutoCompleter{},

			VimMode: false,

			InterruptPrompt: "^C",
			EOFPrompt:       "exit",

			FuncGetWidth: func() int { return width },

			Stdin:       s,
			Stdout:      s,
			StdinWriter: s,
			Stderr:      s.Stderr(),

			UniqueEditLine: false,
		}
		ex, err := readline.NewEx(tcfg)
		if err != nil {
			logger.Warnf("failed to create ex: %v", err)
			fmt.Fprintf(s, "failed to create ex: %v\n", err)
			return
		}
		defer ex.Close()
		ex.Clean()
		for {
			line, err := ex.Readline()
			if err == readline.ErrInterrupt {
			} else if err == io.EOF {
				logger.Infof("connection %s as %s closed", s.RemoteAddr(), s.User())
				return
			} else if err != nil {
				logger.Infof("connection error from %s as %s: %v", s.RemoteAddr(), s.User(), err)
				return
			}

			cmd := parser.ParseCommand(line)

			known := false
			for _, c := range commands {
				if c.Name == cmd.Command {
					s, ok := c.Exec(cmd.Parameters)
					if ok {
						if strings.HasSuffix(s, "\n") {
							s = s[:len(s)-1]
						}
						fmt.Fprintln(ex, s)
					} else {
						fmt.Fprintf(ex, "%s\n", c.usage())
					}
					known = true
					break
				}
			}
			if !known {
				switch cmd.Command {
				case "clear":
					fmt.Fprint(ex, "\033[H")
				case "exit":
					logger.Infof("connection %s as %s closed", s.RemoteAddr(), s.User())
					return
				default:
					fmt.Fprintf(ex, "Unknown command '%s'. Type 'help' for help.\n", cmd.Command)
				}
			}
		}
	})

	logger.Infof("starting ssh server at %s", address)
	options := make([]ssh.Option, 0)
	if HostKeyFile == "" {
		logger.Warnf("no host key file provided, generating a new host key")
	} else if err := util.CheckFile(HostKeyFile); err != nil {
		logger.Warnf("unable to access host key file, generating new host key: %v", err)
	} else {
		options = append(options, ssh.HostKeyFile(HostKeyFile))
	}
	options = append(options, ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
		auth, ok := users[ctx.User()]
		if ok && auth.Password != "" && auth.Password == password {
			return true
		}
		logger.Warnf("failed ssh login attempt from %s as %s using a password", ctx.RemoteAddr(), ctx.User())
		return false
	}), ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
		auth, ok := users[ctx.User()]
		if !ok || len(auth.PubKey) == 0 {
			logger.Warnf("failed ssh login attempt from %s as %s using a public key", ctx.RemoteAddr(), ctx.User())
			return false
		}
		authKey, err := ssh.ParsePublicKey(auth.PubKey)
		if err != nil {
			logger.Warnf("failed to parse ssh key for %s: %v", ctx.User(), err)
			return false
		}
		if ssh.KeysEqual(key, authKey) {
			return true
		}
		logger.Warnf("failed ssh login attempt from %s as %s using a public key", ctx.RemoteAddr(), ctx.User())
		return false
	}))

	err := ssh.ListenAndServe(address, nil, options...)
	logger.Errorf("ssh server at %s stopped: %v", address, err)
}

type sshAutoCompleter struct{}

func (sac *sshAutoCompleter) Do(line []rune, pos int) (newLine [][]rune, length int) {
	newLine = make([][]rune, 0)

	lastChar := '\x00'
	if len(line) != 0 {
		lastChar = line[len(line)-1]
	}

	cmd := parser.ParseCommand(string(line))

	if len(cmd.Parameters) == 0 && lastChar != ' ' {
		for _, c := range commands {
			if strings.HasPrefix(c.Name, cmd.Command) {
				newLine = append(newLine, []rune(c.Name + " ")[pos:])
			}
		}

		for _, c := range []string{"clear", "exit"} {
			if strings.HasPrefix(c, cmd.Command) {
				newLine = append(newLine, []rune(c + " ")[pos:])
			}
		}
	} else {
		for _, c := range commands {
			if c.Name == cmd.Command {
				if c.Options != nil {
					if len(cmd.Parameters) == 0 {
						cmd.Parameters = []string{""}
					}
					argNum := len(cmd.Parameters) - 1
					arg := cmd.Parameters[argNum]
					for _, o := range c.Options(arg, argNum) {
						cmd.Parameters[argNum] = o
						newLine = append(newLine, []rune(cmd.Unparse())[pos:])
					}
					cmd.Parameters[argNum] = arg
				}
			}
		}
	}
	length = pos
	return
}
