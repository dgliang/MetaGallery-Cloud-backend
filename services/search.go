package services

import (
	"MetaGallery-Cloud-backend/models"
	"strings"
	"time"
)

type searchResponse struct {
	Result interface{} `json:"result"`
}

type searchNormalesponse struct {
	Type       string `json:"type"`
	ID         uint   `json:"id"`
	UserID     uint   `json:"user_id"`
	Name       string `json:"name"`
	Path       string `json:"path"`
	IsFavorite bool   `json:"is_favorite"`
	IsShared   bool   `json:"is_shared"`
}

type searchBinResponse struct {
	Type       string    `json:"type"`
	BinID      uint      `json:"bin_id"`
	ID         uint      `json:"id"`
	UserID     uint      `json:"user_id"`
	Name       string    `json:"name"`
	Path       string    `json:"path"`
	IsFavorite bool      `json:"is_favorite"`
	IsShared   bool      `json:"is_shared"`
	DelTime    time.Time `json:"del_time"`
}

func SearchFilesAndFolders(userId uint, rootFolderPath, keyword string) (searchResponse, error) {
	// 使用 UNION 进行查询
	var result []searchNormalesponse

	err := models.DataBase.Raw(`
		(SELECT 
			'FILE' AS type, 
			id, 
			belong_to AS user_id, 
			file_name AS name, 
			path, 
			favorite AS is_favorite, 
			share AS is_shared
		FROM file_data 
		WHERE file_name LIKE ? AND belong_to = ? AND path LIKE ?)
		UNION
		(SELECT 
			'FOLDER' AS type, 
			id, 
			belong_to AS user_id, 
			folder_name AS name, 
			path, 
			favorite AS is_favorite, 
			share AS is_shared
		FROM folder_data 
		WHERE folder_name LIKE ? AND belong_to = ? AND path LIKE ?)
	`, keyword+"%", userId, rootFolderPath+"/%", keyword+"%", userId, rootFolderPath+"/%").Scan(&result).Error
	if err != nil {
		return searchResponse{}, err
	}

	// 去除 path 路径的前缀
	for i := range result {
		result[i].Path = TrimPathPrefix(result[i].Path)
	}

	return searchResponse{
		Result: result,
	}, nil
}

func SearchBinFilesAndFolders(userId uint, keyword string) (searchResponse, error) {
	// 对 keyword 进行预处理
	keyword = keyword + "_bin_"

	// 查询符合条件的记录
	var result []searchBinResponse
	err := models.DataBase.Raw(`
	(
		SELECT
			'FILE' AS type,
			b.id AS bin_id,
			f.id,
			b.user_id AS user_id,
			f.file_name AS name,
			f.path,
			f.favorite AS is_favorite,
			f.share AS is_shared,
			b.deleted_time AS del_time
		FROM file_bins fb
		JOIN bins b ON fb.bin_id = b.id
		JOIN file_data f ON fb.file_id = f.id
		WHERE f.file_name LIKE ? AND b.user_id = ? AND f.deleted_at IS NOT NULL
	)
	UNION ALL
	(
		SELECT 
			'FOLDER' AS type,
			b.id AS bin_id,
			fd.id AS id,
			b.user_id AS user_id,
			fd.folder_name AS name,
			fd.path,
			fd.favorite AS is_favorite,
			fd.share AS is_shared,
			b.deleted_time AS del_time
		FROM folder_bins fb
		JOIN bins b ON fb.bin_id = b.id
		JOIN folder_data fd ON fb.folder_id = fd.id
		WHERE fd.folder_name LIKE ? AND b.user_id = ? AND fd.deleted_at IS NOT NULL
	)
	`, keyword+"%", userId, keyword+"%", userId).Scan(&result).Error
	if err != nil {
		return searchResponse{}, err
	}

	// 去除 path 路径的前缀，去除时间戳
	for i := range result {
		fullName, _ := SplitBinTimestamp(result[i].Name)
		result[i].Path = TrimPathPrefix(result[i].Path)
		result[i].Path = strings.ReplaceAll(result[i].Path, result[i].Name, fullName)
		result[i].Name = fullName
	}

	return searchResponse{
		Result: result,
	}, nil
}
