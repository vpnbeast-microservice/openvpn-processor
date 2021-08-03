package scheduler

import (
	"database/sql"
	"fmt"
	commons "github.com/vpnbeast/golang-commons"
	"go.uber.org/zap"
	"openvpn-processor/pkg/metrics"
	"openvpn-processor/pkg/options"
	"strings"
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
	fmt.Print(opts.DbDriver, opts.DbUrl)
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

	rows, err := db.Query(sqlSelectServers)
	if err != nil {
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
		err := rows.Scan(&ip, &proto, &confData, &port)
		if err != nil {
			logger.Fatal("fatal error occurred while scanning database", zap.String("ip", ip),
				zap.String("proto", proto), zap.Int("port", port), zap.String("error", err.Error()))
		}

		if !isServerInsertable(ip, proto, confData, port, opts.DialTcpTimeoutSeconds) {
			removedServerCount++
			removeServers(db, ip, proto, confData, port)
		}
	}

	logger.Info("Ending remove unreachable server operation on database", zap.Int("removedServerCount", removedServerCount),
		zap.Duration("executionTime", time.Since(beforeExecution)))
}

func insertServers(db *sql.DB, vpnServers []vpnServer) {
	var (
		insertedServerCount = 0
		skippedServerCount  = 0
		beforeExecution     = time.Now()
		values              []interface{}
	)

	logger.Info("Starting insert reachable server operation on database")
	for index, server := range vpnServers {
		if !isServerInsertable(server.ip, server.proto, server.confData, server.port, opts.DialTcpTimeoutSeconds) {
			skippedServerCount++
			metrics.SkippedCounter.Inc()
			continue
		}
		insertedServerCount++
		metrics.InsertedCounter.Inc()
		values = append(values, index+1, server.uuid, server.hostname, server.ip, server.port, server.confData,
			server.proto, server.enabled, server.score, server.ping, server.speed, server.countryLong,
			server.countryLong, server.numVpnSessions, server.uptime, server.totalUsers, server.totalTraffic,
			server.createdAt)
	}
	sqlStr := strings.TrimSuffix(sqlReplaceServers, ",")
	stmt, _ := db.Prepare(sqlStr)
	_, err := stmt.Exec(values...)
	if err != nil {
		logger.Fatal("fatal error occurred while executing query on database", zap.String("query", sqlStr),
			zap.String("error", err.Error()))
	}

	logger.Info("Ending insert reachable server operation on database", zap.Int("insertedServerCount", insertedServerCount),
		zap.Int("skippedServerCount", skippedServerCount), zap.Duration("executionTime", time.Since(beforeExecution)))
}

func removeServers(db *sql.DB, ip string, proto string, confData string, port int) {
	del, err := db.Prepare(sqlDeleteServer)
	if err != nil {
		// TODO: do not panic, handle properly
		panic(err)
	}

	_, err = del.Exec(ip, confData, proto, port)
	if err != nil {
		logger.Fatal("fatal error occurred while executing query on database", zap.String("ip", ip),
			zap.String("proto", proto), zap.Int("port", port), zap.String("error", err.Error()))
	}
}
