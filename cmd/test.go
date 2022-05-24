package main

import (
	"github.com/yu1745/bilibili_crawler_master/model"
)

func main() {
	/*dsn := "host=127.0.0.1 port=5432 user=postgres dbname=bilibili password=asdk7788AA sslmode=disable TimeZone=Asia/Shanghai"
	Db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Db.Logger = logger.Default.LogMode(logger.Info)
	var ups []model.Up
	ups = append(ups, model.Up{UID: 1})
	ups = append(ups, model.Up{UID: 123})
	ups = append(ups, model.Up{UID: 124})
	ups = append(ups, model.Up{UID: 125})
	Db.Clauses(clause.OnConflict{DoNothing: true}).Create(&ups)
	println()*/
	videos := make([]model.Video, 10)
	videos = append(videos, model.Video{Avid: 213312})
	println()
}
