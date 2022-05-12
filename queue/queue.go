package queue

import (
	"log"
	"master/util"
)

var Q *util.DurableQueue

func init() {
	var err error
	Q, err = util.NewQueue("nmsl")
	if err != nil {
		log.Fatalln(err)
	}
}
