package options

import (
	"go.uber.org/zap"
	"openvpn-processor/pkg/logging"
)

var (
	logger *zap.Logger
	opts   *OpenvpnProcessorOptions
)

func init() {
	logger = logging.GetLogger()
	opts = newOpenvpnProcessorOptions()
	opts.initOptions()
}

// GetOpenvpnProcessorOptions returns the initialized EncryptionServiceOptions
func GetOpenvpnProcessorOptions() *OpenvpnProcessorOptions {
	return opts
}

// newOpenvpnProcessorOptions creates an OpenvpnProcessorOptions struct with zero values
func newOpenvpnProcessorOptions() *OpenvpnProcessorOptions {
	return &OpenvpnProcessorOptions{}
}

type OpenvpnProcessorOptions struct {
	// api related environment variables
	VpnGateUrl string
	// database related environment variables
	DbUrl                    string
	DbDriver                 string
	TickerIntervalMin        int
	DbMaxOpenConn            int
	DbMaxIdleConn            int
	DbConnMaxLifetimeMin     int
	HealthCheckMaxTimeoutMin int
	DialTcpTimeoutSeconds    int
	// metric server related environment variables
	MetricsPort         int
	MetricsEndpoint     string
	WriteTimeoutSeconds int
	ReadTimeoutSeconds  int
	// health server related environment variables
	HealthPort			int
	HealthEndpoint		string
}

// initOptions initializes EncryptionServiceOptions while reading environment values, sets default values if not specified
func (opo *OpenvpnProcessorOptions) initOptions() {
	opo.VpnGateUrl = getStringEnv("API_URL", "https://www.vpngate.net/api/iphone/")
	opo.DbUrl = getStringEnv("DB_URL", "spring:123asd456@tcp(127.0.0.1:3306)/vpnbeast?parseTime=true&loc=Local")
	opo.DbDriver = getStringEnv("DB_DRIVER", "mysql")
	opo.TickerIntervalMin = getIntEnv("TICKER_INTERVAL_MIN", 10)
	opo.DbMaxOpenConn = getIntEnv("DB_MAX_OPEN_CONN", 25)
	opo.DbMaxIdleConn = getIntEnv("DB_MAX_IDLE_CONN", 25)
	opo.DbConnMaxLifetimeMin = getIntEnv("DB_CONN_MAX_LIFETIME_MIN", 5)
	opo.HealthCheckMaxTimeoutMin = getIntEnv("HEALTHCHECK_MAX_TIMEOUT_MIN", 5)
	opo.DialTcpTimeoutSeconds = getIntEnv("DIAL_TCP_TIMEOUT_SECONDS", 5)
	opo.MetricsPort = getIntEnv("METRICS_PORT", 3001)
	opo.MetricsEndpoint = getStringEnv("METRICS_ENDPOINT", "/metrics")
	opo.WriteTimeoutSeconds = getIntEnv("WRITE_TIMEOUT_SECONDS", 10)
	opo.ReadTimeoutSeconds = getIntEnv("READ_TIMEOUT_SECONDS", 10)
	opo.HealthPort = getIntEnv("HEALTH_PORT", 9290)
	opo.HealthEndpoint = getStringEnv("HEALTH_ENDPOINT", "/health")
}

