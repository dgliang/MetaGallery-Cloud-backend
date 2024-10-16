package main

import "MetaGallery-Cloud-backend/routes"

func main() {
	r := routes.Router()

	r.Run(":8080")
}
