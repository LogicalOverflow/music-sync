// Package schedule contains methods to start different types of clients/servers
package schedule

import (
	"github.com/LogicalOverflow/music-sync/logging"
	"time"
)

var logger = log.GetLogger("shed")

// TimeSyncInterval is the time interval between syncing time to server
var TimeSyncInterval = 10 * time.Minute

// TimeSyncCycles is the number of cycles used to sync time to server
var TimeSyncCycles = 500

// TimeSyncCycleDelay is the delay between cycles in one time sync
var TimeSyncCycleDelay = 10 * time.Millisecond

// StreamChunkSize is the size of one stream chunk in samples
var StreamChunkSize = 44100 * 4

// StreamChunkTime is the duration is takes to play one stream chunk
var StreamChunkTime = 4 * time.Second

// NanBreakSize is the number of nan-samples to insert between songs, which players use to realign playback
var NanBreakSize = 44100 * 1

// StreamStartDelay is the delay before starting the stream
var StreamStartDelay = 5 * time.Second

// StreamDelay is the delay of the stream, which players use to decode chunks
var StreamDelay = 15 * time.Second

// SampleRate is the sample rate of the stream
var SampleRate = 44100
