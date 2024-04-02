package main

import (
	"log"
	"net/http"

	"github.com/LucasBelusso1/go-OTELChallange/cepvalidation/internal/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post(`/`, handlers.ValidateCEPAndDispatch)

	http.ListenAndServe(":8080", r)
	log.Printf("Listening on port %s", "8080")
}
