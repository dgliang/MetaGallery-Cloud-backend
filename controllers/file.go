package controllers

import (
	"MetaGallery-Cloud-backend/models"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
)

type FileController struct{}

type FileJson struct {
	ID         uint   `json:"id"`
	User       uint   `json:"user"`
	FileName   string `json:"file_name"`
	ParentID   uint   `json:"parent_id"`
	Path       string `json:"path"`
	IsFavorite bool   `json:"is_favorite"`
	IsShare    bool   `json:"is_share"`
	IPFSHash   string `json:"ipfs_hash"`
	IsDeleted  bool   `json:"is_deleted"`
	DeleteTime string `json:"delete_time"`
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

	//读取文件内容
	file, _, err := c.Request.FormFile("file_context")
	if err != nil {
		log.Printf("from %s 接收文件失败\n", c.Request.Host)
		ReturnError(c, "FAILED", "接受文件失败")
		return
	}
	defer file.Close()

	PID, err := strconv.ParseUint(parentFolderID, 10, 0) // 10是进制，0是自动推断结果位数
	if err != nil {
		fmt.Println("转换出错:", err)
		return
	}
	// 将 uint64 转为 uint
	uintPID := uint(PID)

	path, err := models.GetFilePath(userID, fileName, uintPID)
	if err != nil {
		fmt.Println("文件路径生成失败:", err)
		ReturnServerError(c, "文件路径生成失败")
		return
	}
	//在本地创建文件

	out, err := os.Create("userFiles" + path)
	//out, err := os.Create(path)
	if err != nil {
		log.Printf("from %s 创建文件失败\n", c.Request.Host)
		ReturnServerError(c, "服务器创建文件失败")
		return
	}
	defer out.Close()

	// 将上传的文件内容写入到本地文件
	_, err = out.ReadFrom(file)
	if err != nil {
		log.Printf("from %s 写入保存文件失败\n", c.Request.Host)
		ReturnServerError(c, "服务器写入保存文件失败")
		return
	}

	newfile, err := models.CreateFileData(userID, fileName, uintPID)
	if err != nil {
		ReturnServerError(c, err.Error())
		return
	}

	fileRes := FileJson{
		ID:         newfile.ID,
		User:       newfile.BelongTo,
		FileName:   newfile.FileName,
		ParentID:   newfile.ParentFolderID,
		Path:       newfile.Path,
		IsFavorite: newfile.Favorite,
		IsShare:    newfile.Share,
		IsDeleted:  newfile.InBin,
	}

	ReturnSuccess(c, "SUCCESS", "", fileRes)

}

func (receiver FileController) RenameFile(c *gin.Context) {
	account := c.DefaultPostForm("account", "")
	fileID := c.DefaultPostForm("file_id", "-1")
	newFileName := c.DefaultPostForm("new_file_name", "")

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

	FID, err := strconv.ParseUint(fileID, 10, 0) // 10是进制，0是自动推断结果位数
	if err != nil {
		fmt.Println("转换出错:", err)
		return
	}
	// 将 uint64 转为 uint
	uintFID := uint(FID)

	models.RenameFileWithFileID(uintFID, newFileName)

	ReturnSuccess(c, "SUCCESS", "", models.FileData{})

}

type getFilesJson struct {
	Account  string `json:"account" binding:"required"`
	FolderID uint   `json:"folder_id" binding:"required"`
}

func (receiver FileController) GetSubFiles(c *gin.Context) {

	var Request getFilesJson
	// 绑定JSON数据
	if err := c.ShouldBindJSON(&Request); err != nil {
		log.Printf("from %s 查询子文件提供的json绑定失败\n", c.Request.Host)
		ReturnServerError(c, "解析 JSON Request："+err.Error())
		return
	}

	if Request.Account == "" {
		log.Printf("from %s 查询子文件提供的账号不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "账号不能为空")
		return
	}
	userID, err := models.GetUserID(Request.Account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	if Request.FolderID == 0 {
		log.Printf("from %s 上传文件提供的文件夹信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "文件夹ID不能为空")
		return
	}

	subfiles, err := models.GetSubFiles(Request.FolderID)
	if err != nil {
		ReturnServerError(c, err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", "", subfiles)
}

type getFileDetailJson struct {
	Account string `json:"account" binding:"required"`
	FileID  uint   `json:"file_id" binding:"required"`
}

func (receiver FileController) GetFileData(c *gin.Context) {

	var Request getFileDetailJson
	// 绑定JSON数据
	if err := c.ShouldBindJSON(&Request); err != nil {
		log.Printf("from %s 查询子文件提供的json绑定失败\n", c.Request.Host)
		ReturnServerError(c, "解析 JSON Request："+err.Error())
		return
	}

	if Request.Account == "" {
		log.Printf("from %s 查询子文件提供的账号不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "账号不能为空")
		return
	}
	userID, err := models.GetUserID(Request.Account)
	if err != nil {
		ReturnServerError(c, "获取 GetUserID: "+err.Error())
		return
	}
	if userID == 0 {
		ReturnError(c, "Failed", "用户不存在")
		return
	}

	if Request.FileID == 0 {
		log.Printf("from %s 上传文件提供的文件夹信息不全\n", c.Request.Host)
		ReturnError(c, "FAILED", "文件夹ID不能为空")
		return
	}

	fileData, err := models.GetFileData(Request.FileID)
	if err != nil {
		ReturnServerError(c, err.Error())
		return
	}

	ReturnSuccess(c, "SUCCESS", "", fileData)
}
