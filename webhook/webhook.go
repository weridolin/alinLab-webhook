package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"time"

	"github.com/weridolin/alinLab-webhook/webhook/internal/config"
	"github.com/weridolin/alinLab-webhook/webhook/internal/handler"
	"github.com/weridolin/alinLab-webhook/webhook/internal/svc"
	"github.com/weridolin/alinLab-webhook/webhook/utils"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var configFile = flag.String("f", "etc/webhook-api.yaml", "the config file")

func GeLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println("get net interface address failed, err = ", err.Error())
		return ""
	}
	for _, addr := range addrs {
		if ip, ok := addr.(*net.IPNet); ok && !ip.IP.IsLoopback() {
			if ip.IP.To4() != nil {
				return ip.IP.String()
			}
		}
	}
	return ""
}

// 将服务注册到etcd上
func RegisterServiceToETCD(conf *config.Config) {

	client, err := clientv3.New(clientv3.Config{
		Endpoints:   conf.Etcd.Hosts,
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}

	kv := clientv3.NewKV(client)
	lease := clientv3.NewLease(client)
	var curLeaseId clientv3.LeaseID = 0

	listenAddr := fmt.Sprintf("%s:%d", GeLocalIP(), conf.Port)
	// uuid := tools.GetUUID()
	key := fmt.Sprintf("%s/%s", conf.Etcd.Key.Rest, conf.UUID)
	wsKey := fmt.Sprintf("%s/%s", conf.Etcd.Key.Ws, conf.UUID)
	socketioKey := fmt.Sprintf("%s/%s", conf.Etcd.Key.SocketIO, conf.UUID)

	for {
		if curLeaseId == 0 {
			leaseResp, err := lease.Grant(context.TODO(), 10)
			if err != nil {
				panic(err)
			}
			fmt.Println("register alinlab webhook rest api ", key, " -> ", listenAddr)
			if _, err := kv.Put(context.TODO(), key, listenAddr, clientv3.WithLease(leaseResp.ID)); err != nil {
				panic(err)
			}
			fmt.Println("register alinlab webhook websocket api ", wsKey, " -> ", listenAddr)
			if _, err := kv.Put(context.TODO(), wsKey, listenAddr, clientv3.WithLease(leaseResp.ID)); err != nil {
				panic(err)
			}
			fmt.Println("register alinlab webhook socketio api ", socketioKey, " -> ", listenAddr)
			if _, err := kv.Put(context.TODO(), socketioKey, listenAddr, clientv3.WithLease(leaseResp.ID)); err != nil {
				panic(err)
			}
			curLeaseId = leaseResp.ID
		} else {
			// 续约租约，如果租约已经过期将curLeaseId复位到0重新走创建租约的逻辑
			if _, err := lease.KeepAliveOnce(context.TODO(), curLeaseId); err == rpctypes.ErrLeaseNotFound {
				fmt.Println("lease not found, reset lease id to 0")
				curLeaseId = 0
				continue
			}
		}
		time.Sleep(time.Duration(1) * time.Second)
	}
}

func main() {
	flag.Parse()

	var c config.Config
	option := conf.UseEnv()
	conf.MustLoad(*configFile, &c, option)
	c.UUID = utils.GetUUID()
	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)

	// WEBSOCKET
	// go ctx.WebsocketManager.Start()

	// // socket.io
	// go func() {
	// 	if err := ctx.SocketIOServer.Serve(); err != nil {
	// 		fmt.Println("socketio listen error: %s\n", err)
	// 	}
	// }()

	// defer ctx.SocketIOServer.Close()
	// go ctx.SocketIOManager.Start()

	// start rabbitmq consumer
	// RabbitMQClient := rabbitmq.NewRabbitMQTopic(
	// 	c.RabbitMq.BroadcastExchange,
	// 	c.RabbitMq.BroadcastTopic,
	// 	c.RabbitMq.MQURI,
	// )
	// go func() {
	// 	RabbitMQClient.SubscribeTopic(func(msg []byte) {
	// 		ctx.WebsocketManager.MQMsgReceive <- msg
	// 		ctx.SocketIOManager.MQMsgReceive <- msg
	// 	})
	// }()

	// HTTP
	if env := os.Getenv("K8S"); env != "1" {
		go RegisterServiceToETCD(&c)
	}

	server.Start()
}
