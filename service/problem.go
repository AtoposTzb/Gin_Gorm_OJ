package service

import (
	"Gin_Gorm_OJ/define"
	"Gin_Gorm_OJ/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// @Summary 获取问题列表
// @Description 获取问题列表
// @Param page query int false "page" "当前页码"
// @Param size query int false "size" "每页数量"
// @Param keyword query string false "keyword" "搜索关键词"
// @Accept json
// @Produce json
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /problem [get]
func GetProblemList(c *gin.Context) {
	// 从请求参数中获取分页参数
	// DefaultQuery获取查询参数page,默认值为define.DefaultPage,参数解释：page为当前页码,默认值为define.DefaultPage
	//size为每页数量,page为当前页码
	size, _ := strconv.Atoi(c.DefaultQuery("size", define.DefaultSize))
	page, err := strconv.Atoi(c.DefaultQuery("page", define.DefaultPage))
	if err != nil {
		log.Println("转换分页参数失败", err)
	}
	//计算页面偏移量
	offset := (page - 1) * size
	//查询数据库
	var count int64

	keyword := c.Query("keyword")

	data := make([]*models.Problem, 0)
	tx := models.GetProblemList(keyword)
	err = tx.Count(&count).Offset(offset).Limit(size).Find(&data).Error
	if err != nil {
		log.Println("查询问题列表失败", err)
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"count": count,
			"data":  data,
		},
	})

}

/*
测试swagger 接口
// GetProblemList 获取问题列表
// @Summary 获取问题列表
// @Description 获取问题列表
// @Tags 问题
// @Accept json
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /problem [get]
func GetProblemList(c *gin.Context) {
	c.String(http.StatusOK, "success")
	c.JSON(http.StatusOK, gin.H{
		"data": models.GetProblemList(),
	})
}
*/
