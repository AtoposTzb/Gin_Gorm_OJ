/*
	问题表的模型
	名称	类型	长度	小数点	不是 Null	虚拟	键	虚拟类型	表达式	枚举值	默认值	注释	存储	列格式	字符集	排序规则	键长度	键排序	永远生成	根据当前时间戳更新	二进制	自动递增	无符号	填充零

identity	varchar	36		false	false	false				NULL				utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
title	varchar	255		false	false	false				NULL	问题的题目			utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
content	text			false	false	false				NULL	问题的正文描述			utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
max_mem	int			false	false	false				NULL	最大的运行内存							false	false	false	false	false	false
max_runtime	int			false	false	false				NULL	最大的运行时间							false	false	false	false	false	false
*/
package models

import (
	"gorm.io/gorm"
)

type ProblemBasic struct {
	gorm.Model        //问题表的模型,继承gorm.Model,有ID,CreatedAt,UpdatedAt,DeletedAt字段
	Identity   string `gorm:"column:identity;type:varchar(36);" json:"identity"`   //问题表的唯一标识符
	Title      string `gorm:"column:title;type:varchar(255);" json:"title"`        //问题表的标题
	Content    string `gorm:"column:content;type:text;" json:"content"`            //题目描述
	MaxMem     int    `gorm:"column:max_mem;type:int(11);" json:"max_mem"`         //最大的运行内存
	MaxRuntime int    `gorm:"column:max_runtime;type:int(11);" json:"max_runtime"` //最大的运行时间
	// many to many 关联
	Categories []*CategoryBasic `gorm:"many2many:problem_category;joinForeignKey:problem_id;joinReferences:category_id" json:"categories"`
}

func (*ProblemBasic) TableName() string {
	return "problem_basic"
}

// 获取问题列表
func GetProblemList(keyword, categoryIdentity string) *gorm.DB {
	tx := DB.Model(new(ProblemBasic)).Preload("Categories")

	if categoryIdentity != "" {
		// 有分类过滤时，只预加载匹配的分类，并过滤问题
		tx.Preload("Categories", "identity = ?", categoryIdentity).
			Joins("INNER JOIN problem_category pc ON pc.problem_id = problem_basic.id").
			Joins("INNER JOIN category_basic cb ON cb.id = pc.category_id").
			Where("cb.identity = ?", categoryIdentity)
	}

	return tx.Where("title like ? or content like ?", "%"+keyword+"%", "%"+keyword+"%")
}
