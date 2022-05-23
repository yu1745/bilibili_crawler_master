package main

import (
	"fmt"
	"github.com/yu1745/bilibili_crawler_master/db"
	"github.com/yu1745/bilibili_crawler_master/model"
)

func main() {
	db.Init()
	var cmt model.Comment
	db.Db.Where("rpid > ?", 113765133632).Order("rpid").Limit(1).Find(&cmt)
	fmt.Printf("%+v", cmt)
	var i int
	db.Db.Raw("select max(rpid) from comment where `to`= ?", 3).Scan(&i)
	println(i)
}
