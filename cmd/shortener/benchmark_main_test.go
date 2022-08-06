package main

import (
	"database/sql"
	"github.com/naneri/shortener/cmd/shortener/dto"
	"github.com/naneri/shortener/internal/app/link"
	"github.com/naneri/shortener/internal/migrations"
	"log"
	"testing"
)

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

func BenchmarkMain(b *testing.B) {
	dbRepo := initDB()
	links := make([]dto.BatchLink, 5)

	for _, url := range testUrls {
		batchLink := dto.BatchLink{
			CorrelationID: url.CorrelationID,
			OriginalURL:   url.URL,
		}

		links = append(links, batchLink)
	}

	b.ResetTimer()
	b.Run("addBatchLinks", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			for _, url := range links {
				_, _ = dbRepo.AddLink(url.OriginalURL, 0)
			}

			b.StopTimer()
			_ = dbRepo.DeleteAllLinks()
			b.StartTimer()
		}
	})

	b.Run("deleteLinks", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			ids := make([]string, len(links))
			b.StopTimer()
			for _, url := range links {
				linkID, _ := dbRepo.AddLink(url.OriginalURL, 0)

				ids = append(ids, string(rune(linkID)))
			}

			b.StartTimer()
			_ = dbRepo.DeleteLinks(ids)
		}
	})
}

func initDB() *link.DatabaseRepository {
	var err error
	dbConfig := "host=localhost user=postgres password=mysecretpassword dbname=yandex_test port=5432 sslmode=disable TimeZone=UTC"
	db, err = sql.Open("postgres", dbConfig)

	if err != nil {
		log.Fatal("error initializing the database, " + err.Error())
	}

	err = migrations.RunMigrations(db)

	if err != nil {
		log.Fatal("error running migrations: " + err.Error())
	}

	linkRepo, _ := link.InitDatabaseRepository(db)

	defer func() {
		_ = db.Close()
	}()

	return linkRepo
}
