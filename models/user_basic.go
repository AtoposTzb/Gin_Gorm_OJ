/*
	用户表的模型
	通过可视化MySql 数据库表结构创建好

名称	类型	长度	小数点	不是 Null	虚拟	键	虚拟类型	表达式	枚举值	默认值	注释	存储	列格式	字符集	排序规则	键长度	键排序	永远生成	根据当前时间戳更新	二进制	自动递增	无符号	填充零
name	varchar	255		false	false	false				NULL	名称			utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
password	varchar	32		false	false	false				NULL	密码			utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
phone	varchar	20		false	false	false				NULL	手机号			utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
mail	varchar	100		false	false	false				NULL	邮箱			utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
*/
package models

import "gorm.io/gorm"

type UserBasic struct {
	gorm.Model           //用户表的模型,继承gorm.Model,有ID,CreatedAt,UpdatedAt,DeletedAt字段 MySql记得设置字段
	Identity      string `gorm:"column:identity;type:varchar(36);" json:"identity"`     //用户表的唯一标识符
	Name          string `gorm:"column:name;type:varchar(255);" json:"name"`            //用户表的名称
	Password      string `gorm:"column:password;type:varchar(32);" json:"password"`     //用户表的密码
	Phone         string `gorm:"column:phone;type:varchar(20);" json:"phone"`           //用户表的手机号
	Mail          string `gorm:"column:mail;type:varchar(100);" json:"mail"`            //用户表的邮箱
	SubmitCount   int    `gorm:"column:submit_count;type:int;" json:"submit_count"`     //提交记录数量
	CompleteCount int    `gorm:"column:complete_count;type:int;" json:"complete_count"` //完成题目数量
	IsAdmin       int    `gorm:"column:is_admin;type:int;" json:"is_admin"`             //是否是管理员 0:否 1:是
}

func (*UserBasic) TableName() string {
	return "user_basic"
}
