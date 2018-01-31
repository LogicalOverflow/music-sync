package comm

import (
	"net"
	"sync"

	"github.com/LogicalOverflow/music-sync/util"

	"github.com/golang/protobuf/proto"
)

// MessageSender is an interface used to send protobuf messages to clients/the server.
type MessageSender interface {
	// Sends a message to the server/all clients.
	// This method respects subscribe message from clients.
	SendMessage(m proto.Message) error
}

type multiMessageSender struct {
	connections []net.Conn
	channels    map[net.Conn][]Channel
	mutex       sync.RWMutex
}

func (mms *multiMessageSender) SendMessage(m proto.Message) error {
	mms.mutex.RLock()
	defer mms.mutex.RUnlock()

	chs, hasCh := channelOf(m)

	var errCol util.ErrorCollector
	var wg sync.WaitGroup
	wg.Add(len(mms.connections))

	for _, c := range mms.connections {
		go func(c net.Conn) {
			defer wg.Done()
			if !hasCh || mms.isSubscribed(c, chs) {
				if err := sendWire(m, c); err != nil {
					errCol.Add(err)
				}
			}
		}(c)
	}

	wg.Wait()
	errCol.Wait()

	return errCol.Err("failed to send to %d clients: ")
}

func (mms *multiMessageSender) AddConn(c net.Conn) {
	mms.mutex.Lock()
	defer mms.mutex.Unlock()
	mms.connections = append(mms.connections, c)
}

func (mms *multiMessageSender) DelConn(c net.Conn) {
	mms.mutex.Lock()
	defer mms.mutex.Unlock()

	index := -1
	for i, conn := range mms.connections {
		if conn == c {
			index = i
			break
		}
	}

	if 0 <= index {
		mms.connections[index] = mms.connections[len(mms.connections)-1]
		mms.connections = mms.connections[:len(mms.connections)-1]
	}
}

func (mms *multiMessageSender) Subscribe(c net.Conn, channel Channel) {
	mms.mutex.Lock()
	defer mms.mutex.Unlock()
	if l, ok := mms.channels[c]; ok {
		mms.channels[c] = append(l, channel)
	} else {
		mms.channels[c] = []Channel{channel}
	}
}

func (mms *multiMessageSender) isSubscribed(c net.Conn, channels []Channel) bool {
	if l, ok := mms.channels[c]; ok {
		for _, ch := range l {
			for _, c := range channels {
				if ch == c {
					return true
				}
			}
		}
	}
	return false
}

type singleMessageSender struct{ connection net.Conn }

func (sms *singleMessageSender) SendMessage(m proto.Message) error {
	return sendWire(m, sms.connection)
}
