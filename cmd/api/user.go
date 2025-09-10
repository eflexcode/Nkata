package api

import (
	"net/http"
)



func (api *ApiService) GetByID(w http.ResponseWriter, r *http.Request) {

}

func (api *ApiService) GetByUsername(w http.ResponseWriter, r *http.Request) {

	type j struct{
		Message string `json:"message"`
	}

	writeJson(w, http.StatusOK, j{Message: "auth works"})

}

