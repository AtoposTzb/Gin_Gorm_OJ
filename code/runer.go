package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
)

func main() {
	//go run "e:\GoStudy_itying\Gin_Gorm_OJ\code\code_user\main.go"
	//go run code_user/main.go
	cmd := exec.Command("go", "run", "code_user/main.go") //这是创建一个命令对象
	//
	var out, stderr bytes.Buffer //这是创建一个缓冲区，用于存储命令的输出
	cmd.Stdout = &out            //将命令的输出重定向到缓冲区
	cmd.Stderr = &stderr         //将命令的错误输出重定向到缓冲区

	//根据测试的输入按列进行运行，拿到输出结果和标准的输出结果进行对比
	stdin, err := cmd.StdinPipe() //这是创建一个管道，用于将命令的输入重定向到管道
	if err != nil {
		panic("创建管道失败:" + err.Error())
	}
	defer stdin.Close()
	io.WriteString(stdin, "1 2\n") //这是将测试的输入写入管道 然后运行命令
	if err := cmd.Run(); err != nil {
		panic("运行命令失败:" + err.Error())
	}
	log.Println(out.String()) //这是将缓冲区中的输出打印出来
	fmt.Println(out.String())

	//对比缓存中的输出和标准的输出结果是否一致
	if out.String() != "3\n" {
		log.Println("输出结果不一致")
		return
	}
	log.Println("输出结果一致")
}
