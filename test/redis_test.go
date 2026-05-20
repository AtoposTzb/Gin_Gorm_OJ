package test

import (
	"Gin_Gorm_OJ/models"
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// 全局定义变量不可以使用:=，否则会报错：:= 只能在函数体中使用
var ctx = context.Background()
var rdb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

func TestSetRedis(t *testing.T) {
	// 测试设置redis
	err := rdb.Set(ctx, "key:name", "value:mmc", time.Second*10).Err()
	if err != nil {
		panic(err)
	}
}

func TestGetRedis(t *testing.T) {
	// 测试获取redis
	val, err := rdb.Get(ctx, "key:name").Result()
	if err == redis.Nil {
		t.Log("key:name 不存在")
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("key:name", val)
}

func TestGetRedisOfModel(t *testing.T) {
	// 测试获取redis
	val, err := models.RDB.Get(ctx, "key:name").Result()
	if err == redis.Nil {
		t.Log("key:name 不存在")
		return
	}
	if err != nil {
		panic(err)
	}
	fmt.Println("key:name", val)
}

/*
=== RUN   TestSetRedis
--- PASS: TestSetRedis (0.01s)
PASS
ok  	Gin_Gorm_OJ/test	0.294s

10秒内
=== RUN   TestGetRedis
key:name value:mmc
--- PASS: TestGetRedis (0.00s)
PASS
ok  	Gin_Gorm_OJ/test	0.277s

10秒后 key:name 不存在
=== RUN   TestGetRedis
    e:/GoStudy_itying/Gin_Gorm_OJ/test/redis_test.go:32: key:name 不存在
--- PASS: TestGetRedis (0.00s)
PASS
ok  	Gin_Gorm_OJ/test	0.294s
*/
