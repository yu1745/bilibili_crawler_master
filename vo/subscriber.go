package vo

import (
	"github.com/yu1745/bilibili_crawler_master/db"
	"github.com/yu1745/bilibili_crawler_master/model"
	"gorm.io/gorm/clause"
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
}

func (this *Subs) Store() {
	if len(this.Data.List) == 0 {
		return
	}
	var subs []model.Up
	for _, v := range this.Data.List {
		subs = append(subs, model.Up{UID: v.Mid, LastScanned: time.Unix(946656000, 0)})
	}
	db.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&subs)
}
