package websocket

import (
	"log"
	"os"
	"sync"
)

type clients struct {
	mutex *sync.Mutex
	counter int
	list map[*client]bool
}

func (c *clients) attach(client *client) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _,ok := c.list[client]; ok {
		log.Println("重复保存了Client")
		os.Exit(0)
	}
	c.counter++
	c.list[client] = true
}

func (c *clients) detach(client *client) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if _,ok := c.list[client]; !ok {
		log.Println("移除一个已经不存在的Client")
		os.Exit(0)
	}
	c.counter--
	delete(c.list,client)
}