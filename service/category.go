package service

import (
	"Gin_Gorm_OJ/define"
	"Gin_Gorm_OJ/helper"
	"Gin_Gorm_OJ/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

// GetCategoryList
// @Tags 管理员私有方法
// @Summary 获取分类列表
// @Description 获取分类列表
// @Param authorization header string true "authorization" "管理员token"
// @Param page query int false "page" "当前页码"
// @Param size query int false "size" "每页数量"
// @Param keyword query string false "keyword" "搜索关键词"
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /category-list [get]
func GetCategoryList(c *gin.Context) {
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
	categoryList := make([]*models.CategoryBasic, 0)
	err = models.DB.Model(&models.CategoryBasic{}).Where("name LIKE ?", "%"+keyword+"%").
		Count(&count).Offset(offset).Limit(size).
		Find(&categoryList).Error
	//.Offset(offset).Limit(size)的作用是分页查询,offset是偏移量,size是每页数量,从数据库中查询count条数据
	//返回结果
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "查询分类列表失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": gin.H{
			"count": count,
			"data": map[string]interface{}{
				"count": count,
				"list":  categoryList,
			},
		},
	})
}

// CreateCategory
// @Tags 管理员私有方法
// @Summary 创建分类
// @Description 创建分类
// @Param authorization header string true "authorization"
// @Param name formData string true "name" "分类的名称"
// @Param problem_ids formData int false "problem_ids" "分类的标识列表"
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /category-create [post]
func CreateCategory(c *gin.Context) {
	name := c.PostForm("name")
	problemIds, _ := strconv.Atoi(c.PostForm("problem_ids"))
	categoryBasic := &models.CategoryBasic{
		Identity: helper.GenerateUUID(),
		Name:     name,
		ParentId: problemIds,
	}

	err := models.DB.Create(categoryBasic).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "创建分类失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建分类成功",
	})

}

// UpdateCategory
// @Tags 管理员私有方法
// @Summary 更新分类
// @Description 更新分类
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity" "分类的标识"
// @Param name formData string true "name" "分类的名称"
// @Param ParentId formData int false "ParentId" "分类的父级ID"
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /category-update [put]
func UpdateCategory(c *gin.Context) {
	identity := c.PostForm("identity")
	name := c.PostForm("name")
	parentId, _ := strconv.Atoi(c.PostForm("ParentId"))
	// if identity == "" || name == "" {
	// 	c.JSON(http.StatusOK, gin.H{
	// 		"code": -1,
	// 		"msg":  "分类的标识、名称不能为空",
	// 	})
	// 	return
	// }
	category := &models.CategoryBasic{
		Identity: identity,
		Name:     name,
		ParentId: parentId,
	}

	//err := models.DB.Save(categoryBasic).Error
	err := models.DB.Model(new(models.CategoryBasic)).Where("identity = ?", identity).Updates(category).Error
	//.Save()和.Updates()的区别是,.Save()会更新所有字段,.Updates()会更新指定的字段，在这里我们使用.Updates()更新指定的字段
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "更新分类失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新分类成功",
	})

}

// DeleteCategory
// @Tags 管理员私有方法
// @Summary 删除分类
// @Description 删除分类
// @Param authorization header string true "authorization"
// @Param identity query string true "identity" "分类的标识"
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /category-delete [delete]
func DeleteCategory(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "分类的标识不能为空",
		})
		return
	}
	var cnt int64
	err := models.DB.Model(&models.ProblemCategory{}).Where("category_id = (SELECT id FROM category_basic WHERE identity = ? LIMIT 1)", identity).Model(&models.ProblemCategory{}).Count(&cnt).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "获取分类关联问题失败",
		})
		return
	}
	if cnt > 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"data": "分类下有问题关联，不能删除",
		})
		return
	}
	err = models.DB.Where("identity = ?", identity).Delete(&models.CategoryBasic{}).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"code": -1,
			"msg":  "删除分类失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除分类成功",
	})

}
