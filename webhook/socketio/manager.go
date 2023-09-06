package socketio

import (
	"encoding/json"
	"fmt"

	"github.com/weridolin/alinLab-webhook/webhook/models"
	"github.com/weridolin/alinLab-webhook/webhook/utils"
)

type SocketIOConnectionManager struct {
	Clients      map[string]*SocketIOClient
	Register     chan *SocketIOClient
	Unregister   chan *SocketIOClient
	MQMsgReceive chan []byte
	UUID         string
}

func NewSocketIOConnectionManager() *SocketIOConnectionManager {
	return &SocketIOConnectionManager{
		Clients:      make(map[string]*SocketIOClient),
		Register:     make(chan *SocketIOClient),
		Unregister:   make(chan *SocketIOClient),
		MQMsgReceive: make(chan []byte),
		UUID:         utils.GetUUID(),
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
		case msg := <-manager.MQMsgReceive:
			fmt.Println("ws manager ", manager.UUID, " receive msg from rabbitmq -> ", string(msg))
			var notifyMsg = &models.NotifyMessage{}
			err := json.Unmarshal(msg, notifyMsg)
			if err != nil {
				fmt.Println("rabbit mq message json unmarshal error:", err)
			}
			Client, exist := manager.Clients[notifyMsg.ToUUID]
			if exist {
				jsonBytes, err := json.Marshal(notifyMsg.ResourceCalledHistory)
				if err != nil {
					fmt.Println("call history serialize error -> :", err)
				}
				Client.Conn.Emit("msg", string(jsonBytes))
			} else {
				fmt.Println("socketIO client not connect,ignore..")
			}

		}
	}
}
