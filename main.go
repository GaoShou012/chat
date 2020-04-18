package main

import (
	"fmt"
	"github.com/src/websocket"
)


func main() {
	fmt.Println("test websocket")

	settings := websocket.NewSetting()
	settings.HostAddrPort = "0.0.0.0:8081"
	settings.BindPath = "/"
	callback := &websocket.EventCallback {
		OnOpen: func(pack *websocket.Pack) bool {
			fmt.Println("on Open")
			return true
		},
		OnMessage: func(pack *websocket.Pack) {
			fmt.Println("on message")
			pack.Client.Send(pack.Message)
		},
		OnClose: func(pack *websocket.Pack) {
			fmt.Println("on close")
		},
	}
	websocket.Run(settings,callback)

	fmt.Println("exit")
}