package logic

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/weridolin/alinLab-webhook/webhook/internal/svc"
	"github.com/weridolin/alinLab-webhook/webhook/internal/types"
	"github.com/weridolin/alinLab-webhook/webhook/models"

	"github.com/zeromicro/go-zero/core/logx"
)

type WebhookCalledLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewWebhookCalledLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WebhookCalledLogic {
	return &WebhookCalledLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *WebhookCalledLogic) WebhookCalled(req *types.Request, r *http.Request) (resp *types.Response, err error) {
	// userID := l.ctx.Value("id")

	fmt.Println("host:", r.RemoteAddr)
	header := l.ParseHeader(r)
	fmt.Println("header:", header)
	queryParams := l.ParseQueryParams(r)
	fmt.Println("queryParams:", queryParams)
	// formData := l.ParseFormData(r)
	// fmt.Println("formData:", formData)
	raw := l.ParseRaw(r)
	fmt.Println("raw:", raw)
	l.ParseRaw(r)
	newHistory := models.ResourceCalledHistory{
		Uuid:        req.Uuid,
		Header:      header,
		Raw:         raw,
		QueryParams: queryParams,
		Host:        r.RemoteAddr,
		Method:      r.Method,
		// UserID:      userID.(int),
	}
	err = models.CreateNewResourceCalledHistory(l.svcCtx.DB, &newHistory)
	if err != nil {
		fmt.Println("create history error:", err)
	}
	return
}

func (l *WebhookCalledLogic) ParseHeader(r *http.Request) (header map[string]string) {
	header = make(map[string]string)
	for k, v := range r.Header {
		header[k] = strings.Join(v, ",")
	}
	// fmt.Println("header:", r.Header)
	return header
}

func (l *WebhookCalledLogic) ParseQueryParams(r *http.Request) (queryParams map[string]string) {
	queryParams = make(map[string]string)
	for k, v := range r.URL.Query() {
		queryParams[k] = strings.Join(v, ",")
	}
	return queryParams
}

func (l *WebhookCalledLogic) ParseFormData(r *http.Request) (formData map[string]string) {
	formData = make(map[string]string)
	// r.ParseForm()
	r.ParseMultipartForm(32 << 20) //postform 放在这里处理了
	fmt.Println("r.Form:", r.Form, "r.PostForm:", r.PostForm)
	if len(r.Form) > 0 {
		for k, v := range r.Form {
			formData[k] = strings.Join(v, ",")
		}
		return formData
	} else if len(r.PostForm) > 0 {
		for k, v := range r.PostForm {
			formData[k] = strings.Join(v, ",")
		}
		return formData
	}
	return formData
}

func (l *WebhookCalledLogic) ParseRaw(r *http.Request) string {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return ""
	}
	return string(data)
}
