package router

import (
	_ "Gin_Gorm_OJ/docs"
	"Gin_Gorm_OJ/middlewares"
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
	r.GET("/user-detail", service.GetUserDetail)
	r.POST("/login", service.Login)
	r.POST("/send-code", service.SendCode)
	r.POST("/register", service.Register)
	r.GET("/rank-list", service.GetRankList)

	//提交记录
	r.GET("/submit-list", service.GetSubmitList)

	//私有方法
	//管理员私有方法:创建问题
	r.POST("/problem-create", middlewares.AuthAdminCheck, service.CreateProblem) //创建问题，中间件检查用户是否是管理员

	return r
}
