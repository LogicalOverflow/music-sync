package playback

import (
	"github.com/faiface/beep"
	"math"
	"sync"
)

// Playlist is a array of songs, which can then be streamed.
// After reaching the end of the playlist, playback will resume at the start.
type Playlist struct {
	songs        []string
	songsMutex   sync.RWMutex
	position     int
	low          chan float64
	high         chan float64
	forceNext    chan bool
	nanBreakSize int

	playing     bool
	currentSong string

	sampleIndexRead  uint64
	sampleIndexWrite uint64

	newSongHandler     func(startSampleIndex uint64, filename string, songLength int64)
	pauseToggleHandler func(playing bool, sample uint64)

	playingLast bool
}

// StreamLoop reads the samples of the song into the internal buffer.
// This method blocks forever and must be called exactly once before streaming from the playlist.
func (pl *Playlist) StreamLoop() {
	for {
		filename := pl.nextSong()
		if filename == "" {
			for i := 0; i < 44100; i++ {
				pl.pushSample(math.NaN(), math.NaN())
			}
			continue
		}

		pl.currentSong = filename

		s, err := getStreamer(filename)
		if err != nil {
			logger.Warnf("skipping song %s in playlist: %v", filename, err)
			pl.position++
			continue
		}

		pl.pushStreamer(s)
		pl.pusNanBreak()
	}
}

func (pl *Playlist) nextSong() (song string) {
	pl.songsMutex.RLock()
	defer pl.songsMutex.RUnlock()
	if len(pl.songs) == 0 {
		return ""
	}
	pos := pl.position % len(pl.songs)
	pl.position = pos
	return pl.songs[pos]
}

func (pl *Playlist) pushStreamer(s beep.StreamSeekCloser) {
	buf := make([][2]float64, 512)
	go pl.newSongHandler(pl.sampleIndexWrite, pl.currentSong, int64(s.Len()))
	for {
		n, ok := len(buf), true
		pl.callPauseToggleHandler()
		if pl.playing {
			n, ok = s.Stream(buf)
			pl.pushBuffer(buf[:n])
		} else {
			pl.pushNanSamples(len(buf))
		}

		if 0 < len(pl.forceNext) && <-pl.forceNext {
			break
		}

		if !ok || n < len(buf) {
			pl.position++
			break
		}
	}
}

func (pl *Playlist) pushNanSamples(count int) {
	for i := 0; i < count; i++ {
		pl.pushSample(math.NaN(), math.NaN())
	}
}

func (pl *Playlist) pushBuffer(buffer [][2]float64) {
	for _, sample := range buffer {
		pl.pushSample(sample[0], sample[1])
	}
}

func (pl *Playlist) pusNanBreak() {
	for i := 0; i < pl.nanBreakSize; i++ {
		pl.pushSample(math.NaN(), math.NaN())
	}
}

func (pl *Playlist) pushSample(low, high float64) {
	pl.low <- low
	pl.high <- high
	pl.sampleIndexWrite++
}

func (pl *Playlist) callPauseToggleHandler() {
	if pl.playing != pl.playingLast {
		pl.playingLast = pl.playing
		go pl.pauseToggleHandler(pl.playing, pl.sampleIndexWrite)
	}
}

// SetPos jumps to the song at pos.
func (pl *Playlist) SetPos(pos int) {
	pl.position = pos
	pl.forceNext <- true
}

// Pos returns the position of the song currently being played.
func (pl *Playlist) Pos() int {
	return pl.position
}

// Songs returns all songs in the playlist.
func (pl *Playlist) Songs() []string {
	pl.songsMutex.RLock()
	defer pl.songsMutex.RUnlock()
	r := make([]string, len(pl.songs))
	copy(r, pl.songs)
	return r
}

// AddSong adds a song at the end of the playlist.
func (pl *Playlist) AddSong(song string) {
	pl.songsMutex.Lock()
	defer pl.songsMutex.Unlock()
	pl.songs = append(pl.songs, song)
}

// InsertSong inserts a song into the playlist.
// The index is clipped to the bounds of the playlist.
func (pl *Playlist) InsertSong(song string, index int) {
	pl.songsMutex.Lock()
	defer pl.songsMutex.Unlock()
	if index < 0 {
		index = 0
	}
	if len(pl.songs) < index {
		pl.songs = append(pl.songs, song)
	} else {
		pl.songs = append(pl.songs, "")
		copy(pl.songs[index+1:], pl.songs[index:])
		pl.songs[index] = song
	}
}

// RemoveSong remove the song at index from the playlist and returns the removed song.
// The index is clipped to the bounds of the playlist.
// If the playlist is empty, noting happens and "" is returned.
func (pl *Playlist) RemoveSong(index int) string {
	pl.songsMutex.Lock()
	defer pl.songsMutex.Unlock()
	if len(pl.songs) == 0 {
		return ""
	}
	if index < 0 {
		index = 0
	}
	var removed string
	if len(pl.songs) < index {
		removed = pl.songs[len(pl.songs)-1]
		pl.songs = pl.songs[:len(pl.songs)-1]
	} else {
		removed = pl.songs[index]
		copy(pl.songs[index:], pl.songs[index+1:])
		pl.songs[len(pl.songs)-1] = ""
		pl.songs = pl.songs[:len(pl.songs)-1]
	}
	return removed
}

// Fill reads the samples from the internal buffer and fills low and high with them.
// low and high must have the same length.
// returns the sampleIndex of the first read
func (pl *Playlist) Fill(low []float64, high []float64) uint64 {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); copyFloatChannel(low, pl.low) }()
	go func() { defer wg.Done(); copyFloatChannel(high, pl.high) }()
	wg.Wait()
	defer func() { pl.sampleIndexRead += uint64(len(low)) }()
	return pl.sampleIndexRead
}

// Playing returns true if the playlist is currently playing audio.
func (pl *Playlist) Playing() bool {
	return pl.playing
}

// SetPlaying can set whether or not the playlist should be playing audio.
func (pl *Playlist) SetPlaying(p bool) {
	pl.playing = p
}

// CurrentSong returns the song currently being played
func (pl *Playlist) CurrentSong() string {
	return pl.currentSong
}

// SetNewSongHandler sets the new song handler, which is called every time the playlist begins playing a new song
func (pl *Playlist) SetNewSongHandler(nsh func(startSampleIndex uint64, filename string, songLength int64)) {
	pl.newSongHandler = nsh
}

// SetPauseToggleHandler sets the pause toggle handler, which is called every time the playlist is paused/resumed
func (pl *Playlist) SetPauseToggleHandler(psh func(playing bool, sample uint64)) {
	pl.pauseToggleHandler = psh
}

// NewPlaylist create a new playlist with the given buffer size and songs in it, which
// inserts nanBreakSize nan-samples between songs, which players use to realign playback.
func NewPlaylist(bufferSize int, songs []string, nanBreakSize int) *Playlist {
	return &Playlist{
		songs:            songs,
		position:         0,
		low:              make(chan float64, bufferSize),
		high:             make(chan float64, bufferSize),
		forceNext:        make(chan bool, 2),
		nanBreakSize:     nanBreakSize,
		playing:          false,
		playingLast:      true,
		sampleIndexRead:  0,
		sampleIndexWrite: 0,
	}
}

func copyFloatChannel(dst []float64, src chan float64) {
	for i := range dst {
		dst[i] = <-src
	}
}
