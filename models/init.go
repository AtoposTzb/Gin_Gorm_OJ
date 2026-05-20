package models

import (
	"fmt"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB = Init()       // 数据库连接，全局变量，用于在其他文件中使用
var RDB = InitRedis() // Redis连接，全局变量，用于在其他文件中使用

// 初始化数据库模型
func Init() *gorm.DB {
	dsn := "root:qaz..//@tcp(127.0.0.1:3306)/gin_gorm_oj?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("连接数据库失败", err)
		//t.Fatal("连接数据库失败", err)
	}
	return db

}

func InitRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
}
