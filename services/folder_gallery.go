package services

import (
	"MetaGallery-Cloud-backend/config"
	"MetaGallery-Cloud-backend/models"
	"errors"
	"io"
	"net/http"
	"os"
	"path"

	"github.com/gin-gonic/gin"
)

const (
	PAGE_SIZE = 12 // 每页条数
)

type FolderGallery struct {
	Result interface{} `json:"result"`
}

type sharedFolderResponse struct {
	OwnerAccount UserInfo `json:"owner_account"`
	FolderName   string   `json:"folder_name"`
	IPFSHash     string   `json:"ipfs_hash"`
	Intro        string   `json:"intro"`
	CoverImg     string   `json:"cover_img"`
	PinDate      string   `json:"pin_date"`
}

// 获取用户自己共享的所有文件夹列表
func ListUserSharedFolders(ownerId uint) (FolderGallery, error) {
	var sharedFolders []models.SharedFolder
	if err := models.DataBase.Model(&models.SharedFolder{}).Where("owner_id = ?", ownerId).
		Find(&sharedFolders).Error; err != nil {
		return FolderGallery{}, err
	}

	res := matchSharedFolderModelToResponse(sharedFolders)
	return res, nil
}

// 获取所有用户的共享文件夹列表
func ListAllSharedFolders(pageNum int) (FolderGallery, error) {
	if pageNum <= 0 {
		return FolderGallery{}, errors.New("page number should be greater than 0")
	}

	// 计算偏移量
	offset := (pageNum - 1) * PAGE_SIZE
	var sharedFolders []models.SharedFolder
	if err := models.DataBase.Model(&models.SharedFolder{}).
		Limit(PAGE_SIZE).
		Offset(offset).
		Find(&sharedFolders).Error; err != nil {
		return FolderGallery{}, err
	}

	// 计算总页数
	var totalRecords int64
	if err := models.DataBase.Model(&models.SharedFolder{}).
		Count(&totalRecords).Error; err != nil {
		return FolderGallery{}, err
	}

	totalPage := int((totalRecords + int64(PAGE_SIZE) - 1) / int64(PAGE_SIZE))

	res := matchSharedFolderModelToResponse(sharedFolders, totalPage)
	return res, nil
}

// 将 models.SharedFolde 处理成返回的 JSON，传入 totalPages 参数时进行分页，否则不分页
func matchSharedFolderModelToResponse(folders []models.SharedFolder, totalPages ...int) FolderGallery {
	// 不分页的情况下
	if len(totalPages) == 0 {
		var sharedFolders []sharedFolderResponse
		for _, folder := range folders {
			var ownerAccount models.UserData
			if err := models.DataBase.Where("id = ?", folder.OwnerID).First(&ownerAccount).Error; err != nil {
				return FolderGallery{}
			}

			sharedFolders = append(sharedFolders, sharedFolderResponse{
				OwnerAccount: UserInfo{
					Account: ownerAccount.Account,
					Name:    ownerAccount.UserName,
					Intro:   ownerAccount.BriefIntro,
					Avatar:  ownerAccount.ProfilePhoto,
				},
				FolderName: folder.SharedName,
				IPFSHash:   folder.IPFSHash,
				Intro:      folder.Intro,
				CoverImg:   folder.CoverImg,
				PinDate:    folder.CreatedAt.Format("2006-01-02 15:04:05"),
			})
		}

		return FolderGallery{
			Result: sharedFolders,
		}
	}

	// 分页的情况下
	totalPage := totalPages[0]
	var sharedFolders []sharedFolderResponse
	for _, folder := range folders {
		var ownerAccount models.UserData
		if err := models.DataBase.Where("id = ?", folder.OwnerID).First(&ownerAccount).Error; err != nil {
			return FolderGallery{}
		}

		sharedFolders = append(sharedFolders, sharedFolderResponse{
			OwnerAccount: UserInfo{
				Account: ownerAccount.Account,
				Name:    ownerAccount.UserName,
				Intro:   ownerAccount.BriefIntro,
				Avatar:  ownerAccount.ProfilePhoto,
			},
			FolderName: folder.SharedName,
			IPFSHash:   folder.IPFSHash,
			Intro:      folder.Intro,
			CoverImg:   folder.CoverImg,
			PinDate:    folder.CreatedAt.Format("2006-01-02 15:04:05"),
		})
	}
	return FolderGallery{
		Result: struct {
			Folders   []sharedFolderResponse `json:"folders"`
			TotalPage int                    `json:"total_page"`
		}{
			Folders:   sharedFolders,
			TotalPage: totalPage,
		},
	}
}

type sharedFolderInfoResponse struct {
	OwnerAccount string `json:"owner_account"`
	FolderInfo   folder `json:"folder_info"`
}

func GetSharedFolderByIPFSHash(owerId uint, cid string) (models.SharedFolder, error) {
	var sharedFolder models.SharedFolder
	if err := models.DataBase.Where("owner_id = ? AND ipfs_hash = ?", owerId, cid).First(&sharedFolder).Error; err != nil {
		return sharedFolder, err
	}
	return sharedFolder, nil
}

func GetSharedFolderInfoFromIPFS(ownerAccount, cid string) (sharedFolderInfoResponse, error) {
	var sharedFolderInfo sharedFolderInfoResponse

	folderData, err := GetFolderJsonFromIPFS(cid)
	if err != nil {
		return sharedFolderInfo, err
	}

	sharedFolderInfo = sharedFolderInfoResponse{
		OwnerAccount: ownerAccount,
		FolderInfo:   folderData,
	}
	return sharedFolderInfo, nil
}

// 从 IPFS 远程通过 url 下载文件，同时采用本地缓存的机制减少重复下载
func DownloadSharedFile(c *gin.Context, fileName, ipfsHash string) error {
	url := GenerateIPFSUrl(ipfsHash)

	cacheFilePath := path.Join(config.CACHE_RES_PATH, ipfsHash, fileName)

	// 确保目录存在
	cacheDir := path.Dir(cacheFilePath)
	if err := os.MkdirAll(cacheDir, os.ModePerm); err != nil {
		return err
	}

	// 检查本地缓存文件是否存在
	if _, err := os.Stat(cacheFilePath); os.IsNotExist(err) {
		// 如果文件不存在，从 IPFS 下载文件
		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		// 创建本地缓存文件
		file, err := os.Create(cacheFilePath)
		if err != nil {
			return err
		}
		defer file.Close()

		// 将远程文件内容写入本地缓存文件
		_, err = io.Copy(file, resp.Body)
		if err != nil {
			return err
		}
	}

	// 返回本地缓存文件给客户端
	c.File(cacheFilePath)
	return nil
}
