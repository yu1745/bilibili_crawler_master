package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/yu1745/bilibili_crawler_master/db"
	"github.com/yu1745/bilibili_crawler_master/model"
	"github.com/yu1745/bilibili_crawler_master/queue"
	"github.com/yu1745/bilibili_crawler_master/vo"
	"log"
	"net/url"
	"strconv"
	"time"
)

func Video(c *gin.Context) {
	var vs vo.Video
	if err := c.ShouldBind(&vs); err != nil {
		log.Println(err)
	}
	if vs.Data.Page.Pn == 1 {
		var pageNum int
		if vs.Data.Page.Count&vs.Data.Page.Ps == 0 {
			pageNum = vs.Data.Page.Count / vs.Data.Page.Ps
		} else {
			pageNum = vs.Data.Page.Count/vs.Data.Page.Ps + 1
		}
		for i := 2; i <= pageNum; i++ {
			u, _ := url.Parse("http://api.bilibili.com/x/space/arc/search?order=pubdate&ps=50")
			q := u.Query()
			q.Set("mid", strconv.Itoa(vs.Mid))
			q.Set("pn", strconv.Itoa(pageNum))
			u.RawQuery = q.Encode()
			b, _ := json.Marshal(vo.Task{
				TaskType: vo.GetVideoFromUp,
				Payload:  u.String(),
			})
			queue.Q.Offer(b)
		}
	} else {
		var mvs []model.Video
		for _, v := range vs.Data.List.Vlist {
			mvs = append(mvs, model.Video{
				Avid:        v.Aid,
				LastUpdated: time.Now(),
			})
		}
		db.Db.Create(&mvs)
	}
}
