package worker

import (
	"fmt"
	"strconv"
	"testing"
)

func TestGetAllNames(t *testing.T) {
	names, err := GetAllNames()
	if err != nil {
		t.Error(err)
	}
	for _, v := range names {
		println(v)
	}
}

func TestGetALlUrls(t *testing.T) {
	urls, err := GetALlUrls()
	if err != nil {
		t.Error(err)
	}
	for _, v := range urls {
		fmt.Printf("%+v", v)
	}
}

func TestCreateWorker(t *testing.T) {
	PutCode("/tmp/function.zip", "function.zip")
	err := CreateWorker("test")
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveWorker(t *testing.T) {
	err := RemoveWorker("test")
	if err != nil {
		t.Error(err)
	}
}

func TestRemoveAllWorkers(t *testing.T) {
	err := RemoveAllWorkers()
	if err != nil {
		t.Error(err)
	}
}

func TestPutCode(t *testing.T) {
	PutCode("/tmp/function.zip", "function.zip")
}

func TestGetFunctionUrl(t *testing.T) {
	url := GetFunctionUrl("test")
	println(url)
}

func TestInvoke(t *testing.T) {
	s := `{"TaskType":0,"Payload":"http://api.bilibili.com/x/v2/reply?type=1&oid=596088009&ps=49&pn=1&nohot=1","new":false}`
	invoke, err := Invoke("test", []byte(s))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(string(invoke))
	unquote, err := strconv.Unquote(string(invoke))
	if err != nil {
		t.Error(err)
	}
	fmt.Println(unquote)
}
