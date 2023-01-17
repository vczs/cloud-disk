package handler

import (
	"net/http"

	"cloud-disk/disk/internal/logic"
	"cloud-disk/disk/internal/svc"
	"cloud-disk/disk/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func CompletePartHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CompletePartRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewCompletePartLogic(r.Context(), svcCtx)
		resp, err := l.CompletePart(&req, r.Header.Get("Uid"))
		if err != nil {
			ResponseError(r.Context(), w, err)
		} else {
			Response(r.Context(), w, resp.Code, resp.Data)
		}
	}
}
