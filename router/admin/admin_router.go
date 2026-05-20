package admin

import (
	"Gin_Gorm_OJ/middlewares"
	"Gin_Gorm_OJ/service"

	"github.com/gin-gonic/gin"
)

func InitAdminRouter(r *gin.Engine) {
	//管理员路由
	adminGroup := r.Group("/admin", middlewares.AuthAdminCheck)
	//私有方法
	//管理员私有方法:创建问题
	adminGroup.POST("/problem-create", service.CreateProblem) //创建问题，中间件检查用户是否是管理员
	//问题修改
	adminGroup.PUT("/problem-update", service.UpdateProblem)
	//分类列表
	adminGroup.GET("/category-list", service.GetCategoryList)
	//创建分类
	adminGroup.POST("/category-create", service.CreateCategory)
	//更新分类
	adminGroup.PUT("/category-update", service.UpdateCategory)
	//删除分类
	adminGroup.DELETE("/category-delete", service.DeleteCategory)

}
