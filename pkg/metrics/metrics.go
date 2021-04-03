package metrics

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
	"net/http"
	"openvpn-processor/pkg/config"
	"openvpn-processor/pkg/logging"
	"time"
)

var (
	logger *zap.Logger
	metricsPort, writeTimeoutSeconds, readTimeoutSeconds int
	SkippedCounter prometheus.Counter
	InsertedCounter prometheus.Counter
)

func init() {
	logger = logging.GetLogger()
	metricsPort = config.GetIntEnv("METRICS_PORT", 3001)
	writeTimeoutSeconds = config.GetIntEnv("WRITE_TIMEOUT_SECONDS", 10)
	readTimeoutSeconds = config.GetIntEnv("READ_TIMEOUT_SECONDS", 10)
	InsertedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "inserted_server_count",
		Help: "Counts processed server count on last scheduled execution",
	})
	SkippedCounter = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "skipped_server_count",
		Help: "Counts skipped server count on last scheduled execution",
	})
}

// TODO: Generate custom metrics, check below:
// https://prometheus.io/docs/guides/go-application/
// https://www.robustperception.io/prometheus-middleware-for-gorilla-mux

func RunMetricsServer() {
	defer func() {
		err := logger.Sync()
		if err != nil {
			panic(err)
		}
	}()

	router := mux.NewRouter()
	metricServer := &http.Server{
		Handler: router,
		Addr: fmt.Sprintf(":%d", metricsPort),
		WriteTimeout: time.Duration(int32(writeTimeoutSeconds)) * time.Second,
		ReadTimeout:  time.Duration(int32(readTimeoutSeconds)) * time.Second,
	}
	router.Handle("/metrics", promhttp.Handler())
	prometheus.MustRegister(InsertedCounter)
	prometheus.MustRegister(SkippedCounter)

	logger.Info("metric server is up and running", zap.Int("port", metricsPort))
	panic(metricServer.ListenAndServe())
}