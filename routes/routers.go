package routes

import (
	"MetaGallery-Cloud-backend/controllers"

	"github.com/gin-gonic/gin"
)

func Router(r *gin.Engine) {

	// 请求接口都在 “/api” 的目录中
	api := r.Group("/api")
	{
		api.POST("/register", controllers.UerController{}.Register)
		api.POST("/login", controllers.UerController{}.Login)

		// 除了注册登录外，其余接口都要进行 jwt 验证
		// api.Use(middlewares.TokenAuthMiddleware())

		// 账号管理
		api.GET("/getUserInfo", controllers.UerController{}.GetUserInfo)
		api.POST("/updatePassword", controllers.UerController{}.UpdateUserPassword)
		api.POST("/updateProfile", controllers.UerController{}.UpdateUserInfo)

		// 文件夹管理
		api.GET("/getRootFolder", controllers.FolderController{}.GetRootFolder)
		api.POST("/createFolder", controllers.FolderController{}.CreateFolder)
		api.GET("/loadFolder/getChildrenInfo", controllers.FolderController{}.GetChildFolders)
		api.POST("/renameFolder", controllers.FolderController{}.RenameFolder)
		api.POST("/favoriteFolder", controllers.FolderController{}.FavoriteFolder)
		api.GET("/loadFolder/getFolderInfo", controllers.FolderController{}.GetFolderInfo)
		api.DELETE("/removeFolder", controllers.BinController{}.RemoveFolder)
		api.DELETE("/deleteFolder", controllers.BinController{}.DeleteFolder)
		api.GET("/listBinFolder", controllers.BinController{}.ListBinFolder)
		api.POST("/recoverBinFolder", controllers.BinController{}.RecoverBinFolder)
		api.POST("/shareFolder", controllers.FolderShareController{}.SetFolderShared)

		// 文件管理
		api.POST("/uploadFile", controllers.FileController{}.UploadFile)
		api.POST("/renameFile", controllers.FileController{}.RenameFile)
		api.POST("/favoriteFile", controllers.FileController{}.FavoriteFile)
		api.GET("/loadFolder/getSubFileinfo", controllers.FileController{}.GetSubFiles)
		api.GET("/getFileData", controllers.FileController{}.GetFileData)
		api.DELETE("/removeFile", controllers.FileController{}.RemoveFile)
		api.GET("/listBinFile", controllers.FileController{}.GetBinFiles)
		api.POST("/recoverBinFile", controllers.FileController{}.RecoverFile)
	}
}
