package main

import (
	"os"
	"strconv"
	"time"

	"github.com/Core-Mouse/cm-backend/docs"
	"github.com/Core-Mouse/cm-backend/internal/config"
	"github.com/Core-Mouse/cm-backend/internal/controllers"
	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func configureSwagger(gin *gin.Engine, path string) {
	docs.SwaggerInfo.Title = "PC Core Backend"
	docs.SwaggerInfo.Host = path
	docs.SwaggerInfo.Version = "0.0.1"

	swagger := controllers.NewSwaggerController(gin)
	swagger.ApplyRoutes()
}

func main() {
	r := gin.Default()

	config, err := config.ParseConfig("../../cfg.yml")

	if err != nil {
		panic(err)
	}

	r.Use(cors.New(cors.Config{
		AllowOrigins:     config.AllowCors,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	db, err := database.NewDbController(config.DbDriver, os.Getenv("POSTGRES_IBYTE_CONN"))

	if err != nil {
		panic(err)
	}

	if gin.Mode() == gin.DebugMode {
		configureSwagger(r, config.Addr)
	}

	uc := controllers.NewUserController(r, db)
	lc := controllers.NewLaptopController(r, db)
	pc := controllers.NewProductController(r, db)

	uc.ApplyRoutes()
	lc.ApplyRoutes()
	pc.ApplyRoutes()

	r.Run(config.Addr + ":" + strconv.Itoa(config.Port))
}
