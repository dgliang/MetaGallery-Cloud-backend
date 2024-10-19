package main

import (
	"MetaGallery-Cloud-backend/dao"
	"MetaGallery-Cloud-backend/routes"
)

func main() {
	r := routes.Router()

	dao.DataBaseInit() // 确保调用初始化函数

	r.Run(":8080")

}
