package controllers

import (
	_ "github.com/PC-Core/pc-core-backend/docs"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type SwaggerController struct {
	engine *gin.Engine
}

func NewSwaggerController(engine *gin.Engine) *SwaggerController {
	return &SwaggerController{
		engine,
	}
}

func (c *SwaggerController) ApplyRoutes() {
	c.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
}
