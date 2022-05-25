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

const maxPageNumber = 5

type Subs struct {
	Code int `json:"code"`
	Data struct {
		List []struct {
			Mid int `json:"mid"`
		} `json:"list"`
		Total int `json:"total"`
	} `json:"data"`
	Meta
}

func (this *Subs) Next() {
	if this.Pn == 1 && this.Data.Total > this.Ps {
		num := 1
		if this.Data.Total%this.Ps == 0 {
			num = this.Data.Total / this.Ps
		} else {
			num = this.Data.Total/this.Ps + 1
		}
		//取两个之中的较小值
		if maxPageNumber < num {
			num = maxPageNumber
		}
		for i := 2; i <= num; i++ {
			u, err := url.Parse(this.Task.Payload)
			if err != nil {
				log.Println(err)
			}
			q := u.Query()
			q.Set("pn", strconv.Itoa(i))
			u.RawQuery = q.Encode()
			task := Task{TaskType: GetSubscribers, Payload: u.String()}
			var buf bytes.Buffer
			e := json.NewEncoder(&buf)
			e.SetEscapeHTML(false)
			err = e.Encode(&task)
			C.Q.Offer(buf.Bytes())
		}
	}
}

func (this *Subs) Store() {
	if this.Pn == 1 {
		//更新数据库中上次扫面时间字段
		C.Db.Save(&model.User{UID: this.Mid, LastScanned: time.Now()})
	}
	if len(this.Data.List) == 0 {
		return
	}
	var subs []model.Up
	if !this.Task.New {
		//不是第一次扫
		//每页检查一下是否扫到了上次已经扫了的部分
		var uids []int
		for _, v := range this.Data.List {
			uids = append(uids, v.Mid)
		}
		if int(C.Db.Find(&subs, uids).RowsAffected) == len(this.Data.List) {
			this.HasNext = -1
		}
		subs = make([]model.Up, 0)
	}
	for _, v := range this.Data.List {
		if C.Db.Limit(1).Find(&model.Up{UID: v.Mid}).RowsAffected == 0 {
			subs = append(subs, model.Up{UID: v.Mid, LastScanned: time.Unix(946656000, 0)})
			C.Q.Offer(NewInitTask(GetVideoFromUp, strconv.Itoa(v.Mid), false).Encode())
		}
	}
	if len(subs) > 0 {
		C.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&subs)
	}
}
