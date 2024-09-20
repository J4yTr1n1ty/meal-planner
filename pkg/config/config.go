package config

import "time"

const (
	EnvPort         = "PORT"
	EnvSqliteDBFile = "SQLITE_DB_FILE"
	EnvRedisAddr    = "REDIS_ADDR"

	SessionKey    = "sessionId"
	SessionExpire = 5 * time.Hour

	SessionRedisKeyFormat = "session/%s"
)
