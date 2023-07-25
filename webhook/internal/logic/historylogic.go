package logic

import (
	"context"

	"github.com/weridolin/alinLab-webhook/webhook/internal/svc"
	"github.com/weridolin/alinLab-webhook/webhook/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type HistoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewHistoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *HistoryLogic {
	return &HistoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *HistoryLogic) History(req *types.PaginationParams) (resp *types.QueryHistoryResponse, err error) {
	// todo: add your logic here and delete this line

	return
}
