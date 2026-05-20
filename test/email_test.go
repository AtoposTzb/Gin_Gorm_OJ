package test

import (
	"crypto/tls"
	"net/smtp"
	"testing"

	"github.com/jordan-wright/email"
)

func TestSendEmail(t *testing.T) {
	// 测试发送验证码
	e := email.NewEmail()
	e.From = "OJ测试网站 <t1912160135@163.com>"
	e.To = []string{"3488447218@qq.com"}
	// e.Bcc = []string{"test_bcc@example.com"}
	// e.Cc = []string{"test_cc@example.com"}
	e.Subject = "验证码发送测试"
	// e.Text = []byte("")
	e.HTML = []byte("您的验证码<b>123456</b>")
	//e.Send("smtp.163.com:465", smtp.PlainAuth("", "t1912160135@163.com", "网易授权码", "smtp.163.com"))
	//返回EOF时,关闭SSL，使用下面的这个：
	err := e.SendWithTLS("smtp.163.com:465", smtp.PlainAuth("", "t1912160135@163.com", "网易授权码", "smtp.163.com"),
		&tls.Config{InsecureSkipVerify: true, ServerName: "smtp.163.com"})
	if err != nil {
		t.Errorf("发送邮件失败: %v", err)
	}
}

/*
这里使用网易邮箱和qq邮箱，需要准备邮箱服务账号和密码
发信服务器和收信服务器自己去邮件网站查看即可
我这里准备了网易的：
smtp.163.com:465
密码：省略，这里应该使用授权码，而不是登录邮箱的密码

测试结果：
=== RUN   TestSendEmail
--- PASS: TestSendEmail (2.23s)
PASS
ok  	Gin_Gorm_OJ/test	4.314s
qq邮箱收到验证码邮件
*/
