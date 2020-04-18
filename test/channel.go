package test

import (
	"fmt"
	"time"
)

type Router struct {
	cache		chan int
}

func main() {

	cache := make(chan int,100)
	r := &Router{cache:cache}

	go func() {
		for {
			select {
			case data := <- r.cache:
				fmt.Println(data)
			}
		}
	}()

	time.Sleep(5*time.Second)
}