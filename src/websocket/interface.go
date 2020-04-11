package websocket

import "net/http"

type OnOpen func(c *Client,w http.ResponseWriter, r *http.Request) bool
type OnMessage func(c *Client,message []byte)
type OnClose func(c *Client)

