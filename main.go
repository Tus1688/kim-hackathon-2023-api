package main

import (
	"log"
	"net/http"

	"github.com/Tus1688/kim-hackathon-2023-api/authutil"
	"github.com/Tus1688/kim-hackathon-2023-api/controllers"
	"github.com/Tus1688/kim-hackathon-2023-api/database"
	"github.com/Tus1688/kim-hackathon-2023-api/middlwwares"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	err := database.InitMysql()
	if err != nil {
		log.Fatal("unable to connect to mysql", err)
	}
	log.Print("successfully connected to mysql")

	err = database.InitRedis()
	if err != nil {
		log.Fatal("unable to connect to redis", err)
	}
	log.Print("successfully connected to redis")

	err = authutil.InitializeJWTKey()
	if err != nil {
		log.Fatal("unable to initialize jwt key", err)
	}
	log.Print("successfully initialized jwt key")

	err = database.InitAdmin()
	if err != nil {
		log.Fatal("unable to migrate admin account", err)
	}
	log.Print("successfully migrated admin account")

	log.Print("server running on port 3000")
	r := initRouter()

	err = http.ListenAndServe(":3000", r)
	if err != nil {
		log.Fatal(err)
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
							r.Post("/login", controllers.Login)
							r.Get("/refresh", controllers.GetRefreshToken)
						},
					)

					// protected routes for admin
					r.Group(
						func(r chi.Router) {
							r.Use(middlwwares.EnforceAuthentication([]string{"admin"}, 3))

							r.Get("/user", controllers.GetUser)
							r.Post("/user", controllers.CreateUser)
							r.Patch("/user", controllers.ModifyUser)
							r.Delete("/user", controllers.DeleteUser)
						},
					)
				},
			)
		},
	)

	return r
}
