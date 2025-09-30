package api

import (
	"database/sql"
	"errors"
	"main/database"
	"main/internal/evn"
	"math/rand"
	"net/http"
	"strconv"
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

type OtpPayloadLogin struct {
	Email string `json:"email"`
	Otp   int    `json:"otp"`
}

type OtpPayloadReset struct {
	Email    string `json:"email"`
	Otp      int    `json:"otp"`
	Password string `json:"password"`
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

type OtpPayload struct {
	Otp int64 `json:"otp"`
}

var otpPurposeLogin string = "Login"
var otpPurposeResetPassword string = "Reset"
var otpPurposeAddEmail string = "AddEmail"

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

	err := apiService.database.CreateUser(ctx, user)

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

	user, err := apiService.database.GetByUsername(ctx, payload.Username)

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

// sign in with email must verify email
func (api *ApiService) SignInEmail(w http.ResponseWriter, r *http.Request) {

	var payload LoginEmailePayload

	if err := readJson(w, r, &payload); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	user, err := api.database.GetUserByEmail(ctx, payload.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			unauthorized(w, r, err)
			return
		}
		internalServer(w, r, errors.New("somthing went wrong"))
		return
	}

	//send email with your smtp provider
	otpToken := rand.Intn(9000000) + 100000

	err = api.database.InsertOtp(ctx, user.Username, user.Email, otpPurposeLogin, int64(otpToken))

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  http.StatusOK,
		Message: "otp " + strconv.Itoa(otpToken) + " sent to " + user.Email,
	}

	writeJson(w, http.StatusOK, s)

}

func (api *ApiService) VerifySignInEmailOtp(w http.ResponseWriter, r *http.Request) {

	var payload OtpPayloadLogin

	if err := readJson(w, r, &payload); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	otp, err := api.database.GetOtp(ctx, int64(payload.Otp))

	if err != nil {
		unauthorized(w, r, errors.New("otp is invalid"))
		return
	}

	if otp.Purpose != otpPurposeLogin || otp.Email != payload.Email {
		unauthorized(w, r, errors.New("otp is invalid"))
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

	claims := jwt.MapClaims{
		"username": otp.Username,
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

func (api *ApiService) SendResetPasswordOtp(w http.ResponseWriter, r *http.Request) {

	var payload EmailPayload

	if err := readJson(w, r, &payload); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	user, err := api.database.GetUserByEmail(ctx, payload.Email)

	if err != nil {
		if err == sql.ErrNoRows {
			unauthorized(w, r, err)
			return
		}
		internalServer(w, r, errors.New("somthing went wrong"))
		return
	}

	//send email with your smtp provider
	otpToken := rand.Intn(9000000) + 100000

	err = api.database.InsertOtp(ctx, user.Username, user.Email, otpPurposeResetPassword, int64(otpToken))

	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  http.StatusOK,
		Message: "otp " + strconv.Itoa(otpToken) + " sent to " + user.Email,
	}

	writeJson(w, http.StatusOK, s)

}

func (api *ApiService) VerifyResetPasswordOtp(w http.ResponseWriter, r *http.Request) {

	var payload OtpPayloadReset

	if err := readJson(w, r, &payload); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()

	otp, err := api.database.GetOtp(ctx, int64(payload.Otp))

	if err != nil {
		if err == sql.ErrNoRows {
			unauthorized(w, r, err)
			return
		}
		internalServer(w, r, errors.New("somthing went wrong"))
		return
	}

	if otp.Email != payload.Email || otp.Purpose != otpPurposeResetPassword {
		unauthorized(w, r, err)
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

	err = api.database.UpdateUserPassword(ctx, payload.Password, otp.Email)
	
	if err != nil {
		internalServer(w, r, err)
		return
	}

	s := StandardResponse{
		Status:  http.StatusOK,
		Message: "password reset succefully procced to login",
	}

	writeJson(w, http.StatusOK, s)
}

func (apiService *ApiService) CheackUsernameAvailability(w http.ResponseWriter, r *http.Request) {

	var payload UsernamePayload

	if err := readJson(w, r, payload); err != nil {
		badRequest(w, r, err)
		return
	}

	ctx := r.Context()
	exist := apiService.database.CheackUsernameAvailability(ctx, payload.Username)

	returnPayload := BoolPayload{
		Exist: exist,
	}

	writeJson(w, http.StatusOK, returnPayload)

}
