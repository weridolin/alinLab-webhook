package svc

import (
	s "github.com/googollee/go-socket.io"
	"github.com/weridolin/alinLab-webhook/webhook/internal/config"
	"github.com/weridolin/alinLab-webhook/webhook/models"
	socketio "github.com/weridolin/alinLab-webhook/webhook/socketio"
	"github.com/weridolin/alinLab-webhook/webhook/ws"
	"github.com/zeromicro/go-zero/core/logx"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ServiceContext struct {
	Config           config.Config
	DB               *gorm.DB
	WebsocketManager *ws.WebSocketManager
	SocketIOServer   *s.Server
	SocketIOManager  *socketio.SocketIOConnectionManager
}

func NewServiceContext(c config.Config) *ServiceContext {
	db, err := gorm.Open(mysql.Open(c.DBUri), &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			// TablePrefix:   "auth_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	})
	if err != nil {
		logx.Error(err)
	}
	db.AutoMigrate(&models.ResourceCalledHistory{})
	var socketIOManager = socketio.NewSocketIOConnectionManager()
	return &ServiceContext{
		Config:           c,
		DB:               db,
		WebsocketManager: ws.NewWebSocketManager(),
		SocketIOServer:   socketio.NewSocketIoServer(socketIOManager),
		SocketIOManager:  socketIOManager,
	}
}
