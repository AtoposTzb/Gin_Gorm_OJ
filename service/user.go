package service

import (
	"Gin_Gorm_OJ/helper"
	"Gin_Gorm_OJ/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
		return
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

// Login 登录
// @Tags 公共方法
// @Summary 登录
// @Description 登录
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /login [post]
func Login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")
	if username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "用户名或密码不能为空",
		})
		return
	}
	//md5 加密密码
	password = helper.GeTMd5(password)
	//fmt.Println(username, password)

	//
	data := &models.UserBasic{}
	err := models.DB.Where("name = ? AND password = ?",
		username, password).
		First(&data).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"data": "用户名或密码错误",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "Get UserBasic error" + err.Error(),
		})
		return
	}
	//生成token
	tokenString, err := helper.GenerateToken(data.Identity, data.Name)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "GenerateToken error" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"token": tokenString,
		},
	})
}

// SendCode 发送验证码
// @Tags 公共方法
// @Summary 发送验证码
// @Description 发送验证码
// @Param email formData string true "邮箱"
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /send-code [post]
func SendCode(c *gin.Context) {
	email := c.PostForm("email")
	if email == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "邮箱不能为空",
		})
		return
	}
	code := "123456" //这里写死验证码，实际应用中应该随机生成
	err := helper.SendEmail(email, code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "SendEmail error" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"code": "验证码发送成功:" + code,
		},
	})
}
