package vo

import (
	C "github.com/yu1745/bilibili_crawler_master/constant"
	"github.com/yu1745/bilibili_crawler_master/model"
	"gorm.io/gorm/clause"
	"log"
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
				task := Task{
					TaskType: GetCommentsFromVideo,
					Payload:  u.String(),
					New:      true,
				}
				log.Printf("[GetCommentsFromVideo] id=%d page=%d\n", this.Mid, i)
				C.Q.Offer(task.Encode())
			}
		}
	} else {
		//只扫新的
		if this.HasNext != -1 {
			var pageNum int
			if this.Data.Page.Count%this.Data.Page.Size == 0 {
				pageNum = this.Data.Page.Count / this.Data.Page.Size
			} else {
				pageNum = this.Data.Page.Count/this.Data.Page.Size + 1
			}
			if this.Data.Page.Num < pageNum {
				u, err := url.Parse(this.Task.Payload)
				if err != nil {
					log.Println(err)
				}
				q := u.Query()
				q.Set("pn", strconv.Itoa(this.Data.Page.Num+1))
				log.Printf("[GetCommentsFromVideo] id=%d page=%d\n", this.Mid, this.Data.Page.Num+1)
				u.RawQuery = q.Encode()
				task := Task{
					TaskType: GetCommentsFromVideo,
					Payload:  u.String(),
					New:      false,
				}
				C.Q.Offer(task.Encode())
			}
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
	if !this.Task.New {
		//不是第一次扫
		//每页检查一下是否扫到了上次已经扫了的部分
		var keys []int
		for _, v := range this.Data.Replies {
			keys = append(keys, v.Rpid)
		}
		if int(C.Db.Find(&cmts, keys).RowsAffected) == len(this.Data.Replies) {
			this.HasNext = -1
		}
		cmts = make([]model.Comment, 0)
	}
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
	C.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&cmts)
	var users []model.User
	for _, v := range this.Data.Replies {
		if C.Db.Limit(1).Find(&model.User{UID: v.Mid}).RowsAffected == 0 {
			users = append(users, model.User{
				UID:         v.Mid,
				LastScanned: time.Unix(946656000, 0),
			})
			C.Q.Offer(NewInitTask(GetSubscribers, strconv.Itoa(v.Mid), false).Encode())
		}
	}
	if len(users) > 0 {
		C.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&users)
	}
}
