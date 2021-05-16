package scheduler

const (
	sqlSelectServers  = "SELECT ip, proto, conf_data, port FROM servers"
	sqlReplaceServers = "REPLACE INTO servers(id, uuid, hostname, ip, port, conf_data, proto, enabled, score, ping, " +
		"speed, country_long, country_short, num_vpn_sessions, uptime, total_users, total_traffic, created_at) VALUES " +
		"(?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?),"
	sqlDeleteServer = "DELETE FROM servers WHERE ip=? AND conf_data=? AND proto=? AND port=?"
)
