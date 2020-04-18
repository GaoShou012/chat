package websocket

/*
消息格式
c 消息源的终端，message消息
*/
type Pack struct {
	Hub 		*hub
	Client 		*client
	Message 	[]byte
}

type EventCallback struct {
	OnOpen func(pack *Pack) bool
	OnMessage func(pack *Pack)
	OnClose func(pack *Pack)
}