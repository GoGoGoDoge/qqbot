package main

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"regexp"
	"strings"
	"time"

	qqbot "github.com/GoGoGoDoge/qqbot"
	"github.com/tencent-connect/botgo/dto"
	"golang.org/x/text/unicode/norm"
)

func main() {
	client := qqbot.NewClient(
		os.Getenv("QQBOT_APP_ID"),
		os.Getenv("QQBOT_APP_SECRET"),
	)

	// 处理 C2C 消息
	client.Gateway.OnC2CMessage(func(e qqbot.C2CMessageEvent) {
		log.Printf("收到C2C消息: %s (from: %s) with msg id(%s)", e.Content, e.Author.ID, e.ID)

		// 检查是否有附件（图片、视频、文件等）
		if len(e.Attachments) > 0 {
			handleAttachments(client, e.Author.ID, e.ID, e.Attachments)
			return
		}

		switch e.Content {
		case "图片":
			// 方式1: 先上传图片，再发送
			sendImageByUpload(client, e.Author.ID, e.ID)

		case "图片2":
			// 方式2: 直接发送图片URL（会占用主动消息频率）
			sendImageDirect(client, e.Author.ID, e.ID)

		case "语音":
			// 发送语音（仅支持 silk 格式）
			sendVoice(client, e.Author.ID, e.ID)

		case "视频":
			// 发送视频
			sendVideo(client, e.Author.ID, e.ID)

		case "文件":
			// 上传本地文件（FileType=4，使用 base64）
			// 修改为你的本地文件路径
			sendLocalFile(client, e.Author.ID, e.ID, "/Users/marco/Downloads/web4_report_updated.pdf")

		case "本地图片":
			// 上传本地图片（FileType=1，使用 base64）
			sendLocalImage(client, e.Author.ID, e.ID, "/Users/marco/Downloads/test.jpg")

		case "本地视频":
			// 上传本地视频（FileType=2，使用 base64）
			sendLocalVideo(client, e.Author.ID, e.ID, "/Users/marco/Downloads/test.mp4")

		default:
			// 普通文本回复
			client.API.SendC2CMessage(e.Author.ID, "Echo: "+e.Content, e.ID)
		}
	})

	// 处理群消息
	client.Gateway.OnGroupMessage(func(e qqbot.GroupMessageEvent) {
		log.Printf("收到群消息: %s (group: %s)", e.Content, e.GroupOpenID)

		// 检查是否有附件
		if len(e.Attachments) > 0 {
			handleAttachments(client, e.GroupOpenID, e.ID, e.Attachments)
			return
		}

		if e.Content == "图片" {
			sendGroupImage(client, e.GroupOpenID, e.ID)
		}
	})

	if err := client.Connect(); err != nil {
		log.Fatalf("连接失败: %v", err)
	}

	log.Println("Bot started - 发送 '图片', '图片2', '语音', '视频', '文件', '本地图片', '本地视频' 测试富媒体消息")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
	client.Gateway.Close()
}

// sendImageByUpload 先上传图片，再发送（推荐方式，不占用主动消息频率）
func sendImageByUpload(client *qqbot.Client, userID, msgID string) {
	// 步骤1: 上传图片文件
	richMediaMsg := &dto.RichMediaMessage{
		EventID:    msgID,
		FileType:   1,                                                      // 1=图片
		URL:        "https://img.cdn1.vip/i/69aed1f68f255_1773064694.webp", // 图片URL
		SrvSendMsg: false,                                                  // false=仅上传，不直接发送
	}

	msg, err := client.API.GetClient().PostC2CMessage(client.API.GetContext(), userID, richMediaMsg)
	if err != nil {
		log.Printf("❌ 上传图片失败: %v", err)
		return
	}

	log.Printf("✅ 图片上传成功，file_info: %s", string(msg.FileInfo))

	// 步骤2: 使用 file_info 发送消息（消息类型7）
	msgToCreate := &dto.MessageToCreate{
		Content: " ", // 富媒体消息 content 必须有值
		MsgType: 7,   // 7=富媒体消息
		MsgID:   msgID,
		Media: &dto.MediaInfo{
			FileInfo: msg.FileInfo, // 使用上传返回的 file_info
		},
	}

	_, err = client.API.GetClient().PostC2CMessage(client.API.GetContext(), userID, msgToCreate)
	if err != nil {
		log.Printf("❌ 发送图片消息失败: %v", err)
		return
	}

	log.Printf("✅ 图片消息发送成功")
}

// sendImageDirect 直接发送图片（会占用主动消息频率）
func sendImageDirect(client *qqbot.Client, userID, msgID string) {
	richMediaMsg := &dto.RichMediaMessage{
		EventID:    msgID,
		FileType:   1, // 1=图片
		URL:        "https://tse2-mm.cn.bing.net/th/id/OIP-C.afBDBOGj8IrGkUkZKXWvygHaEK?w=310&h=180&c=7&r=0&o=7&pid=1.7&rm=3",
		SrvSendMsg: true, // true=直接发送到用户
	}

	_, err := client.API.GetClient().PostC2CMessage(client.API.GetContext(), userID, richMediaMsg)
	if err != nil {
		log.Printf("❌ 发送图片失败: %v", err)
		return
	}

	log.Printf("✅ 图片发送成功")
}

// sendVoice 发送语音（仅支持 silk 格式）
func sendVoice(client *qqbot.Client, userID, msgID string) {
	// 步骤1: 上传语音文件
	richMediaMsg := &dto.RichMediaMessage{
		EventID:    msgID,
		FileType:   3, // 3=语音（仅支持silk格式）
		URL:        "https://example.com/voice.silk",
		SrvSendMsg: false,
	}

	msg, err := client.API.GetClient().PostC2CMessage(client.API.GetContext(), userID, richMediaMsg)
	if err != nil {
		log.Printf("❌ 上传语音失败: %v", err)
		return
	}

	// 步骤2: 发送语音消息
	msgToCreate := &dto.MessageToCreate{
		Content: " ",
		MsgType: 7,
		MsgID:   msgID,
		Media: &dto.MediaInfo{
			FileInfo: msg.FileInfo,
		},
	}

	_, err = client.API.GetClient().PostC2CMessage(client.API.GetContext(), userID, msgToCreate)
	if err != nil {
		log.Printf("❌ 发送语音消息失败: %v", err)
		return
	}

	log.Printf("✅ 语音消息发送成功")
}

// sendVideo 发送视频
func sendVideo(client *qqbot.Client, userID, msgID string) {
	// 步骤1: 上传视频文件
	richMediaMsg := &dto.RichMediaMessage{
		EventID:    msgID,
		FileType:   2, // 2=视频
		URL:        "https://download.samplelib.com/mp4/sample-5s.mp4",
		SrvSendMsg: false,
	}

	msg, err := client.API.GetClient().PostC2CMessage(client.API.GetContext(), userID, richMediaMsg)
	if err != nil {
		log.Printf("❌ 上传视频失败: %v", err)
		return
	}

	// 步骤2: 发送视频消息
	msgToCreate := &dto.MessageToCreate{
		Content: " ",
		MsgType: 7,
		MsgID:   msgID,
		Media: &dto.MediaInfo{
			FileInfo: msg.FileInfo,
		},
	}

	_, err = client.API.GetClient().PostC2CMessage(client.API.GetContext(), userID, msgToCreate)
	if err != nil {
		log.Printf("❌ 发送视频消息失败: %v", err)
		return
	}

	log.Printf("✅ 视频消息发送成功")
}

// sendGroupImage 发送群图片
func sendGroupImage(client *qqbot.Client, groupID, msgID string) {
	// 步骤1: 上传图片
	richMediaMsg := &dto.RichMediaMessage{
		EventID:    msgID,
		FileType:   1,
		URL:        "https://example.com/image.jpg",
		SrvSendMsg: false,
	}

	msg, err := client.API.GetClient().PostGroupMessage(client.API.GetContext(), groupID, richMediaMsg)
	if err != nil {
		log.Printf("❌ 上传群图片失败: %v", err)
		return
	}

	// 步骤2: 发送图片消息
	msgToCreate := &dto.MessageToCreate{
		Content: " ",
		MsgType: 7,
		MsgID:   msgID,
		Media: &dto.MediaInfo{
			FileInfo: msg.FileInfo,
		},
	}

	_, err = client.API.GetClient().PostGroupMessage(client.API.GetContext(), groupID, msgToCreate)
	if err != nil {
		log.Printf("❌ 发送群图片消息失败: %v", err)
		return
	}

	log.Printf("✅ 群图片消息发送成功")
}

// handleAttachments 处理接收到的附件（图片、视频、文件等）
func handleAttachments(client *qqbot.Client, targetID, msgID string, attachments []qqbot.Attachment) {
	for _, att := range attachments {
		log.Printf("收到附件: 文件名=%s, 类型=%s, 大小=%d bytes, URL=%s",
			att.Filename, att.ContentType, att.Size, att.URL)

		// 构建回复消息
		var reply string
		if att.Filename != "" {
			reply = fmt.Sprintf("收到文件:\n文件名: %s\n大小: %s\n类型: %s",
				att.Filename,
				formatFileSize(att.Size),
				att.ContentType)
		} else {
			// 如果没有文件名，可能是图片
			if att.Width > 0 && att.Height > 0 {
				reply = fmt.Sprintf("收到图片:\n尺寸: %dx%d\n大小: %s\n类型: %s",
					att.Width, att.Height,
					formatFileSize(att.Size),
					att.ContentType)
			} else {
				reply = fmt.Sprintf("收到媒体文件:\n大小: %s\n类型: %s",
					formatFileSize(att.Size),
					att.ContentType)
			}
		}

		// 发送回复
		_, err := client.API.SendC2CMessage(targetID, reply, msgID)
		if err != nil {
			log.Printf("❌ 发送附件信息失败: %v", err)
		} else {
			log.Printf("✅ 已回复附件信息")
		}
	}
}

// formatFileSize 格式化文件大小
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

// sendLocalFile 上传本地文件（FileType=4，使用 file_data 字段，支持 20MB 以下）
func sendLocalFile(client *qqbot.Client, userID, msgID, filePath string) {
	sendLocalMedia(client, userID, msgID, filePath, 4, "文件")
}

// sendLocalImage 上传本地图片（FileType=1，使用 file_data 字段，支持 20MB 以下）
func sendLocalImage(client *qqbot.Client, userID, msgID, filePath string) {
	sendLocalMedia(client, userID, msgID, filePath, 1, "图片")
}

// sendLocalVideo 上传本地视频（FileType=2，使用 file_data 字段，支持 20MB 以下）
func sendLocalVideo(client *qqbot.Client, userID, msgID, filePath string) {
	sendLocalMedia(client, userID, msgID, filePath, 2, "视频")
}

// sendLocalMedia 通用的本地媒体文件上传函数
func sendLocalMedia(client *qqbot.Client, userID, msgID, filePath string, fileType int, mediaTypeName string) {
	// 检查文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		log.Printf("❌ 文件不存在: %v", err)
		return
	}

	// 检查文件大小（base64 编码后会增加 33%，所以限制为 5MB）
	const maxSize = 5 * 1024 * 1024 // 5MB
	if fileInfo.Size() > maxSize {
		log.Printf("❌ 文件太大: %s (file_data 方式最大支持 5MB)", formatFileSize(int(fileInfo.Size())))
		client.API.SendC2CMessage(userID, fmt.Sprintf("%s太大: %s，file_data 方式最大支持 5MB（base64 编码后会增加 33%%）", mediaTypeName, formatFileSize(int(fileInfo.Size()))), msgID)
		return
	}

	log.Printf("准备上传%s: %s, 大小: %s", mediaTypeName, fileInfo.Name(), formatFileSize(int(fileInfo.Size())))

	// 读取文件内容
	fileData, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("❌ 读取文件失败: %v", err)
		return
	}

	// Base64 编码
	base64Data := base64.StdEncoding.EncodeToString(fileData)
	log.Printf("%s base64 编码完成，长度: %d", mediaTypeName, len(base64Data))

	// 步骤1: 上传文件（使用原始 HTTP API，因为 botgo 不支持 file_data）
	var fileName string
	if fileType == 4 { // 只有文件类型需要文件名
		fileName = sanitizeFileName(fileInfo.Name())
	}
	uploadResp, err := uploadMediaWithFileData(client, userID, base64Data, fileType, fileName)
	if err != nil {
		log.Printf("❌ 上传%s失败: %v", mediaTypeName, err)
		client.API.SendC2CMessage(userID, fmt.Sprintf("%s上传失败: %v", mediaTypeName, err), msgID)
		return
	}

	log.Printf("✅ %s上传成功，file_info 长度: %d", mediaTypeName, len(uploadResp.FileInfo))

	// 步骤2: 使用 file_info 发送消息（消息类型7）
	msgSeq := getMsgSeq(msgID)
	err = sendMediaMessage(client, userID, uploadResp.FileInfo, msgID, msgSeq)
	if err != nil {
		log.Printf("❌ 发送%s消息失败: %v", mediaTypeName, err)
		client.API.SendC2CMessage(userID, fmt.Sprintf("%s发送失败: %v", mediaTypeName, err), msgID)
		return
	}

	log.Printf("✅ %s消息发送成功: %s (%s)", mediaTypeName, fileInfo.Name(), formatFileSize(int(fileInfo.Size())))
}

// uploadMediaWithFileData 使用原始 HTTP API 上传媒体文件（file_data 字段）
func uploadMediaWithFileData(client *qqbot.Client, userID, base64Data string, fileType int, fileName string) (*FileUploadResponse, error) {
	// 获取 token
	token, err := client.API.GetTokenSource().Token()
	if err != nil {
		return nil, fmt.Errorf("获取 token 失败: %w", err)
	}

	// 构建请求体
	requestBody := map[string]interface{}{
		"file_type":    fileType,   // 1=图片, 2=视频, 3=语音, 4=文件
		"file_data":    base64Data, // base64 编码的文件数据
		"srv_send_msg": false,      // false=仅上传
	}

	// 添加文件名（仅文件类型需要）
	if fileType == 4 && fileName != "" {
		requestBody["file_name"] = fileName
	}

	bodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	// 发送 HTTP 请求
	url := fmt.Sprintf("https://api.sgroup.qq.com/v2/users/%s/files", userID)
	req, err := http.NewRequest("POST", url, bytes.NewReader(bodyBytes))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	req.Header.Set("Authorization", "QQBot "+token.AccessToken)
	req.Header.Set("Content-Type", "application/json")

	log.Printf("发送媒体上传请求: %s (fileType=%d)", url, fileType)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	// 解析响应
	var uploadResp FileUploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&uploadResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("上传失败: status=%d, response=%+v", resp.StatusCode, uploadResp)
	}

	return &uploadResp, nil
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	FileUUID string `json:"file_uuid"`
	FileInfo string `json:"file_info"`
	TTL      int    `json:"ttl"`
}

// sendMediaMessage 发送富媒体消息（直接使用 HTTP API，避免 file_info 被二次 base64 编码）
func sendMediaMessage(client *qqbot.Client, userID, fileInfo, msgID string, msgSeq int) error {
	// 获取 token
	token, err := client.API.GetTokenSource().Token()
	if err != nil {
		return fmt.Errorf("获取 token 失败: %w", err)
	}

	// 构建请求体（file_info 必须是字符串，不能是 base64 编码的 []byte）
	requestBody := map[string]interface{}{
		"msg_type": 7,
		"media": map[string]string{
			"file_info": fileInfo, // 直接使用字符串
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

	// 发送 HTTP 请求
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

	// 解析响应
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("解析响应失败: %w", err)
	}

	if resp.StatusCode != 200 {
		return fmt.Errorf("发送失败: status=%d, response=%+v", resp.StatusCode, result)
	}

	return nil
}

// getMsgSeq 获取消息序号（范围 0 ~ 65535）
func getMsgSeq(msgID string) int {
	timePart := time.Now().UnixMilli() % 100000000
	random := rand.Intn(65536)
	return int((timePart ^ int64(random)) % 65536)
}

// sanitizeFileName 规范化文件名为 QQ Bot API 要求的 UTF-8 编码格式
// 处理:
// 1. Unicode NFC 规范化（将 macOS NFD 分解形式合并为 NFC 组合形式）
// 2. 去除 ASCII 控制字符（0x00-0x1F, 0x7F）
// 3. 去除首尾空白
// 4. 对 percent-encoded 的文件名尝试 URI 解码
func sanitizeFileName(name string) string {
	if name == "" {
		return name
	}

	result := strings.TrimSpace(name)

	// 尝试 URI 解码（处理 URL 中 percent-encoded 的中文文件名）
	// 例如 %E4%B8%AD%E6%96%87.txt → 中文.txt
	if strings.Contains(result, "%") {
		if decoded, err := url.QueryUnescape(result); err == nil {
			result = decoded
		}
		// 解码失败（非合法 percent-encoding），保留原始值
	}

	// Unicode NFC 规范化：将 macOS NFD 分解形式合并为标准 NFC 组合形式
	result = norm.NFC.String(result)

	// 去除 ASCII 控制字符（保留所有可打印字符和非 ASCII Unicode 字符）
	controlChars := regexp.MustCompile(`[\x00-\x1F\x7F]`)
	result = controlChars.ReplaceAllString(result, "")

	return result
}
