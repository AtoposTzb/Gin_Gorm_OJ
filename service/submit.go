package service

import (
	"Gin_Gorm_OJ/define"
	"Gin_Gorm_OJ/helper"
	"Gin_Gorm_OJ/models"
	"bytes"
	"errors"
	"io"
	"log"
	"net/http"
	"os/exec"
	"runtime"
	"strconv"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
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
// @Router /user/submit-list [get]
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

// SubmitProblem 提交问题
// @Tags 用户私有方法
// @Summary 提交问题
// @Description 提交问题
// @Param authorization header string true "authorization"
// @Param problem_identity query string true "问题的标识"
// @Param code body string true "代码"
// @Success 200 {string} json "{\"code\":200,\"data\":{\"count\":0,\"data\":[]}\""
// @Failure 500 {object} map[string]interface{}
// @Router /user/submit-problem [post]
func SubmitProblem(c *gin.Context) {
	problemIdentity := c.Query("problem_identity")
	//ioutil已经弃用，使用io.ReadAll 请求体
	code, err := io.ReadAll(c.Request.Body) //
	if err != nil {
		//log.Println("读取请求体失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "read request body failed" + err.Error(),
		})
		return
	}
	//代码保存
	path, err := helper.SaveCode(code)
	if err != nil {
		log.Println("保存代码失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "save code failed" + err.Error(),
		})
		return
	}
	// 保存到数据库
	userIdentity, _ := c.Get("userClaim")
	userClaims := userIdentity.(*helper.UserClaims) // 从上下文获取用户标识
	submits := &models.SubmitBasic{
		Identity:        helper.CreateCode(),
		ProblemIdentity: problemIdentity,
		Path:            path,
		UserIdentity:    userClaims.Identity, // 从token中获取用户标识
	}

	//代码判断
	pbs := new(models.ProblemBasic)
	err = models.DB.Where("identity = ?", problemIdentity).Preload("TestCases").First(pbs).Error //查询问题是否存在,并预加载测试用例
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "problem not found" + err.Error(),
		})
		return
	}
	//答案错误的channel(管道)
	WA := make(chan int, len(pbs.TestCases)) //len(pbs.TestCases) 是测试用例的数量
	//运行超内存的channel(管道)
	OM := make(chan int, len(pbs.TestCases))
	//编译错误的channel(管道)
	CE := make(chan int, len(pbs.TestCases))
	//答案正确的channel(管道)
	AC := make(chan int, len(pbs.TestCases))
	//通过的个数
	passCount := 0
	var lock sync.Mutex // 用于保护passCount的互斥锁
	var msg string      //提示信息，给前端展示的提示信息

	//遍历测试用例
	for _, testCase := range pbs.TestCases {
		testCaces := testCase
		go func() {
			//该协程用来执行测试-参考code/runer.go
			cmd := exec.Command("go", "run", path) //这是创建一个命令对象
			//
			var out, stderr bytes.Buffer //这是创建一个缓冲区，用于存储命令的输出
			cmd.Stdout = &out            //将命令的输出重定向到缓冲区
			cmd.Stderr = &stderr         //将命令的错误输出重定向到缓冲区

			//根据测试的输入按列进行运行，拿到输出结果和标准的输出结果进行对比
			stdin, err := cmd.StdinPipe() //这是创建一个管道，用于将命令的输入重定向到管道
			if err != nil {
				panic("创建管道失败:" + err.Error())
			}
			defer stdin.Close()
			io.WriteString(stdin, testCaces.Input) //这是将测试的输入写入管道 然后运行命令

			var startMem runtime.MemStats
			runtime.ReadMemStats(&startMem) // 记录运行前的内存统计信息
			if err := cmd.Run(); err != nil {
				log.Println(err, stderr.String())
				if err.Error() == "exit status 2" {
					msg = stderr.String() // 编译错误信息
					CE <- 1
					return
				}
			}
			var endMem runtime.MemStats
			runtime.ReadMemStats(&endMem) // 记录运行后的内存统计信息
			//答案错误
			if out.String() != testCaces.Output {
				msg = "答案错误"
				WA <- 1
				return
			}
			//运行超内存
			if endMem.Alloc/1024-(startMem.Alloc/1024) > uint64(pbs.MaxMem) {
				msg = "运行超内存"
				OM <- 1
				return
			}
			//没有其他问题就是AC
			lock.Lock()
			passCount++
			msg = "答案正确"
			lock.Unlock()
			AC <- 1

		}()
	}

	//select等待所有测试用例执行完成
	//状态-[0-待判断，1-答案正确，2-答案错误，3-运行超时，4-运行超内存,5-编译错误]
	select {
	case <-WA:
		submits.Status = 2
	case <-OM:
		submits.Status = 4
	case <-CE:
		submits.Status = 5
	case <-AC:
		submits.Status = 1
	case <-time.After(time.Millisecond * time.Duration(pbs.MaxRuntime)):
		if passCount == len(pbs.TestCases) { // 所有测试用例都通过
			submits.Status = 1
		} else { // 有测试用例未通过
			submits.Status = 3
		}
	}

	//事务
	if err = models.DB.Transaction(func(tx *gorm.DB) error {
		err = tx.Create(submits).Error
		if err != nil {
			return errors.New("submit failed: " + err.Error())
		}
		//针对用户提交的代码更新 user_basic和problem_basic
		m := make(map[string]interface{})
		m["submit_count"] = gorm.Expr("submit_count + ?", 1)
		if submits.Status == 1 {
			m["complete_count"] = gorm.Expr("complete_count + ?", 1)
		}
		//更新user_basic
		err = tx.Model(new(models.UserBasic)).Where("identity = ?", userClaims.Identity).Updates(m).Error
		if err != nil {
			return errors.New("更新user_basic Error" + err.Error())
		}
		//更新problem_basic
		err = tx.Model(new(models.ProblemBasic)).Where("identity = ?", problemIdentity).Updates(m).Error
		if err != nil {
			return errors.New("更新problem_basic Error" + err.Error())
		}

		return nil
	}); err != nil {
		log.Println("保存提交记录失败", err)
		c.JSON(http.StatusOK, gin.H{
			"code": -1,
			"msg":  "save submit failed" + err.Error(),
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"data": map[string]interface{}{
			"identity": submits.Identity,
			"status":   submits.Status,
			"msg":      msg,
		},
	})

}
