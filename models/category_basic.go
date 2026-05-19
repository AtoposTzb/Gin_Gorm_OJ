/*
	分类表的模型
	分类表的模型,继承gorm.Model,有ID,CreatedAt,UpdatedAt,DeletedAt字段
	名称	类型	长度	小数点	不是 Null	虚拟	键	虚拟类型	表达式	枚举值	默认值	注释	存储	列格式	字符集	排序规则	键长度	键排序	永远生成	根据当前时间戳更新	二进制	自动递增	无符号	填充零

name	varchar	100		false	false	false				NULL	分类名称			utf8mb4	utf8mb4_0900_ai_ci			false	false	false	false	false	false
parent_id	int	 11		false	false	false				NULL	父级ID							false	false	false	false	false	false
*/
package models

import "gorm.io/gorm"

type CategoryBasic struct {
	gorm.Model        //分类表的模型,继承gorm.Model,有ID,CreatedAt,UpdatedAt,DeletedAt字段
	Identity   string `gorm:"column:identity;type:varchar(36);" json:"identity"` //分类表的唯一标识符
	Name       string `gorm:"column:name;type:varchar(100);" json:"name"`        //分类表的名称
	ParentId   int    `gorm:"column:parent_id;type:int(11);" json:"parent_id"`   //分类表的父级ID
	//many to many 关联
	Categories []*CategoryBasic `gorm:"many2many:category_basic;" json:"categories"`
}

func (*CategoryBasic) TableName() string {
	return "category_basic"
}
