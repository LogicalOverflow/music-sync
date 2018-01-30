// This package is the main package of the music-sync player
package main

import (
	"fmt"
	"github.com/LogicalOverflow/music-sync/cmd"
	"github.com/LogicalOverflow/music-sync/comm"
	"github.com/LogicalOverflow/music-sync/playback"
	"github.com/LogicalOverflow/music-sync/schedule"
	"github.com/LogicalOverflow/music-sync/timing"
	"github.com/urfave/cli"
	"net"
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
	sender, err := comm.ConnectToServer(server, newPlayerPackageHandler())
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	go schedule.Player(sender)

	cmd.WaitForInterrupt()
	return nil
}

type playerPackageHandler struct{}

func (c playerPackageHandler) HandleTimeSyncResponse(tsr *comm.TimeSyncResponse, _ net.Conn) {
	clientRecv := timing.GetRawTime()
	timing.UpdateOffset(tsr.ClientSendTime, tsr.ServerRecvTime, tsr.ServerSendTime, clientRecv)
}

func (c playerPackageHandler) HandleQueueChunkRequest(qsr *comm.QueueChunkRequest, _ net.Conn) { playback.QueueChunk(qsr.StartTime, qsr.ChunkId, playback.CombineSamples(qsr.SampleLow, qsr.SampleHigh)) }
func (c playerPackageHandler) HandleSetVolumeRequest(svr *comm.SetVolumeRequest, _ net.Conn)   { playback.SetVolume(svr.Volume) }
func (c playerPackageHandler) HandlePingMessage(_ *comm.PingMessage, conn net.Conn)            { comm.PingHandler(conn) }

func (c playerPackageHandler) HandleTimeSyncRequest(*comm.TimeSyncRequest, net.Conn)                 {}
func (c playerPackageHandler) HandlePongMessage(*comm.PongMessage, net.Conn)                         {}
func (c playerPackageHandler) HandleSubscribeChannelRequest(*comm.SubscribeChannelRequest, net.Conn) {}
func (c playerPackageHandler) HandleNewSongInfo(*comm.NewSongInfo, net.Conn)                         {}
func (c playerPackageHandler) HandleChunkInfo(*comm.ChunkInfo, net.Conn)                             {}
func (c playerPackageHandler) HandlePauseInfo(*comm.PauseInfo, net.Conn)                             {}

func newPlayerPackageHandler() comm.TypedPackageHandler {
	return comm.TypedPackageHandler{playerPackageHandler{}}
}
