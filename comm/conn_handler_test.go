package comm

import (
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func TestHandleConnection(t *testing.T) {
	log.DefaultCutoffLevel = log.LevelOff
	ch, cr := newPipeConnPair()

	tph := new(testPackageHandler)

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		defer wg.Done()
		handleConnection(ch, tph)
	}()

	comm := make(chan bool, 1)

	go func() {
		defer wg.Done()
		for _, p := range testPackages {
			comm <- true
			result := tph.WaitForPackage()
			assert.Equal(t, p, result, "sending packages %v did not call package handler", p)
		}
	}()

	go func() {
		defer wg.Done()
		for _, p := range testPackages {
			<-comm
			sendWire(p, cr)
		}
		cr.Close()
	}()

	wg.Wait()
}
