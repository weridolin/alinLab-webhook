package ws

import "fmt"

type WebSocketManager struct {
	Clients map[string]*WsClient
	// accept a new websocket client register
	Register chan *WsClient
	// accept a new websocket client unregister
	Unregister chan *WsClient
}

func NewWebSocketManager() *WebSocketManager {
	return &WebSocketManager{
		Clients:    make(map[string]*WsClient),
		Register:   make(chan *WsClient),
		Unregister: make(chan *WsClient),
	}
}

func (manager *WebSocketManager) Start() {
	for {
		select {
		case client := <-manager.Register:
			fmt.Println("register client -> ", client.Id)
			manager.Clients[client.Id] = client
		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client.Id]; ok {
				delete(manager.Clients, client.Id)
				close(client.Send)
			}
		}
	}
}
