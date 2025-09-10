package api

import (
	"database/sql"
	"log"
	"main/database"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type ApiService struct {
	userRpo *database.UserRepository
}

func NewRepos(userRepo *database.UserRepository) *ApiService {
	return &ApiService{userRpo: userRepo}
}

func IntiApi(db *sql.DB) {

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(90 * time.Second))

	uRepo := database.NewUserRepository(db)

	apiService := NewRepos(uRepo)

	r.Route("/v1", func(r chi.Router) {

		r.Get("/ping",func(w http.ResponseWriter, r *http.Request) {

			type ping struct{
				Message string `json:"message"`
			} 

			writeJson(w,http.StatusOK,ping{Message: "pined"})
		})

		r.Route("/user",func(r chi.Router) {
			r.Use(HandleJWTAuth)
			r.Get("/",apiService.GetByUsername)
		})

		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", apiService.RegisterUser)
			r.Post("/sign-in-with-username", apiService.SignInUsername)
			r.Get("/check-username", apiService.CheackUsernameAvailability)
		})

	})

	log.Printf("/**\n" +
		"* ·····························································\n" +
		"* : _   _ _         _          ____                           :\n" +
		"* :| \\ | | | ____ _| |_ __ _  / ___|  ___ _ ____   _____ _ __ :\n" +
		"* :|  \\| | |/ / _` | __/ _` | \\___ \\ / _ \\ '__\\ \\ / / _ \\ '__|:\n" +
		"* :| |\\  |   < (_| | || (_| |  ___) |  __/ |   \\ V /  __/ |   :\n" +
		"* :|_| \\_|_|\\_\\__,_|\\__\\__,_| |____/ \\___|_|    \\_/ \\___|_|   :\n" +
		"* ·····························································\n" +
		"*/")

	log.Printf("Nkata server started on port: 5557")

	err:= http.ListenAndServe(":5557", r)

	if err != nil {
		log.Printf("Nkata server failed to start")
	}

	
	
}
