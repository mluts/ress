package main

import (
	"io"
	"net/http"
)

func main() {
	http.HandleFunc("/feeds", func(w http.ResponseWriter, h *http.Request) {
		io.WriteString(w, "FEEDS!")
		w.WriteHeader(200)
	})
	http.ListenAndServe(":8080", nil)
}
