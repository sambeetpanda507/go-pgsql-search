package routers

import (
	"net/http"

	"github.com/sambeetpanda507/advance-search/controllers"
	"github.com/sambeetpanda507/advance-search/middlewares"
)

func NewsRoutes(mux *http.ServeMux, c controllers.NewController) {
	mux.Handle("/api/news/from-file", middlewares.CORS(http.HandlerFunc(c.GetNewsFromFile)))
}
