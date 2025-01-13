package services

import (
	"MetaGallery-Cloud-backend/config"
	"MetaGallery-Cloud-backend/models"
	"fmt"
	"path"

	// "github.com/unidoc/unipdf/v3/model"

	"github.com/gin-gonic/gin"
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

	return fmt.Errorf("该格式不支持预览")

}

func DerectlyReturnFileURL(fullFilePath string) (string, error) {

	fileURL := path.Join(config.HOST_URL, fullFilePath)

	return fileURL, nil
}
