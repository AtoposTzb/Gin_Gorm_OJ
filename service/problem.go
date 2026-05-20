package service

import (
	"Gin_Gorm_OJ/define"
	"Gin_Gorm_OJ/helper"
	"Gin_Gorm_OJ/models"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// GetProblemList
// @Tags 公共方法
// @Summary 获取问题列表
// @Description 获取问题列表
// @Param page query int false "page" "当前页码"
// @Param size query int false "size" "每页数量"
// @Param keyword query string false "keyword" "搜索关键词"
// @Param category_identity query string false "category_identity" "分类的标识"
// @Accept json
// @Produce json
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /problem-list [get]
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
	categoryIdentity := c.Query("category_identity")

	data := make([]*models.ProblemBasic, 0)
	tx := models.GetProblemList(keyword, categoryIdentity)
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

// GetProblemDetail
// @Tags 公共方法
// @Summary 获取问题详情
// @Description 获取问题详情
// @Param identity query string false "problem identity" "问题的标识"
// @Accept json
// @Produce json
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /problem-detail [get]
func GetProblemDetail(c *gin.Context) {
	identity := c.Query("identity")
	if identity == "" {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "identity is required",
		})
		return
	}
	problemBasic := new(models.ProblemBasic)
	err := models.DB.Where("identity = ?", identity).Preload("Categories").
		First(&problemBasic).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusOK, gin.H{
				"code": -1,
				"msg":  "问题不存在或已被删除",
			})
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "Get problem detail failed" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": problemBasic,
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

// CreateProblem
// @Tags 管理员私有方法
// @Summary 创建问题
// @Description 创建问题
// @Param authorization header string true "authorization"
// @Param title formData string true "title" "问题的标题"
// @Param content formData string true "content" "问题的描述"
// @Param max_mem formData int false "max_mem" "最大的运行内存"
// @Param max_runtime formData int false "max_runtime" "最大的运行时间"
// @Param category_ids formData []string false "category_ids" collectionForm("multi")
// @Param test_cases formData []string true "test_cases" collectionForm("multi")
// @Accept multipart/form-data
// @Produce json
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /admin/problem-create [post]
func CreateProblem(c *gin.Context) {
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	if title == "" || content == "" || maxMem == 0 || maxRuntime == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "参数不能为空",
		})
		return
	}
	identity := helper.GenerateUUID()
	data := models.ProblemBasic{
		Identity:   identity,
		Title:      title,
		Content:    content,
		MaxMem:     maxMem,
		MaxRuntime: maxRuntime,
	}
	//处理分类
	CateGoryBasics := make([]*models.ProblemCategory, 0)
	for _, id := range categoryIds {
		categoryID, _ := strconv.Atoi(id)
		CateGoryBasics = append(CateGoryBasics, &models.ProblemCategory{
			ProblemID:  int(data.ID),
			CategoryID: categoryID,
		})
	}
	data.ProblemCategories = CateGoryBasics
	//处理测试用例
	testCaseBasics := make([]*models.TestCase, 0)
	for _, tc := range testCases {
		//例子{"input":"1,2\n","output":"3\n"}
		caseMap := map[string]string{}
		err := json.Unmarshal([]byte(tc), &caseMap) //这里需要转换为byte类型 ，然后解码为map[string]string
		if err != nil {
			log.Println("解析测试用例失败", err)
			continue
		}
		if _, ok := caseMap["input"]; !ok {
			log.Println("测试用例格式错误", tc)
			continue
		}
		if _, ok := caseMap["output"]; !ok {
			log.Println("测试用例格式错误", tc)
			continue
		}
		testCaseBasics = append(testCaseBasics, &models.TestCase{
			Identity:        helper.GenerateUUID(),
			ProblemIdentity: identity,
			Input:           caseMap["input"],
			Output:          caseMap["output"],
		})
	}
	data.TestCases = testCaseBasics //往问题中添加测试用例

	//创建问题
	err := models.DB.Create(&data).Error
	if err != nil {
		log.Println("创建问题失败", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"identity": identity,
		},
	})
}

// UpdateProblem
// @Tags 管理员私有方法
// @Summary 更新问题
// @Description 更新问题
// @Param authorization header string true "authorization"
// @Param identity formData string true "identity" "问题的标识"
// @Param title formData string true "title" "问题的标题"
// @Param content formData string true "content" "问题的描述"
// @Param max_mem formData int false "max_mem" "最大的运行内存"
// @Param max_runtime formData int false "max_runtime" "最大的运行时间"
// @Param category_ids formData []string false "category_ids" collectionForm("multi")
// @Param test_cases formData []string true "test_cases" collectionForm("multi")
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /admin/problem-update [put]
func UpdateProblem(c *gin.Context) {
	identity := c.PostForm("identity")
	title := c.PostForm("title")
	content := c.PostForm("content")
	maxMem, _ := strconv.Atoi(c.PostForm("max_mem"))
	maxRuntime, _ := strconv.Atoi(c.PostForm("max_runtime"))
	categoryIds := c.PostFormArray("category_ids")
	testCases := c.PostFormArray("test_cases")
	log.Printf("categoryIds: %v, testCases: %v", categoryIds, testCases)
	if identity == "" || title == "" || content == "" || maxMem == 0 || maxRuntime == 0 || len(categoryIds) == 0 || len(testCases) == 0 {
		c.JSON(http.StatusOK, gin.H{
			"code": 400,
			"msg":  "参数不能为空",
		})
		return
	}
	if err := models.DB.Transaction(func(tx *gorm.DB) error { //事务操作-这样可以确保所有操作要么都成功，要么都失败
		//问题信息保存 problem_basic
		problemBasic := &models.ProblemBasic{
			Identity:   identity,
			Title:      title,
			Content:    content,
			MaxMem:     maxMem,
			MaxRuntime: maxRuntime,
		}
		err := tx.Where("identity=?", identity).Updates(problemBasic).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "更新问题失败" + err.Error(),
			})
			return err
		}
		//查询问题详情
		err = tx.Where("identity=?", identity).Find(&problemBasic).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "查询问题详情失败" + err.Error(),
			})
			return err
		}
		//关联问题分类的保存 problem_category
		//1.删除已存在的关联关系
		err = tx.Delete(&models.ProblemCategory{}, "problem_id=?", problemBasic.ID).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "删除问题分类失败" + err.Error(),
			})
			return err
		}
		//2.添加新的关联关系
		CateGoryBasics := make([]*models.ProblemCategory, 0)
		for _, id := range categoryIds {
			categoryID, _ := strconv.Atoi(id)
			CateGoryBasics = append(CateGoryBasics, &models.ProblemCategory{
				ProblemID:  int(problemBasic.ID),
				CategoryID: categoryID, //分类的标识,遍历categoryIds,将每个分类的标识转换为整数
			})
		}
		problemBasic.ProblemCategories = CateGoryBasics
		//3.写入数据库
		if len(CateGoryBasics) > 0 { //ORM 的 Create 方法不接受空切片，如果 categoryIds 解析后 CateGoryBasics 为空就会报错 "empty slice found"。
			err = tx.Create(&CateGoryBasics).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 400,
					"msg":  "添加问题分类失败" + err.Error(),
				})
				return err
			}
		}
		//关联问题测试用例的保存 test_case
		//1.删除已经存在的关联关系
		err = tx.Where("problem_identity=?", identity).Delete(new(models.TestCase)).Error
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code": 400,
				"msg":  "删除测试用例失败" + err.Error(),
			})
			return err
		}
		//2.增加新的关联关系
		tcs := make([]*models.TestCase, 0)

		// 处理 testCases：Swagger 可能发送单个 JSON 数组字符串或多个独立的 JSON 对象
		var testCaseMaps []map[string]string

		for _, tc := range testCases {
			// 优先尝试解析为 JSON 数组（处理 Swagger 发送的单个字符串）
			var arr []map[string]string
			if err := json.Unmarshal([]byte(tc), &arr); err == nil {
				testCaseMaps = append(testCaseMaps, arr...)
				continue
			}

			// 如果失败，尝试添加中括号后解析（Swagger 可能发送缺少外层 [] 的字符串）
			wrapped := "[" + tc + "]"
			if err := json.Unmarshal([]byte(wrapped), &arr); err == nil {
				testCaseMaps = append(testCaseMaps, arr...)
				continue
			}

			// 失败则尝试解析为单个 JSON 对象
			caseMap := map[string]string{}
			if err := json.Unmarshal([]byte(tc), &caseMap); err == nil {
				testCaseMaps = append(testCaseMaps, caseMap)
				continue
			}

			log.Printf("无法解析测试用例: %q", tc)
		}

		for _, caseMap := range testCaseMaps {
			if _, ok := caseMap["input"]; !ok {
				return errors.New("缺少测试输入参数")
			}
			if _, ok := caseMap["output"]; !ok {
				return errors.New("缺少测试输出参数")
			}
			tcs = append(tcs, &models.TestCase{
				Identity:        helper.GenerateUUID(),
				ProblemIdentity: identity,
				Input:           caseMap["input"],
				Output:          caseMap["output"],
			})
		}

		//3.写入数据库
		if len(tcs) > 0 {
			err = tx.Create(&tcs).Error
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"code": 400,
					"msg":  "添加测试用例失败" + err.Error(),
				})
				return err
			}
		}

		return nil

	}); err != nil {
		log.Println("更新问题失败", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新问题成功",
	})
}
