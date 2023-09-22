package main

import (
	"log"

	"github.com/Tus1688/kim-hackathon-2023-api/authutil"
	"github.com/Tus1688/kim-hackathon-2023-api/database"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := database.InitMysql()
	if err != nil {
		log.Fatal("unable to connect to mysql", err)
	}
	err = database.InitRedis()
	if err != nil {
		log.Fatal("unable to connect to redis", err)
	}
	err = authutil.InitializeJWTKey()
	if err != nil {
		log.Fatal("unable to initialize jwt key", err)
	}

}

func initRouter() *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Heartbeat("/ping"))

	r.Route(
		"/api/v1", func(r chi.Router) {
			r.Use(middleware.Compress(5, "application/json"))

			r.Route(
				"/auth", func(r chi.Router) {
					// unprotected routes
					r.Group(
						func(r chi.Router) {

						},
					)
				},
			)
		},
	)

	return r
}
