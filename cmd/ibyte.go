package main

import (
	"os"
	"strconv"

	"github.com/Core-Mouse/cm-backend/internal/controllers"
	"github.com/Core-Mouse/cm-backend/internal/database"
	"github.com/gin-gonic/gin"
)

func main() {
	gin := gin.Default()

	db, err := database.NewDbController("postgres", os.Getenv("POSTGRES_IBYTE_CONN"))

	if err != nil {
		panic("Connection to the database failed")
	}

	uc := controllers.NewUserController(gin, db)

	uc.ApplyRoutes()

	gin.Run("127.0.0.1:" + strconv.Itoa(8080))
}
