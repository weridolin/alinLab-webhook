package handler

import (
	"net/http"

	"github.com/weridolin/alinLab-webhook/webhook/internal/logic"
	"github.com/weridolin/alinLab-webhook/webhook/internal/svc"
	"github.com/weridolin/alinLab-webhook/webhook/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func WebhookCalledHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.Request
		if err := httpx.ParsePath(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewWebhookCalledLogic(r.Context(), svcCtx)
		resp, err := l.WebhookCalled(&req, r)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
