package api

import (
	"github.com/gin-gonic/gin"
	"github.com/yu1745/bilibili_crawler_master/queue"
)

func PutQueue(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		return
	}
	queue.Q.Offer(data)
}

func GetQueue(c *gin.Context) {
	poll, err := queue.Q.Poll()
	if err != nil {
		return
	}
	_, _ = c.Writer.Write(poll)
}
