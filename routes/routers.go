package routes

import (
	"MetaGallery-Cloud-backend/controllers"
	"MetaGallery-Cloud-backend/middlewares"
	"github.com/gin-gonic/gin"
)

func Router() *gin.Engine {
	r := gin.Default()

	// 请求接口都在 “/api” 的目录中
	api := r.Group("/api")
	{
		api.POST("/register", controllers.UerController{}.Register)
		api.POST("/login", controllers.UerController{}.Login)

		// 除了注册登录外，其余接口都要进行 jwt 验证
		api.Use(middlewares.TokenAuthMiddleware())

		api.GET("/getUserInfo", controllers.UerController{}.GetUserInfo)
		api.POST("/updatePassword", controllers.UerController{}.UpdateUserPassword)
		api.POST("/updateProfile", controllers.UerController{}.UpdateUserInfo)
	}

	return r
}
