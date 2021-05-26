package scheduler

import (
	"go.uber.org/zap"
	"time"
)

// RunBackground does the heavy lifting, continuously repeats the application logic
func RunBackground() {
	logger.Info("Starting scheduler execution")
	beforeMainExecution := time.Now()
	csvContent := getCsvContent(opts.VpnGateUrl)
	vpnServers := createStructsFromCsv(csvContent)
	checkUnreachableServersOnDB(db)
	insertServers(db, vpnServers)
	logger.Info("Ending scheduler execution", zap.Duration("executionTime", time.Since(beforeMainExecution)))
}
