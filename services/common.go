package services

import (
	"MetaGallery-Cloud-backend/config"
	"regexp"
	"time"
)

func GenerateBinTimestamp(str string, t time.Time) string {
	timestamp := t.Format("20060102_150405")
	return str + "_bin_" + timestamp
}

func SplitBinTimestamp(str string) (string, string) {
	re := regexp.MustCompile(`^(.*)_bin_\d{8}_\d{6}$`)
	return re.FindStringSubmatch(str)[1], re.FindStringSubmatch(str)[0]
}

func GenerateIPFSUrl(cid string) string {
	return config.PINATA_HOST_URL + cid + config.PINATA_GATEWAY_KEY
}

func TrimPathPrefix(path string) string {
	re := regexp.MustCompile(`^/\d+`)
	res := re.ReplaceAllString(path, "")
	return res
}
