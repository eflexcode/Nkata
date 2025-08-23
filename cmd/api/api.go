package api

import (
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func IntiApi() {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(90 * time.Second))

	log.Printf("/**\n" +
		"* ·····························································\n" +
		"* : _   _ _         _          ____                           :\n" +
		"* :| \\ | | | ____ _| |_ __ _  / ___|  ___ _ ____   _____ _ __ :\n" +
		"* :|  \\| | |/ / _` | __/ _` | \\___ \\ / _ \\ '__\\ \\ / / _ \\ '__|:\n" +
		"* :| |\\  |   < (_| | || (_| |  ___) |  __/ |   \\ V /  __/ |   :\n" +
		"* :|_| \\_|_|\\_\\__,_|\\__\\__,_| |____/ \\___|_|    \\_/ \\___|_|   :\n" +
		"* ·····························································\n" +
		"*/")
	log.Printf("Nkata server started on port :3000")

	r.Route("/v1", func(r chi.Router) {

		r.Route("/auth", func(r chi.Router) {
			r.Post("/sign-up", CreateUser)
			r.Post("/sign-in", SignIn)
		})

	})

	http.ListenAndServe(":3000", r)
}
