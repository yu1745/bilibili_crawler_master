package main

import (
	"github.com/yu1745/bilibili_crawler_master/db"
	"github.com/yu1745/bilibili_crawler_master/model"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"time"
)

func main() {
	db.Init()
	db.Db.Logger = logger.Default.LogMode(logger.Info)
	db.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.Up{UID: 361713503, LastScanned: time.Now()})
}
