package api

import (
	"context"
	"errors"
	"image"
	"io"
	"math/rand"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
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

// GetUser
// @Summary Get-user
// @Description Responds with json
// @Tags User
// @Produce json
// @Success 200 {object} database.User
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Security ApiKeyAuth
// @Router /v1/user/ [get]
func (api *ApiService) GetByUsername(w http.ResponseWriter, r *http.Request) {

	username, err := getUsernameFromCtx(r.Context())

	if err != nil {
		internalServer(w, r, err)
		return
	}

	user, err := api.database.GetByUsername(r.Context(), username)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	writeJson(w, http.StatusOK, user)

}


// GetUser
// @Summary Search User
// @Description Responds with json
// @Tags User
// @Produce json
// @Success 200 {object} database.User
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Security ApiKeyAuth
// @Router /v1/user/search/{username} [get]
func (api *ApiService) GetByUsernameSearch(w http.ResponseWriter, r *http.Request) {


	username := chi.URLParam(r, "username")

	user, err := api.database.GetByUsername(r.Context(), username)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	writeJson(w, http.StatusOK, user)

}


// UploadProfilPic
// @Summary Upload Profil Pic
// @Description Responds with json
// @Tags User
// @Accept multipart/form-data
// @Produce json
// @Param img formData file true "upload send png jpeg and gif"
// @Success 200 {object} StandardResponse
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Security ApiKeyAuth
// @Router /v1/user/upload-profile-picture [post]
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

	_, _, err = image.Decode(file)

	if err != nil {
		badRequest(w, r, err)
		return
	}

	defer file.Close()

	currentTime := time.Now().UnixMilli()

	currentTimeString := strconv.Itoa(int(currentTime)) + filepath.Ext(fileHeader.Filename)

	destinationFile, err := os.Create("C:\\Users\\5557\\Desktop\\nkata_uploads\\profile\\" + currentTimeString)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	defer destinationFile.Close()

	_, err = file.Seek(0, io.SeekStart)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	_, err = io.Copy(destinationFile, file)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	url := "localhost:5557/v1/media/profiles/" + currentTimeString

	err = api.database.UpdateProfilePicUrl(ctx, username, url)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  http.StatusOK,
		Message: "user profile picture updated successfuly",
	}

	writeJson(w, http.StatusOK, s)
}

// LoadProfilPic
// @Summary Download Profil Pic
// @Description Responds with json
// @Tags Media
// @Param img_name path string true "file name"
// @Produce octet-stream
// @Success 200 {file} file
// @Failure 404 {object} errorslope
// @Router /v1/media/profiles/{img_name} [get]
func (api *ApiService) LoadProfilPic(w http.ResponseWriter, r *http.Request) {

	filename := chi.URLParam(r, "img_name")
	url := "C:\\Users\\5557\\Desktop\\nkata_uploads\\profile\\" + filename
	file, err := os.Open(url)

	if err != nil {
		notFound(w, r, errors.New("the system cannot find the file specified"))
		return
	}

	defer file.Close()

	w.Header().Set("Content-Disposition", "attachment; filename= "+filename)
	w.Header().Set("Content-Type", "application/octet-stream")

	http.ServeContent(w, r, filename, time.Time{}, file)
}

// Update
// @Summary Update user display_name or bio
// @Description Responds with json
// @Tags User
// @Param payload body UpdatePayload true "you can sent both or either"
// @Produce json
// @Accept json
// @Success 200 {object} StandardResponse
// @Failure 404 {object} errorslope
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Router /v1/user/update [put]
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

	err = api.database.UpdateUser(r.Context(), username, update.DisplayName, update.Bio)

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

// AddEmdil
// @Summary Add email to user. endpoint sends otp
// @Description Responds with json
// @Tags User
// @Param payload body EmailPayload true "valid email"
// @Produce json
// @Accept json
// @Success 200 {object} StandardResponse
// @Failure 404 {object} errorslope
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Router /v1/user/add-email [post]
func (api *ApiService) AddEmail(w http.ResponseWriter, r *http.Request) {

	username, err := getUsernameFromCtx(r.Context())

	if err != nil {
		internalServer(w, r, err)
		return
	}

	var emailp EmailPayload

	if err := readJson(w, r, &emailp); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	//send email with your smtp provider

	otpToken := rand.Intn(9000000) + 100000

	err = api.database.InsertOtp(ctx, username, emailp.Email, otpPurposeAddEmail, int64(otpToken))

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  http.StatusOK,
		Message: "otp " + strconv.Itoa(otpToken) + " sent to " + emailp.Email,
	}

	writeJson(w, http.StatusOK, s)

}

// AddEmailVerify
// @Summary Send otp sent to email
// @Description Responds with json
// @Tags User
// @Param payload body OtpPayload true "valid otp"
// @Produce json
// @Accept json
// @Success 200 {object} StandardResponse
// @Failure 404 {object} errorslope
// @Failure 400 {object} errorslope
// @Failure 500 {object} errorslope
// @Router /v1/user//add-email-verify [post]
func (api *ApiService) AddEmailVerify(w http.ResponseWriter, r *http.Request) {

	username, err := getUsernameFromCtx(r.Context())

	if err != nil {
		internalServer(w, r, err)
		return
	}

	var otpP OtpPayload

	if err := readJson(w, r, &otpP); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	otp, err := api.database.GetOtp(ctx, otpP.Otp)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	if otp.Purpose != otpPurposeAddEmail {
		unauthorized(w, r, errors.New("user does not have permission to perform this action"))
		return
	}

	if otp.Username != username {
		unauthorized(w, r, errors.New("user does not have permission to perform this action"))
		return
	}

	now := time.Now()
	exp, err := time.Parse(time.RFC1123Z, otp.Exp)

	if err != nil {
		internalServer(w, r, err)
		return
	}

	if exp.Before(now) {
		unauthorized(w, r, errors.New("otp expired"))
		return
	}

	err = api.database.UpdateUserEmail(ctx, username, otp.Email)
	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  http.StatusOK,
		Message: "Email updated succesfully",
	}

	writeJson(w, http.StatusOK, s)

}
