package models

import "gorm.io/gorm"

type ProblemCategory struct {
	gorm.Model
	ProblemID     int            `gorm:"column:problem_id;type:int;" json:"problem_id"`   //问题的ID
	CategoryID    int            `gorm:"column:category_id;type:int;" json:"category_id"` //分类的ID
	CategoryBasic *CategoryBasic `gorm:"foreignKey:ID;references:CategoryID"`             //分类的详细信息
}

func (*ProblemCategory) TableName() string {
	return "problem_category"
}
