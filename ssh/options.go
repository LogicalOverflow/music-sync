package ssh

import (
	"github.com/LogicalOverflow/music-sync/util"
	"github.com/gliderlabs/ssh"
)

func getSSHOptions(users map[string]UserAuth) []ssh.Option {
	options := []ssh.Option{createPasswordAuthOption(users), createPublicKeyAuthOption(users)}
	if hostKeyOption := getHostKeyOption(); hostKeyOption != nil {
		options = append(options, hostKeyOption)
	}
	return options
}

func getHostKeyOption() ssh.Option {
	if HostKeyFile == "" {
		logger.Warnf("no host key file provided, generating a new host key")
	} else if err := util.CheckFile(HostKeyFile); err != nil {
		logger.Warnf("unable to access host key file, generating new host key: %v", err)
	} else {
		return ssh.HostKeyFile(HostKeyFile)
	}
	return nil
}

func createPasswordAuthOption(users map[string]UserAuth) ssh.Option {
	return ssh.PasswordAuth(func(ctx ssh.Context, password string) bool {
		auth, ok := users[ctx.User()]
		if ok && auth.Password != "" && auth.Password == password {
			return true
		}
		logger.Warnf("failed ssh login attempt from %s as %s using a password", ctx.RemoteAddr(), ctx.User())
		return false
	})
}

func createPublicKeyAuthOption(users map[string]UserAuth) ssh.Option {
	return ssh.PublicKeyAuth(func(ctx ssh.Context, key ssh.PublicKey) bool {
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
	})
}
