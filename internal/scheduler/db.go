package scheduler

import (
	"database/sql"
	commons "github.com/vpnbeast/golang-commons"
	"go.uber.org/zap"
	"openvpn-processor/internal/metrics"
	"openvpn-processor/internal/options"
	"time"
)

var (
	logger *zap.Logger
	opts   *options.OpenvpnProcessorOptions
	db     *sql.DB
	err    error
)

func init() {
	logger = commons.GetLogger()
	opts = options.GetOpenvpnProcessorOptions()
	db = initDb()
}

// GetDatabase returns the initialized sql.DB
func GetDatabase() *sql.DB {
	return db
}

// initDb gets parameters and initiate database connection, returns connection then
func initDb() *sql.DB {
	db, err = sql.Open(opts.DbDriver, opts.DbUrl)
	if err != nil {
		logger.Fatal("fatal error occurred while opening database connection", zap.String("error", err.Error()))
	}
	tuneDbPooling(db)
	return db
}

// Read on https://www.alexedwards.net/blog/configuring-sqldb for detailed explanation
func tuneDbPooling(db *sql.DB) {
	// Set the maximum number of concurrently open connections (in-use + idle)
	// to 5. Setting this to less than or equal to 0 will mean there is no
	// maximum limit (which is also the default setting).
	db.SetMaxOpenConns(opts.DbMaxOpenConn)
	// Set the maximum number of concurrently idle connections to 5. Setting this
	// to less than or equal to 0 will mean that no idle connections are retained.
	db.SetMaxIdleConns(opts.DbMaxIdleConn)
	// Set the maximum lifetime of a connection to 1 hour. Setting it to 0
	// means that there is no maximum lifetime and the connection is reused
	// forever (which is the default behavior).
	db.SetConnMaxLifetime(time.Duration(int32(opts.DbConnMaxLifetimeMin)) * time.Minute)
}

func checkUnreachableServersOnDB(db *sql.DB) {
	logger.Info("starting remove unreachable server operation on database")
	var (
		removedServerCount  = 0
		port                int
		ip, confData, proto string
		beforeExecution     = time.Now()
	)

	var rows *sql.Rows
	if rows, err = db.Query(sqlSelectServers); err != nil {
		logger.Fatal("fatal error occurred while querying database", zap.String("query", sqlSelectServers),
			zap.String("error", err.Error()))
	}

	defer func() {
		err := rows.Close()
		if err != nil {
			panic(err)
		}
	}()

	for rows.Next() {
		if err := rows.Scan(&ip, &proto, &confData, &port); err != nil {
			logger.Fatal("fatal error occurred while scanning database", zap.String("ip", ip),
				zap.String("proto", proto), zap.Int("port", port), zap.String("error", err.Error()))
		}

		if !isServerInsertable(ip, proto, confData, port, opts.DialTcpTimeoutSeconds) {
			removedServerCount++
			removeServers(db, ip, proto, confData, port)
		}
	}

	logger.Info("Ending remove unreachable server operation on database", zap.Int("serverCount", removedServerCount),
		zap.Duration("executionTime", time.Since(beforeExecution)))
}

func insertServers(db *sql.DB, vpnServers []vpnServer) {
	var (
		insertedServerCount = 0
		skippedServerCount  = 0
		failedServerCount   = 0
		beforeExecution     = time.Now()
		values              []interface{}
		stmt                *sql.Stmt
		err                 error
	)

	logger.Info("Starting insert reachable server operation on database", zap.Int("serverCount", len(vpnServers)))
	for index, server := range vpnServers {
		if !isServerInsertable(server.ip, server.proto, server.confData, server.port, opts.DialTcpTimeoutSeconds) {
			skippedServerCount++
			metrics.SkippedServerCounter.Inc()
			continue
		}

		values = append(values, index+1, server.uuid, server.hostname, server.ip, server.port, server.confData,
			server.proto, server.enabled, server.score, server.ping, server.speed, server.countryLong,
			server.countryLong, server.numVpnSessions, server.uptime, server.totalUsers, server.totalTraffic,
			server.createdAt)

		if stmt, err = db.Prepare(sqlReplaceServer); err != nil {
			logger.Fatal("fatal error occured while preparing statement", zap.String("query", sqlReplaceServer))
			return
		}

		var res sql.Result
		if res, err = stmt.Exec(values...); err != nil {
			logger.Error("an error occurred while executing query on database", zap.String("server", server.hostname),
				zap.String("query", sqlReplaceServer), zap.String("error", err.Error()))
			failedServerCount++
			metrics.FailedServerCounter.Inc()
			continue
		}

		if row, _ := res.RowsAffected(); row == 1 {
			insertedServerCount++
			metrics.AvailableServerCounter.Inc()
		}

		// clear the slice after all
		values = nil
	}

	logger.Info("Ending insert reachable server operation on database", zap.Int("insertedServerCount", insertedServerCount),
		zap.Int("skippedServerCount", skippedServerCount), zap.Int("failedServerCount", failedServerCount),
		zap.Duration("executionTime", time.Since(beforeExecution)))
}

func removeServers(db *sql.DB, ip string, proto string, confData string, port int) {
	stmt, err := db.Prepare(sqlDeleteServer)
	if err != nil {
		// TODO: do not panic, handle properly
		panic(err)
	}

	if _, err = stmt.Exec(ip, confData, proto, port); err != nil {
		logger.Fatal("fatal error occurred while executing query on database", zap.String("ip", ip),
			zap.String("proto", proto), zap.Int("port", port), zap.String("error", err.Error()))
	}
}
