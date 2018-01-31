package timing

import (
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestResetOffsets(t *testing.T) {
	for i := 0; i < 64; i++ {
		ResetOffsets(i)
		assert.Equal(t, i, cap(offsets), "ResetOffset(%d) created offsets slice with wrong capacity", i)
		assert.Equal(t, 0, len(offsets), "ResetOffset(%d) created offsets slice with wrong length", i)
	}
}

func TestUpdateOffset(t *testing.T) {
	log.DefaultCutoffLevel = log.LevelOff
	for targetOffset := 0; targetOffset < 128; targetOffset += 8 {
		ResetOffsets(32)
		initialOffset := offset
		for i := 0; i < 32; i++ {
			assert.Equal(t, initialOffset, offset, "after updating %d times, offset changed", i+1)
			UpdateOffset(int64(i), int64(i+targetOffset), int64(i+targetOffset), int64(i))
			assert.Equal(t, i+1, len(offsets), "after updating offsets %d times, offsets slice has the wrong length", i+1)
		}
		assert.Equal(t, int64(targetOffset), offset, "after completing time sync, offset is incorrect")
	}
}
