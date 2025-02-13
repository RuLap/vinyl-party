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

	spotifyService, err := service.NewSpotifyService(cfg.ProxyServer, cfg.SpotifyCredentials)
	if err != nil {
		fmt.Println(err)
	}

	userRepo := repository.NewUserRepository(storage.Database())
	albumRepo := repository.NewAlbumRepository(storage.Database())
	partyRepo := repository.NewPartyRepository(storage.Database())
	ratingRepo := repository.NewRatingRepository(storage.Database())
	partyRoleRepo := repository.NewPartyRoleRepository(storage.Database())
	participantRepo := repository.NewParticipantRepository(storage.Database())

	userService := service.NewUserService(userRepo)
	albumService := service.NewAlbumService(albumRepo)
	partyService := service.NewPartyService(partyRepo)
	ratingService := service.NewRatingService(ratingRepo)
	partyRoleService := service.NewPartyRoleService(partyRoleRepo)
	participantService := service.NewParticipantService(participantRepo)

	userHandler := api.NewUserHandler(userService)
	albumHandler := api.NewAlbumHandler(albumService)
	ratingHandler := api.NewRatingHandler(ratingService)
	partyHandler := api.NewPartyHandler(userService, albumService, partyService, ratingService, spotifyService, partyRoleService, participantService)

	router.Post("/users", userHandler.Register)
	router.Post("/login", userHandler.Login)

	router.Post("/ratings", ratingHandler.CreateRating)
	router.Get("/ratings/{id}", ratingHandler.GetRatingByID)

	router.Post("/albums", albumHandler.CreateAlbum)
	router.Get("/albums/{id}", albumHandler.GetAlbumByID)
	router.Delete("/albums/{id}", albumHandler.DeleteAlbum)
	router.Post("/albums/{id}/ratings", albumHandler.AddRating)

	router.Get("/parties", partyHandler.GetAllParties)
	router.Post("/parties", partyHandler.CreateParty)
	router.Get("/parties/{id}", partyHandler.GetPartyInfo)
	router.Post("/parties/{id}/albums", partyHandler.AddAlbumToParty)
	router.Post("/parties/{id}/participants", partyHandler.AddParticipantToParty)

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
