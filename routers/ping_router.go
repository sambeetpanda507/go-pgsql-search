package routers

import (
	"fmt"
	"net/http"

	"github.com/sambeetpanda507/advance-search/controllers"
	"github.com/sambeetpanda507/advance-search/middlewares"
)

func Ping(mux *http.ServeMux) {
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ok")
	})

	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Ok")
	})

	mux.Handle("/api/ping", middlewares.CORS(http.HandlerFunc(controllers.PingHandler)))
}
