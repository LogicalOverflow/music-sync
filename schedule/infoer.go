package schedule

import (
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/timing"
	"os"
	"time"
)

// Infoer start a music-sync client in infoer mode, using sender to communicate with the server
func Infoer(sender comm.MessageSender) {
	go func() {
		syncTime(sender)
		for range time.Tick(TimeSyncInterval) {
			syncTime(sender)
		}
	}()

	go func() {
		if err := sender.SendMessage(&comm.SubscribeChannelRequest{Channel: comm.Channel_META}); err != nil {
			logger.Errorf("failed to subscribe to meta channel")
			os.Exit(1)
		}
	}()
}

func syncTime(sender comm.MessageSender) {
	logger.Infof("syncing time")
	timing.ResetOffsets(TimeSyncCycles)
	for i := 0; i < TimeSyncCycles; i++ {
		if err := sender.SendMessage(&comm.TimeSyncRequest{ClientSend: timing.GetRawTime()}); err != nil {
			logger.Warnf("failed to send sync time request: %v", err)
		}
		time.Sleep(TimeSyncCycleDelay)
	}
}
