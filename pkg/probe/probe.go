package probe

import (
	"context"
	"database/sql"
	"github.com/etherlabsio/healthcheck"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
	"net/http"
	"openvpn-processor/pkg/logging"
	"time"
)

var logger *zap.Logger

func init() {
	logger = logging.GetLogger()
}

func RunHealthProbe(db *sql.DB, healthCheckMaxTimeoutMin int) {
	router := mux.NewRouter()
	router.Handle("/health", healthcheck.Handler(
		// WithTimeout allows you to set a max overall timeout.
		healthcheck.WithTimeout(time.Duration(int32(healthCheckMaxTimeoutMin)) * time.Second),

		healthcheck.WithChecker(
			"database", healthcheck.CheckerFunc(
				func(ctx context.Context) error {
					return db.PingContext(ctx)
				},
			),
		),
	))

	err := http.ListenAndServe(":9290" , router)
	if err != nil {
		logger.Fatal("fatal error occured while spinning up router", zap.String("addr", ":9290"),
			zap.String("error", err.Error()))
	}

	logger.Info("listening on port 9290 for health probes")
}