package comm

import (
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/timing"
	"github.com/golang/protobuf/proto"
	"net"
)

type TypedPackageHandler struct {
	TypedPackageHandlerInterface
}

type TypedPackageHandlerInterface interface {
	HandleTimeSyncRequest(*TimeSyncRequest, net.Conn)
	HandleTimeSyncResponse(*TimeSyncResponse, net.Conn)
	HandleQueueSongRequest(*QueueSongRequest, net.Conn)
	HandlePingMessage(*PingMessage, net.Conn)
	HandlePongMessage(*PongMessage, net.Conn)
	HandleSetVolumeRequest(*SetVolumeRequest, net.Conn)
}

func (t TypedPackageHandler) Handle(message proto.Message, sender net.Conn) {
	switch message.(type) {
	case *TimeSyncRequest:
		go t.HandleTimeSyncRequest(message.(*TimeSyncRequest), sender)
	case *TimeSyncResponse:
		go t.HandleTimeSyncResponse(message.(*TimeSyncResponse), sender)
	case *QueueSongRequest:
		go t.HandleQueueSongRequest(message.(*QueueSongRequest), sender)
	case *PingMessage:
		go t.HandlePingMessage(message.(*PingMessage), sender)
	case *PongMessage:
		go t.HandlePongMessage(message.(*PongMessage), sender)
	case *SetVolumeRequest:
		go t.HandleSetVolumeRequest(message.(*SetVolumeRequest), sender)
	}
}

type MasterPackageHandler struct{}

func (m MasterPackageHandler) HandleTimeSyncRequest(tsr *TimeSyncRequest, c net.Conn) {
	serverRecv := timing.GetRawTime()
	response := &TimeSyncResponse{ClientSendTime: tsr.ClientSend, ServerRecvTime: serverRecv, ServerSendTime: timing.GetRawTime()}
	if err := sendWire(response, c); err != nil {
		logger.Warnf("failed to send handle time sync response: %v", err)
	}
}

func (m MasterPackageHandler) HandlePingMessage(_ *PingMessage, c net.Conn) {
	if err := sendWire(&PongMessage{}, c); err != nil {
		logger.Warnf("failed to send ping response: %v", err)
	}
}

func (m MasterPackageHandler) HandleTimeSyncResponse(*TimeSyncResponse, net.Conn) {}
func (m MasterPackageHandler) HandleQueueSongRequest(*QueueSongRequest, net.Conn) {}
func (m MasterPackageHandler) HandlePongMessage(*PongMessage, net.Conn)           {}
func (m MasterPackageHandler) HandleSetVolumeRequest(*SetVolumeRequest, net.Conn) {}

func NewMasterPackageHandler() TypedPackageHandler {
	return TypedPackageHandler{MasterPackageHandler{}}
}

type SlavePackageHandler struct{}

func (s SlavePackageHandler) HandleTimeSyncRequest(*TimeSyncRequest, net.Conn) {}
func (s SlavePackageHandler) HandlePongMessage(*PongMessage, net.Conn)         {}

func (s SlavePackageHandler) HandleTimeSyncResponse(tsr *TimeSyncResponse, _ net.Conn) {
	clientRecv := timing.GetRawTime()
	timing.UpdateOffset(tsr.ClientSendTime, tsr.ServerRecvTime, tsr.ServerSendTime, clientRecv)
}

func (s SlavePackageHandler) HandleQueueSongRequest(qsr *QueueSongRequest, _ net.Conn) { playback.QueueSong(qsr.StartTime, qsr.ChunkId, playback.CombineSamples(qsr.SampleLow, qsr.SampleHigh)) }
func (s SlavePackageHandler) HandleSetVolumeRequest(svr *SetVolumeRequest, _ net.Conn) { playback.SetVolume(svr.Volume) }

func (s SlavePackageHandler) HandlePingMessage(_ *PingMessage, c net.Conn) {
	if err := sendWire(&PongMessage{}, c); err != nil {
		logger.Warnf("failed to send ping response: %v", err)
	}
}

func NewSlavePackageHandler() TypedPackageHandler {
	return TypedPackageHandler{SlavePackageHandler{}}
}
