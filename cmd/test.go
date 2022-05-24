package main

import (
	C "github.com/yu1745/bilibili_crawler_master/constant"
	"github.com/yu1745/bilibili_crawler_master/model"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/logger"
	"time"
)

func main() {
	C.InitDB()
	C.Db.Logger = logger.Default.LogMode(logger.Info)
	C.Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&model.Up{UID: 361713503, LastScanned: time.Now()})
}
