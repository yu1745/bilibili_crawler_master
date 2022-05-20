package vo

import "encoding/base64"

type TaskType int

const (
	GetCommentsFromVideo TaskType = iota
	GetVideoFromUp
	//GetFollowers //主要通过视频抓取评论，非up爬了也没用
	GetSubscribers
)

type Task struct {
	TaskType TaskType `json:"TaskType"`
	Payload  string   `json:"Payload"`
}

func (this *Task) getPayload() string {
	bytes, err := base64.StdEncoding.DecodeString(this.Payload)
	if err != nil {
		return "null"
	}
	return string(bytes)
}

func (this *Task) setPayload(s string) {
	this.Payload = base64.StdEncoding.EncodeToString([]byte(s))
}
