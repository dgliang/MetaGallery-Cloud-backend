package services

import (
	"MetaGallery-Cloud-backend/config"
	"MetaGallery-Cloud-backend/models"
	"fmt"
	"image/gif"
	"image/jpeg"
	"image/png"
	"os"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

func GetPreview(c *gin.Context, fileID uint) error {
	fileData, err := models.GetFileData(fileID)
	if err != nil {
		return err
	}
	if fileData.ID == 0 {
		return fmt.Errorf("ID not exists")
	}

	filePath := fileData.Path
	fullFilePath := path.Join(config.FileResPath, filePath)

	switch fileData.FileType {
	case ".jpg", ".jpeg":
		return jpegPreview(c, fullFilePath)
	case ".png":
		return pngPreview(c, fullFilePath)
	case ".gif":
		return gifPreview(c, fullFilePath)
	}

	return fmt.Errorf("该格式不支持预览")
}

func GetPreviewURL(c *gin.Context, fileID uint) (string, error) {
	fileData, err := models.GetFileData(fileID)
	if err != nil {
		return "", err
	}
	if fileData.ID == 0 {
		return "", fmt.Errorf("ID not exists")
	}

	filePath := fileData.Path
	fullFilePath := path.Join(config.FileResPath, filePath)

	switch fileData.FileType {
	case ".jpg", ".jpeg", ".png", ".gif", ".pdf":
		return DerectlyReturnFileURL(fullFilePath)
	}

	return "", fmt.Errorf("该格式不支持预览")
}

func gifPreview(c *gin.Context, fullFilePath string) error {
	file, err := os.Open(fullFilePath)
	if err != nil {

		return err
	}
	defer file.Close()

	img, err := gif.Decode(file)
	if err != nil {

		return err
	}

	// 调整图片大小
	previewIMG := resize.Resize(512, 0, img, resize.Lanczos3) // 100表示宽度，高度设为0表示按比例调整

	c.Header("Content-Type", "image/gif")
	png.Encode(c.Writer, previewIMG)

	return nil
}

func jpegPreview(c *gin.Context, fullFilePath string) error {

	file, err := os.Open(fullFilePath)
	if err != nil {

		return err
	}
	defer file.Close()

	img, err := jpeg.Decode(file)
	if err != nil {

		return err
	}

	// 调整图片大小
	previewIMG := resize.Resize(512, 0, img, resize.Lanczos3) // 100表示宽度，高度设为0表示按比例调整

	c.Header("Content-Type", "image/jpeg")
	png.Encode(c.Writer, previewIMG)
	// fmt.Println("缩略图已生成：thumbnail.jpg")

	return nil
}

func pngPreview(c *gin.Context, fullFilePath string) error {
	file, err := os.Open(fullFilePath)
	if err != nil {

		return err
	}
	defer file.Close()
	// log.Println("open file success")
	// 解码PNG图片
	img, err := png.Decode(file)
	if err != nil {

		return err
	}
	// log.Println("decode file success")
	// 调整图片大小
	previewIMG := resize.Resize(512, 0, img, resize.Lanczos3) // 100表示宽度，高度设为0表示按比例调整

	// fmt.Println("缩略图已生成：thumbnail.png")
	c.Header("Content-Type", "image/png")
	png.Encode(c.Writer, previewIMG)

	return nil
}

func DerectlyReturnFileURL(fullFilePath string) (string, error) {

	fileURL := path.Join(config.HostURL, fullFilePath)

	return fileURL, nil
}
