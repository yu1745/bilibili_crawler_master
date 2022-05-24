package vo

import (
	"bytes"
	"encoding/json"
	"github.com/yu1745/bilibili_crawler_master/db"
	"github.com/yu1745/bilibili_crawler_master/model"
	"github.com/yu1745/bilibili_crawler_master/queue"
	"gorm.io/gorm/clause"
	"log"
	"math"
	"net/url"
	"strconv"
	"time"
)

type MainComment struct {
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
	Meta
}

func (this *MainComment) Next() {
	if this.Task.New {
		if this.Data.Page.Num == 1 {
			//没扫过，直接并行生成
			var pageNum int
			if this.Data.Page.Count%this.Data.Page.Size == 0 {
				pageNum = this.Data.Page.Count / this.Data.Page.Size
			} else {
				pageNum = this.Data.Page.Count/this.Data.Page.Size + 1
			}
			for i := 2; i < pageNum; i++ {
				u, err := url.Parse(this.Task.Payload)
				if err != nil {
					log.Println(err)
				}
				q := u.Query()
				q.Set("pn", strconv.Itoa(i))
				u.RawQuery = q.Encode()
				this.Task.Payload = u.String()
				this.Task.New = true
				//批量生成的任务，不用检验是否有下一页
				this.HasNext = -1
				var buf bytes.Buffer
				e := json.NewEncoder(&buf)
				e.SetEscapeHTML(false)
				err = e.Encode(&this.Task)
				if err != nil {
					log.Println(err)
				}
				b := buf.Bytes()
				log.Printf("[%v] page %d\n", this.Task.TaskType, i)
				queue.Q.Offer(b)
			}
		}
	} else {
		//只扫新的
		if this.HasNext != -1 {
			u, err := url.Parse(this.Task.Payload)
			if err != nil {
				log.Println(err)
			}
			q := u.Query()
			q.Set("pn", strconv.Itoa(this.Data.Page.Num+1))
			log.Printf("[%v] page %d\n	", this.Task.TaskType, this.Data.Page.Num+1)
			u.RawQuery = q.Encode()
			this.Task.Payload = u.String()
			this.Task.New = false
			var buf bytes.Buffer
			e := json.NewEncoder(&buf)
			e.SetEscapeHTML(false)
			err = e.Encode(&this.Task)
			if err != nil {
				log.Println(err)
			}
			b := buf.Bytes()
			queue.Q.Offer(b)
		}
	}
}

//func (this *MainComment) HasNextPage() bool {
//	/*if this.HasNext == -1 {
//		return false
//	}else {
//		return true
//	}*/
//	return !(this.HasNext == -1)
//	/*var pageNum int
//	if this.Data.Page.Count%this.Data.Page.Size == 0 {
//		pageNum = this.Data.Page.Count / this.Data.Page.Size
//	} else {
//		pageNum = this.Data.Page.Count/this.Data.Page.Size + 1
//	}
//	return this.Data.Page.Num < pageNum*/
//}

func (this *MainComment) Store() {
	if len(this.Data.Replies) == 0 {
		this.HasNext = -1
		return
	}
	var cmts []model.Comment
	minRpid := math.MaxInt
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
		if v.Rpid < minRpid {
			minRpid = v.Rpid
		}
	}
	if !this.Task.New {
		//不是第一次扫
		//逐页检验
		var dbMaxRpid int
		db.Db.Raw(`select max(rpid) from comment where "to" = ?`, cmts[0].To).Scan(&dbMaxRpid)
		if !(minRpid > dbMaxRpid) {
			this.HasNext = -1
		}
	}
	/*d := */ db.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&cmts)
	/*if int(d.RowsAffected) != len(this.Data.Replies) {
		this.HasNext = -1
		log.Println("insert conflict")
	}*/
	var ups []model.Up
	for _, v := range this.Data.Replies {
		ups = append(ups, model.Up{
			UID:         v.Mid,
			LastScanned: time.Unix(946656000, 0),
		})
	}
	db.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&ups)
}
