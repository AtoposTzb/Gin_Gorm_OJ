package test

import (
	"Gin_Gorm_OJ/models"
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func TestGormMySQL(t *testing.T) {
	// 测试数据库连接
	dsn := "root:qaz..//@tcp(127.0.0.1:3306)/gin_gorm_oj?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接数据库失败", err)
		//t.Fatal("连接数据库失败", err)
	}
	var data = make([]*models.ProblemBasic, 0)
	err = db.Find(&data).Error
	if err != nil {
		t.Fatal("查询数据库失败", err)
	}
	for _, v := range data {
		//t.Log(v)
		fmt.Println(v)
	}

	/*=== RUN   TestGormMySQL
	&{{1 2026-05-19 14:44:30 +0800 CST 2026-05-19 14:44:33 +0800 CST {0001-01-01 00:00:00 +0000 UTC false}} 1 1 测试 测试 1024 1000}
	--- PASS: TestGormMySQL (0.01s)
	PASS
	ok  	Gin_Gorm_OJ/test	2.017s
	*/
}
