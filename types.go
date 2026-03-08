package qqbot

type Attachment struct {
	ContentType string `json:"content_type"`
	Filename    string `json:"filename,omitempty"`
	Height      int    `json:"height,omitempty"`
	Width       int    `json:"width,omitempty"`
	Size        int    `json:"size,omitempty"`
	URL         string `json:"url"`
}

type C2CMessageEvent struct {
	Author struct {
		ID         string `json:"id"`
		UserOpenID string `json:"user_openid"`
	} `json:"author"`
	Content     string       `json:"content"`
	ID          string       `json:"id"`
	Timestamp   string       `json:"timestamp"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type GroupMessageEvent struct {
	Author struct {
		ID           string `json:"id"`
		MemberOpenID string `json:"member_openid"`
	} `json:"author"`
	Content     string       `json:"content"`
	ID          string       `json:"id"`
	Timestamp   string       `json:"timestamp"`
	GroupOpenID string       `json:"group_openid"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type ChannelMessageEvent struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	GuildID   string `json:"guild_id"`
	Content   string `json:"content"`
	Timestamp string `json:"timestamp"`
	Author    struct {
		ID       string `json:"id"`
		Username string `json:"username,omitempty"`
		Bot      bool   `json:"bot,omitempty"`
	} `json:"author"`
	Attachments []Attachment `json:"attachments,omitempty"`
}

type WSPayload struct {
	Op int         `json:"op"`
	D  interface{} `json:"d,omitempty"`
	S  int         `json:"s,omitempty"`
	T  string      `json:"t,omitempty"`
}

type MessageResponse struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
}

type FileUploadResponse struct {
	FileInfo string `json:"file_info"`
}
