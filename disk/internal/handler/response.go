package handler

import (
	"cloud-disk/disk/internal/config"
	"context"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}

func Response(c context.Context, w http.ResponseWriter, code int, data interface{}) {
	body := Body{Code: code, Msg: config.GetMessage(code)}
	if code == config.SUCCESS {
		body.Data = data
		httpx.OkJsonCtx(c, w, body)
	} else {
		httpx.OkJsonCtx(c, w, body)
	}
}

func ResponseError(c context.Context, w http.ResponseWriter, err error) {
	httpx.OkJsonCtx(c, w, Body{Code: -1, Msg: err.Error()})
}
