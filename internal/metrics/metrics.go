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
	// SkippedServerCounter keeps track of skipped vpn servers on last scheduled execution
	SkippedServerCounter prometheus.Counter
	// AvailableServerCounter keeps track of available vpn servers on last scheduled execution
	AvailableServerCounter prometheus.Counter
	// FailedServerCounter keeps track of failed vpn servers on last scheduled execution
	FailedServerCounter prometheus.Counter
)

func init() {
	logger = commons.GetLogger()
	opts = options.GetOpenvpnProcessorOptions()
	AvailableServerCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "available_server_count",
		Help: "Counts processed server count on last scheduled execution",
	})
	SkippedServerCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "skipped_server_count",
		Help: "Counts skipped server count on last scheduled execution",
	})
	FailedServerCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "failed_server_count",
		Help: "Counts failed server count on last scheduled execution",
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
	prometheus.MustRegister(AvailableServerCounter)
	prometheus.MustRegister(SkippedServerCounter)
	prometheus.MustRegister(FailedServerCounter)

	logger.Info("metric server is up and running", zap.Int("port", opts.MetricsPort))
	panic(metricServer.ListenAndServe())
}
