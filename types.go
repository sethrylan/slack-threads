package main

type Cursor struct {
	NextCursor string `json:"next_cursor"`
}

type CursorResponseMetadata struct {
	ResponseMetadata Cursor `json:"response_metadata"`
}

type Attachment struct {
	ID   int    `json:"id"`
	Text string `json:"text"`
}

type Message struct {
	User        string       `json:"user"`
	BotID       string       `json:"bot_id"`
	Text        string       `json:"text"`
	Attachments []Attachment `json:"attachments"`
	TS          string       `json:"ts"`
	Type        string       `json:"type"`
	ReplyCount  int          `json:"reply_count"`
}

type HistoryResponse struct {
	CursorResponseMetadata

	Ok       bool      `json:"ok"`
	HasMore  bool      `json:"has_more"`
	Messages []Message `json:"messages"`
}
