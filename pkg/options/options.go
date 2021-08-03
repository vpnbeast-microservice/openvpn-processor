package options

import (
	"fmt"
	"github.com/spf13/viper"
	commons "github.com/vpnbeast/golang-commons"
	"go.uber.org/zap"
	"net/http"
)

var (
	logger *zap.Logger
	opts   *OpenvpnProcessorOptions
)

func init() {
	logger = commons.GetLogger()
	opts = newOpenvpnProcessorOptions()
	err := opts.initOptions()
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

// initOptions initializes EncryptionServiceOptions while reading environment values, sets default values if not specified
func (opo *OpenvpnProcessorOptions) initOptions() error {
	activeProfile := commons.GetStringEnv("ACTIVE_PROFILE", "local")
	appName := commons.GetStringEnv("APP_NAME", "openvpn-processor")
	if activeProfile == "unit-test" {
		logger.Info("active profile is unit-test, reading configuration from static file")
		// TODO: better approach for that?
		viper.AddConfigPath("./../../config")
		viper.SetConfigName("unit_test")
		viper.SetConfigType("yaml")
		if err := viper.ReadInConfig(); err != nil {
			return err
		}
	} else {
		configHost := commons.GetStringEnv("CONFIG_SERVER_HOST", "localhost")
		configPort := commons.GetIntEnv("CONFIG_SERVER_PORT", 8888)
		logger.Info("loading configuration from remote server", zap.String("host", configHost),
			zap.Int("port", configPort), zap.String("appName", appName),
			zap.String("activeProfile", activeProfile))
		confAddr := fmt.Sprintf("http://%s:%d/%s-%s.yaml", configHost, configPort, appName, activeProfile)
		resp, err := http.Get(confAddr)
		if err != nil {
			return err
		}

		defer func() {
			err := resp.Body.Close()
			if err != nil {
				panic(err)
			}
		}()

		viper.SetConfigName("application")
		viper.SetConfigType("yaml")
		if err = viper.ReadConfig(resp.Body); err != nil {
			return err
		}
	}

	if err := commons.UnmarshalConfig(appName, opo); err != nil {
		return err
	}

	return nil
}
