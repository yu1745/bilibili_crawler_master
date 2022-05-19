package router

import (
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	go Init()
	time.Sleep(time.Hour)
}
