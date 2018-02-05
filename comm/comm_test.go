package comm

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

var testPackages = []proto.Message{
	&QueueChunkRequest{
		StartTime:        1,
		ChunkId:          2,
		SampleLow:        make([]float64, 2),
		SampleHigh:       make([]float64, 2),
		FirstSampleIndex: 3,
	},
	&PauseInfo{
		Playing:           true,
		ToggleSampleIndex: 1,
	},
	&TimeSyncResponse{
		ClientSendTime: 1,
		ServerRecvTime: 2,
		ServerSendTime: 3,
	},
	&SetVolumeRequest{
		Volume: 1.2,
	},
}

var testPackageChannels = [][]Channel{{Channel_AUDIO}, {Channel_META}, {}, {Channel_AUDIO, Channel_META}}

type bufferConn struct {
	*bytes.Buffer
	name string
}

func (bufferConn) LocalAddr() net.Addr                { return nil }
func (bufferConn) RemoteAddr() net.Addr               { return nil }
func (bufferConn) SetDeadline(t time.Time) error      { return nil }
func (bufferConn) SetReadDeadline(t time.Time) error  { return nil }
func (bufferConn) SetWriteDeadline(t time.Time) error { return nil }
func (bufferConn) Close() error                       { return nil }

func (b bufferConn) assertData(t *testing.T, expected []byte, shouldSend bool, p proto.Message) {
	if shouldSend {
		assert.True(t, bytes.Equal(expected, b.Bytes()), "multiMessageSender sendMessage did not write toWire to the connection %s for package %v", b.name, p)
	} else {
		assert.Zero(t, len(b.Bytes()), "multiMessageSender sendMessage did write to the connection %s for package %v", b.name, p)
	}
}

func newBufferConn() bufferConn {
	return bufferConn{Buffer: bytes.NewBuffer([]byte{})}
}

func newNamedBufferConn(name string) bufferConn {
	return bufferConn{Buffer: bytes.NewBuffer([]byte{}), name: name}
}

func newBufferConnWithData(data []byte) bufferConn {
	return bufferConn{Buffer: bytes.NewBuffer(data)}
}
