package main

import "Gin_Gorm_OJ/router"

func main() {

	r := router.InitRouter()
	r.Run(":8080")
}
