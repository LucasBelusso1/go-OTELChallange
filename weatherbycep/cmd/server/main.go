package main

import (
	"log"
	"net/http"

	config "github.com/LucasBelusso1/go-OTELChallange/weatherbycep/configs"
	"github.com/LucasBelusso1/go-OTELChallange/weatherbycep/internal/webserver/handlers"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func init() {
	config.LoadConfig()
}

func main() {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get(`/{cep}`, handlers.GetTemperatureByZipCode)

	http.ListenAndServe(":8081", r)
	log.Printf("Listening on port %s", ":8081")
}
