package app

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type errorResponse struct {
	Error string
}

func decodeJSONRequest(v interface{}, r *http.Request) error {
	dec := json.NewDecoder(r.Body)
	err := dec.Decode(v)
	return err
}

func jsonError(w http.ResponseWriter, msg string, status int) {
	b, err := json.Marshal(errorResponse{msg})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	io.Copy(w, bytes.NewReader(b))
}

func jsonResponse(w http.ResponseWriter, v interface{}) {
	b, err := json.Marshal(v)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	io.Copy(w, bytes.NewReader(b))
}
