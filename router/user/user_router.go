package user

import (
	"Gin_Gorm_OJ/middlewares"
	"Gin_Gorm_OJ/service"

	"github.com/gin-gonic/gin"
)

func InitUserRouter(r *gin.Engine) {
	//用户路由
	userGroup := r.Group("/user")

	//用户路由
	userGroup.GET("/detail", service.GetUserDetail)
	//登录
	userGroup.POST("/login", service.Login)
	//发送验证码
	userGroup.POST("/send-code", service.SendCode)
	//注册
	userGroup.POST("/register", service.Register)
	//排名列表
	userGroup.GET("/rank-list", service.GetRankList)
	//提交记录
	userGroup.GET("/submit-list", service.GetSubmitList)
	//提交问题
	userGroup.POST("/submit-problem", middlewares.AuthUserCheck, service.SubmitProblem)
}
