package timing

import (
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/aristanetworks/goarista/monotime"
	"math/big"
	"sync"
)

var logger = log.GetLogger("time")

var offset int64 = 0
var offsets = make([]int64, 0)
var offsetsMutex sync.Mutex

func GetSyncedTime() int64 {
	return int64(monotime.Now()) + offset
}

func GetRawTime() int64 {
	return int64(monotime.Now())
}

func ResetOffsets(cap int) {
	offsetsMutex.Lock()
	offsets = make([]int64, 0, cap)
	offsetsMutex.Unlock()
}

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
