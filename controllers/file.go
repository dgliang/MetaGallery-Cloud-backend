package controllers

import (
	"MetaGallery-Cloud-backend/models"
	"MetaGallery-Cloud-backend/services"
	"fmt"
	"log"
	"path/filepath"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileController struct{}

type FileBriefJson struct {
	ID         uint   `json:"id"`
	User       uint   `json:"user"`
	FileName   string `json:"file_name"`
	ParentID   uint   `json:"parent_id"`
	Path       string `json:"path"`
	IsFavorite bool   `json:"is_favorite"`
	IsShare    bool   `json:"is_share"`
	IPFSHash   string `json:"ipfs_hash"`
	Deleted    string `json:"is_deleted"`
}

func (receiver FileController) UploadFile(c *gin.Context) {
	account := c.DefaultPostForm("account", "")
	parentFolderID := c.DefaultPostForm("parent_id", "-1")
	fileName := c.DefaultPostForm("file_name", "")

	if account == "" {
		log.Printf("from %s 上传文件提供的信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "上传者账号不能为空")
		return
	}
	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	if parentFolderID == "-1" {
		log.Printf("from %s 上传文件提供的信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "父文件夹不能为空")
		return
	}
	if fileName == "" {
		log.Printf("from %s 上传文件提供的信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "文件名不能为空")
		return
	}

	PID, err := strconv.ParseUint(parentFolderID, 10, 0) // 10是进制，0是自动推断结果位数
	if err != nil {
		fmt.Println("转换出错:", err)
		return
	}
	// 将 uint64 转为 uint
	uintPID := uint(PID)

	if services.FileExist(userID, uintPID, fileName) {
		ReturnError(c, "FAILED", "文件夹下文件重名")
		return
	}

	//读取文件内容
	file, _, err := c.Request.FormFile("file")
	if err != nil {
		log.Printf("from %s 接收文件失败\n", c.Request.Host)
		ReturnError(c, "FAILED", "接收文件失败")
		return
	}
	defer file.Close()

	newfile, err := models.CreateFileData(userID, fileName, uintPID, filepath.Ext(fileName))
	if err != nil {
		ReturnServerError(c, err.Error())
		return
	}

	//在本地创建文件
	saveFileError := services.SaveFileByID(userID, uintPID, newfile.ID, file)
	if saveFileError != nil {
		ReturnServerError(c, saveFileError.Error())
		models.UnscopedDeleteFileData(newfile.ID)
		return
	}

	fileRes := models.FileBrief{
		ID:       newfile.ID,
		FileName: newfile.FileName,
		FileType: filepath.Ext(newfile.FileName),
		Favorite: newfile.Favorite,
		Share:    newfile.Share,
		InBin:    newfile.DeletedAt.Time,
	}

	ReturnSuccess(c, "SUCCESS", "上传文件成功", fileRes)
}

func (receiver FileController) DownloadFile(c *gin.Context) {
	account := c.Query("account")
	fileID := c.Query("file_id")

	if account == "" {
		log.Printf("from %s 上传文件提供的信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "上传者账号不能为空")
		return
	}
	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	if fileID == "" {
		log.Printf("from %s 查询子文件提供的文件夹信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "文件ID不能为空")
		return
	}
	// %isfilebelongto

	FID, err := strconv.ParseUint(fileID, 10, 0)
	if err != nil {
		fmt.Println("转换出错:", err)
		return
	}
	uintFID := uint(FID)

	if !services.IsFileBelongto(userID, uintFID) {
		c.JSON(403, gin.H{
			"error":   "FORBIDDEN",
			"message": "访问禁止",
		})
		c.Abort()
		return
	}

	fileData, err := models.GetFileData(uintFID)
	if err != nil {
		ReturnServerError(c, err.Error())
		return
	}
	if fileData.ID == 0 {
		ReturnError(c, "FAILED", "文件不存在")
		return
	}

	services.DownloadFile(c, userID, uintFID)

}

func (receiver FileController) RenameFile(c *gin.Context) {
	// 从http请求中获取参数
	account := c.DefaultPostForm("account", "")
	fileID := c.DefaultPostForm("file_id", "-1")
	newFileName := c.DefaultPostForm("new_file_name", "")

	// 验证参数合法性
	if account == "" {
		log.Printf("from %s 上传文件提供的信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "上传者账号不能为空")
		return
	}
	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}
	if newFileName == "" {
		log.Printf("from %s 上传文件提供的信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "新文件名不能为空")
		return
	}

	// 将参数类型转换便于后续函数操作
	FID, err := strconv.ParseUint(fileID, 10, 0)
	if err != nil {
		fmt.Println("转换出错:", err)
		return
	}
	uintFID := uint(FID)

	// 操作合法性检查
	Belongto := services.IsFileBelongto(userID, uintFID)
	if !Belongto {
		c.JSON(403, gin.H{
			"error":   "FORBIDDEN",
			"message": "访问禁止",
		})
		c.Abort()
		return
	}

	// 进行操作
	if err := services.RenameFile(userID, uintFID, newFileName); err != nil {
		ReturnError(c, "FAILED", "重命名失败:"+err.Error())
		return
	}

	// 返回成功信息
	ReturnSuccess(c, "SUCCESS", "重命名文件成功", nil)
}

// type getFilesJson struct {
// 	Account  string `json:"account" binding:"required"`
// 	FolderID uint   `json:"folder_id" binding:"required"`
// }

func (receiver FileController) GetSubFiles(c *gin.Context) {

	account := c.Query("account")

	folderID := c.Query("folder_id")
	if folderID == "" {
		log.Printf("from %s 查询子文件提供的文件夹id不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "文件夹id不能为空")
		return
	}

	if account == "" {
		log.Printf("from %s 查询子文件提供的账号不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "账号不能为空")
		return
	}

	FID, err := strconv.ParseUint(folderID, 10, 0)
	if err != nil {
		fmt.Println("转换出错:", err)
		return
	}
	// 将 uint64 转为 uint
	uintFID := uint(FID)
	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	subfiles, err := services.GetSubFiles(uintFID)
	if err != nil {
		ReturnServerError(c, err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", "", subfiles)
}

// type getFileDetailJson struct {
// 	Account string `json:"account" binding:"required"`
// 	FileID  uint   `json:"file_id" binding:"required"`
// }

func (receiver FileController) GetFileData(c *gin.Context) {

	account := c.Query("account")

	fileID := c.Query("file_id")

	if account == "" {
		log.Printf("from %s 查询文件提供的账号不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "账号不能为空")
		return
	}
	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	if fileID == "" {
		log.Printf("from %s 查询子文件提供的文件夹信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "文件ID不能为空")
		return
	}

	FID, err := strconv.ParseUint(fileID, 10, 0)
	if err != nil {
		fmt.Println("转换出错:", err)
		return
	}
	uintFID := uint(FID)

	fileData, err := services.GetFileDetail(uintFID)
	if err != nil {
		ReturnServerError(c, err.Error())
		return
	}
	if fileData.ID == 0 {
		ReturnError(c, "FAILED", "ID不存在")
		return
	}

	ReturnSuccess(c, "SUCCESS", "获取文件详细信息成功", fileData)
}

type favoriteFileRequest struct {
	Account    string `json:"account" binding:"required"`
	FileId     uint   `json:"file_id" binding:"required"`
	IsFavorite int    `json:"is_favorite" binding:"required"`
}

func (receiver FileController) FavoriteFile(c *gin.Context) {
	req, _ := c.Get("jsondata")
	// fmt.Println(req)
	jsondata, ok := req.(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{
			"status": "ERROR",
			"msg":    "服务端获取请求体错误",
		})
		c.Abort()
		return
	}

	isFavorite := jsondata["is_favorite"].(float64)

	account := jsondata["account"].(string)

	fileID := jsondata["file_id"].(float64)
	uintFID := uint(fileID)

	var favoriteStatus bool
	if isFavorite == 1 {
		favoriteStatus = false
	} else if isFavorite == 2 {
		favoriteStatus = true
	} else {
		ReturnError(c, "FAILED", "is_favorite 的取值只能是 1 或者 2")
		return
	}
	fmt.Println("favoriteStatus:", favoriteStatus)

	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "GetUserID"+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "FAILED", "提供的用户不存在")
		return
	}

	Belongto := services.IsFileBelongto(userID, uintFID)
	if !Belongto {
		c.JSON(403, gin.H{
			"error":   "FORBIDDEN",
			"message": "访问禁止",
		})
		c.Abort()
		return
	}

	fileData, err1 := models.GetFileData(uintFID)
	if err1 != nil {
		ReturnServerError(c, "GetFileData: "+err1.Error())
		return
	}
	if fileData.ID == 0 {
		ReturnError(c, "FAILED", "文件不存在")
		return
	}

	// 更新文件的收藏状态
	if favoriteStatus {
		models.SetFileFavorite(uintFID)
	} else {
		models.CancelFileFavorite(uintFID)
	}

	ReturnSuccess(c, "SUCCESS", fmt.Sprintf("成功将 %s 的 %d 文件收藏状态改为 %t",
		account, uintFID, favoriteStatus))
}

func (receiver FileController) GetFavorFiles(c *gin.Context) {

	account := c.Query("account")

	if account == "" {
		log.Printf("from %s 查询收藏文件提供的账号不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "账号不能为空")
		return
	}
	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	favorFiles, err := services.GetAllFavorFiles(userID)
	if err != nil {
		ReturnServerError(c, err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", "", favorFiles)
}

type deleteOrRecoverFilejson struct {
	Account string `json:"account" binding:"required"`
	Fileid  uint   `json:"file_id" binding:"required"`
}

func (receiver FileController) RemoveFile(c *gin.Context) {
	req, _ := c.Get("jsondata")
	jsondata, ok := req.(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{
			"status": "ERROR",
			"msg":    "服务端获取请求体错误",
		})
		c.Abort()
		return
	}

	account := jsondata["account"].(string)

	fileID := jsondata["file_id"].(float64)
	uintFID := uint(fileID)

	if account == "" {
		log.Printf("from %s 将文件移入回收站提供的账号不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "账号不能为空")
		return
	}
	if fileID == 0 {
		log.Printf("from %s 将文件移入回收提供的文件信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "文件ID不能为空")
		return
	}

	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	Belongto := services.IsFileBelongto(userID, uintFID)
	if !Belongto {
		c.JSON(403, gin.H{
			"error":   "FORBIDDEN",
			"message": "访问禁止",
		})
		c.Abort()
		return
	}

	if err := services.RemoveFile(uintFID); err != nil {
		ReturnError(c, "FAILED", err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", "文件移入回收站成功", nil)
}

func (receiver FileController) GetBinFiles(c *gin.Context) {

	account := c.Query("account")

	if account == "" {
		log.Printf("from %s 查询已回收文件提供的账号不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "账号不能为空")
		return
	}

	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	binFiles, err := services.GetBinFiles(userID)
	if err != nil {
		ReturnServerError(c, err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", "该账户已回收以下文件", binFiles)
}

func (receiver FileController) RecoverFile(c *gin.Context) {
	req, _ := c.Get("jsondata")
	jsondata, ok := req.(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{
			"status": "ERROR",
			"msg":    "服务端获取请求体错误",
		})
		c.Abort()
		return
	}

	account := jsondata["account"].(string)

	fileID := jsondata["file_id"].(float64)
	uintFID := uint(fileID)

	if account == "" {
		log.Printf("from %s 将文件移出回收站提供的账号不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "账号不能为空")
		return
	}
	if uintFID == 0 {
		log.Printf("from %s 将文件移出回收提供的文件信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "文件ID不能为空")
		return
	}

	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	Belongto := services.IsFileBelongto(userID, uintFID)
	if !Belongto {
		c.JSON(403, gin.H{
			"error":   "FORBIDDEN",
			"message": "访问禁止",
		})
		c.Abort()
		return
	}

	if err := services.RecoverFile(userID, uintFID); err != nil {
		ReturnError(c, "FAILED", err.Error())
	}
	ReturnSuccess(c, "SUCCESS", "文件移出回收站成功", nil)
}

func (receiver FileController) ActuallyDeleteFile(c *gin.Context) {
	req, _ := c.Get("jsondata")
	jsondata, ok := req.(map[string]interface{})
	if !ok {
		c.JSON(400, gin.H{
			"status": "ERROR",
			"msg":    "服务端获取请求体错误",
		})
		c.Abort()
		return
	}

	account := jsondata["account"].(string)

	fileID := jsondata["file_id"].(float64)
	uintFID := uint(fileID)

	if account == "" {
		log.Printf("from %s 将文件彻底删除提供的账号不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "账号不能为空")
		return
	}
	if fileID == 0 {
		log.Printf("from %s 将文件移彻底删除提供的文件信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "文件ID不能为空")
		return
	}

	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	Belongto := services.IsFileBelongto(userID, uintFID)
	if !Belongto {
		c.JSON(403, gin.H{
			"error":   "FORBIDDEN",
			"message": "访问禁止",
		})
		c.Abort()
		return
	}

	err2 := services.ActuallyDeleteFile(uintFID)
	if err2 != nil {
		ReturnError(c, "FAILED", err2.Error())
		return
	}
	ReturnSuccess(c, "SUCCESS", "文件彻底删除成功", nil)
}

func (receiver FileController) PreviewFile(c *gin.Context) {

	account := c.Query("account")

	fileID := c.Query("file_id")

	if account == "" {
		log.Printf("from %s 上传文件提供的信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "上传者账号不能为空")
		return
	}
	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	if fileID == "" {
		log.Printf("from %s 查询子文件提供的文件夹信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "文件ID不能为空")
		return
	}

	FID, err := strconv.ParseUint(fileID, 10, 0)
	if err != nil {
		fmt.Println("转换出错:", err)
		return
	}
	uintFID := uint(FID)

	err2 := services.GetPreview(c, uintFID)
	if err2 != nil {
		log.Printf("生成预览失败 ：%s", err2)
		ReturnError(c, "FAILED", err2.Error())
		return
	}

}

func (receiver FileController) GetBinFileData(c *gin.Context) {

	account := c.Query("account")

	fileID := c.Query("file_id")

	if account == "" {
		log.Printf("from %s 上传文件提供的信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "上传者账号不能为空")
		return
	}
	userID, err := models.GetUserID(account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	if fileID == "" {
		log.Printf("from %s 查询子文件提供的文件夹信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "文件ID不能为空")
		return
	}

	FID, err := strconv.ParseUint(fileID, 10, 0)
	if err != nil {
		fmt.Println("转换出错:", err)
		return
	}
	uintFID := uint(FID)

	fileData, err2 := services.GetBinFileData(uintFID)
	if err2 != nil {
		ReturnError(c, "FAILED", err2.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", "获取文件信息成功", fileData)

}
