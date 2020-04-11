package websocket

type HubStruct struct {
	// 连接计数器，当连接一旦握手成功，就立即计数+1
	connectionCounter int

	// Attached clients.
	clients 		map[*Client]bool
	// Attach requests from the clients.
	attach			chan *Client
	// Detach requests from the clients.
	detach			chan *Client
	// broadcast message to all clients
	broadcast 		chan []byte

	// callback's events
	onOpen				OnOpen
	onMessage 			OnMessage
	onClose 			OnClose

	// Users
	Users 				*Users
}