package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/garyburd/redigo/redis"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"
)

func main() {
	redisHost := os.Getenv("REDIS_HOST")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	if redisHost == "" {
		redisHost = "localhost:6379"
	}
	if redisPassword == "" {
		redisPassword = ""
	}
	pool := newPool(redisHost, redisPassword)

	token := os.Getenv("TOKEN")

	broadcaster := NewSatelliteHandler(pool, token)

	router := httprouter.New()
	router.HandlerFunc("GET", "/:channel", broadcaster)
	router.HandlerFunc("POST", "/:channel", broadcaster)

	port := strings.TrimSpace(os.Getenv("PORT"))
	if port == "" {
		port = "3001"
	}

	handler := cors.Default().Handler(router)
	err := http.ListenAndServe(":"+port, handler)
	if err != nil {
		log.Fatal(err)
	}
}

func newPool(host, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     2,
		MaxActive:   5,
		IdleTimeout: 5 * time.Second,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", host)
			if err != nil {
				log.Println(err)
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					log.Println(err)
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
	}
}
