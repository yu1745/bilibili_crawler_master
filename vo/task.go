package vo

import (
	"bytes"
	"encoding/json"
	"log"
	"net/url"
	"strconv"
	"strings"
)

type TaskType int

func (t TaskType) Equals(s string) bool {
	return strconv.Itoa(int(t)) == s
}

const (
	GetCommentsFromVideo TaskType = iota
	GetVideoFromUp
	//GetFollowers //主要通过视频抓取评论，非up爬了也没用
	GetSubscribers
)

type Task struct {
	TaskType   TaskType `json:"TaskType"`
	Payload    string   `json:"Payload"`
	New        bool     `json:"new"`
	Derivation string   `json:"derivation"`
}

/*func (this *Task) getPayload() string {
	bytes, err := base64.StdEncoding.DecodeString(this.Payload)
	if err != nil {
		return "null"
	}
	return string(bytes)
}

func (this *Task) setPayload(s string) {
	this.Payload = base64.StdEncoding.EncodeToString([]byte(s))
}*/

func (this *Task) Encode() []byte {
	var buf bytes.Buffer
	e := json.NewEncoder(&buf)
	e.SetEscapeHTML(false)
	err := e.Encode(this)
	if err != nil {
		log.Fatalln(err)
	}
	return buf.Bytes()
}

func NewInitTask(taskType TaskType, target string, new bool) *Task {
	switch taskType {
	case GetCommentsFromVideo:
		s := `http://api.bilibili.com/x/v2/reply?type=1&ps=49&pn=1&nohot=1`
		u, _ := url.Parse(s)
		q := u.Query()
		q.Set("oid", target)
		u.RawQuery = q.Encode()
		return &Task{
			TaskType: taskType,
			Payload:  u.String(),
			New:      new,
		}
	case GetSubscribers:
		s := `http://api.bilibili.com/x/relation/followings?&pn=1&ps=50`
		u, _ := url.Parse(s)
		q := u.Query()
		q.Set("vmid", target)
		u.RawQuery = q.Encode()
		return &Task{
			TaskType: taskType,
			Payload:  u.String(),
			New:      new,
		}
	case GetVideoFromUp:
		s := `http://api.bilibili.com/x/space/arc/search?order=pubdate&pn=1&ps=49`
		u, _ := url.Parse(s)
		q := u.Query()
		q.Set("mid", target)
		u.RawQuery = q.Encode()
		return &Task{
			TaskType: taskType,
			Payload:  u.String(),
			New:      new,
		}
	default:
		panic("Unrecognized task type")
	}
}

func (this *Task) AllowDerivation() bool {
	if this.Derivation == "" {
		return false
	}
	split := strings.Split(this.Derivation, ";")
	for _, v := range split {
		if this.TaskType.Equals(v) {
			return true
		}
	}
	return false
}
