package schedule

import (
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/metadata"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/ssh"
)

// Server starts a music-sync server, using sender to communicate with all clients
func Server(sender comm.MessageSender) {
	ss := &serverState{}
	ss.sender = sender

	ss.lyricsProvider = metadata.GetLyricsProvider()
	ss.metadataProvider = metadata.GetProvider()

	ss.playlist = playback.NewPlaylist(SampleRate, []string{}, NanBreakSize)
	ss.volume = 0.1

	ss.pauses = make([]*comm.PauseInfo, 0)

	comm.NewClientHandler = ss.createClientHandler()

	go ss.playlist.StreamLoop()

	ss.playlist.SetNewSongHandler(ss.createNewSongHandler())
	ss.playlist.SetPauseToggleHandler(ss.createPauseToggleHandler())

	go ss.streamMusic()

	ssh.RegisterCommand(ss.queueCommand())
	ssh.RegisterCommand(ss.playlistCommand())
	ssh.RegisterCommand(ss.removeCommand())
	ssh.RegisterCommand(ss.jumpCommand())
	ssh.RegisterCommand(ss.volumeCommand())
	ssh.RegisterCommand(ss.pauseCommand())
	ssh.RegisterCommand(ss.resumeCommand())
}
