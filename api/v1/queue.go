package v1

import (
	"github.com/gin-gonic/gin"
	C "github.com/yu1745/bilibili_crawler_master/constant"
)

func PutQueue(c *gin.Context) {
	data, err := c.GetRawData()
	if err != nil {
		return
	}
	C.Q.Offer(data)
}

/*func GetQueue(c *gin.Context) {
	poll, err := C.Q.Poll()
	if err != nil {
		return
	}
	_, _ = c.Writer.Write(poll)
}*/
