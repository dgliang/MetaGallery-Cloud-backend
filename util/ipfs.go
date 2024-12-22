package util

import "MetaGallery-Cloud-backend/config"

// GenerateIPFSUrl 根据 cid 生成文件在 IPFS 的 url
func GenerateIPFSUrl(cid string) string {
	return config.PINATA_HOST_URL + cid + config.PINATA_GATEWAY_KEY
}
