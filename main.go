package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/MrWormHole/url-shortener/api"
	mr "github.com/MrWormHole/url-shortener/repository/mongo"
	rr "github.com/MrWormHole/url-shortener/repository/redis"
	"github.com/MrWormHole/url-shortener/shortener"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
)

func main() {
	repository := pickRepository()
	service := shortener.NewRedirectService(repository)
	handler := api.NewHandler(service)

	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/{hash}", handler.Get)
	router.Post("/", handler.Post)

	errorChannel := make(chan error, 2)
	go func() {
		fmt.Println("Listening on port 8000")
		errorChannel <- http.ListenAndServe(pickPort(), router)
	}()

	go func() {
		signalChannel := make(chan os.Signal, 1)
		signal.Notify(signalChannel, syscall.SIGINT)
		errorChannel <- fmt.Errorf("%s", <-signalChannel)
	}()

	fmt.Printf("Halted because of %s", <-errorChannel)

}

func pickPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func pickRepository() shortener.RedirectRepository {
	switch os.Getenv("URL_DB") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repository, err := rr.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repository
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repository, err := mr.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repository
	}
	return nil
}
