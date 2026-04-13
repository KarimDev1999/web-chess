package appconst

import "time"

const (
	DefaultPort       = "8080"
	DefaultDBHost     = "localhost"
	DefaultDBPort     = 5432
	DefaultDBUser     = "chess"
	DefaultDBPassword = "chess"
	DefaultDBName     = "chess"
	DefaultRedisAddr  = "localhost:6379"

	HTTPShutdownTimeout = 10 * time.Second
)
