package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/naneri/shortener/cmd/dto"
	"github.com/naneri/shortener/internal/app/link"
	"io"
	"log"
	"net/http"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseUrl         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
}

var cfg Config
var linkRepository link.Repository

func main() {
	r := mainHandler()

	log.Println("Server started at port " + cfg.ServerAddress)
	http.ListenAndServe(cfg.ServerAddress, r)
}

func mainHandler() *chi.Mux {
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatal(err)
	}

	if flag.Lookup("a") == nil {
		flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "default server Port")
		flag.StringVar(&cfg.BaseUrl, "b", cfg.BaseUrl, "base URL")
		flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")
	}

	flag.Parse()

	linkRepository = link.InitFileRepo(cfg.FileStoragePath)

	r := chi.NewRouter()

	r.Post("/", postUrl)
	r.Post("/api/shorten", shortenUrl)
	r.Get("/{url}", getUrl)

	return r
}

func shortenUrl(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Result string `json:"result"`
	}
	var requestBody dto.ShortenerDto

	body, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(body, &requestBody)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	lastUrlId := linkRepository.AddLink(requestBody.Url)

	responseStruct := response{Result: generateShortLink(lastUrlId)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	fmt.Printf("%+v", responseStruct)
	err = json.NewEncoder(w).Encode(responseStruct)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func getUrl(w http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "url")

	if urlId == "" {
		fmt.Println("URL not found")
		http.Error(w, "The URL not found", http.StatusNotFound)
		return
	}

	if val, err := linkRepository.GetLink(urlId); err == nil {
		w.Header().Set("Location", val)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		fmt.Println("URL not found")
		http.Error(w, "The URL not found", http.StatusNotFound)
		return
	}
}

func postUrl(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	w.Header().Set("content-type", "plain/text")
	w.WriteHeader(http.StatusCreated)

	lastUrlId := linkRepository.AddLink(string(body))

	shortLink := generateShortLink(lastUrlId)
	_, _ = w.Write([]byte(shortLink))
}

func generateShortLink(lastUrlId int) string {
	return fmt.Sprintf("%s/%d", cfg.BaseUrl, lastUrlId)
}
