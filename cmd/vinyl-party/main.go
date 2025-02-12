package main

import (
	"fmt"
	"net/http"
	"vinyl-party/internal/config"
	"vinyl-party/internal/handler/http/api"
	"vinyl-party/internal/repository"
	"vinyl-party/internal/service"
	"vinyl-party/internal/storage/mongodb"
	"vinyl-party/pkg/recovery"

	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := config.MustLoad()

	http.HandleFunc("/panic", recovery.Middleware(panicHandler))

	storage, err := mongodb.New(cfg.MongoURL, cfg.DbName)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(storage)

	router := chi.NewRouter()

	userRepo := repository.NewUserRepository(storage.Database())
	userService := service.NewUserService(userRepo)
	userHandler := api.NewUserHandler(userService)

	ratingRepo := repository.NewRatingRepository(storage.Database())
	ratingService := service.NewRatingService(ratingRepo)
	fmt.Println(ratingService) // TODO: Add handler

	albumRepo := repository.NewAlbumRepository(storage.Database())
	albumService := service.NewAlbumService(albumRepo)
	fmt.Println(albumService) // TODO: Add handler

	partyRepo := repository.NewPartyRepository(storage.Database())
	partySetvice := service.NewPartyService(partyRepo)
	fmt.Println(partySetvice) // TODO: Add handler

	router.Post("/users", userHandler.Register)
	router.Post("/login", userHandler.Login)

	server := &http.Server{
		Addr:         cfg.Address,
		Handler:      router,
		ReadTimeout:  cfg.HTTPServer.Timeout,
		WriteTimeout: cfg.HTTPServer.Timeout,
		IdleTimeout:  cfg.HTTPServer.IdleTimeout,
	}

	if err := server.ListenAndServe(); err != nil {
		fmt.Println(err)
	}

	http.ListenAndServe(cfg.Address, nil)
}

func panicHandler(w http.ResponseWriter, r *http.Request) {
	panic("Something went wrong!")
}
