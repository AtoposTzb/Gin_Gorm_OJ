package helper

import (
	"crypto/md5"
	"crypto/tls"
	"fmt"
	"math/rand"
	"net/smtp"
	"os"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jordan-wright/email"
)

type UserClaims struct {
	Identity string `json:"identity"`
	Name     string `json:"name"`
	IsAdmin  int    `json:"is_admin"`
	jwt.RegisteredClaims
}

func GeTMd5(str string) string {
	return fmt.Sprintf("%x", md5.Sum([]byte(str)))
}

// 生成token
var myKey = []byte("my_secret_key")

func GenerateToken(identity, name string, isAdmin int) (string, error) {
	// 测试JWT生成和验证
	UserClaim := &UserClaims{
		Identity:         identity,
		Name:             name,
		IsAdmin:          isAdmin,
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

// 发送验证码
func SendEmail(toUserEmail, code string) error {
	// 测试发送验证码
	e := email.NewEmail()
	e.From = "OJ测试网站 <t1912160135@163.com>"
	e.To = []string{toUserEmail}
	e.Subject = "验证码已发送"
	e.HTML = []byte("您的验证码<b>" + code + "</b>")
	return e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "t1912160135@163.com", "JDTbq34RedgAkFMu", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
}

// 生成uuid
// 生成唯一标识
func GenerateUUID() string {
	return uuid.New().String()
}

// 随机验证码
func CreateCode() string {
	//生成6位随机数
	code := fmt.Sprintf("%06d", rand.Intn(1000000))
	return code
}

// 代码保存的方法
func SaveCode(code []byte) (string, error) {
	dirName := "code/" + CreateCode() // 生成随机目录名
	path := dirName + "/main.go"      // 生成随机文件名
	err := os.MkdirAll(dirName, 0755) // 创建目录  0755 表示目录权限为 rwxr-xr-x
	if err != nil {
		return "", fmt.Errorf("create dir failed")
	}
	file, err := os.Create(path) // 创建文件 file是文件指针
	if err != nil {
		return "", fmt.Errorf("create file failed")
	}
	file.Write(code) // 写入文件
	if err != nil {
		return "", fmt.Errorf("write file failed")
	}
	defer file.Close()
	return path, nil // 返回文件路径
}
