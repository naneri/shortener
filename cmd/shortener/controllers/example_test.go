package controllers_test

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/naneri/shortener/cmd/shortener/config"
	"github.com/naneri/shortener/cmd/shortener/router"
	"github.com/naneri/shortener/internal/app/link"
	"github.com/naneri/shortener/internal/migrations"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
)

var cfg config.Config
var db *sql.DB

func ExampleMainController_PostURL() {
	appRouter := getRouter()

	// this somehow causes - `panic: runtime error: invalid memory address or nil pointer dereference`
	//defer appRouter.DB.Close()

	postRequest := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(`{"url" : "https://google.com"}`)))

	w := httptest.NewRecorder()
	r := appRouter.GetHandler()
	h := http.HandlerFunc(r.ServeHTTP)
	// запускаем сервер
	h.ServeHTTP(w, postRequest)
	res := w.Result()

	respBody, err := io.ReadAll(res.Body)

	fmt.Println(string(respBody))
	if err != nil {
		log.Fatalln(err)
	}

	// Output:
	// http://localhost:8080/1
}

func getRouter() router.Router {
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatalf("unable to parse env vars: %v", err)
	}

	linkRepo := initDB()

	appRouter := router.Router{
		Repository: linkRepo,
		Config:     cfg,
	}

	return appRouter
}

func initDB() *link.DatabaseRepository {
	var err error
	dbConfig := "host=localhost user=postgres password=mysecretpassword dbname=yandex_test port=5432 sslmode=disable TimeZone=UTC"
	db, err = sql.Open("postgres", dbConfig)

	if err != nil {
		log.Fatal("error initializing the database, " + err.Error())
	}

	_ = migrations.DropTables(db)

	err = migrations.RunMigrations(db)

	if err != nil {
		log.Fatal("error running migrations: " + err.Error())
	}

	linkRepo, _ := link.InitDatabaseRepository(db)

	return linkRepo
}
