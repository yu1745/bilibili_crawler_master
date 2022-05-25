package service

import (
	"context"
	"encoding/json"
	C "github.com/yu1745/bilibili_crawler_master/constant"
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
			poll, err := C.Q.Poll()
			if err != nil {
				continue
			}
			output, err := worker.Invoke(v.Name, poll)
			if err != nil || len(output) == 0 {
				//panic("rebuild worker")
				//todo 重建worker
				continue
			}
			unquote, err := strconv.Unquote(string(output))
			if err != nil {
				log.Println(err)
				println(string(output))
				continue
			}

			println(unquote)

			var task vo.Task
			_ = json.Unmarshal(poll, &task)
			switch task.TaskType {
			case vo.GetCommentsFromVideo:
				var cmt vo.MainComment
				err = json.Unmarshal([]byte(unquote), &cmt)
				if err != nil {
					log.Fatalln(err)
				}
				cmt.Store()
				cmt.Next()
			case vo.GetSubscribers:
				var subs vo.Subs
				err = json.Unmarshal([]byte(unquote), &subs)
				if err != nil {
					log.Fatalln(err)
				}
				subs.Store()
				subs.Next()
			case vo.GetVideoFromUp:
				var video vo.Video
				err = json.Unmarshal([]byte(unquote), &video)
				if err != nil {
					log.Fatalln(err)
				}
				video.Store()
				video.Next()
			default:
			}

		}
	}
}
