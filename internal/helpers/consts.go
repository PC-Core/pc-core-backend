package helpers

import "time"

const (
	RefreshCookieName = "refresh"
)

const (
	JWTAccessLifeTime  = 15 * time.Minute
	JWTRefreshLifeTime = 24 * 30 * time.Hour
)
