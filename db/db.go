package db

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var Db *gorm.DB

func Init() {
	var err error
	Db, err = gorm.Open(mysql.Open("root:asdk7788AA@tcp(127.0.0.1)/bilibili?parseTime=True&loc=Local"), &gorm.Config{
		//SkipDefaultTransaction: true,
	})
	if err != nil {
		panic("failed to connect database")
	}
	Db.Logger = logger.Default.LogMode(logger.Error)
}
