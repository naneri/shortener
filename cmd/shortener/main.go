package main

import (
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/naneri/shortener/cmd/shortener/config"
	"github.com/naneri/shortener/cmd/shortener/controllers"
	"github.com/naneri/shortener/cmd/shortener/middleware"
	"github.com/naneri/shortener/internal/app/link"
	"log"
	"net/http"
	"os"
)

var cfg config.Config
var linkRepository link.Repository
var mainController controllers.MainController

func main() {
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatalf("unable to parse env vars: %v", err)
	}

	if flag.Lookup("a") == nil {
		flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "default server Port")
		flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base URL")
		flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")
	}

	flag.Parse()

	var file *os.File

	if cfg.FileStoragePath != "" {
		file, err = os.OpenFile(cfg.FileStoragePath, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0777)
		defer func(file *os.File) {
			fileCloseErr := file.Close()
			if fileCloseErr != nil {
				log.Fatal("error when closing file: " + fileCloseErr.Error())
			}
		}(file)

		if err != nil {
			log.Fatal("error opening the file")
		}
	}

	linkRepository, err = link.InitFileRepo(file)
	if err != nil {
		log.Fatal("error reading the links")
	}

	mainController = controllers.MainController{
		Db:     linkRepository,
		Config: cfg,
	}

	r := mainHandler()

	log.Println("Server started at port " + cfg.ServerAddress)
	http.ListenAndServe(cfg.ServerAddress, r)
}

func mainHandler() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.GzipMiddleware)
	r.Post("/", mainController.PostURL)
	r.Post("/api/shorten", mainController.ShortenURL)
	r.Get("/{url}", mainController.GetURL)

	return r
}
