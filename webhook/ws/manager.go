package ws

import (
	"encoding/json"
	"fmt"

	"github.com/weridolin/alinLab-webhook/webhook/internal/config"
	"github.com/weridolin/alinLab-webhook/webhook/models"
	rabbitmq "github.com/weridolin/alinLab-webhook/webhook/mq"
	"github.com/weridolin/alinLab-webhook/webhook/utils"
)

type WebSocketManager struct {
	RabbitMQClient *rabbitmq.RabbitMQ
	Clients        map[string]*WsClient
	// accept a new websocket client register
	Register chan *WsClient
	// accept a new websocket client unregister
	Unregister chan *WsClient
	// accept receive message from rabbitmq
	MQMsgReceive chan []byte

	UUID string
}

func NewWebSocketManager(c config.Config) *WebSocketManager {
	// uuid := tools.GetUUID()
	return &WebSocketManager{
		Clients:      make(map[string]*WsClient),
		Register:     make(chan *WsClient),
		Unregister:   make(chan *WsClient),
		MQMsgReceive: make(chan []byte),
		UUID:         utils.GetUUID(),
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
				fmt.Println("unregister client -> ", client.Id)
				delete(manager.Clients, client.Id)
				close(client.Send)
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
				Client.Send <- jsonBytes
			} else {
				fmt.Println("websocket client not connect,ignore..")
			}

		}
	}
}
