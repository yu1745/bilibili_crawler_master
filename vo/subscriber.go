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

type Subs struct {
	Code int `json:"code"`
	Data struct {
		List []struct {
			Mid int `json:"mid"`
		} `json:"list"`
		Total int `json:"total"`
	} `json:"data"`
	Task
}

func (this *Subs) Next() {
	if this.Data.Total > 50 {
		num := 1
		if this.Data.Total%50 == 0 {
			num = this.Data.Total / 50
		} else {
			num = this.Data.Total/50 + 1
		}
		for i := 2; i <= num; i++ {
			u, err := url.Parse(this.Payload)
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
	if len(this.Data.List) == 0 {
		return
	}
	var subs []model.Up
	for _, v := range this.Data.List {
		subs = append(subs, model.Up{UID: v.Mid, LastScanned: time.Unix(946656000, 0)})
	}
	C.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&subs)
}
