package websocket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type hub struct {
	// 基本配置
	settings 				*Settings
	// 当连接发生事件变化的时候，通过此回调，通知应用程序
	callback				*EventCallback

	// 所有的客户端连接
	clients					*clients

	// 握手协议
	upgrader 				*websocket.Upgrader

	// 数据接收器
	// Client to Server
	recv					chan *Pack
}

// 消息广播
func (h *hub) Broadcast(message []byte) {
	for c := range h.clients.list {
		c.Send(message)
	}
}

func (h *hub) receiver() {
	if h.settings.ProcNum < 1 {
		log.Println("进程数量不能少于1")
		os.Exit(1)
	}
	h.recv = make(chan *Pack)
	for i:=0; i<h.settings.ProcNum; i++ {
		go func(h *hub) {
			for {
				select {
				case pack := <- h.recv:
					pack.Hub = h
					h.callback.OnMessage(pack)
				}
			}
		}(h)
	}
}

func (h *hub) handle(w http.ResponseWriter, r *http.Request) {
	// 完成协议握手
	// 握手程序放入goroutine执行发生异常
	conn,err := h.upgrader.Upgrade(w,r,nil)
	if err != nil {
		fmt.Println("握手失败!")
		fmt.Println(err)
		return
	}

	go func(h *hub,conn *websocket.Conn) {
		// 创建Client实例
		c := &client{enable:true,conn:conn,ConnectionTime:time.Now(),sender: struct {
			enable bool
			mutex  *sync.Mutex
		}{enable: true, mutex: &sync.Mutex{}}}

		// 调用OnOpen回调，检查连接是否有效
		// 如果连接无效，关闭连接
		if !h.callback.OnOpen(&Pack{
			Hub:     h,
			Client:  c,
			Message: nil,
		}) {
			c.Close()
			return
		}

		// 防止连接的时候，有公共消息推送到未校验的连接中，所以需要在 onOpen 后进行添加
		// 需要放置在 rPump 前面，防止还没有进行 attach 就立即执行了 detach 导致最后才执行 attach 留有内存残留
		// hub attach event
		h.clients.attach(c)

		// 收到消息的时候，传递到OnMessage回调
		// 方法会阻塞进程，一直接受消息
		c.receiver(h.recv)

		// hub detach event
		h.clients.detach(c)

		// 回调应用程序
		h.callback.OnClose(&Pack{
			Hub:     h,
			Client:  c,
			Message: nil,
		})
	}(h,conn)
}
