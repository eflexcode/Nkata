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

type RegisterUserStruct struct {
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Password    string `json:"password"`
}

type LoginUsernameStuct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginEmailStuct struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type OptStruct struct {
	Email string `json:"email"`
	Otp   int    `json:"otp"`
}

type JwtJson struct {
	Token string `json:"token"`
}

func (apiService *ApiService) RegisterUser(w http.ResponseWriter, r *http.Request) {

	var payload RegisterUserStruct

	if err := readJson(w, r, payload); err != nil {
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

		internalServer(w, r, err)

	}

	s := StandardResponse{
		Status:  http.StatusCreated,
		Message: "User account created successfully procced to signin.",
	}

	writeJson(w, http.StatusCreated, s)
}

func (apiService *ApiService) SignInUsername(w http.ResponseWriter, r *http.Request) {

	var payload LoginUsernameStuct

	if err := readJson(w, r, payload); err != nil {
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
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(payload.Password))

	if err != nil {
		unauthorized(w, r, err)
		return
	}

	claims := jwt.MapClaims{

		"username":user.Username,
		"exp":time.Now().Add(time.Hour*48),

	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	var secret := env.GetString() 


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

func (api *ApiService) CheackUsernameAvailability(w http.ResponseWriter, r *http.Request) {

}
