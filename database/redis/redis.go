package redis

import (
	"fmt"
	"github.com/go-redis/redis"
	"log"
	"os"
)

type Config struct {
	Addr     string
	Password string
	DB       int
	DBurl    string
}

func New(config *Config) *redis.Client {
	redis := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.DB,
	})
	pong, err := redis.Ping().Result()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	fmt.Println(pong, nil)
	return redis
}
