package logic

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/weridolin/alinLab-webhook/webhook/internal/svc"
	"github.com/weridolin/alinLab-webhook/webhook/internal/types"
	"github.com/weridolin/alinLab-webhook/webhook/utils"

	"github.com/zeromicro/go-zero/core/logx"
)

type RegisterWebsocketLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRegisterWebsocketLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterWebsocketLogic {
	return &RegisterWebsocketLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterWebsocketLogic) RegisterWebsocket(req *types.RegisterWebsocketRequest, r *http.Request) (resp *types.RegisterWebsocketResponse, err error) {
	// queryParams := make(map[string]string)
	// for k, v := range r.URL.Query() {
	// 	queryParams[k] = strings.Join(v, ",")
	// }
	fmt.Println("webhook register websocket uuid -> ", req.Uuid)
	if req.Uuid == "" {
		return nil, errors.New("uuid is required")
	}
	// 从环境变量获取jwt key
	jwtKey := os.Getenv("JWT_KEY")
	if jwtKey == "" {
		jwtKey = "DEBUGJWTKEY"
	}
	token := utils.GenToken(req.Uuid, jwtKey)
	websocket_uri := "wss://" + r.Host + "/ws-endpoint/api/v1/?token=" + token
	return &types.RegisterWebsocketResponse{
		WebsocketUri: websocket_uri,
	}, nil
}
