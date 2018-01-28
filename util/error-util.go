package util

import (
	"fmt"
	"strings"
)

// MultiError is an error containing multiple errors
type MultiError struct {
	BaseMessage string
	Errors      []error
}

// Error implements the error interface
func (m MultiError) Error() string {
	msgs := make([]string, len(m.Errors))
	for i, err := range m.Errors {
		msgs[i] = fmt.Sprintf("%d: %v", i, err)
	}
	return fmt.Sprintf("%s: %s", m.BaseMessage, strings.Join(msgs, "; "))
}

// NewMultiError creates a MultiError from a base message and a slice of errors
func NewMultiError(baseMessage string, errors []error) MultiError {
	return MultiError{BaseMessage: baseMessage, Errors: errors}
}
