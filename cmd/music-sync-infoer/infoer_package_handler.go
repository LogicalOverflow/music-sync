package main

import (
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/metadata"
	"github.com/LogicalOverflow/music-sync/timing"
	"net"
	"sort"
)

type infoerPackageHandler struct {
}

func (i *infoerPackageHandler) HandleTimeSyncResponse(tsr *comm.TimeSyncResponse, _ net.Conn) {
	clientRecv := timing.GetRawTime()
	timing.UpdateOffset(tsr.ClientSendTime, tsr.ServerRecvTime, tsr.ServerSendTime, clientRecv)
}
func (i *infoerPackageHandler) HandleNewSongInfo(newSongInfo *comm.NewSongInfo, _ net.Conn) {
	currentState.SongsMutex.Lock()
	defer currentState.SongsMutex.Unlock()
	lyrics := make([]metadata.LyricsLine, len(newSongInfo.Lyrics))
	for i, l := range newSongInfo.Lyrics {
		atoms := make([]metadata.LyricsAtom, len(l.Atoms))
		for j, a := range l.Atoms {
			atoms[j] = metadata.LyricsAtom{Timestamp: a.Timestamp, Caption: a.Caption}
		}
		lyrics[i] = atoms
	}

	md := metadata.SongMetadata{}
	if newSongInfo.Metadata != nil {
		md.Title = newSongInfo.Metadata.Title
		md.Artist = newSongInfo.Metadata.Artist
		md.Album = newSongInfo.Metadata.Album
	}

	currentState.Songs = append(currentState.Songs, upcomingSong{
		filename:   newSongInfo.SongFileName,
		startIndex: newSongInfo.FirstSampleOfSongIndex,
		length:     newSongInfo.SongLength,
		lyrics:     lyrics,
		metadata:   md,
	})
	sort.Sort(songsByStartIndex(currentState.Songs))
}
func (i *infoerPackageHandler) HandleChunkInfo(chunkInfo *comm.ChunkInfo, _ net.Conn) {
	currentState.ChunksMutex.Lock()
	defer currentState.ChunksMutex.Unlock()
	currentState.Chunks = append(currentState.Chunks, upcomingChunk{
		startTime:  chunkInfo.StartTime,
		startIndex: chunkInfo.FirstSampleIndex,
		size:       chunkInfo.ChunkSize,
	})
	sort.Sort(chunksByStartIndex(currentState.Chunks))
}
func (i *infoerPackageHandler) HandlePauseInfo(pauseInfo *comm.PauseInfo, _ net.Conn) {
	currentState.PausesMutex.Lock()
	defer currentState.PausesMutex.Unlock()
	currentState.Pauses = append(currentState.Pauses, pauseToggle{
		playing:     pauseInfo.Playing,
		toggleIndex: pauseInfo.ToggleSampleIndex,
	})
	sort.Sort(pauseByToggleIndex(currentState.Pauses))
}
func (i *infoerPackageHandler) HandleSetVolumeRequest(svr *comm.SetVolumeRequest, _ net.Conn) {
	currentState.Volume = svr.Volume
}

func (i *infoerPackageHandler) HandlePingMessage(_ *comm.PingMessage, conn net.Conn) {
	comm.PingHandler(conn)
}

func (i *infoerPackageHandler) HandleQueueChunkRequest(*comm.QueueChunkRequest, net.Conn) {}
func (i *infoerPackageHandler) HandleTimeSyncRequest(*comm.TimeSyncRequest, net.Conn)     {}
func (i *infoerPackageHandler) HandlePongMessage(*comm.PongMessage, net.Conn)             {}
func (i *infoerPackageHandler) HandleSubscribeChannelRequest(*comm.SubscribeChannelRequest, net.Conn) {
}

// NewPlayerPackageHandler returns the TypedPackageHandler used by players
func newInfoerPackageHandler() comm.TypedPackageHandler {
	return comm.TypedPackageHandler{TypedPackageHandlerInterface: &infoerPackageHandler{}}
}
