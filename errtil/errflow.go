package errtil

import (
	"github.com/go-redis/redis/v8"
	"time"
)

type weMsg struct {
	Location string    `json:"location"`
	Content  string    `json:"content"`
	Time     time.Time `json:"time"`
	Level    int       `json:"level"`
}

func New(location string, err error, level int) {

}

func In(cli *redis.Client, wm weMsg) {

}
