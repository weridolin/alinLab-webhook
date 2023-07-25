package main

import (
	"flag"
	"fmt"

	"github.com/weridolin/alinLab-webhook/webhook/internal/config"
	"github.com/weridolin/alinLab-webhook/webhook/internal/handler"
	"github.com/weridolin/alinLab-webhook/webhook/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/webhook-api.yaml", "the config file")

// func init() {
// 	// 读取.env环境变量
// 	err := godotenv.Load(".env")
// 	if err != nil {
// 		panic("读取.env环境变量失败")
// 	}
// }

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
