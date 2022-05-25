package vo

import (
	"bytes"
	"encoding/json"
	C "github.com/yu1745/bilibili_crawler_master/constant"
	"github.com/yu1745/bilibili_crawler_master/model"
	"gorm.io/gorm/clause"
	"log"
	"net/url"
	"strconv"
	"time"
)

type Video struct {
	Code int `json:"code"`
	Data struct {
		List struct {
			Vlist []struct {
				Aid int `json:"aid"`
			} `json:"vlist"`
		} `json:"list"`
		Page struct {
			Pn    int `json:"pn"`
			Ps    int `json:"ps"`
			Count int `json:"count"`
		} `json:"page"`
	} `json:"data"`
	Meta
}

func (this *Video) Next() {
	if this.Task.New {
		if this.Data.Page.Pn == 1 {
			//没扫过，直接并行生成
			var pageNum int
			if this.Data.Page.Count%this.Data.Page.Ps == 0 {
				pageNum = this.Data.Page.Count / this.Data.Page.Ps
			} else {
				pageNum = this.Data.Page.Count/this.Data.Page.Ps + 1
			}
			for i := 2; i < pageNum; i++ {
				u, err := url.Parse(this.Task.Payload)
				if err != nil {
					log.Println()
				}
				q := u.Query()
				q.Set("pn", strconv.Itoa(i))
				u.RawQuery = q.Encode()
				task := Task{TaskType: GetVideoFromUp, Payload: u.String(), New: true}
				var buf bytes.Buffer
				e := json.NewEncoder(&buf)
				e.SetEscapeHTML(false)
				err = e.Encode(&task)
				if err != nil {
					log.Println(err)
				}
				C.Q.Offer(buf.Bytes())
			}
		}
	} else {
		if this.HasNext != -1 {
			var pageNum int
			if this.Data.Page.Count%this.Data.Page.Ps == 0 {
				pageNum = this.Data.Page.Count / this.Data.Page.Ps
			} else {
				pageNum = this.Data.Page.Count/this.Data.Page.Ps + 1
			}
			if this.Data.Page.Pn < pageNum {
				u, err := url.Parse(this.Task.Payload)
				if err != nil {
					log.Println(err)
				}
				q := u.Query()
				q.Set("pn", strconv.Itoa(this.Data.Page.Pn+1))
				//println("pn:", strconv.Itoa(this.Data.Page.Pn+1))
				u.RawQuery = q.Encode()
				task := Task{TaskType: GetVideoFromUp, Payload: u.String(), New: false}
				var buf bytes.Buffer
				e := json.NewEncoder(&buf)
				e.SetEscapeHTML(false)
				err = e.Encode(&task)
				if err != nil {
					log.Println(err)
				}
				C.Q.Offer(buf.Bytes())
			}
		}
	}
}

func (this *Video) Store() {
	if this.Pn == 1 {
		//更新数据库中上次扫面时间字段
		C.Db.Save(&model.Video{Avid: this.Mid, LastUpdated: time.Now()})
	}
	if len(this.Data.List.Vlist) == 0 {
		return
	}
	var videos []model.Video
	if !this.Task.New {
		//不是第一次扫
		//每页检查一下是否扫到了上次已经扫了的部分
		var avids []int
		for _, v := range this.Data.List.Vlist {
			avids = append(avids, v.Aid)
		}
		if int(C.Db.Find(&videos, avids).RowsAffected) == this.Data.Page.Ps {
			this.HasNext = -1
		}
		videos = make([]model.Video, 0)
	}
	for _, v := range this.Data.List.Vlist {
		if C.Db.Limit(1).Find(&model.Video{Avid: v.Aid}).RowsAffected == 0 {
			videos = append(videos, model.Video{Avid: v.Aid, LastUpdated: time.Unix(946656000, 0)})
			C.Q.Offer(NewInitTask(GetVideoFromUp, strconv.Itoa(v.Aid), false).Encode())
		}
	}
	if len(videos) > 0 {
		C.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&videos)
	}
}
