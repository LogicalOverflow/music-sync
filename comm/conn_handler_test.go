package comm

import (
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
	"time"
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

const invalidLastChannel = -32

func TestServerConnectionAcceptor(t *testing.T) {
	fl := newFakeListener()

	var lastChan Channel = invalidLastChannel
	NewClientHandler = func(channel Channel, conn MessageSender) { lastChan = channel }

	var wg sync.WaitGroup

	wg.Add(2)

	mms := &multiMessageSender{connections: make([]net.Conn, 0), channels: make(map[net.Conn][]Channel)}

	go func() {
		defer wg.Done()
		serverConnectionAcceptor(mms, fl)
	}()

	testConnsServer, testConnsClient := newPipeConnPairs(8)

	go func() {
		defer wg.Done()
		for i := range testConnsServer {
			fl.NewConn(testConnsServer[i])
			time.Sleep(10 * time.Millisecond)
			assert.True(t, containsConn(mms.connections, testConnsServer[i]), "serverConnectionAcceptor did not AddConn (new) to mms (%d conns in mms)", len(mms.connections))

			assert.Equal(t, Channel(-1), lastChan, "serverConnectionAcceptor did not call NewClientHandler with the correct channel")
			lastChan = invalidLastChannel

			testConnsClient[i].Close()
			time.Sleep(10 * time.Millisecond)
			assert.False(t, containsConn(mms.connections, testConnsServer[i]), "serverConnectionAcceptor did not DelConn (closed) from mms (%d conns in mms)", len(mms.connections))
		}
		fl.Close()
	}()

	wg.Wait()
}

func containsConn(cs []net.Conn, c net.Conn) bool {
	for _, co := range cs {
		if co == c {
			return true
		}
	}
	return false
}
