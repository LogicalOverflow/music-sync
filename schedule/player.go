package schedule

import (
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/playback"
	"os"
	"time"
)

// Player starts a music-sync player, using sender to communicate with the server
func Player(sender comm.MessageSender) {
	go func() {
		if err := playback.Init(SampleRate); err != nil {
			logger.Fatalf("failed to initialized playback: %v", err)
			os.Exit(1)
		}
	}()

	go func() {
		syncTime(sender)
		for range time.Tick(TimeSyncInterval) {
			syncTime(sender)
		}
	}()

	go func() {
		if err := sender.SendMessage(&comm.SubscribeChannelRequest{Channel: comm.Channel_AUDIO}); err != nil {
			logger.Errorf("failed to subscribe to audio channel")
			os.Exit(1)
		}
	}()
}
