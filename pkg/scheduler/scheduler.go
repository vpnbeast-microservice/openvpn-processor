package scheduler

import (
	"database/sql"
	"go.uber.org/zap"
	"openvpn-processor/pkg/logging"
	"time"
)

var logger *zap.Logger

func init() {
	logger = logging.GetLogger()
}

// RunBackground does the heavy lifting, continuously repeats the application logic
func RunBackground(db *sql.DB, vpnGateUrl string, dialTcpTimeoutSeconds int) {
	logger.Info("Starting scheduler execution")
	beforeMainExecution := time.Now()
	csvContent := getCsvContent(vpnGateUrl)
	vpnServers := createStructsFromCsv(csvContent)
	checkUnreachableServersOnDB(db, dialTcpTimeoutSeconds)
	insertServers(db, vpnServers, dialTcpTimeoutSeconds)
	logger.Info("Ending scheduler execution", zap.Duration("executionTime", time.Since(beforeMainExecution)))
}
