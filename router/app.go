package router

import (
	_ "Gin_Gorm_OJ/docs"
	"Gin_Gorm_OJ/router/admin"
	"Gin_Gorm_OJ/router/user"
	"Gin_Gorm_OJ/service"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// gin-swagger middleware
// swagger embed files

func InitRouter() *gin.Engine {
	r := gin.Default()
	//swagger路由
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	//公有方法
	//问题路由
	r.GET("/ping", service.Ping)
	r.GET("/problem-list", service.GetProblemList)
	r.GET("/problem-detail", service.GetProblemDetail)

	//用户路由
	user.InitUserRouter(r)

	//管理员路由
	admin.InitAdminRouter(r)

	return r
}
