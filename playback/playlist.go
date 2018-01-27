package playback

import (
	"math"
	"sync"
)

type Playlist struct {
	songs        []string
	songsMutex   sync.RWMutex
	position     int
	low          chan float64
	high         chan float64
	forceNext    chan bool
	nanBreakSize int
	playing      bool
	currentSong  string
}

func (pl *Playlist) StreamLoop() {
	for {
		pl.songsMutex.RLock()
		if len(pl.songs) == 0 {
			for i := 0; i < 4410; i++ {
				pl.low <- math.NaN()
				pl.high <- math.NaN()
			}
			pl.songsMutex.RUnlock()
			continue
		}

		pos := pl.position % len(pl.songs)
		pl.position = pl.position % len(pl.songs)
		filename := pl.songs[pos]
		pl.songsMutex.RUnlock()

		pl.currentSong = filename

		s, err := GetStreamer(filename)
		if err != nil {
			logger.Warnf("skipping song %s in playlist: %v", filename, err)
		}

		buf := make([][2]float64, 512)
		for {
			n, ok := len(buf), true
			if pl.playing {
				n, ok = s.Stream(buf)

				for _, sample := range buf {
					pl.low <- sample[0]
					pl.high <- sample[1]
				}
			} else {
				for range buf {
					pl.low <- math.NaN()
					pl.high <- math.NaN()
				}
			}

			if 0 < len(pl.forceNext) && <-pl.forceNext {
				break
			}
			if !ok || n < len(buf) {
				if pos == pl.position {
					pl.position++
				}
				break
			}
		}
		for i := 0; i < pl.nanBreakSize; i++ {
			pl.low <- math.NaN()
			pl.high <- math.NaN()
		}
	}
}

func (pl *Playlist) SetPos(pos int) {
	pl.position = pos
	pl.forceNext <- true
}

func (pl *Playlist) Pos() int {
	return pl.position
}

func (pl *Playlist) Songs() []string {
	pl.songsMutex.RLock()
	defer pl.songsMutex.RUnlock()
	r := make([]string, len(pl.songs))
	copy(r, pl.songs)
	return r
}

func (pl *Playlist) AddSong(song string) {
	pl.songsMutex.Lock()
	defer pl.songsMutex.Unlock()
	pl.songs = append(pl.songs, song)
}

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

func (pl *Playlist) RemoveSong(index int) string {
	pl.songsMutex.Lock()
	defer pl.songsMutex.Unlock()
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

func (pl *Playlist) Fill(low []float64, high []float64) {
	var wg sync.WaitGroup
	wg.Add(2)
	go func() { defer wg.Done(); copyFloatChannel(low, pl.low) }()
	go func() { defer wg.Done(); copyFloatChannel(high, pl.high) }()
	wg.Wait()
}

func (pl *Playlist) Playing() bool {
	return pl.playing
}

func (pl *Playlist) SetPlaying(p bool) {
	pl.playing = p
}

func (pl *Playlist) CurrentSong() string {
	return pl.currentSong
}

func NewPlaylist(bufferSize int, songs []string, nanBreakSize int) *Playlist {
	return &Playlist{
		songs:        songs,
		position:     0,
		low:          make(chan float64, bufferSize),
		high:         make(chan float64, bufferSize),
		forceNext:    make(chan bool, 2),
		nanBreakSize: nanBreakSize,
		playing:      true,
	}
}

func copyFloatChannel(dst [] float64, src chan float64) {
	for i := range dst {
		dst[i] = <-src
	}
}
