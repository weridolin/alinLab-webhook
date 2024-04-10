package logic

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"

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

	var formData = make(map[string]string)
	urlEncodeForm, err := l.ParsePostForm(r)
	if err != nil {
		fmt.Println("parse post form error:", err)
	} else {
		for k, v := range urlEncodeForm {
			formData[k] = strings.Join(v, ",")
		}
	}
	err = r.ParseMultipartForm(32 << 20)
	if err != nil {
		fmt.Println("parse multipart form error:", err)
	} else {
		for k, v := range r.MultipartForm.Value {
			formData[k] = strings.Join(v, ",")
		}
	}
	fmt.Println("formData:", formData)
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
		FormData:    formData,
		// UserID:      0,
		// UpdatedAt:  time.,
	}
	err = models.CreateNewResourceCalledHistory(l.svcCtx.DB, &newHistory)
	if err != nil {
		fmt.Println("create history error:", err)
	}

	now := time.Now()
	// 延后一天
	tomorrow := now.Add(24 * time.Hour)
	// 转换为 Unix 时间戳（秒）
	tomorrowTimestamp := tomorrow.Unix()

	notifyMsg := struct {
		FromApp     string `json:"from_app"`
		WebsocketId string `json:"websocket_id"`
		exp         int64  `json:"exp"`
		models.ResourceCalledHistory
	}{
		FromApp:               "site.alinlab.webhook",
		WebsocketId:           req.Uuid,
		exp:                   tomorrowTimestamp,
		ResourceCalledHistory: newHistory,
	}
	// 推送广播消息到 mq
	jsonBytes, err := json.Marshal(notifyMsg)
	if err != nil {
		fmt.Println("JSON编码失败:", err)
	}

	topic := fmt.Sprintf("%s.%s", l.svcCtx.Config.RabbitMq.BroadcastTopic, req.Uuid)
	l.svcCtx.RabbitMqPublisher.PublishTopic(topic, jsonBytes)

	// 如果socket io已经建立连接，则发送这次调用结果
	// client, exist := l.svcCtx.SocketIOManager.Clients[req.Uuid]
	// if exist {
	// 	// 将结构体转换为JSON格式的字节切片
	// 	jsonBytes, err := json.Marshal(newHistory)
	// 	if err != nil {
	// 		fmt.Println("JSON编码失败:", err)
	// 	}
	// 	// 将JSON格式的字节切片转换为字符串
	// 	fmt.Println("send call record to socket io client", newHistory.Uuid)
	// 	jsonString := string(jsonBytes)
	// 	client.Conn.Emit("msg", jsonString)
	// } else {
	// 	fmt.Println("socket io client  not connect,ignore..")
	// }

	// 如果websocket已经建立连接，则发送这次调用结果
	// wsClient, exist := l.svcCtx.WebsocketManager.Clients[req.Uuid]
	// if exist {
	// 	// 将结构体转换为JSON格式的字节切片
	// 	jsonBytes, err := json.Marshal(newHistory)
	// 	if err != nil {
	// 		fmt.Println("JSON编码失败:", err)
	// 	}
	// 	// 将JSON格式的字节切片转换为字符串
	// 	fmt.Println("send call record to socket io client", newHistory.Uuid)
	// 	wsClient.Send <- jsonBytes
	// } else {
	// 	fmt.Println("websocket client not connect,ignore..")
	// }

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

type maxBytesReader struct {
	w   http.ResponseWriter
	r   io.ReadCloser // underlying reader
	n   int64         // max bytes remaining
	err error         // sticky error
}

func (l *maxBytesReader) Read(p []byte) (n int, err error) {
	if l.err != nil {
		return 0, l.err
	}
	if len(p) == 0 {
		return 0, nil
	}
	// If they asked for a 32KB byte read but only 5 bytes are
	// remaining, no need to read 32KB. 6 bytes will answer the
	// question of the whether we hit the limit or go past it.
	if int64(len(p)) > l.n+1 {
		p = p[:l.n+1]
	}
	n, err = l.r.Read(p)

	if int64(n) <= l.n {
		l.n -= int64(n)
		l.err = err
		return n, err
	}

	n = int(l.n)
	l.n = 0

	// The server code and client code both use
	// maxBytesReader. This "requestTooLarge" check is
	// only used by the server code. To prevent binaries
	// which only using the HTTP Client code (such as
	// cmd/go) from also linking in the HTTP server, don't
	// use a static type assertion to the server
	// "*response" type. Check this interface instead:
	type requestTooLarger interface {
		requestTooLarge()
	}
	if res, ok := l.w.(requestTooLarger); ok {
		res.requestTooLarge()
	}
	l.err = errors.New("http: request body too large")
	return n, l.err
}

func (l *maxBytesReader) Close() error {
	return l.r.Close()
}

func (l *WebhookCalledLogic) ParsePostForm(r *http.Request) (vs url.Values, err error) {
	if r.Body == nil {
		err = errors.New("missing form body")
		return
	}
	ct := r.Header.Get("Content-Type")
	// RFC 7231, section 3.1.1.5 - empty type
	//   MAY be treated as application/octet-stream
	if ct == "" {
		ct = "application/octet-stream"
	}
	ct, _, err = mime.ParseMediaType(ct)
	switch {
	case ct == "application/x-www-form-urlencoded":
		// fmt.Println(">>>>> application/x-www-form-urlencoded")
		var reader io.Reader = r.Body
		maxFormSize := int64(1<<63 - 1)
		if _, ok := r.Body.(*maxBytesReader); !ok {
			maxFormSize = int64(10 << 20) // 10 MB is a lot of text.
			reader = io.LimitReader(r.Body, maxFormSize+1)
		}
		b, e := io.ReadAll(reader)
		if e != nil {
			if err == nil {
				err = e
			}
			break
		}
		if int64(len(b)) > maxFormSize {
			err = errors.New("http: POST too large")
			return
		}
		vs, e = url.ParseQuery(string(b))
		if err == nil {
			err = e
		}
	case ct == "multipart/form-data":
		// handled by ParseMultipartForm (which is calling us, or should be)
		// TODO(bradfitz): there are too many possible
		// orders to call too many functions here.
		// Clean this up and write more tests.
		// request_test.go contains the start of this,
		// in TestParseMultipartFormOrder and others.
	}
	return
}
