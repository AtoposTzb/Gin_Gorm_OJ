package router

import (
	_ "Gin_Gorm_OJ/docs"
	"Gin_Gorm_OJ/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// gin-swagger middleware
// swagger embed files

func InitRouter() *gin.Engine {
	r := gin.Default()
	//路由
	r.GET("/ping", service.Ping)
	r.GET("/problem", service.GetProblemList)
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	return r
}
