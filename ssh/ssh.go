// Package ssh contains the ssh control interface
package ssh

import (
	"encoding/json"
	"fmt"
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/gliderlabs/ssh"
	"os"
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
	ssh.Handle(sessionHandler)

	options := getSSHOptions(users)
	logger.Infof("starting ssh server at %s", address)
	err := ssh.ListenAndServe(address, nil, options...)
	logger.Errorf("ssh server at %s stopped: %v", address, err)
}
