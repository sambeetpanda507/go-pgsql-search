package routers

import (
	"encoding/json"
	"fmt"
	"net/http"
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

	mux.HandleFunc("/api/ping", func(w http.ResponseWriter, r *http.Request) {
		resp := struct {
			Message string `json:"message"`
		}{
			Message: "Pong",
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	})
}
