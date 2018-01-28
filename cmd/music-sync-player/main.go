// This package is the main package of the music-sync player
package main

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/cmd"
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/schedule"
	"github.com/urfave/cli"
	"os"
)

const usage = "run a music-sync player, which connects to a server and plays music"

func main() {
	app := cmd.NewApp(usage)
	app.Action = run
	app.Flags = cmd.AddLoggingFlags([]cli.Flag{
		cmd.ServerAddressFlag,
		cmd.ServerPortFlag,

		cmd.SampleRateFlag,
	})

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx *cli.Context) error {
	cmd.HandleLoggingFlags(ctx)
	var (
		serverAddress = ctx.String(cmd.FlagKey(cmd.ServerAddressFlag))
		serverPort    = ctx.Int(cmd.FlagKey(cmd.ServerPortFlag))

		sampleRate = ctx.Int(cmd.FlagKey(cmd.SampleRateFlag))
	)

	schedule.SampleRate = sampleRate

	server := fmt.Sprintf("%s:%d", serverAddress, serverPort)
	sender, err := comm.ConnectToServer(server)
	if err != nil {
		cli.NewExitError(err, 1)
	}

	go schedule.Player(sender)

	cmd.WaitForInterrupt()
	return nil
}
