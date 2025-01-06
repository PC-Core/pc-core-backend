package main

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/Core-Mouse/cm-backend/docs"
	"github.com/Core-Mouse/cm-backend/internal/auth/jwt"
	"github.com/Core-Mouse/cm-backend/internal/config"
	"github.com/Core-Mouse/cm-backend/internal/controllers"
	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/Core-Mouse/cm-backend/internal/helpers"
	"github.com/Core-Mouse/cm-backend/internal/middlewares"
	inredis "github.com/Core-Mouse/cm-backend/internal/redis"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
)

const (
	ENV_POSTGRES       = "POSTGRES_IBYTE_CONN"
	ENV_JWT_KEY        = "PCCORE_JWT_KEY"
	ENV_REDIS_PASSWORD = "PCCORE_REDIS_PASSWORD"
)

func configureSwagger(gin *gin.Engine, path string) {
	docs.SwaggerInfo.Title = "PC Core Backend"
	docs.SwaggerInfo.Host = path
	docs.SwaggerInfo.Version = "0.0.1"

	swagger := controllers.NewSwaggerController(gin)
	swagger.ApplyRoutes()
}

func setupCors(r *gin.Engine, cfg *config.Config) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowCors,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
}

func loadJWTAuth(path string) (*jwt.JWTAuth, error) {
	key, err := os.ReadFile(path)

	if err != nil {
		return nil, err
	}

	return jwt.NewJWTAuth(key), nil
}

func setupRedis(cfg *config.Config) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.RedisConn.Addr, cfg.RedisConn.Port),
		Password: os.Getenv(ENV_REDIS_PASSWORD),
	})
}

func main() {
	r := gin.Default()

	config, err := config.ParseConfig("../../cfg.yml")

	if err != nil {
		panic(err)
	}

	setupCors(r, config)

	db, err := database.NewDbController(config.DbDriver, "postgres://postgres:1234@localhost:5432/pccore?sslmode=disable")

	if err != nil {
		panic(err)
	}

	if gin.Mode() == gin.DebugMode {
		configureSwagger(r, config.Addr)
	}

	auth, err := loadJWTAuth("/home/okinai/pccore_secret.key")

	if err != nil {
		panic(err)
	}

	redis := inredis.NewRedisController(setupRedis(config))

	uc := controllers.NewUserController(r, db, redis, auth)
	lc := controllers.NewLaptopController(r, db, middlewares.JWTAuthorize(auth), helpers.JWTRoleCast)
	pc := controllers.NewProductController(r, db)
	ct := controllers.NewCategoryController(r, db)
	jc := controllers.NewJWTController(r, db, auth)
	cc := controllers.NewCartController(r, db, redis, helpers.JWTPublicUserCaster(auth), middlewares.JWTAuthorize(auth))

	uc.ApplyRoutes()
	lc.ApplyRoutes()
	pc.ApplyRoutes()
	ct.ApplyRoutes()
	jc.ApplyRoutes()
	cc.ApplyRoutes()

	r.Run(config.Addr + ":" + strconv.Itoa(config.Port))
}
