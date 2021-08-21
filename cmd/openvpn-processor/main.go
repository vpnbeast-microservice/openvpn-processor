package main

import (
	_ "github.com/go-sql-driver/mysql"
	commons "github.com/vpnbeast/golang-commons"
	"go.uber.org/zap"
	"openvpn-processor/internal/metrics"
	"openvpn-processor/internal/options"
	"openvpn-processor/internal/probe"
	"openvpn-processor/internal/scheduler"
	"time"
)

var (
	opts   *options.OpenvpnProcessorOptions
	logger *zap.Logger
)

func init() {
	opts = options.GetOpenvpnProcessorOptions()
	logger = commons.GetLogger()
}

func main() {
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	go func() {
		probe.RunHealthProbe()
	}()

	go metrics.RunMetricsServer()

	go func() {
		// calling for the instant run before ticker ticks
		scheduler.RunBackground()
		ticker := time.NewTicker(time.Duration(int32(opts.TickerIntervalMin)) * time.Minute)
		for range ticker.C {
			scheduler.RunBackground()
		}
	}()
	select {}
}
