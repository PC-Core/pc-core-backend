package main

import (
	"os"
	"strconv"

	"github.com/Core-Mouse/cm-backend/internal/config"
	"github.com/Core-Mouse/cm-backend/internal/controllers"
	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

func main() {
	gin := gin.Default()

	config, cerr := config.ParseConfig("../cfg.yml")

	if cerr != nil {
		panic(cerr)
	}

	db, err := database.NewDbController(config.DbDriver, os.Getenv("POSTGRES_IBYTE_CONN"))

	if err != nil {
		panic(err)
	}

	uc := controllers.NewUserController(gin, db)
	lc := controllers.NewLaptopController(gin, db)
	pc := controllers.NewProductController(gin, db)

	uc.ApplyRoutes()
	lc.ApplyRoutes()
	pc.ApplyRoutes()

	gin.Run(config.Addr + ":" + strconv.Itoa(config.Port))
}
