package helpers

import "time"

const (
	RefreshCookieName = "refr-tk"
)

const (
	JWTAccessLifeTime  = 15 * time.Minute
	JWTRefreshLifeTime = 24 * 30 * time.Hour
)