package controllers_test

import (
	"bytes"
	"database/sql"
	"github.com/caarlos0/env/v6"
	"github.com/naneri/shortener/cmd/shortener/config"
	"github.com/naneri/shortener/cmd/shortener/router"
	"github.com/naneri/shortener/internal/app/link"
	"github.com/naneri/shortener/internal/migrations"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var db *sql.DB
var cfg config.Config

var testUrls = []struct {
	CorrelationID string
	URL           string
}{
	{
		"1",
		"https://google.com",
	},
	{
		"1",
		"https://yandex.ru",
	},
	{
		"1",
		"https://ya.ru",
	},
	{
		"1",
		"https://amazon.com",
	},
	{
		"1",
		"https://facebook.com",
	},
}

func TestPostUrl(t *testing.T) {
	appRouter := getRouter()

	postRequest := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(`{"url" : "https://google.com"}`)))
	w := httptest.NewRecorder()
	r := appRouter.GetHandler()
	h := http.HandlerFunc(r.ServeHTTP)
	// запускаем сервер
	h.ServeHTTP(w, postRequest)
	res := w.Result()
	respBody, _ := io.ReadAll(res.Body)

	if string(respBody) != "http://localhost:8080/1" {
		t.Errorf("error posting a url")
	}

	db.Close()
}

func TestShortenBatch(t *testing.T) {
	appRouter := getRouter()

	postRequest := httptest.NewRequest("POST", "/api/shorten/batch", bytes.NewReader([]byte(`[{"correlation_id": "1", "original_url": "https://3dnews.ru"}, {"correlation_id": "2", "original_url": "https://amazon.com"}]`)))

	w := httptest.NewRecorder()
	r := appRouter.GetHandler()
	h := http.HandlerFunc(r.ServeHTTP)
	// запускаем сервер
	h.ServeHTTP(w, postRequest)
	res := w.Result()

	if res.StatusCode != 201 {
		t.Errorf("error shortening URLs in batch")
	}

	db.Close()
}

func getRouter() router.Router {
	dbRepo := initDB()

	err := env.Parse(&cfg)

	if err != nil {
		log.Fatalf("unable to parse env vars: %v", err)
	}

	return router.Router{
		Repository: dbRepo,
		Config:     cfg,
	}
}

func initDB() *link.DatabaseRepository {
	var err error
	dbConfig := "host=localhost user=postgres password=mysecretpassword dbname=yandex_test port=5432 sslmode=disable TimeZone=UTC"
	db, err = sql.Open("postgres", dbConfig)

	if err != nil {
		log.Fatal("error initializing the database, " + err.Error())
	}

	dropErr := migrations.DropTables(db)
	if dropErr != nil {
		log.Fatal("error dropping the DB")
	}
	err = migrations.RunMigrations(db)

	if err != nil {
		log.Fatal("error running migrations: " + err.Error())
	}

	linkRepo, _ := link.InitDatabaseRepository(db)

	return linkRepo
}
