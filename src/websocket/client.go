package websocket

import (
	"bytes"
	"github.com/gorilla/websocket"
	"time"
)

const (
	// Timeout when the client has no any response
	pongWait = 60 * time.Second

	// Maximum message size allowed from peer.
	maxMessageSize = 256
)

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	// Point to the Hub
	Hub 			*HubStruct

	// The websocket connection.
	conn 			*websocket.Conn

	// The message sender of the client
	sender 			struct {
		isOpen bool
		reference int
		channel chan[]byte
	}
}

func (c *Client) Send(message []byte) {
	c.sender.reference++

	if c.sender.isOpen {
		c.sender.channel <- message
	}

	c.sender.reference--
}


func (c *Client) wPump() {
	defer func(c *Client) {
		c.sender.isOpen = true
	}(c)

	c.sender = struct {
		isOpen bool
		reference int
		channel chan[]byte
	}{isOpen:false,reference:0,channel:make(chan []byte,maxClients)}

	go func(c *Client) {
		for {
			select {
			case message,ok := <- c.sender.channel:
				if ok {
					c.conn.WriteMessage(websocket.TextMessage,message)
				}else{
					return
				}
			}
		}
	}(c)
}

func (c *Client) rPump(onMessage OnMessage) {
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })

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
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		if bytes.Equal(message,[]byte("ping")) != true {
			if onMessage != nil { onMessage(c,message) }
		}
	}
}

func (c *Client) close() {
	// 如果已经关闭，直接退出
	if c.sender.isOpen == false { return }

	// 关闭发送器
	c.sender.isOpen = false

	// 关闭通道
	// 会触发，rPump产生错误而退出
	c.conn.Close()

	// 等待发送中的消息完成
	// 主要防止channel关闭，导致其他goroutine向channel发送消息时，发生异常
	go func(c *Client) {
		for {
			if c.sender.reference == 0 {
				// 关闭发送通道
				close(c.sender.channel)
				return
			}
			time.Sleep(5*time.Millisecond)
		}
	}(c)
}