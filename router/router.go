package router

import (
	"user-system/controller"
	"user-system/middleware"

	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/register", controller.Register)

	r.POST("/login", controller.Login)

	r.Static("/uploads", "./uploads")

	auth := r.Group("/user")
	auth.Use(middleware.JWTAuth())
	{
		auth.GET("/info", controller.UserInfo)

		auth.GET("/profile", controller.UserProfile)

		auth.GET("/users", controller.UserList)

		auth.PUT("/update", controller.UpdateUser)

		auth.PUT("/password", controller.UpdatePassword)

		auth.POST("/avatar", controller.UploadAvatar)

		auth.POST("/logout", controller.Logout)

		auth.POST("/updateemail", controller.Updateemial)
	}
	return r
}
