package scheduler

import (
	"database/sql"
	"log"
	"time"
)

func RunBackground(db *sql.DB, vpnGateUrl string, dialTcpTimeoutSeconds int) {
	log.Println("Starting scheduler execution...")
	beforeMainExecution := time.Now()
	csvContent := getCsvContent(vpnGateUrl)
	vpnServers := createStructsFromCsv(csvContent)
	checkUnreachableServersOnDB(db, dialTcpTimeoutSeconds)
	insertServers(db, vpnServers, dialTcpTimeoutSeconds)
	log.Println("Ending scheduler execution, took", time.Now().Sub(beforeMainExecution))
}