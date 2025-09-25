package api

import (
	"database/sql"
	"errors"
	"main/database"
	"main/internal/evn"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type RegisterUserPayload struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

type LoginUsernamePayload struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginEmailePayload struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OptPayload struct {
	Email string `json:"email"`
	Otp   int    `json:"otp"`
}

type UsernamePayload struct {
	Username string `json:"username"`
}

type BoolPayload struct {
	Exist bool `json:"exist"`
}

type JwtJson struct {
	Token string `json:"token"`
}

type EmailPayload struct {
	Email string `json:"email"`
}

func (apiService *ApiService) RegisterUser(w http.ResponseWriter, r *http.Request) {

	var payload RegisterUserPayload

	if err := readJson(w, r, &payload); err != nil {
		badRequest(w, r, err)
		return
	}

	if len(payload.Username) > 10 {
		err := errors.New("username max character is 10")
		badRequest(w, r, err)
		return
	}

	if len(payload.Password) < 8 {
		err := errors.New("password min character is 8")
		badRequest(w, r, err)
		return
	}

	user := &database.User{
		Username:    payload.Username,
		DisplayName: payload.DisplayName,
		Password:    payload.Password,
	}

	ctx := r.Context()

	err := apiService.userRpo.CreateUser(ctx, user)

	if err != nil {

		if err.Error() == `pq: duplicate key value violates unique constraint "users_username_key"` {
			err := errors.New("user with username " + user.Username + " already exist")
			conflict(w, r, err)
		} else {
			internalServer(w, r, err)
		}

		return
	}

	s := StandardResponse{
		Status:  http.StatusCreated,
		Message: "User account created successfully procced to signin.",
	}

	writeJson(w, http.StatusCreated, s)
}

func (apiService *ApiService) SignInUsername(w http.ResponseWriter, r *http.Request) {

	var payload LoginUsernamePayload

	if err := readJson(w, r, &payload); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	user, err := apiService.userRpo.GetByUsername(ctx, payload.Username)

	if err != nil {
		if err == sql.ErrNoRows {
			unauthorized(w, r, err)
			return
		}
		internalServer(w, r, errors.New("somthing went wrong"))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))

	if err != nil {
		unauthorized(w, r, err)
		return
	}

	claims := jwt.MapClaims{

		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 48).Unix(),
	}

	var secret_words string = "A request for a long text message: Search results showIf this is your intent, please clarify the context and what you want the text to be about."

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	secret := evn.GetString(secret_words, "JWT_SERECT")

	tokenString, err := token.SignedString([]byte(secret))

	if err != nil {
		err := errors.New("failed to generate token")
		internalServer(w, r, err)
		return
	}

	tokenResponse := JwtJson{Token: tokenString}

	writeJson(w, http.StatusAccepted, tokenResponse)
}

func (api *ApiService) AddEmail(w http.ResponseWriter, r *http.Request) {

	var emailp EmailPayload

	if err := readJson(w, r, &emailp); err != nil {
		badRequest(w, r, err)
		return
	}

	// ctx = r.Context()

}

// sign in with email must verify email
func (api *ApiService) SignInEmail(w http.ResponseWriter, r *http.Request) {

}

func (api *ApiService) VerifySignInEmailOtpEmail(w http.ResponseWriter, r *http.Request) {

}

func (api *ApiService) SendResetPasswordOtpEmail(w http.ResponseWriter, r *http.Request) {

}

func (api *ApiService) VerifyResetPasswordOtpEmail(w http.ResponseWriter, r *http.Request) {

}

func (api *ApiService) ResetPassword(w http.ResponseWriter, r *http.Request) {

}

func (apiService *ApiService) CheackUsernameAvailability(w http.ResponseWriter, r *http.Request) {

	var payload UsernamePayload

	if err := readJson(w, r, payload); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()
	exist := apiService.userRpo.CheackUsernameAvailability(ctx, payload.Username)

	returnPayload := BoolPayload{
		Exist: exist,
	}

	writeJson(w, http.StatusOK, returnPayload)

}
