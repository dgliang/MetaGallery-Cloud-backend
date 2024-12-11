package main

import (
	"MetaGallery-Cloud-backend/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default() // 创建默认的 gin 引擎

	r.Use(cors.New(cors.Config{
		AllowAllOrigins: true,                                                // 允许所有来源
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},            // 允许的 HTTP 方法
		AllowHeaders:    []string{"Origin", "Content-Type", "Authorization"}, // 允许的请求头
	}))

	routes.Router(r)
	// 提供静态文件服务
	r.Static("/resources/img", "./resources/img")

	// 添加静态文件服务：假设 public 文件夹中包含 register.html 和其他静态文件
	r.Static("/public", "./public")
	r.StaticFile("/favicon.ico", "./public/favicon.ico")

	r.Run(":8080")

}
