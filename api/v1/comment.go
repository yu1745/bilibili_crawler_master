package v1

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/yu1745/bilibili_crawler_master/vo"
	"log"
)

func RootComment(c *gin.Context) {
	var cmt vo.Comment
	if err := c.ShouldBind(&cmt); err != nil {
		log.Println(err)
	}

	marshal, _ := json.Marshal(&cmt)
	println(string(marshal))

	if cmt.Code == 0 {
		if len(cmt.Data.Replies) == 0 {
			return
		}
		/*var pageNum int
		if cmt.Data.Page.Count%cmt.Data.Page.Size == 0 {
			pageNum = cmt.Data.Page.Count / cmt.Data.Page.Size
		} else {
			pageNum = cmt.Data.Page.Count/cmt.Data.Page.Size + 1
		}
		u, _ := url.Parse("http://api.bilibili.com/x/v2/reply?type=1&ps=49&nohot=1")
		query := u.Query()
		query.Set("oid", strconv.Itoa(cmt.Data.Replies[0].Oid))
		query.Set("pn", strconv.Itoa(pageNum))
		u.RawQuery = query.Encode()
		t := &model.Task{
			TaskType: model.GetCommentsFromVideo,
			Payload:  u.String(),
		}
		bytes, err := json.Marshal(t)
		if err != nil {
			log.Fatalln(err)
		}
		queue.Q.Offer(bytes)*/
		/*var cmts []model.Comment
		for _, v := range cmt.Data.Replies {
			cmts = append(cmts, model.Comment{
				Rpid:    v.Rpid,
				From:    v.Mid,
				To:      v.Oid,
				Avid:    v.Oid,
				Like:    v.Like,
				Ctime:   time.Unix(int64(v.Ctime), 0),
				Content: v.Content.Message,
			})
		}
		db.Db.Create(&cmts)*/
		cmt.Store()
		if cmt.HasNextPage() {
			cmt.Next()
		}
		//处理user
		/*for _, v := range cmt.Data.Replies {
			u, _ := url.Parse("http://api.bilibili.com/x/space/arc/search?order=pubdate&pn=1&ps=50")
			q := u.Query()
			q.Set("mid", strconv.Itoa(v.Mid))
			u.RawQuery = q.Encode()
			b, _ := json.Marshal(&model.Task{
				TaskType: model.GetVideoFromUp,
				Payload:  u.String(),
			})
			queue.Q.Offer(b)
		}*/
	} else {
		log.Println("ip blocked")
	}
}
