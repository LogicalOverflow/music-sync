// Package playback contains functions and types to stream and play audio
package playback

//

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
	streamer   *timedMultiStreamer
	bufferSize int
	volume     float64
)

var logger = log.GetLogger("play")

// AudioDir is the directory containing the audio file
var AudioDir string

func getStreamer(filename string) (beep.StreamSeekCloser, error) {
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

// QueueChunk queue a chunk for playback
func QueueChunk(startTime int64, chunkID int64, samples [][2]float64) {
	if streamer == nil {
		logger.Infof("not queuing chunk %d: streamer not ready", chunkID)
		return
	}
	logger.Debugf("queueing chunk %d at %d", chunkID, startTime)

	q := newQueuedStream(startTime, samples)
	streamer.chunks = append(streamer.chunks, q)
	logger.Debugf("chunk %d queued at %d", chunkID, startTime)
}

// CombineSamples combines to []float64 to one [][2]float64,
// such that low[i] == returned[i][0] and high[i] == returned[i][1]
func CombineSamples(low []float64, high []float64) [][2]float64 {
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

// Init prepares a player for playback with the given sample rate
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

	initStreamer(sampleRate)

	go playLoop()

	go streamer.ReadChunks()

	logger.Infof("playback initialized")

	return nil
}

func initStreamer(sampleRate int) {
	streamer = &timedMultiStreamer{
		format:         format,
		streamers:      make([]*queuedStream, 0),
		chunks:         make([]*queuedStream, 0),
		background:     beep.Silence(-1),
		offset:         0,
		sampleDuration: int64(format.SampleRate.D(1) / time.Nanosecond),
		maxCorrection:  10,
		samples:        newTimedSampleQueue(2 * sampleRate),
		syncing:        true,
	}
}

// SetVolume sets the playback volume of the player
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
		samplesToAudioBuf(samples, buf)
		player.Write(buf)
	}
}

func samplesToAudioBuf(samples [][2]float64, buf []byte) {
	for i := range samples {
		for c := range samples[i] {
			buf[i*4+c*2+0], buf[i*4+c*2+1] = convertSampleToBytes(samples[i][c] * volume)
		}
	}
}

func convertSampleToBytes(val float64) (low, high byte) {
	if val < -1 {
		val = -1
	}
	if val > +1 {
		val = +1
	}
	valInt16 := int16(val * (1<<15 - 1))
	low = byte(valInt16)
	high = byte(valInt16 >> 8)
	return
}
