package main

import (
	"user-system/config"
	"user-system/model"
	"user-system/router"

	_ "github.com/gin-gonic/gin"
)

func main() {
	config.InitDB()
	config.InitRedis()
	config.Db.AutoMigrate(&model.User{})
	r := router.SetupRouter()
	r.Run(":9090")
}
