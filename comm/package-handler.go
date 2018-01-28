package comm

import (
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/timing"
	"github.com/golang/protobuf/proto"
	"net"
)

type typedPackageHandler struct {
	typedPackageHandlerInterface
}

type typedPackageHandlerInterface interface {
	HandleTimeSyncRequest(*TimeSyncRequest, net.Conn)
	HandleTimeSyncResponse(*TimeSyncResponse, net.Conn)
	HandleQueueChunkRequest(*QueueChunkRequest, net.Conn)
	HandlePingMessage(*PingMessage, net.Conn)
	HandlePongMessage(*PongMessage, net.Conn)
	HandleSetVolumeRequest(*SetVolumeRequest, net.Conn)
	HandleSubscribeChannelRequest(*SubscribeChannelRequest, net.Conn)
}

func (t typedPackageHandler) Handle(message proto.Message, sender net.Conn) {
	switch message.(type) {
	case *TimeSyncRequest:
		go t.HandleTimeSyncRequest(message.(*TimeSyncRequest), sender)
	case *TimeSyncResponse:
		go t.HandleTimeSyncResponse(message.(*TimeSyncResponse), sender)
	case *QueueChunkRequest:
		go t.HandleQueueChunkRequest(message.(*QueueChunkRequest), sender)
	case *PingMessage:
		go t.HandlePingMessage(message.(*PingMessage), sender)
	case *PongMessage:
		go t.HandlePongMessage(message.(*PongMessage), sender)
	case *SetVolumeRequest:
		go t.HandleSetVolumeRequest(message.(*SetVolumeRequest), sender)
	case *SubscribeChannelRequest:
		go t.HandleSubscribeChannelRequest(message.(*SubscribeChannelRequest), sender)
	}
}

type serverPackageHandler struct {
	sender *multiMessageSender
}

func (s serverPackageHandler) HandleTimeSyncRequest(tsr *TimeSyncRequest, c net.Conn) {
	serverRecv := timing.GetRawTime()
	response := &TimeSyncResponse{ClientSendTime: tsr.ClientSend, ServerRecvTime: serverRecv, ServerSendTime: timing.GetRawTime()}
	if err := sendWire(response, c); err != nil {
		logger.Warnf("failed to send handle time sync response: %v", err)
	}
}

func (s serverPackageHandler) HandlePingMessage(_ *PingMessage, c net.Conn) {
	if err := sendWire(&PongMessage{}, c); err != nil {
		logger.Warnf("failed to send ping response: %v", err)
	}
}
func (s serverPackageHandler) HandleSubscribeChannelRequest(scr *SubscribeChannelRequest, c net.Conn) {
	s.sender.Subscribe(c, scr.Channel)
	NewClientHandler(scr.Channel, &singleMessageSender{c})
}

func (s serverPackageHandler) HandleTimeSyncResponse(*TimeSyncResponse, net.Conn)   {}
func (s serverPackageHandler) HandleQueueChunkRequest(*QueueChunkRequest, net.Conn) {}
func (s serverPackageHandler) HandlePongMessage(*PongMessage, net.Conn)             {}
func (s serverPackageHandler) HandleSetVolumeRequest(*SetVolumeRequest, net.Conn)   {}

func newMasterPackageHandler(sender *multiMessageSender) typedPackageHandler {
	return typedPackageHandler{serverPackageHandler{sender: sender}}
}

type clientPackageHandler struct{}

func (c clientPackageHandler) HandleTimeSyncResponse(tsr *TimeSyncResponse, _ net.Conn) {
	clientRecv := timing.GetRawTime()
	timing.UpdateOffset(tsr.ClientSendTime, tsr.ServerRecvTime, tsr.ServerSendTime, clientRecv)
}

func (c clientPackageHandler) HandleQueueChunkRequest(qsr *QueueChunkRequest, _ net.Conn) { playback.QueueChunk(qsr.StartTime, qsr.ChunkId, playback.CombineSamples(qsr.SampleLow, qsr.SampleHigh)) }
func (c clientPackageHandler) HandleSetVolumeRequest(svr *SetVolumeRequest, _ net.Conn)   { playback.SetVolume(svr.Volume) }

func (c clientPackageHandler) HandlePingMessage(_ *PingMessage, conn net.Conn) {
	if err := sendWire(&PongMessage{}, conn); err != nil {
		logger.Warnf("failed to send ping response: %v", err)
	}
}

func (c clientPackageHandler) HandleTimeSyncRequest(*TimeSyncRequest, net.Conn)                 {}
func (c clientPackageHandler) HandlePongMessage(*PongMessage, net.Conn)                         {}
func (c clientPackageHandler) HandleSubscribeChannelRequest(*SubscribeChannelRequest, net.Conn) {}

func newSlavePackageHandler() typedPackageHandler {
	return typedPackageHandler{clientPackageHandler{}}
}
