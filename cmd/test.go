package main

import (
	"github.com/yu1745/bilibili_crawler_master/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func main() {
	dsn := "host=127.0.0.1 port=5432 user=postgres dbname=bilibili password=asdk7788AA sslmode=disable TimeZone=Asia/Shanghai"
	Db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Db.Logger = logger.Default.LogMode(logger.Info)
	Db.Save(&model.User{
		UID:         524931,
		LastScanned: time.Now(),
	})
	println()
}
