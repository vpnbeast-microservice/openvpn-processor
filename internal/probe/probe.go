package probe

import (
	"fmt"
	"github.com/dimiro1/health"
	"github.com/dimiro1/health/db"
	"github.com/gorilla/mux"
	commons "github.com/vpnbeast/golang-commons"
	"go.uber.org/zap"
	"net/http"
	"openvpn-processor/internal/options"
	"openvpn-processor/internal/scheduler"
)

var (
	logger *zap.Logger
	opts   *options.OpenvpnProcessorOptions
)

func init() {
	logger = commons.GetLogger()
	opts = options.GetOpenvpnProcessorOptions()
}

// RunHealthProbe spins up a router and continuously checks the health of database connection
func RunHealthProbe() {
	router := mux.NewRouter()
	mysql := db.NewMySQLChecker(scheduler.GetDatabase())

	handler := health.NewHandler()
	handler.AddChecker("MySQL", mysql)
	router.Handle(opts.HealthEndpoint, handler)

	logger.Info("health server is up and running", zap.Int("port", 9290))
	panic(http.ListenAndServe(fmt.Sprintf(":%d", opts.HealthPort), router))
}
