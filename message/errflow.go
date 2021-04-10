package message

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	Common    = 0
	Waring    = 1
	Emergency = 2
	MQueue    = "MESSAGE_QUEUE"

	wePushUrl  = "https://qyapi.weixin.qq.com/cgi-bin/webhook/send?key=5fa7e1e9-a337-4474-b706-38ae73648fb9"
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

// -------------------------------------

type PushData struct {
	MsgType  string   `json:"msgtype"`
	Markdown Markdown `json:"markdown"`
}

type Markdown struct {
	Content string `json:"content"`
}

type PushRes struct {
	ErrCode int    `json:"errcode"`
	ErrMsg  string `json:"errmsg"`
}

func New(title string, content string, level int) (wm *WeMsg) {
	return &WeMsg{
		Title:   title,
		Content: content,
		Time:    time.Now(),
		Level:   level,
	}
}

func In(cli *redis.Client, ctx context.Context, wm *WeMsg) {
	wmStr, _ := json.Marshal(wm)
	switch wm.Level {
	case Emergency:
		cli.RPush(ctx, MQueue, wmStr)
	default:
		cli.LPush(ctx, MQueue, wmStr)
	}
}

func Out(cli *redis.Client, ctx context.Context) (wm WeMsg, err error) {

	wmStr, err := cli.RPop(ctx, MQueue).Result()
	if err != nil {
		return wm, err
	}
	err = json.Unmarshal([]byte(wmStr), &wm)
	if err != nil {
		return wm, err
	}
	return
}

func Push(wm WeMsg, l int64) error {
	var (
		status string
		color  string
	)
	switch wm.Level {
	case Common:
		color, status = "comment", "通知  "
		break
	case Waring:
		color, status = "info", "提醒  "
		break
	case Emergency:
		color, status = "warning", "需处理"
		break
	}
	var pushData PushData
	pushData.MsgType = "markdown"
	pushData.Markdown.Content = fmt.Sprintf(wePushTemp, color, status, wm.Title, wm.Content, wm.Time.Format("2006-01-02 15:04:05"), l)

	msgB, _ := json.Marshal(&pushData)

	cli := &http.Client{}
	req, _ := http.NewRequest("POST", wePushUrl, bytes.NewBuffer(msgB))
	req.Header.Set("Content-Type", "application/json")
	res, err := cli.Do(req)
	if err != nil {
		return err
	}
	body, _ := ioutil.ReadAll(res.Body)

	var pushRes PushRes
	json.Unmarshal(body, &pushRes)

	if pushRes.ErrCode != 0 || pushRes.ErrMsg != "ok" {
		return errors.New(fmt.Sprintf("push failed msg: %s, code: %d", pushRes.ErrMsg, pushRes.ErrCode))
	} else {
		return nil
	}
}
