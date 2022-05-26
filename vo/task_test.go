package vo

import (
	"encoding/json"
	"testing"
)

func TestTask(t *testing.T) {
	marshal, err := json.Marshal(&Task{
		TaskType: GetCommentsFromVideo,
		Payload:  "http://api.bilibili.com/x/v2/reply?type=1&oid=981746036&ps=50&pn=1&nohot=1",
	})
	if err != nil {
		t.Error(err)
	}
	println(string(marshal))
}

func TestTask_AllowDerivation(t *testing.T) {
	task := Task{
		TaskType:   GetCommentsFromVideo,
		Derivation: "0",
	}
	println(task.AllowDerivation())
	task.Derivation = "0;1;2"
	println(task.AllowDerivation())
	task.Derivation = "1;2"
	println(task.AllowDerivation())
}
