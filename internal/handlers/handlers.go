package handlers

import "net/http"

func HelloHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = w.Write([]byte("Hello, World"))
}
