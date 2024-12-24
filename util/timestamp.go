package util

import (
	"regexp"
	"time"
)

// GenerateTimestamp 生成移入回收站的文件夹的时间戳
func GenerateBinTimestamp(str string, t time.Time) string {
	timestamp := t.Format("20060102_150405")
	return str + "_bin_" + timestamp
}

// SplitBinTimestamp 从移入回收站的文件夹中提取时间戳
func SplitBinTimestamp(str string) (string, string) {
	re := regexp.MustCompile(`^(.*)(_bin_\d{8}_\d{6})$`)
	matches := re.FindStringSubmatch(str)
	if len(matches) == 3 {
		return matches[1], matches[2]
	}
	return str, ""
}
