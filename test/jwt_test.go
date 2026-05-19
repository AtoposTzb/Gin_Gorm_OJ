package test

import (
	"fmt"
	"log"
	"testing"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

var myKey = []byte("my_secret_key")

// 生成token
func TestGenerateToken(t *testing.T) {
	// 测试JWT生成和验证
	UserClaim := &UserClaims{
		Identity:         "user_1",
		Name:             "Get",
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		t.Fatal("生成token失败", err)
	}
	t.Log(tokenString)

}

/*
测试输出：
=== RUN   TestGenerateToken
    e:/GoStudy_itying/Gin_Gorm_OJ/test/jwt_test.go:30: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.vdX-OsHhxmyPlMpoJsOyvOt4JuCnhCos36LWs1ULTTA
--- PASS: TestGenerateToken (0.00s)
PASS
ok  	Gin_Gorm_OJ/test	2.033s
*/

// 解析token
func TestAnalyseToken(t *testing.T) {
	// 测试JWT解析
	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.vdX-OsHhxmyPlMpoJsOyvOt4JuCnhCos36LWs1ULTTA"
	userClaim := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (any, error) {
		return myKey, nil
	})
	if err != nil {
		log.Fatal("错误:", err)
	}
	if claims, ok := token.Claims.(*UserClaims); ok {
		fmt.Println(claims.Identity, claims.Name)
	} else {
		log.Fatal("unknown claims type, cannot proceed")
	}

}

/*
解析测试输出：
=== RUN   TestAnalyseToken
user_1 Get
--- PASS: TestAnalyseToken (0.00s)
PASS
ok  	Gin_Gorm_OJ/test	2.138s

*/
