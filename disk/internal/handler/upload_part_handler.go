package handler

import (
	"net/http"

	"cloud-disk/disk/internal/logic"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func UploadPartHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.UploadPartRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewUploadPartLogic(r.Context(), svcCtx)
		resp, err := l.UploadPart(&req, r)
		if err != nil {
			ResponseError(r.Context(), w, err)
		} else {
			Response(r.Context(), w, resp.Code, nil)
		}
	}
}
