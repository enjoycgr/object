package handler

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"object/app/address/api/internal/logic"
	"object/app/address/api/internal/svc"
)

func AddressTreeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		l := logic.NewAddressLogic(r.Context(), svcCtx)
		resp, err := l.Address()
		if err != nil {
			httpx.Error(w, err)
		} else {
			httpx.OkJson(w, resp)
		}
	}
}
