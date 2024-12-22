package services

import (
	"MetaGallery-Cloud-backend/config"
	"MetaGallery-Cloud-backend/models"
	"fmt"
	"image/jpeg"
	"image/png"
	"os"
	"path"

	// "github.com/unidoc/unipdf/v3/model"

	"github.com/gin-gonic/gin"
	"github.com/nfnt/resize"
)

var PreviewAvailable = []string{".jpg", ".jpeg", ".png", ".gif", ".svg", ".pdf", ".txt", ".csv", ".json", ".xml"}

func GetPreview(c *gin.Context, fileID uint) error {
	fileData, err := models.GetFileData(fileID)
	if err != nil {
		return err
	}
	if fileData.ID == 0 {
		return fmt.Errorf("ID not exists")
	}

	filePath := fileData.Path
	fullFilePath := path.Join(config.FILE_RES_PATH, filePath)

	for _, filetype := range PreviewAvailable {
		if filetype == fileData.FileType {
			c.File(fullFilePath)
			return nil
		}
	}
	// switch fileData.FileType {
	// case ".jpg", ".jpeg":
	// 	return jpegPreview(c, fullFilePath)
	// case ".png":
	// 	return pngPreview(c, fullFilePath)
	// case ".gif":
	// 	return gifPreview(c, fullFilePath)
	// case ".pdf":
	// 	return pdfPreview(c, fullFilePath)
	// case ".svg":
	// 	return svgPreview(c, fullFilePath)
	// }
	return fmt.Errorf("该格式不支持预览")
	// return fmt.Errorf("该格式不支持预览")
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
	fullFilePath := path.Join(config.FILE_RES_PATH, filePath)

	switch fileData.FileType {
	case ".jpg", ".jpeg", ".png", ".gif", ".pdf":
		return DerectlyReturnFileURL(fullFilePath)
	}

	return "", fmt.Errorf("该格式不支持预览")
}

func gifPreview(c *gin.Context, fullFilePath string) error {
	// file, err := os.Open(fullFilePath)
	// if err != nil {

	// 	return err
	// }
	// defer file.Close()

	// img, err := gif.Decode(file)
	// if err != nil {

	// 	return err
	// }

	// 调整图片大小
	// previewIMG := resize.Resize(512, 0, img, resize.Lanczos3) // 100表示宽度，高度设为0表示按比例调整

	c.Header("Content-Type", "image/gif")
	// png.Encode(c.Writer, previewIMG)
	c.File(fullFilePath)

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
	// 解码PNG图片
	img, err := png.Decode(file)
	if err != nil {

		return err
	}
	// 调整图片大小
	previewIMG := resize.Resize(512, 0, img, resize.Lanczos3) // 100表示宽度，高度设为0表示按比例调整

	// fmt.Println("缩略图已生成：thumbnail.png")
	c.Header("Content-Type", "image/png")
	png.Encode(c.Writer, previewIMG)

	return nil
}

func pdfPreview(c *gin.Context, fullFilePath string) error {
	// file, err := os.Open(fullFilePath)
	// if err != nil {

	// 	return err
	// }
	// defer file.Close()
	// // 读取PDF文件
	// pdfReader, err := model.NewPdfReader(file)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read PDF"})
	// 	return err
	// }

	// // 创建一个新的PDF文档
	// newPdf := model.NewPdfDocument()

	// // 获取PDF的第一页
	// page, err := pdfReader.GetPage(1)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get page 1"})
	// 	return err
	// }

	// // 将第一页添加到新的PDF中
	// err = newPdf.AddPage(page)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add page"})
	// 	return err
	// }

	// // 将裁剪后的PDF写入响应
	// c.Header("Content-Type", "application/pdf")
	// c.Header("Content-Disposition", "inline; filename=first_page.pdf")
	// err = newPdf.Write(c.Writer)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to write PDF"})
	// 	return err
	// }

	c.Header("Content-Type", "application/pdf")
	// c.Header("Content-Disposition", "inline; filename=文件预览")
	c.File(fullFilePath)

	return nil
}

func svgPreview(c *gin.Context, fullFilePath string) error {

	c.Header("Content-Type", "image/svg+xml")
	// c.Header("Content-Disposition", "inline; filename=文件预览")
	c.File(fullFilePath)

	return nil
}

func DerectlyReturnFileURL(fullFilePath string) (string, error) {

	fileURL := path.Join(config.HOST_URL, fullFilePath)

	return fileURL, nil
}
