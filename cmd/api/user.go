package api

import (
	"context"
	"errors"
	"net/http"
)

type UpdatePayload struct {
	DisplayName string `json:"display_name"`
	Bio         string `json:"bio"`
}

func getUsernameFromCtx(ctx context.Context) (string, error) {

	username, ok := ctx.Value("user").(string)

	if !ok {
		err := errors.New("no username found in token")
		return "nil", err
	}

	return username, nil
}

func (api *ApiService) GetByID(w http.ResponseWriter, r *http.Request) {

}

func (api *ApiService) GetByUsername(w http.ResponseWriter, r *http.Request) {

	username, err := getUsernameFromCtx(r.Context())

	if err != nil {
		internalServer(w, r, err)
		return
	}

	user, err := api.userRpo.GetByUsername(r.Context(), username)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	writeJson(w, http.StatusOK, user)

}

func (api *ApiService) Update(w http.ResponseWriter, r *http.Request) {

	var update UpdatePayload

	if err := readJson(w, r, &update); err != nil {
		badRequest(w, r, err)
		return
	}

	username, err := getUsernameFromCtx(r.Context())

	if err != nil {
		internalServer(w, r, err)
		return
	}

	err = api.userRpo.Update(r.Context(),username,update.DisplayName,update.Bio)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status: http.StatusOK,
		Message: "user details updated successfuly",
	}

	writeJson(w,http.StatusOK,s)

}
