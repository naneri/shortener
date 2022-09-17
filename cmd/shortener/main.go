package main

import (
	"context"
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/naneri/shortener/cmd/shortener/config"
	"github.com/naneri/shortener/cmd/shortener/router"
	"github.com/naneri/shortener/internal/app/link"
	"github.com/naneri/shortener/internal/migrations"
	"golang.org/x/crypto/acme/autocert"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
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
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGTERM, syscall.SIGINT, syscall.SIGQUIT)
	idleConnsClosed := make(chan struct{})

	printBuildData()

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
		cfg.BaseURL = fileConfig.BaseURL
		cfg.FileStoragePath = fileConfig.FileStoragePath
		cfg.DatabaseAddress = fileConfig.DatabaseDsn
		cfg.EnableHTTPS = fileConfig.EnableHTTPS
	}

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
		flag.BoolVar(&cfg.EnableHTTPS, "s", cfg.EnableHTTPS, "enable HTTPS")
		flag.StringVar(&cfg.TrustedSubnet, "t", cfg.TrustedSubnet, "trusted Subnet")
	}

	flag.Parse()

	if cfg.TrustedSubnet != "" {
		_, _, parseErr := net.ParseCIDR(cfg.TrustedSubnet)

		if parseErr != nil {
			log.Fatal("error parsing subnet:" + parseErr.Error())
		}
	}

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

	var server *http.Server

	if !cfg.EnableHTTPS {
		server = &http.Server{
			Addr:    cfg.ServerAddress,
			Handler: appRouter.GetHandler(),
		}
	} else {
		manager := &autocert.Manager{
			// директория для хранения сертификатов
			Cache: autocert.DirCache("cache-dir"),
			// функция, принимающая Terms of Service издателя сертификатов
			Prompt: autocert.AcceptTOS,
			// перечень доменов, для которых будут поддерживаться сертификаты
			HostPolicy: autocert.HostWhitelist(cfg.BaseURL),
		}
		// конструируем сервер с поддержкой TLS
		server = &http.Server{
			Addr:    ":443",
			Handler: appRouter.GetHandler(),
			// для TLS-конфигурации используем менеджер сертификатов
			TLSConfig: manager.TLSConfig(),
		}
	}

	go func() {
		// читаем из канала прерываний
		// поскольку нужно прочитать только одно прерывание,
		// можно обойтись без цикла
		<-sigint
		// получили сигнал os.Interrupt, запускаем процедуру graceful shutdown
		if serverShutDownErr := server.Shutdown(context.Background()); serverShutDownErr != nil {
			// ошибки закрытия Listener
			log.Printf("HTTP server Shutdown: %v", serverShutDownErr)
		}
		// сообщаем основному потоку,
		// что все сетевые соединения обработаны и закрыты
		close(idleConnsClosed)
	}()

	if !cfg.EnableHTTPS {
		if listenErr := server.ListenAndServe(); listenErr != http.ErrServerClosed {
			// ошибки старта или остановки Listener
			log.Fatalf("HTTP server ListenAndServe: %v", listenErr)
		}
	} else {
		if listenErr := server.ListenAndServeTLS("", ""); listenErr != http.ErrServerClosed {
			// ошибки старта или остановки Listener
			log.Fatalf("HTTPS server ListenAndServe: %v", listenErr)
		}
	}

	<-idleConnsClosed

	fmt.Println("Server Shutdown gracefully")
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
