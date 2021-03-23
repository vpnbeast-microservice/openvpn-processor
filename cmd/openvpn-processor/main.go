package main

import (
	_ "github.com/go-sql-driver/mysql"
	"openvpn-processor/pkg/config"
	"openvpn-processor/pkg/probe"
	"openvpn-processor/pkg/scheduler"
	"time"
)

var (
	vpnGateUrl, dbUrl, dbDriver string
	tickerIntervalMin, dbMaxOpenConn, dbMaxIdleConn, dbConnMaxLifetimeMin, healthCheckMaxTimeoutMin, dialTcpTimeoutSeconds int
)

func init() {
	vpnGateUrl = config.GetStringEnv("API_URL", "https://www.vpngate.net/api/iphone/")
	dbUrl = config.GetStringEnv("DB_URL", "spring:123asd456@tcp(127.0.0.1:3306)/vpnbeast")
	dbDriver = config.GetStringEnv("DB_DRIVER", "mysql")
	tickerIntervalMin = config.GetIntEnv("TICKER_INTERVAL_MIN", 10)
	dbMaxOpenConn = config.GetIntEnv("DB_MAX_OPEN_CONN", 25)
	dbMaxIdleConn = config.GetIntEnv("DB_MAX_IDLE_CONN", 25)
	dbConnMaxLifetimeMin = config.GetIntEnv("DB_CONN_MAX_LIFETIME_MIN", 5)
	healthCheckMaxTimeoutMin = config.GetIntEnv("HEALTHCHECK_MAX_TIMEOUT_MIN", 5)
	dialTcpTimeoutSeconds = config.GetIntEnv("DIAL_TCP_TIMEOUT_SECONDS", 5)
}

func main() {
	db := scheduler.InitDb(dbDriver, dbUrl, dbMaxOpenConn, dbMaxIdleConn, dbConnMaxLifetimeMin)

	go func() {
		probe.RunHealthProbe(db, healthCheckMaxTimeoutMin)
	}()

	go func() {
		// calling for the instant run before ticker ticks
		scheduler.RunBackground(db, vpnGateUrl, dialTcpTimeoutSeconds)
		ticker := time.NewTicker(time.Duration(int32(tickerIntervalMin)) * time.Minute)
		for range ticker.C {
			scheduler.RunBackground(db, vpnGateUrl, dialTcpTimeoutSeconds)
		}
	}()
	select {}
}
