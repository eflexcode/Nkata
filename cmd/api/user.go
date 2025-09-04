package api

import (
	"database/sql"
	"errors"
	"main/database"
	"net/http"
)

type RegisterUserStruct struct{
	Username     string `json:"username"`
	DisplayName  string `json:"display_name"`
	Password     string `json:"password"`
}

func RegisterUser(w http.ResponseWriter, r *http.Request) {

	var payload RegisterUserStruct

	if err :=readJson(w,r,payload); err != nil{
		badRequest(w,r,err)
		return
	}

	if (len(payload.Username) > 10){
		err := errors.New("username max character is 10")
		badRequest(w,r,err)
		return
	}

	if len(payload.Password) <= 3 {
		err := errors.New("password min is 3")
		badRequest(w,r,err)
		return
	}

	user := &database.User{
		Username: payload.Username,
		DisplayName: payload.DisplayName,
		Password: payload.Password,
	}

	ctx := r.Context()

	uRepo := database.NewUserRepository(d)

	err :=uRepo.CreateUser(ctx,user)


	if err != nil {
		
	}

}

func SignIn(w http.ResponseWriter, r *http.Request) {

}

func GetByID(w http.ResponseWriter, r *http.Request) {

}

func CheackUsernameAvailability(w http.ResponseWriter, r *http.Request) {

}