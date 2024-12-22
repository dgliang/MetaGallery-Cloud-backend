package util

import "regexp"

// TrimPathPrefix 去除文件夹路径前缀，即 "/{userId}"
func TrimPathPrefix(path string) string {
	re := regexp.MustCompile(`^/\d+`)
	res := re.ReplaceAllString(path, "")
	return res
}
