package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	_ "github.com/joho/godotenv/autoload"
	"github.com/naneri/shortener/cmd/grpc/config"
	"github.com/naneri/shortener/cmd/grpc/proto"
	"github.com/naneri/shortener/internal/app/link"
	"github.com/naneri/shortener/internal/migrations"
	"google.golang.org/grpc"
	"log"
	"net"
	"os"
)

var cfg config.Config
var linkRepository link.Repository
var db *sql.DB

func main() {
	if flag.Lookup("c") != nil {
		flag.StringVar(&cfg.FileConfig, "c", cfg.FileConfig, "read config from File")

		var fileConfig config.FileConfig
		fileConfigData, readErr := os.ReadFile(cfg.FileConfig)
		if readErr != nil {
			log.Fatal("error reading the File config: " + readErr.Error())
		}

		unMarshalErr := json.Unmarshal(fileConfigData, &fileConfig)

		if unMarshalErr != nil {
			log.Fatal("error unmarshalling the config: " + unMarshalErr.Error())
		}

		cfg.ServerAddress = fileConfig.ServerAddress
		cfg.FileStoragePath = fileConfig.FileStoragePath
		cfg.DatabaseAddress = fileConfig.DatabaseDsn
		cfg.BaseURL = fileConfig.BaseURL
	}

	err := env.Parse(&cfg)

	if err != nil {
		log.Fatalf("unable to parse env vars: %v", err)
	}

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

	// определяем порт для сервера
	listen, err := net.Listen("tcp", cfg.ServerAddress)
	if err != nil {
		log.Fatal(err)
	}

	// создаём gRPC-сервер без зарегистрированной службы
	s := grpc.NewServer()
	// регистрируем сервис
	proto.RegisterShortenerServiceServer(s, &ShortenerServer{
		LinkRepository: linkRepository,
		Config:         cfg,
	})

	fmt.Println("Сервер gRPC начал работу")
	// получаем запрос gRPC
	if err := s.Serve(listen); err != nil {
		log.Fatal(err)
	}
}
