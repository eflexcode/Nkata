package api

import (
	"context"
	"errors"
	"image"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"

	_ "image/gif" 
	_ "image/jpeg" 
	_ "image/png"  
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

func (api *ApiService) UploadProfilPic(w http.ResponseWriter, r *http.Request) {

	ctx := r.Context()
	username, err := getUsernameFromCtx(ctx)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	r.ParseMultipartForm(30 << 20)

	file, fileHeader, err := r.FormFile("img")

	if err != nil {
		badRequest(w, r, err)
		return
	}

	_,_, err = image.Decode(file)

	if err != nil {
		badRequest(w, r, err)
		return
	}

	defer file.Close()

	destinationFile, err := os.Create("C:\\Users\\5557\\Desktop\\nkata_uploads\\profile" + fileHeader.Filename)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	defer destinationFile.Close()

	_, err = io.Copy(destinationFile, file)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	url := "localhost:5557/v1/media/profiles/" + fileHeader.Filename

	api.userRpo.UpdateProfilePicUrl(ctx, username, url)

	s := StandardResponse{
		Status:  http.StatusOK,
		Message: "user profile picture updated successfuly",
	}

	writeJson(w, http.StatusOK, s)
}

func (api *ApiService) LoadProfilPic(w http.ResponseWriter, r *http.Request) {

	filename := chi.URLParam(r, "img_name")
	url := "C:\\Users\\5557\\Desktop\\nkata_uploads\\profile" + filename
	file, err := os.Open(url)

	if err != nil {
		notFound(w, r, err)
		return
	}
	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment: "+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeContent(w, r, filename, time.Now(), file)
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

	err = api.userRpo.Update(r.Context(), username, update.DisplayName, update.Bio)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  http.StatusOK,
		Message: "user details updated successfuly",
	}

	writeJson(w, http.StatusOK, s)

}
