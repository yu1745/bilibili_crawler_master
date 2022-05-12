package worker

import (
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
		println(v)
	}
}

func TestCreateWorker(t *testing.T) {
	PutCode("/tmp/function2.zip", "function.zip")
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
