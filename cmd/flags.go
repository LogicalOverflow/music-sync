package cmd

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/logging"
	"github.com/urfave/cli"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Default Values for cli flags
const (
	DefaultMusicAddress = "127.0.0.1"
	DefaultMusicPort    = 13333

	DefaultSSHAddress   = DefaultMusicAddress
	DefaultSSHPort      = 13334
	DefaultSSHUsersFile = "users.json"
	DefaultSSHKeyFile   = "id_rsa"

	DefaultAudioDir = "audio"

	DefaultSampleRate = 44100

	DefaultTimeSyncInterval   = 10 * time.Minute
	DefaultTimeSyncCycles     = 500
	DefaultTimeSyncCycleDelay = 10 * time.Millisecond
	DefaultStreamChunkSize    = DefaultSampleRate * 4
	DefaultNanBreakSize       = DefaultSampleRate * 1
	DefaultStreamStartDelay   = 5 * time.Second
	DefaultStreamDelay        = 15 * time.Second
)

// TODO: refine logging
// TODO: check for old commented code

// LoggingFlag is a cli flag which controls log output
type LoggingFlag struct {
	cli.StringFlag
	logName string
}

var loggingFlags = []LoggingFlag{
	newLoggingFlag("comm"),
	newLoggingFlag("play"),
	newLoggingFlag("shed"),
	newLoggingFlag("ssh"),
	newLoggingFlag("time"),
}

// AddLoggingFlags adds the logging cli flags
func AddLoggingFlags(f []cli.Flag) []cli.Flag {
	for _, l := range loggingFlags {
		f = append(f, l)
	}
	f = append(f, cli.StringFlag{
		Name:  "logging",
		Usage: "sets the default logging level",
	})
	return f
}

func newLoggingFlag(name string) LoggingFlag {
	l := make([]string, 0, len(log.LevelNames))
	for _, n := range log.LevelNames {
		l = append(l, strings.ToLower(n.Full))
	}
	levels := strings.Join(l, ", ")

	return LoggingFlag{StringFlag: cli.StringFlag{
		Name:  fmt.Sprintf("%s-logging", name),
		Usage: fmt.Sprintf("sets %s logging level (values: %s)", name, levels),
	}}
}

// HandleLoggingFlags sets the log levels according to the logging cli flags
func HandleLoggingFlags(ctx *cli.Context) {
	for _, f := range loggingFlags {
		levelName := ctx.String(FlagKey(f))
		if levelName != "" {
			level := log.LevelByName(levelName)
			log.CutoffLevels[f.Name] = level
		}
	}
	defaultLevelName := ctx.String("logging")
	if defaultLevelName != "" {
		defaultLevel := log.LevelByName(defaultLevelName)
		log.DefaultCutoffLevel = defaultLevel
	}
}

func init() {

}

// NewApp creates a new cli app with same defaults
func NewApp(usage string) *cli.App {
	app := cli.NewApp()
	app.Name = filepath.Base(os.Args[0])
	app.Version = Version
	app.Author = Author
	app.EnableBashCompletion = true
	app.Usage = usage
	return app
}

var (
	// ServerAddressFlag is a flag for the server's address
	ServerAddressFlag = cli.StringFlag{
		Name:  "address, addr, a",
		Usage: "the master's address",
		Value: DefaultMusicAddress,
	}
	// ServerPortFlag is a flag for the server's port
	ServerPortFlag = cli.IntFlag{
		Name:  "port, p",
		Usage: "the master's port",
		Value: DefaultMusicPort,
	}

	// ListenAddressFlag is a flag for the address the master should listen on
	ListenAddressFlag = cli.StringFlag{
		Name:  "address, addr, a",
		Usage: "the address to listen on",
		Value: DefaultMusicAddress,
	}
	// ListenPortFlag is a flag for the port the master should listen on
	ListenPortFlag = cli.IntFlag{
		Name:  "port, p",
		Usage: "the port to listen on",
		Value: DefaultMusicPort,
	}

	// MusicDirFlag is a flag for the directory containing music files
	MusicDirFlag = cli.StringFlag{
		Name:  "music-dir, dir, d",
		Usage: "the directory containing the music files",
		Value: DefaultAudioDir,
	}

	// SSHAddressFlag is a flag for the master's ssh server address
	SSHAddressFlag = cli.StringFlag{
		Name:  "ssh-address, ssh-addr, sa",
		Usage: "the address to listen for ssh connections on",
		Value: DefaultSSHAddress,
	}
	// SSHPortFlag is a flag for the master's ssh server port
	SSHPortFlag = cli.IntFlag{
		Name:  "ssh-port, sp",
		Usage: "the port to listen for ssh connections on",
		Value: DefaultSSHPort,
	}
	// SSHUsersFlag is a flag for the master's ssh users file
	SSHUsersFlag = cli.StringFlag{
		Name:  "users-file, users, u",
		Usage: "the json file containing the ssh users",
		Value: DefaultSSHUsersFile,
	}
	// SSHKeyFileFlag is a flag for the master's ssh host key file
	SSHKeyFileFlag = cli.StringFlag{
		Name:  "host-key, key, k",
		Usage: "the file containing the host key for the ssh server",
		Value: DefaultSSHKeyFile,
	}

	// TimeSyncIntervalFlag is a flag for the time sync interval
	TimeSyncIntervalFlag = cli.DurationFlag{
		Name:  "time-sync-interval",
		Usage: "interval between synchronizing time to master",
		Value: DefaultTimeSyncInterval,
	}
	// TimeSyncCyclesFlag is a flag for the time sync cycle count
	TimeSyncCyclesFlag = cli.IntFlag{
		Name:  "time-sync-cycles",
		Usage: "number of time synchronization cycles",
		Value: DefaultTimeSyncCycles,
	}
	// TimeSyncCycleDelayFlag is a flag for the time sync cycle delay
	TimeSyncCycleDelayFlag = cli.DurationFlag{
		Name:  "time-sync-cycle-delay",
		Usage: "time to wait between time synchronization cycles",
		Value: DefaultTimeSyncCycleDelay,
	}

	// StreamChunkSizeFlag is a flag for the stream chunk size
	StreamChunkSizeFlag = cli.IntFlag{
		Name:  "stream-chunk-size",
		Usage: "number of samples per stream chunk",
		Value: DefaultStreamChunkSize,
	}
	// StreamStartDelayFlag is a flag for the stream start delay
	StreamStartDelayFlag = cli.DurationFlag{
		Name:  "stream-start-delay",
		Usage: "time to wait before beginning streaming",
		Value: DefaultStreamStartDelay,
	}
	// StreamDelayFlag is a flag for the stream delay
	StreamDelayFlag = cli.DurationFlag{
		Name:  "stream-delay",
		Usage: "delay between streaming a chunk and playing it",
		Value: DefaultStreamDelay,
	}

	// NanBreakSizeFlag is a flag for the nan break size between songs
	NanBreakSizeFlag = cli.IntFlag{
		Name:  "song-break-size",
		Usage: "number of samples between songs",
		Value: DefaultNanBreakSize,
	}

	// SampleRateFlag is a flag for the sample rate
	SampleRateFlag = cli.IntFlag{
		Name:  "sample-rate",
		Usage: "the sample rate of the stream",
		Value: DefaultSampleRate,
	}
)

// FlagKey returns the key to retrieve the flags value from the cli context
func FlagKey(flag cli.Flag) string {
	return strings.Split(flag.GetName(), ",")[0]
}
