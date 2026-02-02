package routers

import (
	"net/http"

	"github.com/sambeetpanda507/advance-search/controllers"
	"github.com/sambeetpanda507/advance-search/middlewares"
)

func NewsRoutes(mux *http.ServeMux, c controllers.NewsController) {
	mux.Handle("/api/news/from-file", middlewares.CORS(http.HandlerFunc(c.GetNewsFromFile)))
	mux.Handle("/api/news", middlewares.CORS(http.HandlerFunc(c.GetAllNews)))
	mux.Handle("/api/news/fill-embedding", middlewares.CORS(http.HandlerFunc(c.HandleFillEmbedding)))
}
