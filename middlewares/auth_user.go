package middlewares

import (
	"Gin_Gorm_OJ/helper"
	"net/http"

	"github.com/gin-gonic/gin"
)

// AuthUserCheck 中间件,检查用户是否登录
func AuthUserCheck(c *gin.Context) {
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
	if userClaim == nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"code": 401,
			"data": "请先登录",
		})
		c.Abort() //中断请求，防止继续执行后续的中间件
		return
	}
	c.Set("userClaim", userClaim) //将用户信息存储到上下文,方便后续使用
	// 继续执行后续的中间件
	c.Next()
}
