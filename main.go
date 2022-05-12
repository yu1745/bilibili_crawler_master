package main

import (
	"bytes"
	"context"
	"github.com/google/uuid"
	"github.com/yu1745/bilibili_crawler_master/queue"
	"github.com/yu1745/bilibili_crawler_master/router"
	"github.com/yu1745/bilibili_crawler_master/util/worker"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var wg sync.WaitGroup

//var lock sync.RWMutex
var cancelFunc context.CancelFunc

func init() {
	//处理ctrl+c
	go func() {
		ch := make(chan os.Signal)
		signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
		<-ch
		cancelFunc()
		time.Sleep(time.Second * 5)
		os.Exit(0)
	}()

}

func main() {
	worker.Init()
	go router.Init()
	//time.Sleep(time.Hour)
	ctx, cancelFunc_ := context.WithCancel(context.Background())
	for _, v := range worker.Urls {
		cancelFunc = cancelFunc_
		wg.Add(1)
		go process(v, ctx)
	}
	wg.Wait()
}

func process(v string, ctx context.Context) {
	defer wg.Done()
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			poll, err := queue.Q.Poll()
			if err != nil {
				continue
			}
			req, err := http.NewRequest("POST", v, bytes.NewReader(poll))
			if err != nil {
				log.Fatalln(err)
			}
			req.Header.Set("Content-Type", "application/json")
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Println(err)
				//重建worker
				err := worker.RemoveWorker(v)
				if err != nil {
					log.Fatalln(err)
				}
				newUUID, err := uuid.NewUUID()
				if err != nil {
					log.Fatalln(err)
				}
				s := "worker-" + newUUID.String()
				err = worker.CreateWorker(s)
				if err != nil {
					log.Fatalln(err)
				}
				u := worker.GetFunctionUrl(s)
				wg.Add(1)
				go process(u, ctx)
				return
			}
			all, _ := ioutil.ReadAll(resp.Body)
			if len(all) > 0 {
				//重建worker
				err := worker.RemoveWorker(v)
				if err != nil {
					log.Fatalln(err)
				}
				newUUID, err := uuid.NewUUID()
				if err != nil {
					log.Fatalln(err)
				}
				s := "worker-" + newUUID.String()
				err = worker.CreateWorker(s)
				if err != nil {
					log.Fatalln(err)
				}
				u := worker.GetFunctionUrl(s)
				wg.Add(1)
				go process(u, ctx)
				return
			}
		}
	}
}
