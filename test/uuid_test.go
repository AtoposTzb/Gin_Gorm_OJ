package test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestUUID(t *testing.T) {
	// 测试生成uuid
	uuid := uuid.New()
	fmt.Println(uuid)
	// 打印uuid的字符串
	fmt.Println(len(uuid))
	fmt.Println(uuid.String())

}

/*
=== RUN   TestUUID
e2387059-e86b-4513-8eec-84540dd7e349
16
e2387059-e86b-4513-8eec-84540dd7e349
--- PASS: TestUUID (0.00s)
PASS
ok  	Gin_Gorm_OJ/test	2.083s
*/
