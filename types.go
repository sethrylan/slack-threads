package main

type Cursor struct {
	NextCursor string `json:"next_cursor"`
}

type CursorResponseMetadata struct {
	ResponseMetadata Cursor `json:"response_metadata"`
}

type Attachment struct {
	ID   int
	Text string
}

type Message struct {
	User        string
	BotID       string `json:"bot_id"`
	Text        string
	Attachments []Attachment
	Ts          string
	Type        string
	ReplyCount  int `json:"reply_count"`
}

type HistoryResponse struct {
	CursorResponseMetadata
	Ok       bool
	HasMore  bool `json:"has_more"`
	Messages []Message
}
