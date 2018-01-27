package main

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/cmd"
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/schedule"
	"github.com/urfave/cli"
	"os"
)

const usage = "run a music-sync slave, playing music"

func main() {
	app := cmd.NewApp(usage)
	app.Action = run
	app.Flags = cmd.AddLoggingFlags([]cli.Flag{
		cmd.MasterAddressFlag,
		cmd.MasterPortFlag,

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
		masterAddress = ctx.String(cmd.FlagKey(cmd.MasterAddressFlag))
		masterPort    = ctx.Int(cmd.FlagKey(cmd.MasterPortFlag))

		sampleRate = ctx.Int(cmd.FlagKey(cmd.SampleRateFlag))
	)

	schedule.SampleRate = sampleRate

	master := fmt.Sprintf("%s:%d", masterAddress, masterPort)
	sender, err := comm.ConnectToMaster(master)
	if err != nil {
		cli.NewExitError(err, 1)
	}

	go schedule.Slave(sender)

	cmd.WaitForInterrupt()
	return nil
}
