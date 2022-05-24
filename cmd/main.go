package main

import (
	"context"
	"flag"
	"github.com/yu1745/bilibili_crawler_master/db"
	"github.com/yu1745/bilibili_crawler_master/queue"
	"github.com/yu1745/bilibili_crawler_master/router"
	"github.com/yu1745/bilibili_crawler_master/service"
	"github.com/yu1745/bilibili_crawler_master/util/worker"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var cancelFunc context.CancelFunc

var (
	num int
)

func init() {
	flag.IntVar(&num, "n", 20, "num of worker")
	flag.Parse()
	go func() {
		//处理ctrl+c
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		println("-------SIGINT-------")
		if cancelFunc != nil {
			cancelFunc()
			time.Sleep(time.Second * 3)
		}
		os.Exit(0)
	}()
}

func main() {
	//log.SetFlags(log.Lshortfile)
	log.SetFlags(^log.Ltime)
	go db.Init()
	go router.Init()
	queue.Init()
	worker.Init(num)
	ctx, cancelFunc_ := context.WithCancel(context.Background())
	cancelFunc = cancelFunc_
	for i, v := range worker.Workers {
		time.Sleep(time.Duration(int64(float64(i)/float64(len(worker.Workers))*1000) * int64(time.Millisecond)))
		go service.Process(v, ctx)
	}
	//永久阻塞
	<-make(chan struct{})
}
