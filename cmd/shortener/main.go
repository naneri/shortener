package main

import (
	"encoding/json"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/naneri/shortener/cmd/dto"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Config struct {
	ServerAddress string `env:"SERVER_ADDRESS" envDefault:"8080"`
	BaseUrl       string `env:"BASE_URL" envDefault:"http://localhost:8080"`
}

const shortLinkHost = "http://localhost:8080"

var lastUrlId int
var storage map[string]string
var cfg Config

func main() {
	r := mainHandler()

	log.Println("Server started at port 8080")
	http.ListenAndServe(":"+cfg.ServerAddress, r)
}

func mainHandler() *chi.Mux {
	err := env.Parse(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	if storage == nil {
		storage = make(map[string]string)
	}

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

	lastUrlId++
	storage[strconv.Itoa(lastUrlId)] = requestBody.Url

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

	if val, ok := storage[urlId]; ok {
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

	lastUrlId++
	storage[strconv.Itoa(lastUrlId)] = string(body)

	shortLink := generateShortLink(lastUrlId)
	_, _ = w.Write([]byte(shortLink))
}

func generateShortLink(lastUrlId int) string {
	return fmt.Sprintf("%s/%d", cfg.BaseUrl, lastUrlId)
}
