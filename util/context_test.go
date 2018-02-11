package util

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsCanceled(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	assert.False(t, IsCanceled(ctx), "IsCanceled returned true before cancel was called")
	cancel()
	assert.True(t, IsCanceled(ctx), "IsCanceled returned false after cancel was called")
}
