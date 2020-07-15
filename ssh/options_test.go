package ssh

import (
	"context"
	"encoding/base64"
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/gliderlabs/ssh"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net"
	"path"
	"sync"
	"testing"
)

type tCtx struct {
	context.Context
	l             *sync.Mutex
	user          string
	sessionID     string
	clientVersion string
	serverVersion string
	remoteAddr    net.Addr
	localAddr     net.Addr
	permissions   *ssh.Permissions
}

func (t tCtx) User() string                    { return t.user }
func (t tCtx) SessionID() string               { return t.sessionID }
func (t tCtx) ClientVersion() string           { return t.clientVersion }
func (t tCtx) ServerVersion() string           { return t.serverVersion }
func (t tCtx) RemoteAddr() net.Addr            { return t.remoteAddr }
func (t tCtx) LocalAddr() net.Addr             { return t.localAddr }
func (t tCtx) Permissions() *ssh.Permissions   { return t.permissions }
func (t tCtx) SetValue(key, value interface{}) {}
func (t tCtx) Lock()                           { t.l.Lock() }
func (t tCtx) Unlock()                         { t.l.Unlock() }

var testPubKeyBytes, pubKeyBytesErr = base64.StdEncoding.DecodeString("AAAAB3NzaC1kc3MAAACBAPZ53Jl0pArkxTrx6kGTUB9FcE5luWaLowbkSybVG5GzqHKTFehDHTP6TQ5tZyRlG5nGeYJbjIbaQ8zALhu4Ubl2THxliFN8pdARSfK9pfThpV0AM1naYqo7qAx3o+6Jc85FZpnMqQckJCxejfeD2BFGJIwCE4D3UyZ5ZdNh4EtBAAAAFQCI/Q2U0TOYjc0noYnXB9DxEixJawAAAIBHVl7srtbQMAjnTN9vpQ4zaHS8XeroDreDW0TGtdBAezqMgWEURwfh32omnEEWct25PcHR2oUS/P6aHa4KPtytJiDGg45LeS14RN5HyKKNFQljGD1eL9Z7KtNoHJYLdzZBz2jgA2uqf1MSR6fRwbqxGYwCi/Bd/OtIgA/XqcUa8QAAAIEA0ct95BoR9+qBLqpwUvJSAPEP0wwQ88c0/YwyOte/sh9homF6qC3ky3XVB0CDImv2LO3GhWvaGtVoKj2WRqr68p3NC6Cy3gbpvcZzSqJJqGO0U2Ai0NoDT5A853oO8rOqVoshxGvG+nyyQ4WGXnxlO/8v5d9DZDWNklPR085qDIU=")
var testPubKey, pubKeyErr = ssh.ParsePublicKey(testPubKeyBytes)

var testUsers = map[string]UserAuth{
	"test-user":     {Password: "test-pw", PubKey: testPubKeyBytes},
	"password-user": {Password: "test-pw"},
	"pub-key-user":  {PubKey: testPubKeyBytes},
}

var passwordCases = []struct {
	username, password string
	success            bool
}{
	{username: "test-user", password: "test-pw", success: true},
	{username: "password-user", password: "test-pw", success: true},
	{username: "password-user", password: "wrong-pw", success: false},
	{username: "wrong-user", password: "a-pw", success: false},
	{username: "wrong-user", password: "", success: false},
	{username: "pub-key-user", password: "a-pub-key", success: false},
	{username: "pub-key-user", password: "", success: false},
}
var pubKeyCases = []struct {
	username string
	pubKey   ssh.PublicKey
	success  bool
}{
	{username: "test-user", pubKey: testPubKey, success: true},
	{username: "test-user", pubKey: nil, success: false},
	{username: "password-user", pubKey: testPubKey, success: false},
	{username: "password-user", pubKey: nil, success: false},
	{username: "wrong-user", pubKey: testPubKey, success: false},
	{username: "wrong-user", pubKey: nil, success: false},
	{username: "pub-key-user", pubKey: testPubKey, success: true},
	{username: "pub-key-user", pubKey: nil, success: false},
}

func TestCreatePasswordAuthOption(t *testing.T) {
	log.DefaultCutoffLevel = log.LevelOff

	pwao := createPasswordAuthOption(testUsers)
	svr := new(ssh.Server)
	pwao(svr)
	require.False(t, svr.PasswordHandler == nil, "PasswordAuthOption did not set password handler")

	for _, c := range passwordCases {
		ok := svr.PasswordHandler(tCtx{user: c.username}, c.password)
		assert.Equal(t, c.success, ok, "PasswordAuthOption did not return correctly for case %v", c)
	}
}

func TestCreatePublicKeyAuthOption(t *testing.T) {
	log.DefaultCutoffLevel = log.LevelOff

	pkao := createPublicKeyAuthOption(testUsers)
	svr := new(ssh.Server)
	pkao(svr)
	require.NotNil(t, testPubKey, "testPubKey is nil: %v, %v: %v", pubKeyErr, pubKeyBytesErr, testPubKeyBytes)
	require.False(t, svr.PublicKeyHandler == nil, "PublicKeyAuthOption did not set public key handler")

	for _, c := range pubKeyCases {
		ok := svr.PublicKeyHandler(tCtx{user: c.username}, c.pubKey)
		assert.Equal(t, c.success, ok, "PublicKeyAuthOption did not return correctly for case %v", c)
	}
}

func TestGetHostKeyOption(t *testing.T) {
	log.DefaultCutoffLevel = log.LevelOff
	oldHostKeyFile := HostKeyFile

	HostKeyFile = ""
	assert.Nil(t, getHostKeyOption(), "HostKeyOption is not nil for HostKeyFile \"\"")

	HostKeyFile = path.Join("_host_key_test_files", "non-existent")
	assert.Nil(t, getHostKeyOption(), "HostKeyOption is not nil for HostKeyFile \"non-existent\"")

	HostKeyFile = path.Join("_host_key_test_files", "id_rsa")
	assert.NotNil(t, getHostKeyOption(), "HostKeyOption is nil for HostKeyFile \"id_rsa\"")

	HostKeyFile = oldHostKeyFile
}

func TestGetSSHOptions(t *testing.T) {
	log.DefaultCutoffLevel = log.LevelOff
	oldHostKeyFile := HostKeyFile

	HostKeyFile = ""
	assert.Equal(t, 2, len(getSSHOptions(testUsers)), "wrong number of options returned by getSSHOptions for HostKeyFile \"\"")

	HostKeyFile = path.Join("_host_key_test_files", "non-existent")
	assert.Equal(t, 2, len(getSSHOptions(testUsers)), "wrong number of options returned by getSSHOptions for HostKeyFile \"non-existent\"")

	HostKeyFile = path.Join("_host_key_test_files", "id_rsa")
	assert.Equal(t, 3, len(getSSHOptions(testUsers)), "wrong number of options returned by getSSHOptions for HostKeyFile \"id_rsa\"")

	HostKeyFile = oldHostKeyFile
}
