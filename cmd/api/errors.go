package api

import (
	"log"
	"net/http"
)

func badRequest(w http.ResponseWriter, r *http.Request, err error) {

	log.Default().Println("forbidden", "method", r.Method, "path", r.URL.Path, "error")

    errorResponse(w, http.StatusBadRequest, err.Error())
}
