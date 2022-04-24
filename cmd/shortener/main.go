package main

import (
	"database/sql"
	"flag"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/naneri/shortener/cmd/shortener/config"
	"github.com/naneri/shortener/cmd/shortener/controllers"
	"github.com/naneri/shortener/cmd/shortener/middleware"
	"github.com/naneri/shortener/internal/app/link"
	"github.com/naneri/shortener/internal/migrations"
	"log"
	"net/http"
	"os"
)

var cfg config.Config
var linkRepository link.Repository
var db *sql.DB

func main() {
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatalf("unable to parse env vars: %v", err)
	}

	// checking the CLI params and filling in case that they are passed
	if flag.Lookup("a") == nil {
		flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "default server Port")
		flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base URL")
		flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")
		flag.StringVar(&cfg.DatabaseAddress, "d", cfg.DatabaseAddress, "database DSN")
	}

	flag.Parse()

	if cfg.DatabaseAddress != "" {
		db, err = sql.Open("pgx", cfg.DatabaseAddress)

		if err != nil {
			log.Fatal("error initializing the database, " + err.Error())
		}

		err = migrations.RunMigrations(db)

		if err != nil {
			log.Fatal("error running migrations: " + err.Error())
		}

		linkRepository, _ = link.InitDatabaseRepository(db)
		defer func() {
			_ = db.Close()
		}()
	}

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

		linkRepository, err = link.InitFileRepo(file)
		if err != nil {
			log.Fatal("error reading the links")
		}
	}

	r := mainHandler()

	log.Println("Server started at port " + cfg.ServerAddress)
	http.ListenAndServe(cfg.ServerAddress, r)
}

func mainHandler() *chi.Mux {
	r := chi.NewRouter()

	// if I don't do this, the main_test.go will fail as it only tests this handler and MainController does need the Repo
	if linkRepository == nil {
		linkRepository, _ = link.InitFileRepo(nil)
	}

	mainController := controllers.MainController{
		LinkRepository: linkRepository,
		Config:         cfg,
	}

	utilityController := controllers.UtilityController{
		DbConnection: db,
	}

	r.Use(middleware.GzipMiddleware)
	r.Use(middleware.IDMiddleware)
	r.Post("/", mainController.PostURL)
	r.Post("/api/shorten", mainController.ShortenURL)
	r.Get("/{url}", mainController.GetURL)
	r.Get("/api/user/urls", mainController.UserUrls)
	r.Get("/ping", utilityController.PingDb)
	r.Post("/api/shorten/batch", mainController.ShortenBatch)

	return r
}
