package playback

import (
	"github.com/LogicalOverflow/music-sync/timing"
	"github.com/faiface/beep"
	"math"
	"time"
)

type timedSample struct {
	sample [2]float64
	time   int64
}

type TimedMultiStreamer struct {
	format         beep.Format
	streamers      []*QueuedStream
	chunks         []*QueuedStream
	background     beep.Streamer
	offset         int64
	sampleDuration int64
	maxCorrection  int
	samples        *timedSampleQueue
	syncing        bool
}

func (tms *TimedMultiStreamer) Stream(samples [][2]float64) {
	var n int
	var drained bool
	now := timing.GetSyncedTime()
	for 0 < len(samples) {
		if tms.syncing {
			n, drained = tms.StreamSync(samples, now)
		} else {
			n, drained = tms.StreamDirect(samples)
		}
		now += tms.samplesDuration(n)
		samples = samples[n:]
		if drained {
			tms.syncing = !tms.syncing
			_, t := tms.samples.Peek()
			logger.Debugf("playback error: %s", time.Duration(t-now)*time.Nanosecond)
		}
	}
}

func (tms *TimedMultiStreamer) StreamDirect(samples [][2]float64) (n int, drained bool) {
	for i := range samples {
		s, _ := tms.samples.Remove()
		if math.IsNaN(s[0]) {
			return i, true
		} else {
			samples[i] = s
		}
	}
	return len(samples), false
}

func (tms *TimedMultiStreamer) StreamSync(samples [][2]float64, now int64) (n int, drained bool) {
	s, t := tms.samples.Peek()
	for math.IsNaN(s[0]) {
		tms.samples.Remove()
		s, t = tms.samples.Peek()
	}
	if now < t+tms.samplesDuration(len(samples)) {
		tms.background.Stream(samples)
		return len(samples), false
	} else if now < t {
		silence := tms.samplesCount(t - now)
		tms.background.Stream(samples[:silence])
		return silence, true
	} else {
		for t < now {
			tms.samples.Remove()
			_, t = tms.samples.Peek()
		}
		return 0, true
	}
}

func (tms *TimedMultiStreamer) ReadChunks() {
	for {
		if 0 < len(tms.chunks) {
			st := tms.chunks[0].startTime
			for i, s := range tms.chunks[0].samples {
				tms.samples.Add(s, st+tms.samplesDuration(i))
			}
			tms.chunks = tms.chunks[1:]
		} else {
			time.Sleep(time.Millisecond)
		}
	}
}

func (tms *TimedMultiStreamer) samplesDuration(n int) int64 {
	return int64(tms.format.SampleRate.D(n) / time.Nanosecond)
}

func (tms *TimedMultiStreamer) samplesCount(n int64) int {
	return tms.format.SampleRate.N(time.Duration(n) * time.Nanosecond)
}

func (tms *TimedMultiStreamer) Err() error { return nil }

type QueuedStream struct {
	startTime int64
	samples   [][2]float64
	sampleN   int
	pos       int
}

func (q *QueuedStream) copySamples(target [][2]float64) (n int) {
	if q.sampleN <= q.pos {
		return 0
	}

	n = copy(target, q.samples[q.pos:])
	q.pos += n
	return
}

func (q *QueuedStream) drained() bool {
	return q.sampleN <= q.pos
}

func NewQueuedStream(startTime int64, samples [][2]float64) *QueuedStream {
	return &QueuedStream{startTime: startTime, samples: samples, sampleN: len(samples)}
}
