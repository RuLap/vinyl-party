package main

import (
	"fmt"
	"net/http"
	"vinyl-party/internal/config"
	"vinyl-party/internal/handler/http/api"
	"vinyl-party/internal/repository"
	"vinyl-party/internal/service"
	"vinyl-party/internal/storage/mongodb"
	"vinyl-party/pkg/auth"
	"vinyl-party/pkg/cors"
	"vinyl-party/pkg/jwt_helper"
	"vinyl-party/pkg/recovery"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()

	if err := jwt_helper.NewJwtHelper(cfg.JWT.Secret); err != nil {
		fmt.Printf("Failed to initialize JWT: %v\n", err)
		return
	}

	storage, err := mongodb.New(cfg.MongoURL, cfg.DbName)
	if err != nil {
		fmt.Println(err)
	}

	router := chi.NewRouter()

	http.HandleFunc("/panic", recovery.Middleware(panicHandler))
	router.Use(cors.Middleware)

	spotifyService, err := service.NewSpotifyService(cfg.ProxyServer, cfg.SpotifyCredentials)
	if err != nil {
		fmt.Println(err)
	}

	userRepo := repository.NewUserRepository(storage.Database())
	albumRepo := repository.NewAlbumRepository(storage.Database())
	partyRepo := repository.NewPartyRepository(storage.Database())
	ratingRepo := repository.NewRatingRepository(storage.Database())

	err = partyRepo.EnsureIndexes()
	if err != nil {
		fmt.Println(err)
	}

	userService := service.NewUserService(userRepo)
	albumService := service.NewAlbumService(albumRepo)
	partyService := service.NewPartyService(partyRepo)
	ratingService := service.NewRatingService(ratingRepo)

	userHandler := api.NewUserHandler(userService)
	albumHandler := api.NewAlbumHandler(albumService)
	partyHandler := api.NewPartyHandler(userService, albumService, partyService, ratingService, spotifyService)

	router.Group(func(r chi.Router) {
		r.Post("/login", userHandler.Login)
		r.Post("/register", userHandler.Register)
	})

	router.Group(func(r chi.Router) {
		r.Use(
			jwt_helper.Middleware,
			auth.Middleware,
		)

		router.Post("/parties", partyHandler.CreateParty)
		router.Get("/users/{id}/parties", partyHandler.GetUserParties)
		//router.Get("/parties/{id}", partyHandler.)
		//router.Post("/parties/{id}/albums", partyHandler)
		//router.Post("/parties/{id}/participants", partyHandler)

		router.Post("/albums", albumHandler.CreateAlbum)
		router.Delete("/albums/{id}", albumHandler.DeleteAlbum)
		router.Post("/albums/{id}/ratings", albumHandler.AddRating)
	})

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("Something went wrong!")
}
