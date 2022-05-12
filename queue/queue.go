package queue

import (
	"github.com/yu1745/bilibili_crawler_master/util"
	"log"
)

var Q *util.DurableQueue

func init() {
	var err error
	Q, err = util.NewQueue("nmsl")
	if err != nil {
		log.Fatalln(err)
	}
}
