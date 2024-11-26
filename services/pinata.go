package services

import (
	"MetaGallery-Cloud-backend/config"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"path"
	"strings"
)

func UploadFileToIPFS(filePath string) (string, error) {
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

func UploadJsonToIPFS(jsonData map[string]interface{}) (string, error) {
	url := "https://api.pinata.cloud/pinning/pinJSONToIPFS"

	payloadData := map[string]interface{}{
		// "pinataOptions": map[string]interface{}{
		// 	"cidVersion": 1,
		// },
		// "pinataMetadata": map[string]interface{}{
		// 	"name": "pinnie.json",
		// },
		"pinataContent": jsonData,
	}

	// 将 payload 转为 JSON
	payload, err := json.Marshal(payloadData)
	if err != nil {
		return "", fmt.Errorf("将 payload 转为 JSON: %v", err)
	}

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(payload))

	req.Header.Add("Authorization", "Bearer "+config.PinataJWT)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	log.Println(res)
	log.Println(string(body))

	// 解析响应体 response
	var response map[string]interface{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return "", fmt.Errorf("解析响应体 response: %v", err)
	}

	// 获取 IPFS 哈希
	if ipfsHash, ok := response["IpfsHash"]; ok {
		return ipfsHash.(string), nil
	} else {
		return "", fmt.Errorf("无法从响应中获取 IPFS 哈希")
	}
}

func CreatePinataGroup(groupName string) (string, error) {
	url := "https://api.pinata.cloud/groups"

	payload := strings.NewReader(fmt.Sprintf("{\n  \"name\": \"%s\"\n}", groupName))

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("Authorization", "Bearer "+config.PinataJWT)
	req.Header.Add("Content-Type", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	log.Println(res)
	log.Println(string(body))

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("创建 Pinata 群组失败，unexpected status code: %d", res.StatusCode)
	}

	// 解析 JSON 响应
	var response map[string]interface{}
	err := json.Unmarshal(body, &response)
	if err != nil {
		return "", err
	}

	if groupId, ok := response["id"]; ok {
		return groupId.(string), nil
	} else {
		return "", fmt.Errorf("无法从响应中获取群组 ID")
	}
}
