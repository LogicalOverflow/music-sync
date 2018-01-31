// This package is the main package of the music-sync player
package main

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/cmd"
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/LogicalOverflow/music-sync/schedule"
	"github.com/LogicalOverflow/music-sync/timing"
	"github.com/gdamore/tcell"
	"github.com/urfave/cli"
	"math"
	"net"
	"os"
	"sort"
	"sync"
	"time"
)

const usage = "run a music-sync client in info mode, which connects to a server and prints information about the current song"

type state struct {
	Songs      []upcomingSong
	SongsMutex sync.RWMutex

	Chunks      []upcomingChunk
	ChunksMutex sync.RWMutex

	Pauses      []pauseToggle
	PausesMutex sync.RWMutex

	Volume float64
}

type pauseByToggleIndex []pauseToggle

func (p pauseByToggleIndex) Len() int           { return len(p) }
func (p pauseByToggleIndex) Less(i, j int) bool { return p[i].toggleIndex < p[j].toggleIndex }
func (p pauseByToggleIndex) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

type chunksByStartIndex []upcomingChunk

func (c chunksByStartIndex) Len() int           { return len(c) }
func (c chunksByStartIndex) Less(i, j int) bool { return c[i].startIndex < c[j].startIndex }
func (c chunksByStartIndex) Swap(i, j int)      { c[i], c[j] = c[j], c[i] }

type songsByStartIndex []upcomingSong

func (s songsByStartIndex) Len() int           { return len(s) }
func (s songsByStartIndex) Less(i, j int) bool { return s[i].startIndex < s[j].startIndex }
func (s songsByStartIndex) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s *state) Info() *playbackInformation {
	now := timing.GetSyncedTime()
	sample := int64(-1)

	s.ChunksMutex.Lock()
	if len(s.Chunks) != 0 {
		passed := 0
		for ; s.Chunks[passed].endTime() < now; passed++ {
			if len(s.Chunks) <= passed {
				goto afterCalcSample
			}
		}
		if passed != 0 {
			copy(s.Chunks, s.Chunks[passed:])
			for i := len(s.Chunks) - passed; i < len(s.Chunks); i++ {
				s.Chunks[i] = upcomingChunk{}
			}
			s.Chunks = s.Chunks[:len(s.Chunks)-passed]
		}
		if len(s.Chunks) != 0 {
			timeInChunk := now - s.Chunks[0].startTime
			sampleInChunk := int64(time.Duration(timeInChunk)*time.Nanosecond) / int64(schedule.SampleRate)
			sample = int64(s.Chunks[0].startIndex) + sampleInChunk
		}
	}

afterCalcSample:
	s.ChunksMutex.Unlock()

	currentSong := upcomingSong{filename: "None", startIndex: 0, length: 0}
	s.SongsMutex.Lock()
	for i := len(s.Songs) - 1; 0 <= i; i-- {
		if int64(s.Songs[i].startIndex) < sample {
			currentSong = s.Songs[i]
			break
		}
	}
	s.SongsMutex.Unlock()

	s.PausesMutex.Lock()
	passed := 0
	for i, p := range s.Pauses {
		if p.playing && p.toggleIndex < currentSong.startIndex {
			passed = i
		}
	}
	if 0 < passed {
		copy(s.Pauses, s.Pauses[passed:])
		for i := len(s.Pauses) - passed; i < len(s.Pauses); i++ {
			s.Pauses[i] = pauseToggle{}
		}
		s.Pauses = s.Pauses[:len(s.Pauses)-passed]
	}

	playing := true
	pauseBegin := int64(0)
	pausesInCurrentSong := int64(0)
	for _, p := range s.Pauses {
		if p.toggleIndex < uint64(sample) && playing != p.playing {
			if p.playing {
				pausesInCurrentSong += int64(p.toggleIndex) - pauseBegin
			} else {
				pauseBegin = int64(p.toggleIndex)
			}

			if p.toggleIndex < currentSong.startIndex {
				if p.playing {
					pausesInCurrentSong = 0
				} else {
					pausesInCurrentSong = int64(p.toggleIndex) - int64(currentSong.startIndex)
				}
			}

			playing = p.playing
		}
	}

	if !playing {
		pausesInCurrentSong += sample - pauseBegin
	}

	s.PausesMutex.Unlock()

	return &playbackInformation{
		CurrentSong:         currentSong,
		CurrentSample:       sample,
		PausesInCurrentSong: pausesInCurrentSong,
		Now:                 now,
		Playing:             playing,
		Volume:              s.Volume,
	}
}

type playbackInformation struct {
	CurrentSong         upcomingSong
	CurrentSample       int64
	PausesInCurrentSong int64
	Now                 int64
	Playing             bool
	Volume              float64
}

var currentState = &state{Songs: make([]upcomingSong, 0), Chunks: make([]upcomingChunk, 0)}

func main() {
	app := cmd.NewApp(usage)
	app.Action = run
	app.Flags = []cli.Flag{
		cmd.ServerAddressFlag,
		cmd.ServerPortFlag,

		cmd.SampleRateFlag,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx *cli.Context) error {
	// disable logging
	log.DefaultCutoffLevel = log.LevelOff
	var (
		serverAddress = ctx.String(cmd.FlagKey(cmd.ServerAddressFlag))
		serverPort    = ctx.Int(cmd.FlagKey(cmd.ServerPortFlag))

		sampleRate = ctx.Int(cmd.FlagKey(cmd.SampleRateFlag))
	)
	tcell.SetEncodingFallback(tcell.EncodingFallbackASCII)
	s, e := tcell.NewScreen()
	if e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}
	if e = s.Init(); e != nil {
		fmt.Fprintf(os.Stderr, "%v\n", e)
		os.Exit(1)
	}

	s.SetStyle(tcell.StyleDefault.Background(tcell.ColorBlack).Foreground(tcell.ColorWhite))
	s.Clear()

	schedule.SampleRate = sampleRate

	server := fmt.Sprintf("%s:%d", serverAddress, serverPort)
	sender, err := comm.ConnectToServer(server, newInfoerPackageHandler())
	if err != nil {
		cli.NewExitError(err, 1)
	}

	go schedule.Infoer(sender)

	drawLoop(s)

	return nil
}

func fmtDuration(duration time.Duration) string {
	if duration < time.Hour {
		return fmt.Sprintf("%d:%02d", duration/time.Minute, duration/time.Second%60)
	}
	return fmt.Sprintf("%d:%02d:%02d", duration/time.Hour, duration/time.Minute%60, duration/time.Second%60)
}

func drawLoop(screen tcell.Screen) {
	w, h := screen.Size()

	running := true
	go func() {
	eventLoop:
		for {
			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				if ev.Key() == tcell.KeyCtrlC {
					running = false
					break eventLoop
				}
			case *tcell.EventResize:
				screen.Sync()
				w, h = screen.Size()
			default:
				panic(fmt.Sprintf("%T", ev))
			}
		}
	}()

	for running {
		screen.Clear()

		info := currentState.Info()
		currentSong := info.CurrentSong
		currentSample := info.CurrentSample

		songLength := time.Duration(0)
		timeInSong := time.Duration(0)
		progressInSong := float64(0)
		if currentSong.startIndex != 0 && int64(currentSong.startIndex) < currentSample {
			sampleInSong := currentSample - int64(currentSong.startIndex) - info.PausesInCurrentSong
			timeInSong = time.Duration(sampleInSong) * time.Second / time.Duration(schedule.SampleRate)
			if 0 < currentSong.length {
				progressInSong = float64(sampleInSong) / float64(currentSong.length)
			}
			songLength = time.Duration(currentSong.length) * time.Second / time.Duration(schedule.SampleRate)
		}

		playState := "Paused"
		if info.Playing {
			playState = "Playing"
		}

		sampleLine := fmt.Sprintf("Current Sample: %d", currentSample)
		songLine := fmt.Sprintf("Current Song (%s): %s (%s/%s)", playState, currentSong.filename, fmtDuration(timeInSong), fmtDuration(songLength))
		volumeLine := fmt.Sprintf("Volume: %06.2f%%", currentState.Volume*100)

		drawString(w-len(volumeLine)-1, h-3, tcell.StyleDefault, volumeLine, screen)
		drawString(1, h-3, tcell.StyleDefault, songLine, screen)

		if len(sampleLine) < w-10 {
			drawString(w-len(sampleLine)-1, h-2, tcell.StyleDefault, sampleLine, screen)
			drawProgress(1, h-2, tcell.StyleDefault, w-2-len(sampleLine), progressInSong, screen)
		} else {
			drawProgress(1, h-2, tcell.StyleDefault, w-2, progressInSong, screen)
		}
		drawBox(0, h-4, w, 4, tcell.StyleDefault, screen)
		screen.Show()
		time.Sleep(500 * time.Millisecond)
	}
	screen.Fini()
}

func drawString(x, y int, style tcell.Style, str string, screen tcell.Screen) {
	for i, r := range str {
		screen.SetContent(x+i, y, r, nil, style)
	}
}

func drawProgress(x, y int, style tcell.Style, length int, progress float64, screen tcell.Screen) {
	head := int(math.Floor(float64(length-1) * progress))
	_, headProgress := math.Modf(float64(length) * progress)
	filledRune := '█'
	emptyRune := ' '
	var headRune rune
	if headProgress < 0.3 {
		headRune = emptyRune
	} else if headProgress < 0.7 {
		headRune = '▌'
	} else {
		headRune = filledRune
	}
	if head == length-1 {
		headRune = filledRune
	}
	for i := 0; i < length; i++ {
		var r rune
		if i < head {
			r = filledRune
		} else if head == i {
			r = headRune
		} else {
			r = emptyRune
		}
		screen.SetContent(x+i, y, r, nil, style)
	}
}

func drawBox(x, y, w, h int, style tcell.Style, screen tcell.Screen) {
	for i := x; i < x+w; i++ {
		screen.SetContent(i, y+h-1, '═', nil, style)
		screen.SetContent(i, y, '═', nil, style)
	}

	for j := y; j < y+h; j++ {
		screen.SetContent(x, j, '║', nil, style)
		screen.SetContent(x+w-1, j, '║', nil, style)
	}
	screen.SetContent(x, y, '╔', nil, style)
	screen.SetContent(x+w-1, y, '╗', nil, style)
	screen.SetContent(x, y+h-1, '╚', nil, style)
	screen.SetContent(x+w-1, y+h-1, '╝', nil, style)
}

type upcomingSong struct {
	filename   string
	startIndex uint64
	length     int64
}

type upcomingChunk struct {
	startTime  int64
	startIndex uint64
	size       uint64
}

func (uc upcomingChunk) endTime() int64 {
	return uc.startTime + uc.length()
}
func (uc upcomingChunk) length() int64 {
	return int64(time.Duration(uc.size)*time.Second) / int64(schedule.SampleRate)
}

type pauseToggle struct {
	playing     bool
	toggleIndex uint64
}

type infoerPackageHandler struct {
}

func (i *infoerPackageHandler) HandleTimeSyncResponse(tsr *comm.TimeSyncResponse, _ net.Conn) {
	clientRecv := timing.GetRawTime()
	timing.UpdateOffset(tsr.ClientSendTime, tsr.ServerRecvTime, tsr.ServerSendTime, clientRecv)
}
func (i *infoerPackageHandler) HandleNewSongInfo(newSongInfo *comm.NewSongInfo, _ net.Conn) {
	currentState.SongsMutex.Lock()
	defer currentState.SongsMutex.Unlock()
	currentState.Songs = append(currentState.Songs, upcomingSong{
		filename:   newSongInfo.SongFileName,
		startIndex: newSongInfo.FirstSampleOfSongIndex,
		length:     newSongInfo.SongLength,
	})
	sort.Sort(songsByStartIndex(currentState.Songs))
}
func (i *infoerPackageHandler) HandleChunkInfo(chunkInfo *comm.ChunkInfo, _ net.Conn) {
	currentState.ChunksMutex.Lock()
	defer currentState.ChunksMutex.Unlock()
	currentState.Chunks = append(currentState.Chunks, upcomingChunk{
		startTime:  chunkInfo.StartTime,
		startIndex: chunkInfo.FirstSampleIndex,
		size:       chunkInfo.ChunkSize,
	})
	sort.Sort(chunksByStartIndex(currentState.Chunks))
}
func (i *infoerPackageHandler) HandlePauseInfo(pauseInfo *comm.PauseInfo, _ net.Conn) {
	currentState.PausesMutex.Lock()
	defer currentState.PausesMutex.Unlock()
	currentState.Pauses = append(currentState.Pauses, pauseToggle{
		playing:     pauseInfo.Playing,
		toggleIndex: pauseInfo.ToggleSampleIndex,
	})
	sort.Sort(pauseByToggleIndex(currentState.Pauses))
}
func (i *infoerPackageHandler) HandleSetVolumeRequest(svr *comm.SetVolumeRequest, _ net.Conn) {
	currentState.Volume = svr.Volume
}

func (i *infoerPackageHandler) HandlePingMessage(_ *comm.PingMessage, conn net.Conn) {
	comm.PingHandler(conn)
}

func (i *infoerPackageHandler) HandleQueueChunkRequest(*comm.QueueChunkRequest, net.Conn) {}
func (i *infoerPackageHandler) HandleTimeSyncRequest(*comm.TimeSyncRequest, net.Conn)     {}
func (i *infoerPackageHandler) HandlePongMessage(*comm.PongMessage, net.Conn)             {}
func (i *infoerPackageHandler) HandleSubscribeChannelRequest(*comm.SubscribeChannelRequest, net.Conn) {
}

// NewPlayerPackageHandler returns the TypedPackageHandler used by players
func newInfoerPackageHandler() comm.TypedPackageHandler {
	return comm.TypedPackageHandler{TypedPackageHandlerInterface: &infoerPackageHandler{}}
}
