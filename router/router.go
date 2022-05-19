package router

import (
	"github.com/gin-gonic/gin"
	"github.com/yu1745/bilibili_crawler_master/api/v1"
	"log"
)

func Init() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		_, _ = c.Writer.Write([]byte("pong\n"))
	})
	r.PUT("/comment", v1.RootComment)
	r.PUT("/video", v1.Video)
	/*r.GET("/", func(c *gin.Context) {
		c.JSON(0, "done")
	})
	err := r.RunTLS(":46512", "/root/fullchain.pem", "/root/privkey.pem")
	if err != nil {
		log.Fatalln(err)
	}*/
	r.PUT("/q", v1.PutQueue)
	r.GET("/q", v1.GetQueue)
	err := r.RunTLS(":8443", "/root/fullchain.pem", "/root/privkey.pem")
	if err != nil {
		log.Fatalln(err)
	}
}
