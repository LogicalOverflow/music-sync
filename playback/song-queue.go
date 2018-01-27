package playback

import (
	"fmt"
	"os"
	"path"
	"time"

	"github.com/LogicalOverflow/music-sync/logging"

	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/hajimehoshi/oto"
)

var (
	player     *oto.Player
	format     beep.Format
	streamer   *TimedMultiStreamer
	bufferSize int
	volume     float64
)

var logger = log.GetLogger("play")
var AudioDir string

func GetStreamer(filename string) (beep.Streamer, error) {
	filename = path.Join(AudioDir, filename)
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %s: %v", filename, err)
	}

	s, _, err := mp3.Decode(f)
	if err != nil {
		return nil, fmt.Errorf("failed to decode file %s: %v", filename, err)
	}

	return s, nil
}

func QueueSong(startTime int64, chunkId int64, samples [][2]float64) {
	if streamer == nil {
		logger.Infof("not queuing chunk %d: streamer not ready", chunkId)
		return
	}
	logger.Debugf("queueing chunk %d at %d", chunkId, startTime)

	q := NewQueuedStream(startTime, samples)
	streamer.chunks = append(streamer.chunks, q)
	logger.Debugf("chunk %d queued at %d", chunkId, startTime)
}

func CombineSamples(low []float64, high []float64) ([][2]float64) {
	e := len(low)
	if len(high) < e {
		e = len(high)
	}
	s := make([][2]float64, e)
	for i := 0; i < e; i++ {
		s[i] = [2]float64{low[i], high[i]}
	}
	return s
}

func Init(sampleRate int) error {
	logger.Infof("initializing playback")
	var err error

	volume = .1

	format = beep.Format{
		SampleRate:  beep.SampleRate(sampleRate),
		NumChannels: 2,
		Precision:   2,
	}

	bufferSize = format.SampleRate.N(time.Second / 10)
	player, err = oto.NewPlayer(int(format.SampleRate), format.NumChannels, format.Precision,
		format.NumChannels*format.Precision*bufferSize)

	if err != nil {
		return fmt.Errorf("failed to initialize speaker: %v", err)
	}
	player.SetUnderrunCallback(func() { logger.Warn("player is underrunning") })

	streamer = &TimedMultiStreamer{
		format:         format,
		streamers:      make([]*QueuedStream, 0),
		chunks:         make([]*QueuedStream, 0),
		background:     beep.Silence(-1),
		offset:         0,
		sampleDuration: int64(format.SampleRate.D(1) / time.Nanosecond),
		maxCorrection:  10,
		samples:        newTimedSampleQueue(2 * sampleRate),
		syncing:        true,
	}

	go playLoop()

	go streamer.ReadChunks()

	logger.Infof("playback initialized")

	return nil
}

func SetVolume(v float64) {
	volume = v
	logger.Infof("volume set to %.3f", v)
}

func playLoop() {
	numBytes := bufferSize * format.NumChannels * format.Precision
	samples := make([][2]float64, bufferSize)
	buf := make([]byte, numBytes)

	for {
		streamer.Stream(samples)

		for i := range samples {
			for c := range samples[i] {
				val := samples[i][c] * volume
				if val < -1 {
					val = -1
				}
				if val > +1 {
					val = +1
				}
				valInt16 := int16(val * (1<<15 - 1))
				low := byte(valInt16)
				high := byte(valInt16 >> 8)
				buf[i*4+c*2+0] = low
				buf[i*4+c*2+1] = high
			}
		}

		player.Write(buf)
	}
}
