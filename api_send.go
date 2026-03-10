package qqbot

import (
	"fmt"

	"github.com/tencent-connect/botgo/dto"
)

func (a *API) SendC2CMessage(openID, content, msgID string) (*MessageResponse, error) {
	msgToCreate := &dto.MessageToCreate{
		Content: content,
		MsgID:   msgID,
	}

	if a.MarkdownMode {
		msgToCreate.Markdown = &dto.Markdown{
			Content: content,
		}
		msgToCreate.MsgType = 2
	}

	msg, err := a.client.PostC2CMessage(a.ctx, openID, msgToCreate)
	if err != nil {
		return nil, err
	}

	return &MessageResponse{
		ID:        msg.ID,
		Timestamp: string(msg.Timestamp),
	}, nil
}

func (a *API) SendGroupMessage(groupOpenID, content, msgID string) (*MessageResponse, error) {
	msgToCreate := &dto.MessageToCreate{
		Content: content,
		MsgID:   msgID,
	}

	if a.MarkdownMode {
		msgToCreate.Markdown = &dto.Markdown{
			Content: content,
		}
		msgToCreate.MsgType = 2
	}

	msg, err := a.client.PostGroupMessage(a.ctx, groupOpenID, msgToCreate)
	if err != nil {
		return nil, err
	}

	return &MessageResponse{
		ID:        msg.ID,
		Timestamp: string(msg.Timestamp),
	}, nil
}

func (a *API) SendChannelMessage(channelID, content, msgID string) (*MessageResponse, error) {
	msgToCreate := &dto.MessageToCreate{
		Content: content,
		MsgID:   msgID,
	}

	msg, err := a.client.PostMessage(a.ctx, channelID, msgToCreate)
	if err != nil {
		return nil, err
	}

	return &MessageResponse{
		ID:        msg.ID,
		Timestamp: string(msg.Timestamp),
	}, nil
}

func (a *API) SendC2CImage(openID, imageURL, msgID string) error {
	msgToCreate := &dto.MessageToCreate{
		Content: " ",
		MsgType: 7,
		Media: &dto.MediaInfo{
			FileInfo: []byte(imageURL),
		},
		MsgID: msgID,
	}

	_, err := a.client.PostC2CMessage(a.ctx, openID, msgToCreate)
	return err
}

func (a *API) SendGroupImage(groupOpenID, imageURL, msgID string) error {
	msgToCreate := &dto.MessageToCreate{
		Content: " ",
		MsgType: 7,
		Media: &dto.MediaInfo{
			FileInfo: []byte(imageURL),
		},
		MsgID: msgID,
	}

	_, err := a.client.PostGroupMessage(a.ctx, groupOpenID, msgToCreate)
	return err
}

func (a *API) SendC2CFile(openID, fileInfo, msgID string) error {
	msgToCreate := &dto.MessageToCreate{
		Content: " ",
		MsgType: 7,
		Media: &dto.MediaInfo{
			FileInfo: []byte(fileInfo),
		},
		MsgID: msgID,
	}

	_, err := a.client.PostC2CMessage(a.ctx, openID, msgToCreate)
	return err
}

func (a *API) UploadC2CFile(openID, filePath string) (string, error) {
	// botgo 库暂不直接支持文件上传，需要使用 Rich Media API
	return "", fmt.Errorf("file upload not implemented yet")
}

func (a *API) UploadMedia(url string, data []byte) (string, error) {
	return "", fmt.Errorf("media upload not implemented yet")
}

func (a *API) SendC2CInputNotify(openID, msgID string, seconds int) error {
	// 输入状态通知暂不支持
	return nil
}
