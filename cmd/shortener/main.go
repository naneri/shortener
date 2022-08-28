package main

import (
	"database/sql"
	"flag"
	"fmt"
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

var (
	buildVersion string
	buildDate    string
	buildCommit  string
)

var cfg config.Config
var linkRepository link.Repository
var db *sql.DB

func main() {
	printBuildData()
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
		defer func(file *os.File) {
			closeErr := file.Close()
			if closeErr != nil {
				fmt.Println("error closing file: " + closeErr.Error())
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

	if linkRepository == nil {
		linkRepository, _ = link.InitFileRepo(nil)
	}

	appRouter := router.Router{
		Repository: linkRepository,
		Config:     cfg,
		DB:         db,
	}

	log.Println("Server started at port " + cfg.ServerAddress)
	go func() {
		log.Println(http.ListenAndServe(":8087", nil))
	}()
	log.Println(http.ListenAndServe(cfg.ServerAddress, appRouter.GetHandler()))

	os.Exit(0)
}

func printBuildData() {
	fmt.Println(generateMessageOrNa("Build version: ", buildVersion))
	fmt.Println(generateMessageOrNa(`Build date: `, buildDate))
	fmt.Println(generateMessageOrNa(`Build commit: `, buildCommit))
}

func generateMessageOrNa(message string, value string) string {
	if value == "" {
		value = "N/A"
	}
	return message + value
}
