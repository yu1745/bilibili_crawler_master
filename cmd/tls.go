package main

import (
	"github.com/gin-gonic/gin"
	"log"
)

func main() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(0, "nmsl")
	})
	err := r.RunTLS(":46512", "/root/fullchain.pem", "/root/privkey.pem")
	if err != nil {
		log.Fatalln(err)
	}
}
