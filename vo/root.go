package vo

import (
	"github.com/yu1745/bilibili_crawler_master/model"
)

type MidAndTask struct {
	Mid     int        `json:"mid"`
	Task    model.Task `json:"task"`
	hasNext int        //-1就是插入发生重复
}

type Paged interface {
	HasNextPage() bool
	NextTask()
	Store()
}
