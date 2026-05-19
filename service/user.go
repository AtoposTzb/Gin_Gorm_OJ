package service

import (
	"Gin_Gorm_OJ/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetUserDetail 获取用户详情
// @Tags 公共方法
// @Summary 获取用户详情
// @Description 获取用户详情
// @Param identity query string false "user identity" "用户的标识"
// @Accept json
// @Produce json
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /user-detail [get]
func GetUserDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "用户唯一标识不能为空",
		})
	}
	//开始查询用户详情
	username := &models.UserBasic{}
	err := models.DB.Omit("password").
		Where("identity = ? ", identity).
		First(&username).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "用户不存在",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": username,
	})

}
