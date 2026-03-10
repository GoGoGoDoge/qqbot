package qqbot

type Author struct {
	ID         string
	UserOpenID string
	Username   string
}

type Attachment struct {
	ContentType string `json:"content_type"`
	Filename    string `json:"filename,omitempty"`
	Height      int    `json:"height,omitempty"`
	Width       int    `json:"width,omitempty"`
	Size        int    `json:"size,omitempty"`
	URL         string `json:"url"`
}

type C2CMessageEvent struct {
	ID          string
	Content     string
	Timestamp   string
	Author      Author
	Attachments []Attachment
}

type GroupMessageEvent struct {
	ID          string
	Content     string
	Timestamp   string
	GroupOpenID string
	Author      Author
	Attachments []Attachment
}

type ChannelMessageEvent struct {
	ID          string
	ChannelID   string
	GuildID     string
	Content     string
	Timestamp   string
	Author      Author
	Attachments []Attachment
}

type MessageResponse struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Code      int    `json:"code"`
	Message   string `json:"message"`
}

type FileUploadResponse struct {
	FileInfo string `json:"file_info"`
}
