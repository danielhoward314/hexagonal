package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	_ "github.com/joho/godotenv/autoload"

	h "github.com/danielhoward314/hexagonal/api"
	es "github.com/danielhoward314/hexagonal/repository/elasticsearch"
	mr "github.com/danielhoward314/hexagonal/repository/mongodb"
	pg "github.com/danielhoward314/hexagonal/repository/postgres"
	rr "github.com/danielhoward314/hexagonal/repository/redis"

	"github.com/danielhoward314/hexagonal/shortener"
)

func main() {
	repoType := flag.String("r", "redis", "sets repo type, must be one of 'redis', 'mongo', 'postgres' 'es'")
	flag.Parse()
	repo := chooseRepo(*repoType)
	service := shortener.NewRedirectService(repo)
	handler := h.NewHandler(service)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)
	errs := make(chan error, 2)
	go func() {
		port := httpPort()
		fmt.Printf("Listening on port %v", port)
		errs <- http.ListenAndServe(port, r)

	}()

	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT)
		errs <- fmt.Errorf("%s", <-c)
	}()
	fmt.Printf("Terminated %s", <-errs)
}

func httpPort() string {
	port := "8080"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func chooseRepo(repoType string) shortener.RedirectRepository {
	switch repoType {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := rr.NewRedisRepository(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGO_URL")
		mongodb := os.Getenv("MONGO_DB_NAME")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGO_TIMEOUT"))
		repo, err := mr.NewMongoRepository(mongoURL, mongodb, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "postgres":
		postgresHost := os.Getenv("POSTGRES_HOST")
		postgresUser := os.Getenv("POSTGRES_USER")
		postgresPassword := os.Getenv("POSTGRES_PASSWORD")
		postgresDBName := os.Getenv("POSTGRES_DB_NAME")
		postgresPort := os.Getenv("POSTGRES_PORT")
		repo, err := pg.NewPostgresRepository(postgresHost, postgresUser, postgresPassword, postgresDBName, postgresPort)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "es":
		repo, err := es.NewEsRepository()
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}
