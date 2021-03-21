package scheduler

import "time"

type vpnServer struct {
	hostname, uuid, ip, proto, countryLong, countryShort, confData             string
	port, score, ping, speed, numVpnSessions, uptime, totalUsers, totalTraffic int
	enabled				                                                       bool
	createdAt                                                                  time.Time
}