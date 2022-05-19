package vo

import (
	"encoding/json"
	"github.com/yu1745/bilibili_crawler_master/db"
	"github.com/yu1745/bilibili_crawler_master/model"
	"github.com/yu1745/bilibili_crawler_master/queue"
	"gorm.io/gorm/clause"
	"log"
	"net/url"
	"strconv"
	"time"
)

type Comment struct {
	Code int `json:"code"`
	Data struct {
		Page struct {
			Num   int `json:"num"`
			Size  int `json:"size"`
			Count int `json:"count"`
		} `json:"page"`
		Replies []struct {
			Rpid    int `json:"rpid"`
			Oid     int `json:"oid"`
			Mid     int `json:"mid"`
			Like    int `json:"like"`
			Ctime   int `json:"ctime"`
			Content struct {
				Message string `json:"message"`
			} `json:"content,omitempty"`
		} `json:"replies"`
	} `json:"data"`
	MidAndTask
}

func (this *Comment) NextTask() {
	u, err := url.Parse(this.Task.Payload)
	if err != nil {
		log.Println(err)
	}
	q := u.Query()
	q.Set("pn", strconv.Itoa(this.Data.Page.Count+1))
	u.RawQuery = q.Encode()
	this.Task.Payload = u.String()
	b, err := json.Marshal(&this.Task)
	if err != nil {
		log.Println(err)
	}
	queue.Q.Offer(b)
}

func (this *Comment) HasNextPage() bool {
	if this.hasNext == -1 {
		return false
	}
	var pageNum int
	if this.Data.Page.Count%this.Data.Page.Size == 0 {
		pageNum = this.Data.Page.Count / this.Data.Page.Size
	} else {
		pageNum = this.Data.Page.Count/this.Data.Page.Size + 1
	}
	return this.Data.Page.Num < pageNum
}

func (this *Comment) Store() {
	var cmts []model.Comment
	for _, v := range this.Data.Replies {
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
	d := db.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&cmts)
	if int(d.RowsAffected) != len(this.Data.Replies) {
		this.hasNext = -1
	}
}
