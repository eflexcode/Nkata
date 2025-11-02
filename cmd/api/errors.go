package api

import (
	"log"
	"net/http"
)

func badRequest(w http.ResponseWriter, r *http.Request, err error) {

	log.Default().Println("badrequet", "method", r.Method, "path", r.URL.Path, "error")

    errorResponse(w, http.StatusBadRequest, err.Error())
}

func notFound(w http.ResponseWriter, r *http.Request, err error) {

	log.Default().Println("badrequet", "method", r.Method, "path", r.URL.Path, "error")

    errorResponse(w, http.StatusNotFound, err.Error())
}

func internalServer(w http.ResponseWriter, r *http.Request, err error) {

	log.Default().Println("internalServer", "method", r.Method, "path", r.URL.Path, "error")

    errorResponse(w, http.StatusInternalServerError, err.Error())
}

func unauthorized(w http.ResponseWriter, r *http.Request, err error) {

	log.Default().Println("Unauthorized", "method", r.Method, "path", r.URL.Path, "error")

    errorResponse(w, http.StatusUnauthorized, err.Error())
}

func conflict(w http.ResponseWriter, r *http.Request, err error) {

	log.Default().Println("Duplicate", "method", r.Method, "path", r.URL.Path, "error")

    errorResponse(w, http.StatusConflict, err.Error())
}

func tooManyRequest(w http.ResponseWriter, r *http.Request, err error) {

	log.Default().Println("TooManyRequests", "method", r.Method, "path", r.URL.Path, "error")

    errorResponse(w, http.StatusTooManyRequests, err.Error())
}