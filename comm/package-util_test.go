package comm

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"net"
	"testing"
	"time"
)

type fakeConn struct {
	*bytes.Buffer
}

func (fakeConn) LocalAddr() net.Addr                { return nil }
func (fakeConn) RemoteAddr() net.Addr               { return nil }
func (fakeConn) SetDeadline(t time.Time) error      { return nil }
func (fakeConn) SetReadDeadline(t time.Time) error  { return nil }
func (fakeConn) SetWriteDeadline(t time.Time) error { return nil }
func (fakeConn) Close() error                       { return nil }

func TestWireConversation(t *testing.T) {
	testPackages := []proto.Message{
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
	}
	for _, tp := range testPackages {
		wire, err := toWire(tp)
		if assert.Nil(t, err, "toWire returned an error for package %v: %v", tp, err) {
			result, err := readWire(fakeConn{Buffer: bytes.NewBuffer(wire)})
			if assert.Nil(t, err, "readWire returned an error for package %v: %v", tp, err) {
				assert.Equal(t, tp.String(), result.String(), "readWire returned a different package then toWire created")
			}
		}

	}
}
