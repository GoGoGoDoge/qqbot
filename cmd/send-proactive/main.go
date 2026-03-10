package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"github.com/GoGoGoDoge/qqbot"
	"github.com/tencent-connect/botgo/dto"
	"golang.org/x/text/unicode/norm"
)

func main() {
	to := flag.String("to", "", "Target openid")
	text := flag.String("text", "", "Message text")
	typ := flag.String("type", "c2c", "Message type (c2c/group)")
	list := flag.Bool("list", false, "List known users")

	// 富媒体参数
	mediaURL := flag.String("media", "", "Media URL (image/video/voice)")
	mediaType := flag.String("media-type", "image", "Media type: image, video, voice, file")
	direct := flag.Bool("direct", false, "Send media directly (占用主动消息频率)")

	// 本地文件参数
	localFile := flag.String("file", "", "Local file path (image/video/file)")

	flag.Parse()

	appID := os.Getenv("QQBOT_APP_ID")
	secret := os.Getenv("QQBOT_APP_SECRET")
	if appID == "" || secret == "" {
		log.Fatal("QQBOT_APP_ID and QQBOT_APP_SECRET required")
	}

	client := qqbot.NewClient(appID, secret)
	store := qqbot.NewProactiveStore(".")

	if *list {
		users := store.List(*typ, "", 0)
		fmt.Printf("Known users (%d):\n", len(users))
		for _, u := range users {
			fmt.Printf("%s\t%s\t%s\n", u.Type, u.OpenID, u.Nickname)
		}
		return
	}

	if *to == "" {
		flag.Usage()
		os.Exit(1)
	}

	// 发送本地文件
	if *localFile != "" {
		if err := sendLocalFile(client, *to, *typ, *localFile, *mediaType); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Local file sent successfully")
		return
	}

	// 发送富媒体消息（URL）
	if *mediaURL != "" {
		if err := sendMedia(client, *to, *typ, *mediaURL, *mediaType, *direct); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Media sent successfully")
		return
	}

	// 发送文本消息
	if *text == "" {
		flag.Usage()
		os.Exit(1)
	}

	var err error
	if *typ == "c2c" {
		_, err = client.API.SendC2CMessage(*to, *text, "")
	} else {
		_, err = client.API.SendGroupMessage(*to, *text, "")
	}

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Message sent successfully")
}

func sendLocalFile(client *qqbot.Client, to, typ, filePath, mediaType string) error {
	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return fmt.Errorf("文件不存在: %w", err)
	}

	// 检查文件大小（base64 编码后会增加 33%，所以限制为 5MB）
	const maxSize = 5 * 1024 * 1024
	if fileInfo.Size() > maxSize {
		return fmt.Errorf("文件太大: %s (file_data 方式最大支持 5MB，base64 编码后会增加 33%%)\n提示：对于大文件，请先上传到公网 URL，然后使用 --media 参数", formatFileSize(int(fileInfo.Size())))
	}

	// 确定文件类型
	var fileType int
	switch mediaType {
	case "image":
		fileType = 1
	case "video":
		fileType = 2
	case "voice":
		fileType = 3
	case "file":
		fileType = 4
	default:
		return fmt.Errorf("unsupported media type: %s (use: image, video, voice, file)", mediaType)
	}

	fmt.Printf("准备上传%s: %s, 大小: %s\n", mediaType, fileInfo.Name(), formatFileSize(int(fileInfo.Size())))

	// 读取文件内容
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("读取文件失败: %w", err)
	}

	// Base64 编码
	base64Data := base64.StdEncoding.EncodeToString(fileData)
	fmt.Printf("%s base64 编码完成，长度: %d\n", mediaType, len(base64Data))

	// 步骤1: 上传文件
	var fileName string
	if fileType == 4 {
		fileName = sanitizeFileName(filepath.Base(filePath))
	}

	uploadResp, err := uploadMediaWithFileData(client, to, base64Data, fileType, fileName)
	if err != nil {
		return fmt.Errorf("上传%s失败: %w", mediaType, err)
	}

	fmt.Printf("✅ %s上传成功，file_info 长度: %d\n", mediaType, len(uploadResp.FileInfo))

	// 步骤2: 发送消息
	msgSeq := getMsgSeq("")
	err = sendMediaMessage(client, to, uploadResp.FileInfo, "", msgSeq)
	if err != nil {
		return fmt.Errorf("发送%s消息失败: %w", mediaType, err)
	}

	fmt.Printf("✅ %s消息发送成功\n", mediaType)
	return nil
}

func sendMedia(client *qqbot.Client, to, typ, mediaURL, mediaType string, direct bool) error {
	// 确定文件类型
	var fileType uint64
	switch mediaType {
	case "image":
		fileType = 1
	case "video":
		fileType = 2
	case "voice":
		fileType = 3
	case "file":
		fileType = 4
	default:
		return fmt.Errorf("unsupported media type: %s (use: image, video, voice, file)", mediaType)
	}

	// 如果是直接发送模式
	if direct {
		richMediaMsg := &dto.RichMediaMessage{
			FileType:   fileType,
			URL:        mediaURL,
			SrvSendMsg: true,
		}

		if typ == "c2c" {
			_, err := client.API.GetClient().PostC2CMessage(client.API.GetContext(), to, richMediaMsg)
			return err
		} else {
			_, err := client.API.GetClient().PostGroupMessage(client.API.GetContext(), to, richMediaMsg)
			return err
		}
	}

	// 两步法：先上传，再发送
	richMediaMsg := &dto.RichMediaMessage{
		FileType:   fileType,
		URL:        mediaURL,
		SrvSendMsg: false,
	}

	var msg *dto.Message
	var err error

	if typ == "c2c" {
		msg, err = client.API.GetClient().PostC2CMessage(client.API.GetContext(), to, richMediaMsg)
	} else {
		msg, err = client.API.GetClient().PostGroupMessage(client.API.GetContext(), to, richMediaMsg)
	}

	if err != nil {
		return fmt.Errorf("upload media failed: %w", err)
	}

	// 发送富媒体消息（使用直接 HTTP API 避免 file_info 被二次编码）
	msgSeq := getMsgSeq("")
	err = sendMediaMessage(client, to, string(msg.FileInfo), "", msgSeq)
	if err != nil {
		return fmt.Errorf("send media message failed: %w", err)
	}

	return nil
}

// uploadMediaWithFileData 使用原始 HTTP API 上传媒体文件
func uploadMediaWithFileData(client *qqbot.Client, userID, base64Data string, fileType int, fileName string) (*FileUploadResponse, error) {
	token, err := client.API.GetTokenSource().Token()
	if err != nil {
		return nil, fmt.Errorf("获取 token 失败: %w", err)
	}

	requestBody := map[string]interface{}{
		"file_type":    fileType,
		"file_data":    base64Data,
		"srv_send_msg": false,
	}

	if fileType == 4 && fileName != "" {
		requestBody["file_name"] = fileName
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	url := fmt.Sprintf("https://api.sgroup.qq.com/v2/users/%s/files", userID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "QQBot "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err2 := io.ReadAll(resp.Body)
	if err2 != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err2)
	}

	// 检查状态码
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("上传失败: status=%d, body=%s", resp.StatusCode, string(bodyBytes))
	}

	// 解析 JSON 响应
	var uploadResp FileUploadResponse
	if err := json.Unmarshal(bodyBytes, &uploadResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w, body=%s", err, string(bodyBytes))
	}

	return &uploadResp, nil
}

// sendMediaMessage 发送富媒体消息
func sendMediaMessage(client *qqbot.Client, userID, fileInfo, msgID string, msgSeq int) error {
	token, err := client.API.GetTokenSource().Token()
	if err != nil {
		return fmt.Errorf("获取 token 失败: %w", err)
	}

	requestBody := map[string]interface{}{
		"msg_type": 7,
		"media": map[string]string{
			"file_info": fileInfo,
		},
		"msg_seq": msgSeq,
	}

	if msgID != "" {
		requestBody["msg_id"] = msgID
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %w", err)
	}

	url := fmt.Sprintf("https://api.sgroup.qq.com/v2/users/%s/messages", userID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "QQBot "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("发送失败: status=%d, response=%+v", resp.StatusCode, result)
	}

	return nil
}

type FileUploadResponse struct {
	FileUUID string `json:"file_uuid"`
	FileInfo string `json:"file_info"`
	TTL      int    `json:"ttl"`
}

func getMsgSeq(msgID string) int {
	timePart := time.Now().UnixMilli() % 100000000
	random := rand.Intn(65536)
	return int((timePart ^ int64(random)) % 65536)
}

func sanitizeFileName(name string) string {
	if name == "" {
		return name
	}

	result := strings.TrimSpace(name)

	if strings.Contains(result, "%") {
		if decoded, err := url.QueryUnescape(result); err == nil {
			result = decoded
		}
	}

	result = norm.NFC.String(result)

	controlChars := regexp.MustCompile(`[\x00-\x1F\x7F]`)
	result = controlChars.ReplaceAllString(result, "")

	return result
}

func formatFileSize(size int) string {
	const (
		KB = 1024
		MB = 1024 * KB
		GB = 1024 * MB
	)

	switch {
	case size >= GB:
		return fmt.Sprintf("%.2f GB", float64(size)/float64(GB))
	case size >= MB:
		return fmt.Sprintf("%.2f MB", float64(size)/float64(MB))
	case size >= KB:
		return fmt.Sprintf("%.2f KB", float64(size)/float64(KB))
	default:
		return fmt.Sprintf("%d bytes", size)
	}
}
