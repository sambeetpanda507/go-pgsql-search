package routers

import (
	"net/http"

	"github.com/sambeetpanda507/advance-search/controllers"
	"github.com/sambeetpanda507/advance-search/middlewares"
)

func AIRouter(mux *http.ServeMux, c controllers.AIController) {
	mux.Handle("POST /api/ai/embedding", middlewares.CORS(http.HandlerFunc(c.HandleEmbedding)))
}
