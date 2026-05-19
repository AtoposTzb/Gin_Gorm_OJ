package helper

import (
	"crypto/md5"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	jwt.RegisteredClaims
}

func GeTMd5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// 生成token
var myKey = []byte("my_secret_key")

func GenerateToken(identity, name string) (string, error) {
	// 测试JWT生成和验证
	UserClaim := &UserClaims{
		Identity:         identity,
		Name:             name,
		RegisteredClaims: jwt.RegisteredClaims{},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, UserClaim)
	tokenString, err := token.SignedString(myKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

// 解析token
func AnalyseToken(tokenString string) (*UserClaims, error) {
	// 测试JWT解析
	//	tokenString := "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZGVudGl0eSI6InVzZXJfMSIsIm5hbWUiOiJHZXQifQ.vdX-OsHhxmyPlMpoJsOyvOt4JuCnhCos36LWs1ULTTA"
	userClaim := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, userClaim, func(token *jwt.Token) (any, error) {
		return myKey, nil
	})
	if err != nil {
		return nil, err
	} else if claims, ok := token.Claims.(*UserClaims); ok {
		return claims, nil
	} else {
		return nil, fmt.Errorf("unknown claims type")
	}

}
