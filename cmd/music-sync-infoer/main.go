// This package is the main package of the music-sync player
package main

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/cmd"
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/LogicalOverflow/music-sync/schedule"
	"github.com/gdamore/tcell"
	"github.com/urfave/cli"
	"os"
	"time"
)

const usage = "run a music-sync client in info mode, which connects to a server and prints information about the current song"

func main() {
	app := cmd.NewApp(usage)
	app.Action = run
	app.Flags = []cli.Flag{
		cmd.ServerAddressFlag,
		cmd.ServerPortFlag,

		cmd.SampleRateFlag,
		cmd.LyricsHistorySizeFlag,
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var lyricsHistorySize int

func run(ctx *cli.Context) error {
	// disable logging
	log.DefaultCutoffLevel = log.LevelOff
	var (
		serverAddress = ctx.String(cmd.FlagKey(cmd.ServerAddressFlag))
		serverPort    = ctx.Int(cmd.FlagKey(cmd.ServerPortFlag))

		sampleRate = ctx.Int(cmd.FlagKey(cmd.SampleRateFlag))
	)
	lyricsHistorySize = int(ctx.Uint(cmd.FlagKey(cmd.LyricsHistorySizeFlag)))

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

	tcellLoop(s)

	return nil
}

func fmtDuration(duration time.Duration) string {
	if duration < time.Hour {
		return fmt.Sprintf("%d:%02d", duration/time.Minute, duration/time.Second%60)
	}
	return fmt.Sprintf("%d:%02d:%02d", duration/time.Hour, duration/time.Minute%60, duration/time.Second%60)
}

func tcellLoop(screen tcell.Screen) {
	w, h := screen.Size()

	running := true
	go func() {
		eventLoop(screen, &w, &h)
		running = false
	}()

	for range time.Tick(200 * time.Millisecond) {
		if !running {
			break
		}
		redraw(screen, w, h)
	}

	screen.Fini()
}

func eventLoop(screen tcell.Screen, w, h *int) {
	for {
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			if ev.Key() == tcell.KeyCtrlC {
				return
			}
		case *tcell.EventResize:
			screen.Sync()
			*w, *h = screen.Size()
		default:
			panic(fmt.Sprintf("%T", ev))
		}
	}
}

func redraw(screen tcell.Screen, w, h int) {
	screen.Clear()

	info := currentState.Info()
	currentSong := info.CurrentSong
	currentSample := info.CurrentSample

	songLength := time.Duration(0)
	timeInSong := time.Duration(0)
	progressInSong := float64(0)
	if currentSong.startIndex != 0 && int64(currentSong.startIndex) < currentSample {
		sampleInSong := currentSample - int64(currentSong.startIndex) - info.PausesInCurrentSong
		timeInSong = time.Duration(sampleInSong) * time.Second / time.Duration(schedule.SampleRate) / time.Nanosecond
		if 0 < currentSong.length {
			progressInSong = float64(sampleInSong) / float64(currentSong.length)
		}
		songLength = time.Duration(currentSong.length) * time.Second / time.Duration(schedule.SampleRate) / time.Nanosecond
	}

	playState := "Paused"
	if info.Playing {
		playState = "Playing"
	}

	songLineName := ""
	songLineArtistAlbum := ""

	timeLine := fmt.Sprintf("%s/%s", fmtDuration(timeInSong), fmtDuration(songLength))

	if currentSong.metadata.Title != "" {
		songLineName = currentSong.metadata.Title
	} else {
		songLineName = currentSong.filename
	}

	if currentSong.metadata.Artist != "" {
		songLineArtistAlbum = currentSong.metadata.Artist
		if currentSong.metadata.Album != "" {
			songLineArtistAlbum += " - " + currentSong.metadata.Album
		}
	}

	volumeLine := fmt.Sprintf("Volume: %06.2f%%", currentState.Volume*100)

	drawString(w-len(volumeLine)-1, h-4, tcell.StyleDefault, volumeLine, screen)
	drawString(1, h-4, tcell.StyleDefault, songLineName, screen)
	drawString(w-len(timeLine)-1, h-3, tcell.StyleDefault, timeLine, screen)
	drawString(1, h-3, tcell.StyleDefault, songLineArtistAlbum, screen)

	drawProgress(1, h-2, tcell.StyleDefault, w-2, progressInSong, screen)

	drawBox(0, h-5, w, 5, tcell.StyleDefault, screen)
	drawString(2, h-5, tcell.StyleDefault, playState, screen)

	lyricsHeight := lyricsHistorySize
	if h < lyricsHeight+7 {
		lyricsHeight = h - 7
	}
	if 0 < lyricsHeight {
		if currentSong.lyrics != nil && 0 < len(currentSong.lyrics) {
			nextLine := 0
			for ; nextLine < len(currentSong.lyrics); nextLine++ {
				l := currentSong.lyrics[nextLine]
				if l != nil && 0 < len(l) && int64(timeInSong/time.Millisecond) < l[0].Timestamp {
					break
				}
			}

			lines := make([]string, lyricsHeight)

			for i := range lines {
				lines[i] = ""
				if 0 <= nextLine-i-1 {
					for _, atom := range currentSong.lyrics[nextLine-i-1] {
						if atom.Timestamp < int64(timeInSong/time.Millisecond)+100 {
							lines[i] += atom.Caption
						}
					}
				}
			}

			for i, l := range lines {
				drawString(1, h-7-i, tcell.StyleDefault, l, screen)
			}
		}
		drawBox(0, h-7-lyricsHeight, w, lyricsHeight+2, tcell.StyleDefault, screen)
		drawString(2, h-7-lyricsHeight, tcell.StyleDefault, "Lyrics", screen)
	}

	screen.Show()
}
