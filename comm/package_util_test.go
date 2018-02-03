package comm

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestWireConversation(t *testing.T) {
	for _, tp := range testPackages {
		wire, err := toWire(tp)
		if assert.Nil(t, err, "toWire returned an error for package %v: %v", tp, err) {
			result, err := readWire(newBufferConnWithData(wire))
			if assert.Nil(t, err, "readWire returned an error for package %v: %v", tp, err) {
				assert.Equal(t, tp.String(), result.String(), "readWire returned a different package then toWire created")
			}
		}

	}
}
