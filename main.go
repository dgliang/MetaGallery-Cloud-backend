package main

import (
	"MetaGallery-Cloud-backend/dao"
	"MetaGallery-Cloud-backend/routes"
)

func main() {
	r := routes.Router()

	dao.DataBaseInit() // 确保调用初始化函数

	// 提供静态文件服务
	r.Static("/resources", "./resources")

	r.Run(":8080")

}
