package routes

import (
	"MetaGallery-Cloud-backend/controllers"
	"MetaGallery-Cloud-backend/middlewares"
	"fmt"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {
	//r := gin.Default()
	fmt.Println(1) // 打印数字 1
	// 添加静态文件服务：假设 public 文件夹中包含 register.html 和其他静态文件
	// r.Static("/public", "./public")
	// r.StaticFile("/favicon.ico", "./public/favicon.ico")

	// r.POST("/api/register", func(c *gin.Context) { // 测试是否能命中 /api/register
	// 	fmt.Println("Static /api/register endpoint hit")
	// 	c.String(200, "Static test successful")
	// })
	// 请求接口都在 “/api” 的目录中
	api := r.Group("/api")
	{
		fmt.Println(2)
		api.POST("/register", controllers.UerController{}.Register)
		fmt.Println(2)
		api.POST("/login", controllers.UerController{}.Login)
		fmt.Println(2)
		// 除了注册登录外，其余接口都要进行 jwt 验证
		api.Use(middlewares.TokenAuthMiddleware())

		api.GET("/getUserInfo", controllers.UerController{}.GetUserInfo)
		api.POST("/updatePassword", controllers.UerController{}.UpdateUserPassword)
		api.POST("/updateProfile", controllers.UerController{}.UpdateUserInfo)
	}
	fmt.Println(2)
	//return r
}
