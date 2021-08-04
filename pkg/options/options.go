package options

import (
	commons "github.com/vpnbeast/golang-commons"
	"go.uber.org/zap"
)

var (
	logger *zap.Logger
	opts   *OpenvpnProcessorOptions
)

func init() {
	logger = commons.GetLogger()
	opts = newOpenvpnProcessorOptions()
	err := commons.InitOptions(opts, "openvpn-processor")
	if err != nil {
		logger.Fatal("fatal error occured while initializing options", zap.Error(err))
	}
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
	VpnGateUrl string `env:"VPNGATE_URL"`
	// database related environment variables
	DbUrl                    string `env:"DB_URL"`
	DbDriver                 string `env:"DB_DRIVER"`
	TickerIntervalMin        int    `env:"TICKER_INTERVAL_MIN"`
	DbMaxOpenConn            int    `env:"DB_MAX_OPEN_CONN"`
	DbMaxIdleConn            int    `env:"DB_MAX_IDLE_CONN"`
	DbConnMaxLifetimeMin     int    `env:"DB_CONN_MAX_LIFETIME_MIN"`
	HealthCheckMaxTimeoutMin int    `env:"HEALTH_CHECK_MAX_TIMEOUT_MIN"`
	DialTcpTimeoutSeconds    int    `env:"DIAL_TCP_TIMEOUT_SECONDS"`
	// metric server related environment variables
	MetricsPort         int    `env:"METRICS_PORT"`
	MetricsEndpoint     string `env:"METRICS_ENDPOINT"`
	WriteTimeoutSeconds int    `env:"WRITE_TIMEOUT_SECONDS"`
	ReadTimeoutSeconds  int    `env:"READ_TIMEOUT_SECONDS"`
	// health server related environment variables
	HealthPort     int    `env:"HEALTH_PORT"`
	HealthEndpoint string `env:"HEALTH_ENDPOINT"`
}