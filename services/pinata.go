package services

import (
	"MetaGallery-Cloud-backend/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
)

func UploadFileToPinata(filePath string) (string, error) {
	url := "https://api.pinata.cloud/pinning/pinFileToIPFS"

	// 对 filepath 进行预处理
	workingDir, _ := os.Getwd()
	workingDir = strings.ReplaceAll(workingDir, "\\", "/")
	log.Printf("Working directory: %s", workingDir)
	filePath = strings.ReplaceAll(filePath, "\\", "/")

	abosultePath := path.Join(workingDir, config.FileResPath, filePath)
	log.Printf("Absolute path: %s", abosultePath)

	// 创建一个 buffer 和 multipart writer
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加文件字段
	file, err := os.Open(abosultePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	part, err := writer.CreateFormFile("file", abosultePath)
	if err != nil {
		return "", err
	}
	io.Copy(part, file)

	// // 添加元数据字段
	// writer.WriteField("pinataMetadata", `{"name": "Pinnie.json"}`)

	// // 添加选项字段
	// writer.WriteField("pinataOptions", `{"cidVersion": 1}`)

	// 关闭 writer，完成请求体构建
	writer.Close()

	// 创建请求
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		panic(err)
	}
	req.Header.Add("Authorization", "Bearer "+config.PinataJWT)
	req.Header.Add("Content-Type", writer.FormDataContentType())

	// 发送请求
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	// 读取响应
	responseBody, err := io.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	fmt.Println("Status:", res.Status)
	fmt.Println("Response:", string(responseBody))

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("文件上传 IPFS 失败，unexpected status code: %d", res.StatusCode)
	}

	// 解析 JSON 响应
	var response map[string]interface{}
	err = json.Unmarshal(responseBody, &response)
	if err != nil {
		return "", err
	}

	if ipfsHash, ok := response["IpfsHash"]; ok {
		return ipfsHash.(string), nil
	} else {
		return "", fmt.Errorf("无法从响应中获取 IPFS 哈希")
	}
}
