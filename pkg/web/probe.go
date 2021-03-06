package web

import (
	"context"
	"database/sql"
	"github.com/etherlabsio/healthcheck"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
)

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
	log.Println("Listening on port 9290 for health probes...")
	log.Fatal(http.ListenAndServe(":9290" , router))
}