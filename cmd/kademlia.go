package cmd

import (
	"errors"
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/dewebprotocol/IBS/cmd/flags"
	"github.com/dewebprotocol/IBS/network"
	"github.com/dewebprotocol/IBS/simulator"
	"github.com/urfave/cli/v2"
)

var Kademlia = &cli.Command{
	Name: "kademlia",
	Subcommands: []*cli.Command{
		kadLatency,
		kadCoverage,
	},
}

func datetime() string {
	t := strings.Split(time.Now().String(), ".")[0]
	return strings.ReplaceAll(t, ":", ".")
}

var kadLatency = &cli.Command{
	Name: "latency",
	Flags: []cli.Flag{
		flags.K,
		flags.Beta,
		flags.BroadcastPerNode,
		flags.BroadcastInterval,
		flags.NetSize,
		flags.WithNE,
	},
	Action: func(ctx *cli.Context) error {
		netSize := ctx.Int("net_size")
		broadcastPerNode := ctx.Int("broadcast_per_node")
		nMessage := broadcastPerNode * netSize
		broadcastInterval := ctx.Int("broadcast_interval")
		k := ctx.Int("k")
		beta := ctx.Int("beta")
		var crashInterval = 3000000_000_000 // 3000000s
		withNE := ctx.Bool("with_ne")

		folder := fmt.Sprintf(
			"kad_latency_%s(beta=%d,interval=%d)",
			datetime(), beta, broadcastInterval,
		)
		if withNE {
			folder = "ne" + folder
		}
		fmt.Println("output will be write to path:", folder)
		os.Mkdir(folder, 0777)
		logFile, err := os.Create(folder + "/output_log.txt")
		if err != nil {
			return err
		}
		configText := "===== kademlia =====\n"
		if withNE {
			configText = "===== ne-kademlia =====\n"
		}
		configText += fmt.Sprintf("beta: %d, bucket size: %d\n", beta, k)
		configText += fmt.Sprintf("NetSize: %d, NMessage: %d, BroadcastInterval: %d(μs)\n",
			netSize, nMessage, broadcastInterval)
		logFile.Write([]byte(configText))
		fmt.Println(configText)
		var InitNet func(size int, k int, beta int) network.Network
		InitNet = network.NewKadcastNet
		if withNE {
			InitNet = network.NewNecastNet
		}
		net := InitNet(netSize, k, beta)
		sim := simulator.New(net, nMessage, netSize, crashInterval, broadcastInterval)
		sim.InitState()
		sim.Run(true, false)
		outputText := sim.Statistic()
		logFile.Write([]byte(outputText))
		fmt.Println(outputText)
		sim.OutputNodes(folder)
		sim.OutputLatency(folder)
		return nil
	},
}

var kadCoverage = &cli.Command{
	Name: "coverage",
	Flags: []cli.Flag{
		flags.K,
		flags.Beta,
		flags.BroadcastPerNode,
		flags.BroadcastInterval,
		flags.CrashInterval,
		flags.NetSize,
		flags.WithNE,
		flags.MaliciousNode,
	},
	Action: func(ctx *cli.Context) error {
		netSize := ctx.Int("net_size")
		broadcastPerNode := ctx.Int("broadcast_per_node")
		nMessage := broadcastPerNode * netSize
		broadcastInterval := ctx.Int("broadcast_interval")
		k := ctx.Int("k")
		beta := ctx.Int("beta")
		crashInterval := ctx.Int("crash_interval")
		withNE := ctx.Bool("with_ne")
		malicious := ctx.Bool("malicious")

		if malicious && crashInterval < math.MaxInt {
			return errors.New("bad instruction")
		}

		folder := fmt.Sprintf(
			"kad_coverage_%s(beta=%d,interval=%d)",
			datetime(), beta, broadcastInterval,
		)
		if withNE {
			folder = "ne" + folder
		}
		fmt.Println("output will be write to path:", folder)
		os.Mkdir(folder, 0777)
		logFile, err := os.Create(folder + "/output_log.txt")
		if err != nil {
			return err
		}
		configText := "===== kademlia =====\n"
		if withNE {
			configText = "===== ne-kademlia =====\n"
		}
		configText += fmt.Sprintf("beta: %d, bucket size: %d\n", beta, k)
		configText += fmt.Sprintf("NetSize: %d, NMessage: %d, BroadcastInterval: %d(μs)\n",
			netSize, nMessage, broadcastInterval)
		fmt.Println(configText)
		var InitNet func(size int, k int, beta int) network.Network
		InitNet = network.NewKadcastNet
		if withNE {
			InitNet = network.NewNecastNet
		}
		net := InitNet(netSize, k, beta)
		if malicious {
			cntInfest := net.Infest(1)
			fmt.Println("infested: ", cntInfest)
		} else {
			cntCrash := net.Churn(1, !(crashInterval < math.MaxInt))
			fmt.Println("first crashed: ", cntCrash)
		}
		sim := simulator.New(net, nMessage, netSize, crashInterval, broadcastInterval)
		sim.InitState()
		sim.Run(false, false)
		outputText := sim.Statistic()
		logFile.Write([]byte(configText))
		logFile.Write([]byte(outputText))
		fmt.Println(outputText)
		sim.OutputNodes(folder)
		sim.OutputCoverage(folder)
		return nil
	},
}
