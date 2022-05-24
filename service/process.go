package service

import (
	"context"
	"encoding/json"
	"github.com/yu1745/bilibili_crawler_master/queue"
	"github.com/yu1745/bilibili_crawler_master/util/worker"
	"github.com/yu1745/bilibili_crawler_master/vo"
	"log"
	"strconv"
	"time"
)

func Process(v worker.Worker, ctx context.Context) {
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
			output, err := worker.Invoke(v.Name, poll)
			if err != nil {
				panic("rebuild worker")
				//todo 重建worker
			}
			unquote, err := strconv.Unquote(string(output))
			if err != nil {
				log.Fatalln(err)
			}
			var task vo.Task
			_ = json.Unmarshal(poll, &task)
			var paged vo.Paged
			switch task.TaskType {
			case vo.GetCommentsFromVideo:
				var cmt vo.MainComment
				err = json.Unmarshal([]byte(unquote), &cmt)
				if err != nil {
					log.Fatalln(err)
				}
				paged = &cmt
			default:
			}
			paged.Store()
			paged.Next()
		}
	}
}
