package cmd

import (
	"github.com/stretchr/testify/require"
	"os"
	"syscall"
	"testing"
)

func TestWaitForInterrupt(t *testing.T) {
	comm := make(chan bool, 1)

	go func() {
		WaitForInterrupt()
		comm <- true
	}()

	select {
	case <-comm:
		require.Fail(t, "WaitForInterrupt did not block")
	default:
		break
	}

	p, err := os.FindProcess(syscall.Getpid())
	require.Nil(t, err, "failed to find process to send interrupt signal")
	require.Nil(t, p.Signal(os.Interrupt), "failed to send interrupt signal")

	select {
	case <-comm:
		break
	default:
		require.Fail(t, "WaitForInterrupt blocks after sending interrupt")
	}
	<-comm
}
