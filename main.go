package main

import (
	"MetaGallery-Cloud-backend/routes"
)

func main() {
	r := routes.Router()

	// 提供静态文件服务
	r.Static("/resources", "./resources")

	r.Run(":8080")
}
