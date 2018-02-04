package comm

import (
	"github.com/LogicalOverflow/music-sync/timing"
	"github.com/golang/protobuf/proto"
	"net"
)

// TypedPackageHandler calls the handle functions of TypedPackageHandlerInterface
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
	HandlePauseInfo(*PauseInfo, net.Conn)
}

// Handle forwards the message and sender to the matching Handle function of TypedPackageHandlerInterface
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
	case *PauseInfo:
		go t.HandlePauseInfo(message.(*PauseInfo), sender)
	}
}

// TODO: move this is the server cmd?

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
func (s serverPackageHandler) HandlePauseInfo(*PauseInfo, net.Conn)                 {}

func newServerPackageHandler(sender *multiMessageSender) TypedPackageHandler {
	return TypedPackageHandler{serverPackageHandler{sender: sender}}
}

// PingHandler handle a PingMessage
func PingHandler(conn net.Conn) {
	if err := sendWire(&PongMessage{}, conn); err != nil {
		logger.Warnf("failed to send ping response: %v", err)
	}
}
