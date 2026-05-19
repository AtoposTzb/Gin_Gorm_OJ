package service

import (
	"Gin_Gorm_OJ/define"
	"Gin_Gorm_OJ/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetSubmitList 获取提交记录
// @Tags 公共方法
// @Summary 获取提交记录
// @Description 获取提交记录
// @Param page query int false "page" "当前页码"
// @Param size query int false "size" "每页数量"
// @Param problem_identity query string false "problem_identity" "问题的标识"
// @Param user_identity query string false "user_identity" "用户的标识"
// @Param status query int false "status" "状态"
// @Accept json
// @Produce json
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /submit-list [get]
func GetSubmitList(c *gin.Context) {
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("转换分页参数失败", err)
	}
	//计算页面偏移量
	offset := (page - 1) * size
	//查询数据库
	var count int64
	data := make([]models.SubmitBasic, 0)

	problemIdentity := c.Query("problem_identity")
	userIdentity := c.Query("user_identity")
	status, _ := strconv.Atoi(c.Query("Status"))
	tx := models.GetSubmitList(problemIdentity, userIdentity, status)
	err = tx.Count(&count).Offset(offset).Limit(size).Find(&data).Error
	if err != nil {
		log.Println("查询提交记录失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get submit list failed" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"count": count,
			"data":  data,
		},
	})

}
