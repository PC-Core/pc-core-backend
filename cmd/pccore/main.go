package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/PC-Core/pc-core-backend/docs"
	"github.com/PC-Core/pc-core-backend/internal/auth/jwt"
	"github.com/PC-Core/pc-core-backend/internal/config"
	"github.com/PC-Core/pc-core-backend/internal/controllers"
	"github.com/PC-Core/pc-core-backend/internal/database"
	"github.com/PC-Core/pc-core-backend/internal/helpers"
	"github.com/PC-Core/pc-core-backend/internal/middlewares"
	inredis "github.com/PC-Core/pc-core-backend/internal/redis"
	"github.com/PC-Core/pc-core-backend/internal/static"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/swaggo/swag"
)

const (
	ENV_POSTGRES       = "PCCORE_POSTGRES_CONN"
	ENV_JWT_KEY        = "PCCORE_JWT_KEY"
	ENV_REDIS_PASSWORD = "PCCORE_REDIS_PASSWORD"
	ENV_MINIO_ACCESS   = "MINIO_ACCESS"
	ENV_MINIO_SECRET   = "MINIO_SECRET"
)

const SWAGGER_KEY = "swagger"

func configureSwagger(gin *gin.Engine, path string) {
	docs.SwaggerInfo.Title = "PC Core Backend"
	docs.SwaggerInfo.Host = path
	docs.SwaggerInfo.Version = "0.0.1"

	swaggerRegisterCheck()

	swagger := controllers.NewSwaggerController(gin)
	swagger.ApplyRoutes()
}

func swaggerRegisterCheck() {
	if swag.GetSwagger(SWAGGER_KEY) == nil {
		swag.Register(SWAGGER_KEY, docs.SwaggerInfo)
	}
}

func setupCors(r *gin.Engine, cfg *config.Config) {
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowCors,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
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

func setupMinio(config *config.Config) *static.MinIOClient {
	client, err := static.NewMinIOClient(config.MinIOConn.Ep, os.Getenv(ENV_MINIO_ACCESS), os.Getenv(ENV_MINIO_SECRET), config.MinIOConn.Secure, config.MinIOConn.Bucket)

	if err != nil {
		panic(err)
	}

	return client
}

func setupWorkingDir() *string {
	dir := flag.String("working-dir", "./", "The directory containing config files.")

	flag.Parse()

	return dir
}

func main() {
	wd := setupWorkingDir()

	if wd == nil {
		panic("No working directory specified")
	}

	err := godotenv.Load(fmt.Sprintf("%s/.env", *wd))

	if err != nil {
		log.Fatal("Error loading .env file!")
	}

	r := gin.Default()

	config, err := config.ParseConfig(fmt.Sprintf("%s/cfg.yml", *wd))

	if err != nil {
		panic(err)
	}

	setupCors(r, config)

	db, err := database.NewDPostgresDbController(config.DbDriver, os.Getenv(ENV_POSTGRES))

	if err != nil {
		panic(err)
	}

	if gin.Mode() == gin.DebugMode {
		configureSwagger(r, config.Addr)
	}

	auth, err := loadJWTAuth(os.Getenv(ENV_JWT_KEY))

	if err != nil {
		panic(err)
	}

	staticDataController := setupMinio(config)

	redis := inredis.NewRedisController(setupRedis(config))

	uc := controllers.NewUserController(r, db, redis, auth)
	lc := controllers.NewLaptopController(r, db, middlewares.JWTAuthorize(auth), helpers.JWTRoleCast)
	pc := controllers.NewProductController(r, db)
	ct := controllers.NewCategoryController(r, db)
	jc := controllers.NewJWTController(r, db, auth)
	cc := controllers.NewCartController(r, db, redis, helpers.JWTPublicUserCaster(auth), middlewares.JWTAuthorize(auth))
	prc := controllers.NewProfileController(r, helpers.JWTPublicUserCaster(auth), middlewares.JWTAuthorize(auth))
	mc := controllers.NewStaticController(r, staticDataController)

	uc.ApplyRoutes()
	lc.ApplyRoutes()
	pc.ApplyRoutes()
	ct.ApplyRoutes()
	jc.ApplyRoutes()
	cc.ApplyRoutes()
	prc.ApplyRoutes()
	mc.ApplyRoutes()

	r.Run(config.Addr + ":" + strconv.Itoa(config.Port))
}
