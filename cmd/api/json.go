package api

import (
	"encoding/json"
	"net/http"
)

type StandardResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type errorslope struct {
		Error  string `json:"error"`
		Status int    `json:"status"`
	}

func writeJson(w http.ResponseWriter, status int, data any) error {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	return json.NewEncoder(w).Encode(data)
}

func readJson(w http.ResponseWriter, r *http.Request, data any) error {

	maxPayLoadSize := 1_048_578
	r.Body = http.MaxBytesReader(w, r.Body, int64(maxPayLoadSize))

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	return decoder.Decode(data)
}

func errorResponse(w http.ResponseWriter, status int, message string) error {
	return writeJson(w, status, &errorslope{Error: message, Status: status})
}
