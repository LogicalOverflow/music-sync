package util

import (
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNewMultiError(t *testing.T) {
	msg := "base message"
	errs := []error{errors.New("msg-1"), errors.New("msg-2")}
	mr := NewMultiError(msg, errs)
	assert.Equal(t, msg, mr.BaseMessage, "MultiError has wrong BaseMessage")
	assert.Equal(t, len(errs), len(mr.Errors), "MultiError has wrong number of Errors")

	for i := 0; i < len(errs); i++ {
		assert.Equal(t, errs[i], mr.Errors[i], "MultiError has wrong error at index %d", i)
	}
}

func TestMultiError_Error(t *testing.T) {
	msg := "base message"
	errs := []error{errors.New("msg-1"), errors.New("msg-2")}
	mr := NewMultiError(msg, errs)

	err := mr.Error()
	assert.Equal(t, "base message: 0: msg-1; 1: msg-2", err, "MultiError has wrong Error() message")
}