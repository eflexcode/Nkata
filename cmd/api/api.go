package api

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

func IntiApi() {
	r := chi.NewRouter()

	r.Get("/test", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("test")
	})

	log.Printf("Nkata http server started on port :3000")
	http.ListenAndServe(":3000", r)
}
