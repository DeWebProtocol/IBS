package cmd

import (
	"errors"
	"fmt"
	"math"
	"os"

	"github.com/dewebprotocol/IBS/cmd/flags"
	"github.com/dewebprotocol/IBS/network"
	"github.com/dewebprotocol/IBS/simulator"
	"github.com/urfave/cli/v2"
)

var Flood = &cli.Command{
	Name: "flood",
	Subcommands: []*cli.Command{
		floodLatency,
		floodCoverage,
	},
}

var floodLatency = &cli.Command{
	Name: "latency",
	Flags: []cli.Flag{
		flags.TableSize,
		flags.Degree,
		flags.BroadcastPerNode,
		flags.BroadcastInterval,
		flags.CrashInterval,
		flags.NetSize,
		flags.WithNE,
	},
	Action: func(ctx *cli.Context) error {
		tableSize := ctx.Int("table_size")
		degree := ctx.Int("degree")

		folder := fmt.Sprintf(
			"flood_latency_%s(tableSize=%d,degree=%d)",
			datetime(), tableSize, degree,
		)
		return floodProcess(ctx, false, folder)
	},
}

var floodCoverage = &cli.Command{
	Name: "coverage",
	Flags: []cli.Flag{
		flags.TableSize,
		flags.Degree,
		flags.BroadcastPerNode,
		flags.BroadcastInterval,
		flags.CrashInterval,
		flags.NetSize,
		flags.WithNE,
		flags.MaliciousNode,
		flags.OutputPacket,
		flags.NeLearn,
	},
	Action: func(ctx *cli.Context) error {
		tableSize := ctx.Int("table_size")
		degree := ctx.Int("degree")

		folder := fmt.Sprintf(
			"flood_coverage_%s(tableSize=%d,degree=%d)",
			datetime(), tableSize, degree,
		)
		return floodProcess(ctx, true, folder)
	},
}

func floodProcess(ctx *cli.Context, disturbNet bool, folder string) error {
	netSize := ctx.Int("net_size")
	//netSize := 30
	broadcastPerNode := ctx.Int("broadcast_per_node")
	nMessage := broadcastPerNode * netSize
	//nMessage := netSize * 3
	broadcastInterval := ctx.Int("broadcast_interval")
	tableSize := ctx.Int("table_size")
	degree := ctx.Int("degree")
	crashInterval := ctx.Int("crash_interval")
	withNE := ctx.Bool("with_ne")
	learn := ctx.Bool("ne_learn")
	outputPacket := ctx.Bool("output_packet")

	if !withNE && learn {
		return errors.New("bad instruction")
	}

	if disturbNet {
		malicious := ctx.Bool("malicious")
		if malicious {
			folder = folder + "_mal"
		}
		if malicious && crashInterval < math.MaxInt {
			return errors.New("bad instruction")
		}
	}

	if withNE {
		folder = "ne" + folder
	}

	fmt.Println("output will be write to path:", folder)
	os.Mkdir(folder, 0777)
	logFile, err := os.Create(folder + "/output_log.txt")
	if err != nil {
		return err
	}
	configText := "===== flood =====\n"
	if withNE {
		configText = "===== ne-flood =====\n"
	}
	configText += fmt.Sprintf("table size: %d, degree: %d\n", tableSize, degree)
	configText += fmt.Sprintf("NetSize: %d, NMessage: %d, BroadcastInterval: %d(μs)\n",
		netSize, nMessage, broadcastInterval)
	logFile.Write([]byte(configText))
	fmt.Println(configText)
	//net := network.NewNeFloodNet(netSize, tableSize, degree)
	var InitNet func(size int, k int, beta int) network.Network
	InitNet = network.NewFloodNet
	if withNE {
		InitNet = network.NewNeFloodNet
	}
	net := InitNet(netSize, tableSize, degree)
	var sim *simulator.Simulator
	if learn {
		sim = simulator.New(net, 10*nMessage, netSize, crashInterval, broadcastInterval)
	} else {
		sim = simulator.New(net, nMessage, netSize, crashInterval, broadcastInterval)
	}
	sim.InitState()
	if disturbNet {
		malicious := ctx.Bool("malicious")
		if malicious {
			cntInfest := net.Infest(1)
			logFile.Write([]byte(fmt.Sprintf("infested: %d\n", cntInfest)))
			fmt.Println("infested: ", cntInfest)
		} else {
			cntCrash := net.Churn(1, !(crashInterval < math.MaxInt))
			logFile.Write([]byte(fmt.Sprintf("first crashed: %d\n", cntCrash)))
			fmt.Println("first crashed: ", cntCrash)
		}
	}
	sim.Run(false, outputPacket && !learn)
	if learn {
		sim.InitState()
		sim.ResetNMsg(nMessage)
		sim.Run(false, outputPacket)
	}
	outputText := sim.Statistic()
	logFile.Write([]byte(outputText))
	fmt.Println(outputText)
	sim.OutputNodes(folder)
	sim.OutputLatency(folder)
	sim.OutputReceived(folder)
	if outputPacket {
		fmt.Println("output packet to:", folder)
		sim.OutputPackets(folder)
	}
	return nil
}
