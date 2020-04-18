package websocket

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"os"
	"sync"
)

type Settings struct {
	// 监听 地址：端口
	HostAddrPort 			string
	// 路径
	BindPath 				string
	// 消息处理协程数量
	// 在收到客户端消息的时候，同时多少个协程在处理这些消息
	ProcNum					int

	// 连接设置
	// 读：Client To Server
	// 写：Server To Client
	Connection struct {
		ReadBufferSize		int
		WriteBufferSize		int
	}
}

func NewSetting() *Settings {
	tmp := &Settings{
		HostAddrPort: "0.0.0.0:8080",
		BindPath:     "/",
		ProcNum:      100,
		Connection: struct {
			ReadBufferSize  int
			WriteBufferSize int
		}{ReadBufferSize:1024,WriteBufferSize:1024},
	}
	return tmp
}

func Run(settings *Settings,callback *EventCallback) {
	hub := &hub{
		settings:settings,
		callback:callback,
		upgrader:&websocket.Upgrader{
			//HandshakeTimeout:  0,
			ReadBufferSize:    settings.Connection.ReadBufferSize,
			WriteBufferSize:   settings.Connection.WriteBufferSize,
			//WriteBufferPool:   nil,
			//Subprotocols:      nil,
			//Error:             nil,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			//EnableCompression: false,
		},
		clients:&clients{mutex:&sync.Mutex{},list:make(map[*client]bool)},
	}
	hub.receiver()

	http.HandleFunc(hub.settings.BindPath,hub.handle)
	err := http.ListenAndServe(hub.settings.HostAddrPort,nil)
	if err != nil {
		log.Println("监听服务失败!")
		log.Printf("%s %s\n",hub.settings.HostAddrPort,err)
		os.Exit(1)
	}
}
