package main

import (
	"database/sql"
	"flag"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/naneri/shortener/cmd/shortener/config"
	"github.com/naneri/shortener/cmd/shortener/router"
	"github.com/naneri/shortener/internal/app/link"
	"github.com/naneri/shortener/internal/migrations"
	"log"
	"net/http"
	_ "net/http/pprof"
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
		db, err = sql.Open("postgres", cfg.DatabaseAddress)

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
		defer file.Close()

		if err != nil {
			log.Fatal("error opening the file")
		}

		linkRepository, err = link.InitFileRepo(file)
		if err != nil {
			log.Fatal("error reading the links")
		}
	}

	if linkRepository == nil {
		linkRepository, _ = link.InitFileRepo(nil)
	}

	appRouter := router.Router{
		Repository: linkRepository,
		Config:     cfg,
		Db:         db,
	}

	log.Println("Server started at port " + cfg.ServerAddress)
	go func() {
		log.Println(http.ListenAndServe(":8087", nil))
	}()
	log.Println(http.ListenAndServe(cfg.ServerAddress, appRouter.GetHandler()))
}

//func mainHandler() *chi.Mux {
//	r := chi.NewRouter()
//
//	// if I don't do this, the main_test.go will fail as it only tests this handler and MainController does need the Repo
//	if linkRepository == nil {
//		linkRepository, _ = link.InitFileRepo(nil)
//	}
//
//	mainController := controllers.MainController{
//		LinkRepository: linkRepository,
//		Config:         cfg,
//	}
//
//	utilityController := controllers.UtilityController{
//		DBConnection: db,
//	}
//
//	r.Use(middleware.GzipMiddleware)
//	r.Use(middleware.IDMiddleware)
//	r.Post("/", mainController.PostURL)
//	r.Post("/api/shorten", mainController.ShortenURL)
//	r.Get("/{url}", mainController.GetURL)
//	r.Get("/api/user/urls", mainController.UserUrls)
//	r.Delete("/api/user/urls", mainController.DeleteUserUrls)
//	r.Get("/ping", utilityController.PingDB)
//	r.Post("/api/shorten/batch", mainController.ShortenBatch)
//
//	return r
//}
