package handlers

import (
	"busha/repository"
	"github.com/go-redis/redis"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	"log"
)

type Handler struct {
	DB      *gorm.DB
	Repo    *repository.Repository
	RedisDB *redis.Client
}

func New(DB *gorm.DB, redisDB *redis.Client) *Handler {
	repo := repository.New(DB)
	return &Handler{
		DB:      DB,
		Repo:    repo,
		RedisDB: redisDB,
	}
}

func (h *Handler) AddRoutes(route *mux.Router) {
	route.HandleFunc("/movies", h.GetMovies).Methods("GET")
	route.HandleFunc("/comment", h.AddComment).Methods("POST")
	route.HandleFunc("/comments/{id}", h.GetMovieComment).Methods("GET")
	route.HandleFunc("/characters/{id}", h.GetMovieCharacters).Methods("GET")
	log.Println("Loaded routes")
}
