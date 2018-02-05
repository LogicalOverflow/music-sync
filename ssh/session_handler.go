package ssh

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/ssh/parser"
	"github.com/chzyer/readline"
	"github.com/gliderlabs/ssh"
	"io"
	"strings"
)

type session struct {
	ssh.Session
	cfg *readline.Config
	ex  *readline.Instance

	width, height int
}

func (s *session) preparePty() bool {
	pty, resize, ok := s.Pty()
	if !ok {
		logger.Infof("no pty req from %s as %s", s.RemoteAddr(), s.User())
		return false
	}

	s.width, s.height = pty.Window.Width, pty.Window.Height
	go func() {
		for w := range resize {
			s.width, s.height = w.Width, w.Height
		}
	}()

	return true
}

func (s *session) init() bool {
	if !s.preparePty() {
		return false
	}

	s.createReadlineConfig()

	if err := s.createEx(); err != nil {
		logger.Warnf("failed to create ex: %v", err)
		fmt.Fprintf(s, "failed to create ex: %v\n", err)
		return false
	}

	return true
}

func (s *session) createEx() (err error) {
	s.ex, err = readline.NewEx(s.cfg)
	return
}

func (s *session) createReadlineConfig() {
	s.cfg = &readline.Config{
		Prompt: "\033[31mmusic-sync»\033[0m ",

		AutoComplete: &sshAutoCompleter{},

		VimMode: false,

		InterruptPrompt: "^C",
		EOFPrompt:       "exit",

		FuncGetWidth: func() int { return s.width },

		Stdin:       s,
		Stdout:      s,
		StdinWriter: s,
		Stderr:      s.Stderr(),

		UniqueEditLine: false,
	}
}

func (s *session) handleLine(line string) bool {
	cmd := parser.ParseCommand(line)

	if c := commandByName(cmd.Command); c != nil {
		s.execCommand(*c, cmd.Parameters)
	} else {
		switch cmd.Command {
		case "clear":
			fmt.Fprint(s.ex, "\033[H")
		case "exit":
			logger.Infof("connection %s as %s closed", s.RemoteAddr(), s.User())
			return false
		default:
			fmt.Fprintf(s.ex, "Unknown command '%s'. Type 'help' for help.\n", cmd.Command)
		}
	}
	return true
}

func (s *session) execCommand(c Command, args []string) {
	msg, ok := c.Exec(args)
	if ok {
		if strings.HasSuffix(msg, "\n") {
			msg = msg[:len(msg)-1]
		}
		fmt.Fprintln(s.ex, msg)
	} else {
		fmt.Fprintf(s.ex, "%s\n", c.usage())
	}
}

func (s *session) readLoop() {
	for {
		line, err := s.ex.Readline()
		if err == readline.ErrInterrupt {
		} else if err == io.EOF {
			logger.Infof("connection %s as %s closed", s.RemoteAddr(), s.User())
			return
		} else if err != nil {
			logger.Infof("connection error from %s as %s: %v", s.RemoteAddr(), s.User(), err)
			return
		}

		if !s.handleLine(line) {
			return
		}
	}
}

func sessionHandler(sshSession ssh.Session) {
	s := &session{Session: sshSession}
	defer s.Close()

	logger.Infof("new ssh connection from %s as %s", s.RemoteAddr(), s.User())

	if !s.init() {
		return
	}

	defer s.ex.Close()
	s.ex.Clean()
	s.readLoop()
}
