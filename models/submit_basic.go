/*
提交表的模型
名称	类型	长度	小数点	不是 Null	虚拟	键	虚拟类型	表达式	枚举值	默认值	注释	存储	列格式	字符集	排序规则	键长度	键排序	永远生成	根据当前时间戳更新	二进制	自动递增	无符号	填充零
identity	varchar	36		false	false	false				NULL				utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
problem_identity	varchar	36		false	false	false				NULL	问题的唯一标识			utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
user_identity	varchar	36		false	false	false				NULL	用户的唯一标识			utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
path	varchar	255		false	false	false				NULL	代码路径			utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
status	tinyint	1		false	false	false				NULL	状态-[-1--待判断，1-答案正确，2-答案错误，3-运行超时，4-运行超内存]							false	false	false	false	false	false
*/
package models

import "gorm.io/gorm"

type SubmitBasic struct {
	gorm.Model                    //一样的包含ID,CreatedAt,UpdatedAt,DeletedAt字段
	Identity        string        `gorm:"column:identity;type:varchar(36);" json:"identity"`                 //提交表的唯一标识符
	ProblemIdentity string        `gorm:"column:problem_identity;type:varchar(36);" json:"problem_identity"` //问题的唯一标识
	ProblemBasic    *ProblemBasic `gorm:"foreignKey:identity;references:problem_identity"`                   //关联问题表
	UserIdentity    string        `gorm:"column:user_identity;type:varchar(36);" json:"user_identity"`       //用户的唯一标识
	UserBasic       *UserBasic    `gorm:"foreignKey:identity;references:user_identity"`                      //关联用户表
	Path            string        `gorm:"column:path;type:varchar(255);" json:"path"`                        //代码路径
	Status          int           `gorm:"column:status;type:tinyint(1);" json:"status"`                      //状态-[0-待判断，1-答案正确，2-答案错误，3-运行超时，4-运行超内存]
}

func (*SubmitBasic) TableName() string {
	return "submit_basic"
}

//查询
func GetSubmitList(problemIdentity, userIdentity string, status int) *gorm.DB {
	tx := DB.Model(&SubmitBasic{}).Preload("ProblemBasic", func(db *gorm.DB) *gorm.DB {
		return db.Omit("content") //排除问题表的content字段
	}).Preload("UserBasic")
	if problemIdentity != "" {
		tx = tx.Where("problem_identity = ?", problemIdentity)
	}
	if userIdentity != "" {
		tx = tx.Where("user_identity = ?", userIdentity)
	}
	if status != 0 {
		tx = tx.Where("status = ?", status)
	}
	return tx
}
