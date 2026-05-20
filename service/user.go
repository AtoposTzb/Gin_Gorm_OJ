package service

import (
	"Gin_Gorm_OJ/define"
	"Gin_Gorm_OJ/helper"
	"Gin_Gorm_OJ/models"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
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
	tokenString, err := helper.GenerateToken(data.Identity, data.Name, data.IsAdmin)
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
	//code := "123456" //这里写死验证码，实际应用中应该随机生成
	code := helper.CreateCode() //随机生成验证码
	err := models.RDB.Set(c, email, code, time.Second*300).Err()
	if err != nil {
		log.Println(err) // 写入redis失败，记录日志，用户看不到
		return
	}

	err = helper.SendEmail(email, code)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "验证码发送失败:" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"code": "验证码发送成功:",
		},
	})
}

// 用户注册
// @Tags 公共方法
// @Summary 用户注册
// @Description 用户注册
// @Param email formData string true "邮箱"
// @Param code formData string true "验证码"
// @Param username formData string true "用户名"
// @Param password formData string true "密码"
// @Param phone formData string false "手机号"
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /register [post]
func Register(c *gin.Context) {
	email := c.PostForm("email")
	userCode := c.PostForm("code")
	username := c.PostForm("username")
	password := c.PostForm("password")
	phone := c.PostForm("phone")
	if email == "" || userCode == "" || username == "" || password == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "邮箱、验证码、用户名、密码不能为空",
		})
		return
	}
	//验证码是否正确
	sysCode, err := models.RDB.Get(c, email).Result()
	if err == redis.Nil {
		log.Println(err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "验证码不正确，请重新发送验证码",
		})
		return
	}
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "验证码获取失败",
		})
		return
	}
	if sysCode != userCode {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "验证码错误",
		})
		return
	}
	//写入数据之前需要判断email是否存在
	var cnt int64
	err = models.DB.Where("mail = ?", email).Model(&models.UserBasic{}).Count(&cnt).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "查询邮箱是否存在失败:" + err.Error(),
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "邮箱已注册",
		})
		return
	}

	//数据写入数据库
	userIdenetity := helper.GenerateUUID()
	data := &models.UserBasic{
		Identity: userIdenetity,
		Name:     username,
		Password: helper.GeTMd5(password),
		Mail:     email,
		Phone:    phone,
	}
	err = models.DB.Create(data).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "注册失败:" + err.Error(),
		})
		return
	}
	//生成token
	tokenString, err := helper.GenerateToken(userIdenetity, username, data.IsAdmin)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "生成token失败:" + err.Error(),
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

// GetRankList 获取排名列表
// @Tags 公共方法
// @Summary 获取排名列表
// @Description 获取排名列表
// @Param page query int false "page" "当前页码"
// @Param size query int false "size" "每页数量"
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /rank-list [get]
func GetRankList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("转换分页参数失败", err)
	}
	//计算页面偏移量
	offset := (page - 1) * size
	//查询数据库
	var count int64
	list := make([]models.UserBasic, 0)
	err = models.DB.Model(&models.UserBasic{}).Count(&count).
		Order("complete_count desc,submit_count asc"). //按完成题目数量降序，提交记录数量升序
		Offset(offset).Limit(size).Find(&list).Error
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "查询排名失败" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"count": count,
			"data":  list,
		},
	})

}
