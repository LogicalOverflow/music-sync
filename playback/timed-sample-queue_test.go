package playback

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

const testQueueSize = 16

func TestFIFO(t *testing.T) {
	q := newTestQueue()

	for i := 0; i < testQueueSize; i++ {
		q.Add([2]float64{float64(i), float64(i)}, int64(i))
	}

	for i := 0; i < testQueueSize; i++ {
		sam, ti := q.Remove()
		assert.Equal(t, float64(i), sam[0], "%d-th remove did not yield the element added %d-th (sample[0])", i+1, i+1)
		assert.Equal(t, float64(i), sam[1], "%d-th remove did not yield the element added %d-th (sample[1])", i+1, i+1)
		assert.Equal(t, int64(i), ti, "%d-th remove did not yield the element added %d-th (time)", i+1, i+1)
	}
}

func TestPeek(t *testing.T) {
	q := newTestQueue()
	q.Add([2]float64{1, 1}, 1)
	q.Add([2]float64{2, 2}, 2)

	sam, ti := q.Peek()
	assert.Equal(t, float64(1), sam[0], "peek did not yield the correct element (sample[0])")
	assert.Equal(t, float64(1), sam[1], "peek did not yield the correct element (sample[1])")
	assert.Equal(t, int64(1), ti, "peek did not yield the correct element (time)")

	sam2, ti2 := q.Peek()
	assert.Equal(t, sam[0], sam2[0], "a second peek did not yield the same element as the first (sample[0])")
	assert.Equal(t, sam[1], sam2[1], "a second peek did not yield the same element as the first (sample[1])")
	assert.Equal(t, ti, ti2, "a second peek did not yield the same element as the first (time)")
}

func TestCap(t *testing.T) {
	for i := 0; i < testQueueSize; i++ {
		assert.Equal(t, i, newTimedSampleQueue(i).Cap(), "queue created with size %d has wrong capacity", i)
	}
}

func TestLen(t *testing.T) {
	q := newTestQueue()

	assert.Equal(t, 0, q.Len(), "an empty queue does not have length 0")

	for i := 0; i < testQueueSize; i++ {
		q.Add([2]float64{float64(i), float64(i)}, int64(i))
		assert.Equal(t, i+1, q.Len(), "after adding %d elements, queue length is incorrect", i+1)
	}

	for i := 0; i < testQueueSize; i++ {
		q.Remove()
		assert.Equal(t, testQueueSize-i-1, q.Len(), "after adding %d elements and removing %d again, queue length is incorrect", testQueueSize, i+1)
	}

	q = newTestQueue()
	for i := 0; i < testQueueSize/2; i++ {
		q.Add([2]float64{float64(i), float64(i)}, int64(i))
	}
	for i := 0; i < 2*testQueueSize; i++ {
		q.Add([2]float64{float64(i), float64(i)}, int64(i))
		q.Remove()
		assert.Equal(t, testQueueSize/2, q.Len(), "after adding %d elements and the adding and removing %d elements, length is incorrect", testQueueSize/2, i+1)
	}
}

func TestFull(t *testing.T) {
	q := newTestQueue()

	added := 0
	removed := 0

	assert.False(t, q.full(), "newly created queue of size %d claims to be full", testQueueSize)

	for i := 0; i < testQueueSize-1; i++ {
		q.Add([2]float64{float64(i), float64(i)}, int64(i))
		added++
		assert.False(t, q.full(), "after adding %d elements, queue of size %d claims to be full", added, testQueueSize)
	}

	q.Add([2]float64{float64(testQueueSize - 1), float64(testQueueSize - 1)}, int64(testQueueSize-1))
	added++
	assert.True(t, q.full(), "after adding %d elements, queue of size %d does not claim to be full", added, testQueueSize)

	for i := 0; i < testQueueSize; i++ {
		q.Remove()
		removed++
		assert.False(t, q.full(), "after adding %d elements and removing %d again, queue of size %d claims to be full", added, removed, testQueueSize)
	}
}

func TestEmpty(t *testing.T) {
	q := newTestQueue()

	added := 0
	removed := 0

	assert.True(t, q.empty(), "newly created queue of size %d does not claim to be empty", testQueueSize)

	for i := 0; i < testQueueSize; i++ {
		q.Add([2]float64{float64(i), float64(i)}, int64(i))
		added++
		assert.False(t, q.empty(), "after adding %d elements, queue of size %d claims to be empty", added, testQueueSize)
	}

	for i := 0; i < testQueueSize-1; i++ {
		q.Remove()
		removed++
		assert.False(t, q.empty(), "after adding %d elements and removing %d again, queue of size %d claims to be empty", added, removed, testQueueSize)
	}

	q.Remove()
	removed++
	assert.True(t, q.empty(), "after adding %d elements and removing %d again, queue of size %d does not claim to be empty", added, removed, testQueueSize)
}

func TestAddRemoveAsync(t *testing.T) {
	q := newTestQueue()
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 4*testQueueSize; i++ {
			q.Add([2]float64{float64(i), float64(i)}, int64(i))
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 4*testQueueSize; i++ {
			sam, ti := q.Remove()

			assert.Equal(t, float64(i), sam[0], "async removing the %d-th time while first filling did not yield the correct element (sample[0])", i+1)
			assert.Equal(t, float64(i), sam[1], "async removing the %d-th time while first filling did not yield the correct element (sample[1])", i+1)
			assert.Equal(t, int64(i), ti, "async removing the %d-th time while first filling did not yield the correct element (time)", i+1)
		}
	}()

	q = newTestQueue()
	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < 4*testQueueSize; i++ {
			sam, ti := q.Remove()

			assert.Equal(t, float64(i), sam[0], "async removing the %d-th time while first removing did not yield the correct element (sample[0])", i+1)
			assert.Equal(t, float64(i), sam[1], "async removing the %d-th time while first removing did not yield the correct element (sample[1])", i+1)
			assert.Equal(t, int64(i), ti, "async removing the %d-th time while first removing did not yield the correct element (time)", i+1)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 4*testQueueSize; i++ {
			q.Add([2]float64{float64(i), float64(i)}, int64(i))
		}
	}()

	wg.Wait()
}

func newTestQueue() *timedSampleQueue {
	return newTimedSampleQueue(testQueueSize)
}
