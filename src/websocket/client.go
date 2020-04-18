package websocket

import (
	"github.com/gorilla/websocket"
	"sync"
	"time"
)

const (
	// Timeout when the client has no any response
	pongWait = 5 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 256
)

// Client is a middleman between the websocket connection and the hub.
type client struct {
	// The Client toggle
	enable			bool

	// The websocket connection.
	conn 			*websocket.Conn

	// connection time
	ConnectionTime	time.Time

	// The message sender of the client
	sender 			struct {
		enable bool
		mutex *sync.Mutex
	}
}

func (c *client) Send(message []byte) {
	if !c.enable { return }

	c.sender.mutex.Lock()
	defer c.sender.mutex.Unlock()

	if !c.enable { return }
	if !c.sender.enable { return }
	c.conn.WriteMessage(websocket.TextMessage,message)
}

func (c *client) Close() {
	c.enable = false

	c.sender.mutex.Lock()
	defer c.sender.mutex.Unlock()

	if c.sender.enable {
		c.sender.enable = false
	}else{
		return
	}

	// 关闭通道
	// 会触发，receiver 产生错误而退出
	c.conn.Close()
}

func (c *client) receiver(channel chan *Pack) {
	defer c.Close()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPingHandler(func(appData string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_,message,err := c.conn.ReadMessage()
		if err != nil {
			// 如果错误信息，不是以下情况
			// 打印日志
			//if websocket.IsUnexpectedCloseError(err,websocket.CloseGoingAway,websocket.CloseAbnormalClosure) {
			//	log.Printf("error : %v",err)
			//}
			return
		}

		// 重置dead line时间
		c.conn.SetReadDeadline(time.Now().Add(pongWait))

		// 消息存放到channel
		// 等待消息消费进程处理
		channel <- &Pack{Client:c,Message:message}
	}
}
