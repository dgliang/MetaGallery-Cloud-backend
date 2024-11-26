package services

import (
	"MetaGallery-Cloud-backend/models"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func SetFolderShareState(userId, folderId uint, shareState bool) error {
	// 1. 更新数据库的 share 字段
	if err := setShareFolderState(userId, folderId, shareState); err != nil {
		return fmt.Errorf("SetFolderShareState: database: %w", err)
	}

	// workingDir, err := os.Getwd()
	// if err != nil {
	// 	return fmt.Errorf("SetFolderShareState: get working dir: %w", err)
	// }
	// workingDir = strings.ReplaceAll(workingDir, "\\", "/")
	// log.Println(workingDir)

	// absolutePath := path.Join(workingDir, "resources", "./img/B.png")
	// log.Println(absolutePath)

	// ipfsHash, err := UploadFileToPinata(absolutePath)
	// if err != nil {
	// 	return fmt.Errorf("SetFolderShareState: upload file: %w", err)
	// }
	// log.Println(ipfsHash)
	// CreateGroup("")

	// jsonData := map[string]interface{}{
	// 	"folder_name": "my_folder",
	// 	"files": []map[string]interface{}{
	// 		{
	// 			"file_name": "file1.txt",
	// 			"cid":       "Qm...cid_of_file1",
	// 		},
	// 		{
	// 			"file_name": "file2.txt",
	// 			"cid":       "Qm...cid_of_file2",
	// 		},
	// 	},
	// 	"subfolders": []map[string]interface{}{
	// 		{
	// 			"folder_name": "subfolder1",
	// 			"cid":         "Qm...cid_of_subfolder",
	// 		},
	// 	},
	// }
	// UploadJsonToIPFS(jsonData)

	return nil
}

func setShareFolderState(userId, folderId uint, shareState bool) error {
	return models.DataBase.Transaction(func(tx *gorm.DB) error {
		var folder models.FolderData
		if err := tx.First(&folder, "id = ? AND belong_to = ?", folderId, userId).Error; err != nil {
			return err
		}

		if folder.Share == shareState {
			return nil
		}

		folder.Share = shareState
		tmpParentPath := folder.Path

		// 更新父文件夹的 Share 字段
		if err := tx.Save(&folder).Error; err != nil {
			return err
		}

		// 更新所有子文件和子文件夹的 Share 字段
		parentPath := strings.ReplaceAll(strings.TrimSpace(tmpParentPath), "\\", "/")
		var subFolders []models.FolderData
		if err := tx.Where("path LIKE ? AND belong_to = ?", parentPath+"/%", userId).
			Find(&subFolders).Error; err != nil {
			return err
		}

		for _, subFolder := range subFolders {
			subFolder.Share = shareState
			if err := tx.Save(&subFolder).Error; err != nil {
				return err
			}
		}

		// TODO: 更新所有子文件的 Share 字段

		return nil
	})
}

/*
1. 使用 ListChildFolders 检查是否有子文件夹
2. 对每个子文件夹，继续递归，直到所有子文件夹和文件都被处理。
3. 每一层的文件夹都会上传到 IPFS，并为每个文件夹生成一个唯一的 CID（内容标识符）。
4. 子文件夹中的文件也需要上传，并生成对应的 CID。在文件夹上传之前，必须先处理该文件夹中的所有内容（文件和子文件夹）。
	可以将文件夹和文件信息（包括子文件夹的 CID）组织成一个 JSON 数据结构，然后上传该结构到 IPFS，获得文件夹的 CID。
*/

// 递归上传文件夹
// 1. 上传当前文件夹中的文件
// 2. 递归上传子文件夹
// 3. 构建文件夹的元数据
// 4. 将文件夹结构上传到 IPFS
// 5. 返回文件夹的 CID

// 获取文件夹内部子文件夹和子文件的结构，并创建 JSON 格式文件
func getFilesInFolder() map[string]interface{} {
	return map[string]interface{}{}
}
