# Gin_Gorm_OJ
基于Gin、Gorm实现在线练习系统--两个星期学完Go/Gin/Zero/Gorm后，首个练手小项目。

>语言：Go、框架：Gin、ORM：Gorm、数据库：MySQL

# 准备工作
## 环境配置
1. 安装Go语言环境
   - 可以参考[Go官方文档](https://go.dev/doc/install)进行安装。
2. 安装Gin框架
   - 可以参考[gin官方文档](https://gin-gonic.com/docs/installation/)进行安装。
   - 安装完成后，需要在项目根目录下执行`go mod tidy`项目依赖。
   - 这个命令可以自动下载并安装项目依赖的库。后面操作时演示。
3. 安装Gorm ORM
   - 可以参考[Gorm官方文档](https://gorm.io/docs/installation.html)进行安装。
   - 安装完成后，执行`go mod tidy`命令，安装项目依赖。
4. 安装MySQL数据库
   - 本地练手项目，使用的是可视化数据库工具[Navicat Premium](https://www.navicat.com/zh-cn/product/navicat-premium)。
## 项目开发前准备
   - github建仓库-空仓库可以有readme.md文件
   - 本地clone仓库，vscode or traeIDE 打开项目目录
     - 初始化module，执行`go mod init gin_gorm_oj`命令
     - 可视化数据库工具Navicat Premium建表，创建数据库表结构,详见models文件夹-gorm的模型定义。
   - 数据库配置，开始开发项目的实际功能。



