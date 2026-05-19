/*
	问题表的模型
	名称	类型	长度	小数点	不是 Null	虚拟	键	虚拟类型	表达式	枚举值	默认值	注释	存储	列格式	字符集	排序规则	键长度	键排序	永远生成	根据当前时间戳更新	二进制	自动递增	无符号	填充零

identity	varchar	36		false	false	false				NULL				utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
category_id	varchar	255		false	false	false				NULL	以逗号分隔的分类			utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
title	varchar	255		false	false	false				NULL	问题的题目			utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
content	text			false	false	false				NULL	问题的正文描述			utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
max_mem	int			false	false	false				NULL	最大的运行内存							false	false	false	false	false	false
max_runtime	int			false	false	false				NULL	最大的运行时间							false	false	false	false	false	false
*/
package models

import (
	"gorm.io/gorm"
)

type Problem struct {
	gorm.Model        //问题表的模型,继承gorm.Model,有ID,CreatedAt,UpdatedAt,DeletedAt字段
	Identity   string `gorm:"column:identity;type:varchar(36);" json:"identity"`        //问题表的唯一标识符
	CategoryId string `gorm:"column:category_id;type:varchar(255);" json:"category_id"` //问题表的分类标识符
	Title      string `gorm:"column:title;type:varchar(255);" json:"title"`             //问题表的标题
	Content    string `gorm:"column:content;type:text;" json:"content"`                 //题目描述
	MaxMem     int    `gorm:"column:max_mem;type:int(11);" json:"max_mem"`              //最大的运行内存
	MaxRuntime int    `gorm:"column:max_runtime;type:int(11);" json:"max_runtime"`      //最大的运行时间
}

func (*Problem) TableName() string {
	return "problem"
}

// 获取问题列表
func GetProblemList(keyword string) *gorm.DB {
	return DB.Model(new(Problem)).
		Where("title like ? or content like ?", "%"+keyword+"%", "%"+keyword+"%")
}
