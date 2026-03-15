package httpserver

import (
	"net/http"

	"example.com/example-project/internal/app"
)

func NewMux() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", helloHandler)

	return mux
}

func helloHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(app.Greeting("world") + "\n"))
}
