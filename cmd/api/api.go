package api

import (
	"database/sql"
	"log"
	"main/database"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "main/cmd/docs"
)

type ApiService struct {
	database *database.DataRepository
}

func NewRepos(userRepo *database.DataRepository) *ApiService {
	return &ApiService{database: userRepo}
}

// @title Example API
// @version 1.0
// @description This is a sample server using Chi and Swagger.
// @host localhost:8080
// @BasePath /api/v1
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

		r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {

			type ping struct {
				Message string `json:"message"`
			}

			writeJson(w, http.StatusOK, ping{Message: "pined"})
		})
		
		r.Get("/swagger/*", httpSwagger.WrapHandler)
		// r.Get("/swagger/*", httpSwagger.Handler(
		// 	httpSwagger.URL("http://localhost:5557/v1/swagger/doc.json"), 
		// ))

		r.Route("/user", func(r chi.Router) {
			r.Use(HandleJWTAuth)
			r.Get("/", apiService.GetByUsername)
			r.Put("/update", apiService.Update)
			r.Post("/upload-profile-picture", apiService.UploadProfilPic)
		})

		r.Route("/media", func(r chi.Router) {
			r.Get("/profiles/{img_name}", apiService.LoadProfilPic)
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

	err := http.ListenAndServe(":5557", r)

	if err != nil {
		log.Printf("Nkata server failed to start")
	}

}
