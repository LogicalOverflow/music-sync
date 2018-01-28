package comm

import (
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/timing"
	"github.com/golang/protobuf/proto"
	"net"
)

// Typed package handler calls the handle functions of TypedPackageHandlerInterface
type TypedPackageHandler struct {
	TypedPackageHandlerInterface
}

// TypedPackageHandlerInterface has methods to handle all packages received
type TypedPackageHandlerInterface interface {
	HandleTimeSyncRequest(*TimeSyncRequest, net.Conn)
	HandleTimeSyncResponse(*TimeSyncResponse, net.Conn)
	HandleQueueChunkRequest(*QueueChunkRequest, net.Conn)
	HandlePingMessage(*PingMessage, net.Conn)
	HandlePongMessage(*PongMessage, net.Conn)
	HandleSetVolumeRequest(*SetVolumeRequest, net.Conn)
	HandleSubscribeChannelRequest(*SubscribeChannelRequest, net.Conn)
	HandleNewSongInfo(*NewSongInfo, net.Conn)
	HandleChunkInfo(*ChunkInfo, net.Conn)
}

func (t TypedPackageHandler) Handle(message proto.Message, sender net.Conn) {
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
	case *NewSongInfo:
		go t.HandleNewSongInfo(message.(*NewSongInfo), sender)
	case *ChunkInfo:
		go t.HandleChunkInfo(message.(*ChunkInfo), sender)
	}
}

// TODO: move this is the server cmd

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

func (s serverPackageHandler) HandleSubscribeChannelRequest(scr *SubscribeChannelRequest, c net.Conn) {
	s.sender.Subscribe(c, scr.Channel)
	NewClientHandler(scr.Channel, &singleMessageSender{c})
}

func (s serverPackageHandler) HandlePingMessage(_ *PingMessage, c net.Conn) { PingHandler(c) }

func (s serverPackageHandler) HandleTimeSyncResponse(*TimeSyncResponse, net.Conn)   {}
func (s serverPackageHandler) HandleQueueChunkRequest(*QueueChunkRequest, net.Conn) {}
func (s serverPackageHandler) HandlePongMessage(*PongMessage, net.Conn)             {}
func (s serverPackageHandler) HandleSetVolumeRequest(*SetVolumeRequest, net.Conn)   {}
func (s serverPackageHandler) HandleNewSongInfo(*NewSongInfo, net.Conn)             {}
func (s serverPackageHandler) HandleChunkInfo(*ChunkInfo, net.Conn)                 {}

func newServerPackageHandler(sender *multiMessageSender) TypedPackageHandler {
	return TypedPackageHandler{serverPackageHandler{sender: sender}}
}

// TODO: move this in the player cmd

type playerPackageHandler struct{}

func (c playerPackageHandler) HandleTimeSyncResponse(tsr *TimeSyncResponse, _ net.Conn) {
	clientRecv := timing.GetRawTime()
	timing.UpdateOffset(tsr.ClientSendTime, tsr.ServerRecvTime, tsr.ServerSendTime, clientRecv)
}

func (c playerPackageHandler) HandleQueueChunkRequest(qsr *QueueChunkRequest, _ net.Conn) { playback.QueueChunk(qsr.StartTime, qsr.ChunkId, playback.CombineSamples(qsr.SampleLow, qsr.SampleHigh)) }
func (c playerPackageHandler) HandleSetVolumeRequest(svr *SetVolumeRequest, _ net.Conn)   { playback.SetVolume(svr.Volume) }
func (c playerPackageHandler) HandlePingMessage(_ *PingMessage, conn net.Conn)            { PingHandler(conn) }

func (c playerPackageHandler) HandleTimeSyncRequest(*TimeSyncRequest, net.Conn)                 {}
func (c playerPackageHandler) HandlePongMessage(*PongMessage, net.Conn)                         {}
func (c playerPackageHandler) HandleSubscribeChannelRequest(*SubscribeChannelRequest, net.Conn) {}
func (c playerPackageHandler) HandleNewSongInfo(*NewSongInfo, net.Conn)                         {}
func (c playerPackageHandler) HandleChunkInfo(*ChunkInfo, net.Conn)                             {}

// NewPlayerPackageHandler returns the TypedPackageHandler used by players
func NewPlayerPackageHandler() TypedPackageHandler {
	return TypedPackageHandler{playerPackageHandler{}}
}

// PingHandler handle a PingMessage
func PingHandler(conn net.Conn) {
	if err := sendWire(&PongMessage{}, conn); err != nil {
		logger.Warnf("failed to send ping response: %v", err)
	}
}
