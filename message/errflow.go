package message

import (
	"time"
)

const (
	Common     = 0
	Waring     = 1
	Emergency  = 2
	MQueue     = "MESSAGE_QUEUE"
	wePushTemp = `#### <font color="%s">%s</font> **%s** 
###### 内容: %s
###### 发生时间: %s
###### 剩余错误: %d`
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
