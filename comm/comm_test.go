package comm

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"io"
	"net"
	"sync"
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

type testPackageHandler struct {
	packages      []proto.Message
	packagesMutex sync.RWMutex
	packagesCond  *sync.Cond
}

func (tph *testPackageHandler) Handle(message proto.Message, _ net.Conn) {
	tph.packagesMutex.Lock()
	defer tph.packagesMutex.Unlock()
	if tph.packages == nil {
		tph.packages = []proto.Message{}
	}

	tph.packages = append(tph.packages, message)

	if tph.packagesCond != nil {
		tph.packagesCond.Broadcast()
	}
}

func (tph *testPackageHandler) Packages() []proto.Message {
	tph.packagesMutex.RLock()
	defer tph.packagesMutex.RUnlock()
	if tph.packages == nil || len(tph.packages) == 0 {
		return []proto.Message{}
	}

	r := make([]proto.Message, len(tph.packages))
	copy(r, tph.packages)
	return r
}

func (tph *testPackageHandler) Latest() proto.Message {
	tph.packagesMutex.RLock()
	defer tph.packagesMutex.RUnlock()

	if tph.packages == nil || len(tph.packages) == 0 {
		return nil
	}

	return tph.packages[len(tph.packages)-1]
}

func (tph *testPackageHandler) WaitForPackage() proto.Message {
	if tph.packagesCond == nil {
		tph.packagesCond = sync.NewCond(new(sync.Mutex))
	}

	tph.packagesCond.L.Lock()
	tph.packagesCond.Wait()
	tph.packagesCond.L.Unlock()
	return tph.Latest()
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

type pipeConn struct {
	r *io.PipeReader
	w *io.PipeWriter
}

func (p *pipeConn) Read(b []byte) (n int, err error)   { return p.r.Read(b) }
func (p *pipeConn) Write(b []byte) (n int, err error)  { return p.w.Write(b) }
func (p *pipeConn) LocalAddr() net.Addr                { return nil }
func (p *pipeConn) RemoteAddr() net.Addr               { return nil }
func (p *pipeConn) SetDeadline(t time.Time) error      { return nil }
func (p *pipeConn) SetReadDeadline(t time.Time) error  { return nil }
func (p *pipeConn) SetWriteDeadline(t time.Time) error { return nil }
func (p *pipeConn) Close() error                       { p.r.Close(); p.w.Close(); return nil }

func newPipeConnPair() (*pipeConn, *pipeConn) {
	r1, w1 := io.Pipe()
	r2, w2 := io.Pipe()
	return &pipeConn{r: r1, w: w2}, &pipeConn{r: r2, w: w1}
}
