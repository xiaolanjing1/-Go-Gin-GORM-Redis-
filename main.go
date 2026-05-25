// @title User System API
// @version 1.0
// @description Go Gin 用户中心接口文档
// @host localhost:9090
// @BasePath /
package main

import (
	"user-system/config"
	"user-system/model"
	"user-system/router"

	_ "user-system/docs"

	_ "github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	config.InitDB()
	config.InitRedis()
	config.Db.AutoMigrate(&model.User{})
	r := router.SetupRouter()
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	r.Run(":9090")
}
