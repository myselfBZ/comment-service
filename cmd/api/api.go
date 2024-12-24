package main

import (
	"github.com/myselfBZ/comment-service/internal/kafka"
	"github.com/myselfBZ/comment-service/internal/storage"
	"log"
	"net/http"
)

var (
	redisStore    *storage.RedisStore
	kafkaProducer *kafka.CommentProducer
)

func init() {
	redisClient, err := storage.InitRedisClient()
	if err != nil {
		log.Fatal("error initializing a connection with redis: ", err)
	}
	redisStore = storage.NewRedistStore(redisClient)
	log.Println("Redis is up")
	conn := kafka.InitCommentProducer()
	kafkaProducer = kafka.NewCommentProducer(conn, "comments")
}

type app struct {
	kafkaProducer *kafka.CommentProducer
	store         storage.CommentStorage
}

func NewApp() *app {
	return &app{store: redisStore, kafkaProducer: kafkaProducer}
}

func (a *app) mount() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /comments/{id}", a.getCommentByID)
	mux.HandleFunc("POST /comments", a.createComment)
	return mux
}

func (a *app) run(port string) error {
	return http.ListenAndServe(port, a.mount())
}
