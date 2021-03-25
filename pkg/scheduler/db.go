package scheduler

import (
	"database/sql"
	"go.uber.org/zap"
	"strings"
	"time"
)

func InitDb(dbDriver, dbUrl string, dbMaxOpenConn, dbMaxIdleConn, dbConnMaxLifetimeMin int) *sql.DB {
	db, err := sql.Open(dbDriver, dbUrl)
	if err != nil {
		logger.Fatal("fatal error occured while opening database connection", zap.String("error", err.Error()))
	}
	tuneDbPooling(db, dbMaxOpenConn, dbMaxIdleConn, dbConnMaxLifetimeMin)
	return db
}

// Read on https://www.alexedwards.net/blog/configuring-sqldb for detailed explanation
func tuneDbPooling(db *sql.DB, dbMaxOpenConn int, dbMaxIdleConn int, dbConnMaxLifetimeMin int) {
	// Set the maximum number of concurrently open connections (in-use + idle)
	// to 5. Setting this to less than or equal to 0 will mean there is no
	// maximum limit (which is also the default setting).
	db.SetMaxOpenConns(dbMaxOpenConn)
	// Set the maximum number of concurrently idle connections to 5. Setting this
	// to less than or equal to 0 will mean that no idle connections are retained.
	db.SetMaxIdleConns(dbMaxIdleConn)
	// Set the maximum lifetime of a connection to 1 hour. Setting it to 0
	// means that there is no maximum lifetime and the connection is reused
	// forever (which is the default behavior).
	db.SetConnMaxLifetime(time.Duration(int32(dbConnMaxLifetimeMin)) * time.Minute)
}

func checkUnreachableServersOnDB(db *sql.DB, dialTcpTimeoutSeconds int) {
	logger.Info("starting remove unreachable server operation on database")
	var (
		removedServerCount = 0
		port int
		ip, confData, proto string
		beforeExecution = time.Now()
	)
	query := "SELECT ip, proto, conf_data, port FROM servers"
	rows, err := db.Query(query)
	if err != nil {
		logger.Fatal("fatal error occured while querying database", zap.String("query", query),
			zap.String("error", err.Error()))
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		err := rows.Scan(&ip, &proto, &confData, &port)
		if err != nil {
			logger.Fatal("fatal error occured while scanning database", zap.String("ip", ip),
				zap.String("proto", proto), zap.Int("port", port), zap.String("error", err.Error()))
		}

		if !isServerInsertable(ip, proto, confData, port, dialTcpTimeoutSeconds) {
			removedServerCount++
			removeServers(db, ip, proto, confData, port)
		}
	}

	logger.Info("Ending remove unreachable server operation on database", zap.Int("removedServerCount", removedServerCount),
		zap.Duration("executionTime", time.Since(beforeExecution)))
}

func insertServers(db *sql.DB, vpnServers []vpnServer, dialTcpTimeoutSeconds int) {
	logger.Info("Starting insert reachable server operation on database")
	var (
		insertedServerCount = 0
		skippedServerCount = 0
		beforeExecution = time.Now()
		values []interface{}
	)
	var sqlStr = "REPLACE INTO servers(id, uuid, hostname, ip, port, conf_data, proto, enabled, score, ping, speed, " +
		"country_long, country_short, num_vpn_sessions, uptime, total_users, total_traffic, created_at) VALUES "
	for index, server := range vpnServers {
		if !isServerInsertable(server.ip, server.proto, server.confData, server.port, dialTcpTimeoutSeconds) {
			skippedServerCount++
			continue
		}
		insertedServerCount++
		sqlStr += "(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
		values = append(values, index + 1, server.uuid, server.hostname, server.ip, server.port, server.confData,
			server.proto, server.enabled, server.score, server.ping, server.speed, server.countryLong,
			server.countryLong, server.numVpnSessions, server.uptime, server.totalUsers, server.totalTraffic,
			server.createdAt)
	}
	sqlStr = strings.TrimSuffix(sqlStr, ",")
	stmt, _ := db.Prepare(sqlStr)
	_, err := stmt.Exec(values...)
	if err != nil {
		logger.Fatal("fatal error occured while executing query on database", zap.String("query", sqlStr),
			zap.String("error", err.Error()))
	}

	logger.Info("Ending insert reachable server operation on database", zap.Int("insertedServerCount", insertedServerCount),
		zap.Int("skippedServerCount", skippedServerCount), zap.Duration("executionTime", time.Since(beforeExecution)))
}

func removeServers(db *sql.DB, ip string, proto string, confData string, port int) {
	query := "DELETE FROM servers WHERE ip=? AND conf_data=? AND proto=? AND port=?"
	del, err := db.Prepare(query)
	if err != nil {
		panic(err)
	}

	_, err = del.Exec(ip, confData, proto, port)
	if err != nil {
		logger.Fatal("fatal error occured while executing query on database", zap.String("ip", ip),
			zap.String("proto", proto), zap.Int("port", port), zap.String("error", err.Error()))
	}
}