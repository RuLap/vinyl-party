package main

import (
	"fmt"
	"log/slog"
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
	slog.Info("loaded config successfuly")

	if err := jwt_helper.NewJwtHelper(cfg.JWT.Secret); err != nil {
		slog.Error("failed to init JWT", "error", err)
		return
	}
	slog.Info("generated JWT successfuly")

	storage, err := mongodb.New(cfg.MongoURL, cfg.DbName)
	if err != nil {
		slog.Error("failed to init mongodb", "error", err)
		fmt.Println(err)
	}
	slog.Info("init mongodb successfuly")

	router := chi.NewRouter()
	slog.Info("init chi router successfuly")

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
	slog.Info("init repositories successfuly")

	err = partyRepo.EnsureIndexes()
	if err != nil {
		slog.Info("failed to init party repository indexes")
	}

	userService := service.NewUserService(userRepo)
	albumService := service.NewAlbumService(albumRepo, ratingRepo, userRepo, storage.Database().Client())
	partyService := service.NewPartyService(partyRepo, albumRepo, ratingRepo, userRepo, storage.Database().Client())
	slog.Info("init services successfuly")

	userHandler := api.NewUserHandler(userService)
	albumHandler := api.NewAlbumHandler(albumService)
	partyHandler := api.NewPartyHandler(userService, albumService, partyService, spotifyService)
	slog.Info("init handlers successfuly")

	router.Group(func(r chi.Router) {
		r.Post("/login", userHandler.Login)
		r.Post("/register", userHandler.Register)
	})

	router.Group(func(r chi.Router) {
		r.Use(
			jwt_helper.Middleware,
			auth.Middleware,
		)

		router.Get("/users/{id}", userHandler.GetUser)
		router.Post("/parties", partyHandler.CreateParty)
		router.Get("/users/{id}/parties", partyHandler.GetUserParties)
		router.Get("/parties/{id}", partyHandler.GetParty)
		router.Post("/parties/{id}/albums", partyHandler.AddAlbum)
		router.Post("/parties/{id}/participants", partyHandler.AddParticipant)
		router.Post("/albums/{id}/ratings", albumHandler.AddRating)
	})
	slog.Info("init routes successfuly")

	server := &http.Server{
		Addr:         cfg.HTTPServer.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		slog.Info("server error", "error", err)
	}
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("Something went wrong!")
}
