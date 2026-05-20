package middlewares

import (
	"Gin_Gorm_OJ/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthAdminCheck 中间件,检查用户是否是管理员
func AuthAdminCheck(c *gin.Context) {
	// 从请求头中获取token
	token := c.GetHeader("Authorization")
	userClaim, err := helper.AnalyseToken(token) //解析token
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"data": "Unauthorized Authorization",
		})
		c.Abort() //中断请求，防止继续执行后续的中间件
		return
	}
	if userClaim.IsAdmin != 1 || userClaim == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"data": "未授权，非管理员权限",
		})
		c.Abort() //中断请求，防止继续执行后续的中间件
		return
	}
	//	否则,继续执行后续的中间件
	c.Next()
}
