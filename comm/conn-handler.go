// Package comm contains functions and types from communication between music-sync clients and the server
package comm

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
)

type packageHandler interface {
	Handle(message proto.Message, sender net.Conn)
}

func handleConnection(conn net.Conn, h packageHandler) {
	defer conn.Close()
	for {
		m, err := readWire(conn)
		if err != nil {
			logger.Infof("failed to read data from %s: %v. closing connection", conn.RemoteAddr(), err)
			break
		}
		h.Handle(m, conn)
	}
}

// NewClientHandler is called when a new client connects to the server (with channel -1)
// and when a client subscribes to a channel.
var NewClientHandler func(channel Channel, conn MessageSender)

// StartServer starts a music-sync server listening at address and returns a MessageSender to broadcast
// to clients
func StartServer(address string) (MessageSender, error) {
	logger.Infof("starting server at %s", address)

	l, err := net.Listen("tcp", address)
	if err != nil {
		logger.Fatalf("failed to start server at %s: %v", address, err)
		return nil, fmt.Errorf("failed to start server at %s: %v", address, err)
	}
	logger.Infof("server running at %s", address)
	mms := &multiMessageSender{connections: make([]net.Conn, 0), channels: make(map[net.Conn][]Channel)}
	go func() {
		h := newServerPackageHandler(mms)
		for {
			conn, err := l.Accept()
			if err != nil {
				logger.Warnf("failed to accept connection: %v", err)
			}
			go func(conn net.Conn) {
				mms.AddConn(conn)
				handleConnection(conn, h)
				mms.DelConn(conn)
			}(conn)
			if NewClientHandler != nil {
				go NewClientHandler(-1, &singleMessageSender{conn})
			}
		}
	}()

	return mms, nil
}

// ConnectToServer connects to the server at server and returns a MessageSender to communicate with the master
func ConnectToServer(master string, handler TypedPackageHandler) (MessageSender, error) {
	logger.Infof("connecting to master at %s", master)
	conn, err := net.Dial("tcp", master)
	if err != nil {
		logger.Fatalf("could not connect to master at %s: %v", master, err)
		return nil, fmt.Errorf("could not connect to master at %s: %v", master, err)
	}
	logger.Infof("connected to master at %s", master)
	go func() {
		handleConnection(conn, handler)
		logger.Fatalf("connection to master closed")
		os.Exit(1)
	}()
	return &singleMessageSender{connection: conn}, nil
}
