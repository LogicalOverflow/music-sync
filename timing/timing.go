// Package timing provides functions to access raw and synced time with nanosecond precision
package timing

import (
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/aristanetworks/goarista/monotime"
	"math/big"
	"sync"
)

var logger = log.GetLogger("time")

var offset int64
var offsets = make([]int64, 0)
var offsetsMutex sync.Mutex

// GetSyncedTime returns the current time, sync to the sever, with nanosecond precision
func GetSyncedTime() int64 {
	return int64(monotime.Now()) + offset
}

// GetRawTime returns the current time, not sync to the server, with nanosecond precision
func GetRawTime() int64 {
	return int64(monotime.Now())
}

// ResetOffsets clears the slice holding offsets by replacing it with an empty slice with cap as capacity
func ResetOffsets(cap int) {
	offsetsMutex.Lock()
	offsets = make([]int64, 0, cap)
	offsetsMutex.Unlock()
}

// UpdateOffset handles the four timestamps used to synchronize time to the server
func UpdateOffset(clientSend, serverRecv, serverSend, clientRecv int64) {
	logger.Tracef("updating offset: %d, %d, %d, %d", clientSend, serverRecv, serverSend, clientRecv)

	calcOffset := ((serverRecv - clientSend) + (serverSend - clientRecv)) / 2
	offsetsMutex.Lock()
	offsets = append(offsets, calcOffset)
	offsetsMutex.Unlock()

	if len(offsets) == cap(offsets) {
		offsetsMutex.Lock()
		avg := big.NewInt(0)
		for _, o := range offsets {
			avg = new(big.Int).Add(avg, big.NewInt(o))
		}
		avg = new(big.Int).Div(avg, big.NewInt(int64(len(offsets))))
		offset = avg.Int64()
		offsetsMutex.Unlock()
		logger.Debugf("time synced: offset: %d", offset)
	}
}
