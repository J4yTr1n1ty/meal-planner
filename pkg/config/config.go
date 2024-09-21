package config

import "time"

const (
	EnvPort         = "PORT"
	EnvSqliteDBFile = "SQLITE_DB_FILE"
	EnvRedisAddr    = "REDIS_ADDR"

	SessionKey            = "sessionId"
	SessionExpire         = 5 * time.Hour
	SessionCookieName     = "session_id"
	SessionCookieMaxAge   = 5 * 60 * 60
	SessionRedisKeyFormat = "session/%s"

	Password = "4912Essen"
)
