package response

import (
	"context"
	"github.com/sphinx-camfield/utils/rid"
	"net/http"
)

func Error(w http.ResponseWriter, ctx context.Context) {
	Json(w, http.StatusInternalServerError,
		map[string]string{
			"code":       "ERR_SERVER_ERROR",
			"message":    "Server error",
			"request_id": ctx.Value("trace").(*rid.Rid).String(),
		})
}
