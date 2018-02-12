package ssh

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestReadUsersFile(t *testing.T) {
	_, err := ReadUsersFile("non-existent.json")
	assert.NotNil(t, err, "ReadUsersFile did not return an error for a non-existent users file")

	_, err = ReadUsersFile("_test_users_file_error.json")
	assert.NotNil(t, err, "ReadUsersFile did not return an error for an invalid users file")

	u, err := ReadUsersFile("_test_users_file.json")
	assert.Nil(t, err, "ReadUsersFile returned an error for a valid users file")
	assert.Equal(t, map[string]UserAuth{"test-user": {Password: "test-password"}}, u, "ReadUsersFile returned the wrong users")
}
