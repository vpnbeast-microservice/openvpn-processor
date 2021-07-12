package options

import (
	"go.uber.org/zap"
	"openvpn-processor/pkg/logging"
	"os"
	"strconv"
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
	opo.DbUrl = getStringEnv("DB_URL", "spring:123asd456@tcp(127.0.0.1:3306)/vpnbeast")
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

// getStringEnv gets the specific environment variables with default value, returns default value if variable not set
func getStringEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return value
}

// getIntEnv gets the specific environment variables with default value, returns default value if variable not set
func getIntEnv(key string, defaultValue int) int {
	value := os.Getenv(key)
	if len(value) == 0 {
		return defaultValue
	}
	return ConvertStringToInt(value)
}

// ConvertStringToInt converts string environment variables to integer values
func ConvertStringToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		logger.Warn("an error occurred while converting from string to int. Setting it as zero",
			zap.String("error", err.Error()))
		i = 0
	}
	return i
}
