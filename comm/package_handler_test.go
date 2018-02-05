package comm

import (
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/assert"
	"net"
	"sync"
	"testing"
)

type testTypedPackageHandler struct {
	lastPackage proto.Message
	lastType    string
	cond        sync.Cond
}

func (t *testTypedPackageHandler) HandleTimeSyncRequest(p *TimeSyncRequest, _ net.Conn) {
	t.lastPackage = p
	t.lastType = "TimeSyncRequest"
	t.cond.Broadcast()
}
func (t *testTypedPackageHandler) HandleTimeSyncResponse(p *TimeSyncResponse, _ net.Conn) {
	t.lastPackage = p
	t.lastType = "TimeSyncResponse"
	t.cond.Broadcast()
}
func (t *testTypedPackageHandler) HandleQueueChunkRequest(p *QueueChunkRequest, _ net.Conn) {
	t.lastPackage = p
	t.lastType = "QueueChunkRequest"
	t.cond.Broadcast()
}
func (t *testTypedPackageHandler) HandlePingMessage(p *PingMessage, _ net.Conn) {
	t.lastPackage = p
	t.lastType = "PingMessage"
	t.cond.Broadcast()
}
func (t *testTypedPackageHandler) HandlePongMessage(p *PongMessage, _ net.Conn) {
	t.lastPackage = p
	t.lastType = "PongMessage"
	t.cond.Broadcast()
}
func (t *testTypedPackageHandler) HandleSetVolumeRequest(p *SetVolumeRequest, _ net.Conn) {
	t.lastPackage = p
	t.lastType = "SetVolumeRequest"
	t.cond.Broadcast()
}
func (t *testTypedPackageHandler) HandleSubscribeChannelRequest(p *SubscribeChannelRequest, _ net.Conn) {
	t.lastPackage = p
	t.lastType = "SubscribeChannelRequest"
	t.cond.Broadcast()
}
func (t *testTypedPackageHandler) HandleNewSongInfo(p *NewSongInfo, _ net.Conn) {
	t.lastPackage = p
	t.lastType = "NewSongInfo"
	t.cond.Broadcast()
}
func (t *testTypedPackageHandler) HandleChunkInfo(p *ChunkInfo, _ net.Conn) {
	t.lastPackage = p
	t.lastType = "ChunkInfo"
	t.cond.Broadcast()
}
func (t *testTypedPackageHandler) HandlePauseInfo(p *PauseInfo, _ net.Conn) {
	t.lastPackage = p
	t.lastType = "PauseInfo"
	t.cond.Broadcast()
}

var typedPackageHandlerHandleCases = []struct {
	pType string
	p     proto.Message
}{
	{pType: "TimeSyncRequest", p: &TimeSyncRequest{ClientSend: 1}},
	{pType: "TimeSyncResponse", p: &TimeSyncResponse{ClientSendTime: 1, ServerRecvTime: 2, ServerSendTime: 3}},
	{pType: "QueueChunkRequest", p: &QueueChunkRequest{StartTime: 1, ChunkId: 2, FirstSampleIndex: 3}},
	{pType: "PingMessage", p: &PingMessage{}},
	{pType: "PongMessage", p: &PongMessage{}},
	{pType: "SetVolumeRequest", p: &SetVolumeRequest{Volume: 1.2}},
	{pType: "SubscribeChannelRequest", p: &SubscribeChannelRequest{Channel: Channel_AUDIO}},
	{pType: "NewSongInfo", p: &NewSongInfo{FirstSampleOfSongIndex: 1, SongFileName: "abc", SongLength: 2}},
	{pType: "ChunkInfo", p: &ChunkInfo{StartTime: 1, FirstSampleIndex: 2, ChunkSize: 3}},
	{pType: "PauseInfo", p: &PauseInfo{Playing: true, ToggleSampleIndex: 2}},
}

func TestTypedPackageHandler_Handle(t *testing.T) {
	ttph := &testTypedPackageHandler{cond: sync.Cond{L: new(sync.Mutex)}}
	ph := TypedPackageHandler{TypedPackageHandlerInterface: ttph}

	for _, c := range typedPackageHandlerHandleCases {
		done := make(chan bool)
		go func() {
			ttph.cond.L.Lock()
			done <- true
			ttph.cond.Wait()
			done <- true
			ttph.cond.L.Unlock()
		}()
		<-done
		ph.Handle(c.p, nil)
		<-done
		assert.Equal(t, c.pType, ttph.lastType, "handling package of type %s called handler for %s through TypedPackageHandler", c.pType, ttph.lastType)
		assert.Equal(t, c.p, ttph.lastPackage, "handling package of type %s called handler with wrong package through TypedPackageHandler", c.pType)
	}
}
