package service

import (
	"Gin_Gorm_OJ/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

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
