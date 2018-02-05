// This package is the main package of the music-sync server
package main

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/cmd"
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/schedule"
	"github.com/LogicalOverflow/music-sync/ssh"
	"github.com/LogicalOverflow/music-sync/util"
	"github.com/urfave/cli"
	"os"
	"time"
)

const usage = "run a music-sync server for clients to connect to"

func main() {
	app := cmd.NewApp(usage)
	app.Flags = cmd.AddLoggingFlags([]cli.Flag{
		cmd.ListenAddressFlag,
		cmd.ListenPortFlag,
		cmd.MusicDirFlag,
		cmd.SSHAddressFlag,
		cmd.SSHPortFlag,
		cmd.SSHUsersFlag,
		cmd.SSHKeyFileFlag,

		cmd.TimeSyncIntervalFlag,
		cmd.TimeSyncCyclesFlag,
		cmd.TimeSyncCycleDelayFlag,
		cmd.StreamChunkSizeFlag,
		cmd.StreamStartDelayFlag,
		cmd.StreamDelayFlag,
		cmd.NanBreakSizeFlag,
		cmd.SampleRateFlag,
	})
	app.Action = run

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func setScheduleVars(ctx *cli.Context) {
	var (
		timeSyncInterval   = ctx.Duration(cmd.FlagKey(cmd.TimeSyncIntervalFlag))
		timeSyncCycles     = ctx.Int(cmd.FlagKey(cmd.TimeSyncCyclesFlag))
		timeSyncCycleDelay = ctx.Duration(cmd.FlagKey(cmd.TimeSyncCycleDelayFlag))
		streamChunkSize    = ctx.Int(cmd.FlagKey(cmd.StreamChunkSizeFlag))
		streamStartDelay   = ctx.Duration(cmd.FlagKey(cmd.StreamStartDelayFlag))
		streamDelay        = ctx.Duration(cmd.FlagKey(cmd.StreamDelayFlag))
		nanBreakSize       = ctx.Int(cmd.FlagKey(cmd.NanBreakSizeFlag))
		sampleRate         = ctx.Int(cmd.FlagKey(cmd.SampleRateFlag))
	)

	schedule.TimeSyncInterval = timeSyncInterval
	schedule.TimeSyncCycles = timeSyncCycles
	schedule.TimeSyncCycleDelay = timeSyncCycleDelay
	schedule.StreamChunkSize = streamChunkSize
	schedule.StreamChunkTime = time.Duration(streamChunkSize) * time.Second / time.Duration(sampleRate)
	schedule.NanBreakSize = nanBreakSize
	schedule.StreamStartDelay = streamStartDelay
	schedule.StreamDelay = streamDelay
	schedule.SampleRate = sampleRate

}

func run(ctx *cli.Context) error {
	cmd.HandleLoggingFlags(ctx)
	var (
		listenAddress = ctx.String(cmd.FlagKey(cmd.ListenAddressFlag))
		listenPort    = ctx.Int(cmd.FlagKey(cmd.ListenPortFlag))
		musicDir      = ctx.String(cmd.FlagKey(cmd.MusicDirFlag))
		sshAddress    = ctx.String(cmd.FlagKey(cmd.SSHAddressFlag))
		sshPort       = ctx.Int(cmd.FlagKey(cmd.SSHPortFlag))
		sshUsers      = ctx.String(cmd.FlagKey(cmd.SSHUsersFlag))
		sshKeyFile    = ctx.String(cmd.FlagKey(cmd.SSHKeyFileFlag))
	)

	listen := fmt.Sprintf("%s:%d", listenAddress, listenPort)
	sshListen := fmt.Sprintf("%s:%d", sshAddress, sshPort)

	if err := util.CheckDir(musicDir); err != nil {
		return cli.NewExitError(fmt.Sprintf("invalid music dir: %v", err), 1)
	}

	playback.AudioDir = musicDir
	setScheduleVars(ctx)

	sender, err := comm.StartServer(listen)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	users, err := ssh.ReadUsersFile(sshUsers)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	ssh.HostKeyFile = sshKeyFile
	go ssh.StartSSH(sshListen, users)
	go schedule.Server(sender)

	cmd.WaitForInterrupt()
	return nil
}
