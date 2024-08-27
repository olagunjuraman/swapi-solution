package main

import (
	"busha/database/postgres"
	"busha/database/redis"
	"busha/handlers"
	"busha/models"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
	port, exist := os.LookupEnv("PORT")
	if !exist {
		port = "8080"
	}

	_, route := Setup()
	srv := &http.Server{
		Addr:         ":" + port,
		WriteTimeout: time.Second * 60,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      route,
	}

	if err := srv.ListenAndServe(); err != nil {
		log.Println(err)
	}
}

func Setup() (*handlers.Handler, *mux.Router) {
	dbUser, exist := os.LookupEnv("POSTGRES_USER")
	if !exist {
		log.Fatal("POSTGRES_USER not set in .env")
	}

	dbPass, exist := os.LookupEnv("POSTGRES_PASSWORD")
	if !exist {
		log.Fatal("POSTGRES_PASSWORD not set in .env")
	}

	dbHost, exist := os.LookupEnv("POSTGRES_HOST")
	if !exist {
		log.Fatal("POSTGRES_HOST not set in .env")
	}

	dbName, exist := os.LookupEnv("POSTGRES_DB")
	if !exist {
		log.Fatal("POSTGRES_DB not set in .env")
	}

	dbPort, exist := os.LookupEnv("POSTGRES_PORT")
	if !exist {
		log.Fatal("POSTGRES_PORT not set in .env")
	}

	sslMode, exist := os.LookupEnv("POSTGRES_SSLMode")
	if !exist {
		sslMode = "disable"
	}

	db, err := postgres.New(&postgres.Config{
		User:     dbUser,
		Password: dbPass,
		DBName:   dbName,
		Host:     dbHost,
		Port:     dbPort,
		SSLMode:  sslMode,
	})

	if err != nil {
		log.Fatal("Failed to connect to Postgres database", err)
	}

	//Migrate table(s)
	if err = postgres.SetupDatabase(db, &models.Comments{}); err != nil {
		log.Fatal("Failed To setup tables", err)
	}

	redisAddress, exist := os.LookupEnv("REDIS_ADDRESS")
	if !exist {
		log.Fatal("REDIS_ADDRESS not set in .env")
	}

	redisPassword, exist := os.LookupEnv("REDIS_PASSWORD")
	//if !exist {
	//	log.Fatal("REDIS_PASSWORD not set in .env")
	//}
	//Initiate Redis
	redisDB := redis.New(&redis.Config{
		Addr:     redisAddress,
		Password: redisPassword,
	})

	//Initiate Routes
	route := mux.NewRouter()
	h := handlers.New(db, redisDB)
	h.AddRoutes(route)
	return h, route
}
