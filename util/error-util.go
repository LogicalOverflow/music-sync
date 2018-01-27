package util

import (
	"fmt"
	"strings"
)

type MultiError struct {
	BaseMessage string
	Errors      []error
}

func (m MultiError) Error() string {
	msgs := make([]string, len(m.Errors))
	for i, err := range m.Errors {
		msgs[i] = fmt.Sprintf("%d: %v", i, err)
	}
	return fmt.Sprintf("%s: %s", m.BaseMessage, strings.Join(msgs, "; "))
}

func NewMultiError(baseMessage string, errors []error) MultiError {
	return MultiError{BaseMessage: baseMessage, Errors: errors}
}
