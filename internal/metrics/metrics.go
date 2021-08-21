package metrics

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	commons "github.com/vpnbeast/golang-commons"
	"go.uber.org/zap"
	"net/http"
	"openvpn-processor/internal/options"
	"time"
)

var (
	logger *zap.Logger
	opts   *options.OpenvpnProcessorOptions
	// SkippedCounter keeps track of skipped vpn servers for various reasons
	SkippedCounter prometheus.Counter
	// InsertedCounter keeps track of inserted vpn servers to database
	InsertedCounter prometheus.Counter
)

func init() {
	logger = commons.GetLogger()
	opts = options.GetOpenvpnProcessorOptions()
	InsertedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "inserted_server_count",
		Help: "Counts processed server count on last scheduled execution",
	})
	SkippedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "skipped_server_count",
		Help: "Counts skipped server count on last scheduled execution",
	})
}

// RunMetricsServer spins up a router to provide prometheus metrics
func RunMetricsServer() {
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	router := mux.NewRouter()
	metricServer := &http.Server{
		Handler:      router,
		Addr:         fmt.Sprintf(":%d", opts.MetricsPort),
		WriteTimeout: time.Duration(int32(opts.WriteTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(opts.ReadTimeoutSeconds)) * time.Second,
	}
	router.Handle(opts.MetricsEndpoint, promhttp.Handler())
	prometheus.MustRegister(InsertedCounter)
	prometheus.MustRegister(SkippedCounter)

	logger.Info("metric server is up and running", zap.Int("port", opts.MetricsPort))
	panic(metricServer.ListenAndServe())
}
