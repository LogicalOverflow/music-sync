package comm

import (
	"net"
	"sync"

	"github.com/LogicalOverflow/music-sync/util"

	"github.com/golang/protobuf/proto"
)

type MessageSender interface {
	SendMessage(m proto.Message) error
}

type multiMessageSender struct {
	connections []net.Conn
	mutex       sync.RWMutex
}

func (mms *multiMessageSender) SendMessage(m proto.Message) error {
	mms.mutex.RLock()
	defer mms.mutex.RUnlock()

	var errCol util.ErrorCollector
	var wg sync.WaitGroup
	wg.Add(len(mms.connections))

	for _, c := range mms.connections {
		go func(c net.Conn) {
			defer wg.Done()
			if err := sendWire(m, c); err != nil {
				errCol.Add(err)
			}
		}(c)
	}

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

type singleMessageSender struct{ connection net.Conn }

func (sms *singleMessageSender) SendMessage(m proto.Message) error {
	return sendWire(m, sms.connection)
}
