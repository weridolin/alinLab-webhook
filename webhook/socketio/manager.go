package socketio

import "fmt"

type SocketIOConnectionManager struct {
	Clients    map[string]*SocketIOClient
	Register   chan *SocketIOClient
	Unregister chan *SocketIOClient
}

func NewSocketIOConnectionManager() *SocketIOConnectionManager {
	return &SocketIOConnectionManager{
		Clients:    make(map[string]*SocketIOClient),
		Register:   make(chan *SocketIOClient),
		Unregister: make(chan *SocketIOClient),
	}
}

func (manager *SocketIOConnectionManager) Start() {
	for {
		select {
		case client := <-manager.Register:
			fmt.Println("socket io receive a new client ", client.Id)
			manager.Clients[client.Id] = client
		case client := <-manager.Unregister:
			if _, ok := manager.Clients[client.Id]; ok {
				fmt.Println("socket io disconnect a new client", client.Id)
				manager.Clients[client.Id].Conn.Close()
				delete(manager.Clients, client.Id)
			}
		}
	}
}
