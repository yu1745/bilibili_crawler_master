package router

import (
	"github.com/gin-gonic/gin"
	"log"
	"master/api"
)

func Init() {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		_, _ = c.Writer.Write([]byte("pong"))
	})
	r.PUT("/comment", api.RootComment)
	/*r.GET("/", func(c *gin.Context) {
		c.JSON(0, "done")
	})
	err := r.RunTLS(":46512", "/root/fullchain.pem", "/root/privkey.pem")
	if err != nil {
		log.Fatalln(err)
	}*/
	r.PUT("/q", api.PutQueue)
	r.GET("/q", api.GetQueue)
	//err := r.Run(":8080")
	err := r.RunTLS(":8443", "/root/fullchain.pem", "/root/privkey.pem")
	if err != nil {
		log.Fatalln(err)
	}
}
