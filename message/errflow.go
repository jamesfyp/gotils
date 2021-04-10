package message

import (
	"time"
)

const (
	Common     = 0
	Waring     = 1
	Emergency  = 2
	MQueueName = "MESSAGE_QUEUE"
)

type WeMsg struct {
	Title   string    `json:"title"`
	Content string    `json:"content"`
	Time    time.Time `json:"time"`
	Level   int       `json:"level"`
}

func New(title string, content string, level int) (wm *WeMsg) {
	return &WeMsg{
		Title:   title,
		Content: content,
		Time:    time.Now(),
		Level:   level,
	}
}
