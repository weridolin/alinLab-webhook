package logic

import (
	"context"
	"fmt"

	"github.com/weridolin/alinLab-webhook/webhook/internal/svc"
	"github.com/weridolin/alinLab-webhook/webhook/internal/types"
	"github.com/weridolin/alinLab-webhook/webhook/models"

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

func (l *HistoryLogic) History(req *types.QueryHistoryRequest) (resp *types.QueryHistoryResponse, err error) {
	//count
	fmt.Println(req.Uuid, req.Page, req.Size)
	var count int64
	l.svcCtx.DB.Model(&models.ResourceCalledHistory{}).Where(map[string]string{"uuid": req.Uuid}).Count(&count)
	var records []*models.ResourceCalledHistory
	l.svcCtx.DB.Model(&models.ResourceCalledHistory{}).Where(map[string]string{"uuid": req.Uuid}).Offset((req.Page - 1) * req.Size).Limit(req.Size).Order("id desc").Find(&records)
	var resItems []types.HistoryItem
	for _, record := range records {
		resItems = append(resItems, types.HistoryItem{
			Header:      record.Header,
			Raw:         record.Raw,
			QueryParams: record.QueryParams,
			FormData:    record.FormData,
			Host:        record.Host,
			Method:      record.Method,
			UserID:      record.UserID,
			Uuid:        record.Uuid,
		})
	}
	return &types.QueryHistoryResponse{
		Items: resItems,
		Total: count,
	}, nil
}
