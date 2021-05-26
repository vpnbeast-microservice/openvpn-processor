package probe

import (
	"context"
	"github.com/etherlabsio/healthcheck"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"openvpn-processor/pkg/logging"
	"openvpn-processor/pkg/options"
	"openvpn-processor/pkg/scheduler"
	"time"
)

var (
	logger *zap.Logger
	opts   *options.OpenvpnProcessorOptions
)

func init() {
	logger = logging.GetLogger()
	opts = options.GetOpenvpnProcessorOptions()
}

// RunHealthProbe spins up a router and continuously checks the health of database connection
func RunHealthProbe() {
	router := mux.NewRouter()
	router.Handle("/health", healthcheck.Handler(
		healthcheck.WithTimeout(time.Duration(int32(opts.HealthCheckMaxTimeoutMin))*time.Second),
		healthcheck.WithChecker(
			"database", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					return scheduler.GetDatabase().PingContext(ctx)
				},
			),
		),
	))

	// TODO: get variables here in that package from environment variables
	logger.Info("metric server is up and running", zap.Int("port", 9290))
	panic(http.ListenAndServe(":9290", router))
}
