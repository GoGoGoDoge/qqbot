package qqbot

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func (a *API) GetGatewayURL() (string, error) {
	token, err := a.GetAccessToken()
	if err != nil {
		return "", err
	}
	req, _ := http.NewRequest("GET", APIBase+"/gateway", nil)
	req.Header.Set("Authorization", "QQBot "+token)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result struct{ URL string `json:"url"` }
	json.NewDecoder(resp.Body).Decode(&result)
	return result.URL, nil
}

func (a *API) SendC2CMessage(openID, content, msgID string) (*MessageResponse, error) {
	token, _ := a.GetAccessToken()
	var body map[string]interface{}
	if a.MarkdownMode {
		body = map[string]interface{}{"markdown": map[string]string{"content": content}, "msg_type": 2}
	} else {
		body = map[string]interface{}{"content": content, "msg_type": 0}
	}
	if msgID != "" {
		body["msg_id"] = msgID
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v2/users/%s/messages", APIBase, openID), bytes.NewReader(data))
	req.Header.Set("Authorization", "QQBot "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result MessageResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return &result, nil
}

func (a *API) SendGroupMessage(groupOpenID, content, msgID string) (*MessageResponse, error) {
	token, _ := a.GetAccessToken()
	var body map[string]interface{}
	if a.MarkdownMode {
		body = map[string]interface{}{"markdown": map[string]string{"content": content}, "msg_type": 2}
	} else {
		body = map[string]interface{}{"content": content, "msg_type": 0}
	}
	if msgID != "" {
		body["msg_id"] = msgID
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v2/groups/%s/messages", APIBase, groupOpenID), bytes.NewReader(data))
	req.Header.Set("Authorization", "QQBot "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result MessageResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return &result, nil
}

func (a *API) SendChannelMessage(channelID, content, msgID string) (*MessageResponse, error) {
	token, _ := a.GetAccessToken()
	body := map[string]interface{}{"content": content, "msg_type": 0}
	if msgID != "" {
		body["msg_id"] = msgID
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/channels/%s/messages", APIBase, channelID), bytes.NewReader(data))
	req.Header.Set("Authorization", "QQBot "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	var result MessageResponse
	json.NewDecoder(resp.Body).Decode(&result)
	return &result, nil
}

func (a *API) UploadC2CFile(openID, filePath string) (string, error) {
	token, _ := a.GetAccessToken()
	data, _ := os.ReadFile(filePath)
	body := map[string]interface{}{
		"file_type": 1,
		"url":       "https://example.com/" + filePath,
		"srv_send_msg": false,
	}
	jsonData, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v2/users/%s/files", APIBase, openID), bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "QQBot "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	var result FileUploadResponse
	json.NewDecoder(resp.Body).Decode(&result)
	_ = data
	return result.FileInfo, nil
}

func (a *API) SendC2CImage(openID, imageURL, msgID string) error {
	token, _ := a.GetAccessToken()
	body := map[string]interface{}{
		"content": " ",
		"msg_type": 7,
		"media": map[string]string{"file_info": imageURL},
	}
	if msgID != "" {
		body["msg_id"] = msgID
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v2/users/%s/messages", APIBase, openID), bytes.NewReader(data))
	req.Header.Set("Authorization", "QQBot "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (a *API) SendGroupImage(groupOpenID, imageURL, msgID string) error {
	token, _ := a.GetAccessToken()
	body := map[string]interface{}{
		"content": " ",
		"msg_type": 7,
		"media": map[string]string{"file_info": imageURL},
	}
	if msgID != "" {
		body["msg_id"] = msgID
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v2/groups/%s/messages", APIBase, groupOpenID), bytes.NewReader(data))
	req.Header.Set("Authorization", "QQBot "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (a *API) SendC2CFile(openID, fileInfo, msgID string) error {
	token, _ := a.GetAccessToken()
	body := map[string]interface{}{
		"content": " ",
		"msg_type": 7,
		"media": map[string]string{"file_info": fileInfo},
	}
	if msgID != "" {
		body["msg_id"] = msgID
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v2/users/%s/messages", APIBase, openID), bytes.NewReader(data))
	req.Header.Set("Authorization", "QQBot "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}

func (a *API) UploadMedia(url string, data []byte) (string, error) {
	encoded := base64.StdEncoding.EncodeToString(data)
	token, _ := a.GetAccessToken()
	body := map[string]interface{}{"file_type": 1, "url": url, "srv_send_msg": false, "file_data": encoded}
	jsonData, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", url, bytes.NewReader(jsonData))
	req.Header.Set("Authorization", "QQBot "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	respData, _ := io.ReadAll(resp.Body)
	var result FileUploadResponse
	json.Unmarshal(respData, &result)
	return result.FileInfo, nil
}

func (a *API) SendC2CInputNotify(openID, msgID string, seconds int) error {
	token, _ := a.GetAccessToken()
	body := map[string]interface{}{
		"msg_type": 6,
		"input_notify": map[string]int{"input_type": 1, "input_second": seconds},
	}
	if msgID != "" {
		body["msg_id"] = msgID
	}
	data, _ := json.Marshal(body)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/v2/users/%s/messages", APIBase, openID), bytes.NewReader(data))
	req.Header.Set("Authorization", "QQBot "+token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	resp.Body.Close()
	return nil
}
