package main

import (
	"github.com/yu1745/bilibili_crawler_master/util/worker"
	"log"
	"time"
)

func main() {
	t := time.Now()
	s := `{
  "TaskType": 0,
  "Payload": "http://api.bilibili.com/x/v2/reply?type=1&oid=2&ps=49&pn=10&nohot=1"
}`
	invoke, err := worker.Invoke("test", []byte(s))
	if err != nil {
		log.Println(err)
	}
	println(string(invoke))
	println(time.Now().Sub(t).String())
}
