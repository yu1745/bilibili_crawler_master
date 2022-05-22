package v1

import (
	"github.com/gin-gonic/gin"
	"github.com/yu1745/bilibili_crawler_master/vo"
	"log"
)

func RootComment(c *gin.Context) {
	var cmt vo.Comment
	if err := c.ShouldBind(&cmt); err != nil {
		log.Println(err)
	}
	//marshal, _ := json.Marshal(&cmt)
	//println(string(marshal))
	//fmt.Printf("%+v", cmt)
	if cmt.Code == 0 {
		if len(cmt.Data.Replies) == 0 {
			return
		}
		cmt.Store()
		if cmt.HasNextPage() {
			cmt.Next()
		}
	} else {
		log.Println("ip blocked")
	}
}
