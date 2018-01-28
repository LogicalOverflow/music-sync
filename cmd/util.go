package cmd

import (
	"os"
	"os/signal"
)

// WaitForInterrupt blocks until os.Interrupt
func WaitForInterrupt() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for range c {
		return
	}
}
