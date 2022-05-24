package db

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func Init() {
	dsn := "host=127.0.0.1 port=5432 user=postgres dbname=bilibili password=asdk7788AA sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	//Db, err = gorm.Open(mysql.Open("root:asdk7788AA@tcp(127.0.0.1)/bilibili?parseTime=True&loc=Local"), &gorm.Config{})
	Db, err = gorm.Open(postgres.Open(dsn /*"postgres:asdk7788AA@tcp(127.0.0.1)/bilibili?parseTime=True&loc=Local"*/), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	Db.Logger = logger.Default.LogMode(logger.Error)
}
