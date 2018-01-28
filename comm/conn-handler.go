package comm

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"net"
	"os"
)

type PackageHandler interface {
	Handle(message proto.Message, sender net.Conn)
}

func HandleConnection(conn net.Conn, h PackageHandler) {
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

var NewSlaveHandler func(channel Channel, conn MessageSender)

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
		h := NewMasterPackageHandler(mms)
		for {
			conn, err := l.Accept()
			if err != nil {
				logger.Warnf("failed to accept connection: %v", err)
			}
			go func(conn net.Conn) {
				mms.AddConn(conn)
				HandleConnection(conn, h)
				mms.DelConn(conn)
			}(conn)
			if NewSlaveHandler != nil {
				go NewSlaveHandler(-1, &singleMessageSender{conn})
			}
		}
	}()

	return mms, nil
}

func ConnectToMaster(master string) (MessageSender, error) {
	logger.Infof("connecting to master at %s", master)
	conn, err := net.Dial("tcp", master)
	if err != nil {
		logger.Fatalf("could not connect to master at %s: %v", master, err)
		return nil, fmt.Errorf("could not connect to master at %s: %v", master, err)
	}
	logger.Infof("connected to master at %s", master)
	go func() {
		HandleConnection(conn, NewSlavePackageHandler())
		logger.Fatalf("connection to master closed")
		os.Exit(1)
	}()
	return &singleMessageSender{connection: conn}, nil
}
