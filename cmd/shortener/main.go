package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"github.com/naneri/shortener/cmd/dto"
	"github.com/naneri/shortener/cmd/shortener/middleware"
	"github.com/naneri/shortener/internal/app/link"
	"log"
	"net/http"
	"os"
)

type Config struct {
	ServerAddress   string `env:"SERVER_ADDRESS" envDefault:":8080"`
	BaseURL         string `env:"BASE_URL" envDefault:"http://localhost:8080"`
	FileStoragePath string `env:"FILE_STORAGE_PATH" envDefault:""`
}

var cfg Config
var linkRepository link.Repository

func main() {
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatalf("unable to parse env vars: %v", err)
	}

	if flag.Lookup("a") == nil {
		flag.StringVar(&cfg.ServerAddress, "a", cfg.ServerAddress, "default server Port")
		flag.StringVar(&cfg.BaseURL, "b", cfg.BaseURL, "base URL")
		flag.StringVar(&cfg.FileStoragePath, "f", cfg.FileStoragePath, "file storage path")
	}

	flag.Parse()

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
	}

	linkRepository, err = link.InitFileRepo(file)
	if err != nil {
		log.Fatal("error reading the links")
	}

	r := mainHandler()

	log.Println("Server started at port " + cfg.ServerAddress)
	http.ListenAndServe(cfg.ServerAddress, r)
}

func mainHandler() *chi.Mux {
	r := chi.NewRouter()

	if linkRepository == nil {
		linkRepository, _ = link.InitFileRepo(nil)
	}

	r.Use(middleware.GzipMiddleware)
	r.Post("/", postURL)
	r.Post("/api/shorten", shortenURL)
	r.Get("/{url}", getURL)

	return r
}

func shortenURL(w http.ResponseWriter, r *http.Request) {
	type response struct {
		Result string `json:"result"`
	}
	var requestBody dto.ShortenerDto

	body, err := middleware.ReadBody(r)
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
	//if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
	//	log.Printf("io.ReadAll: %v\n", err)
	//	http.Error(w, "unable to read request body", http.StatusBadRequest)
	//	return
	//}

	lastURLID, err := linkRepository.AddLink(requestBody.URL)

	if err != nil {
		log.Print("error when adding a link:" + err.Error())
		http.Error(w, "error shortening the link.", http.StatusInternalServerError)
	}

	responseStruct := response{Result: generateShortLink(lastURLID)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(responseStruct)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getURL(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "url")

	if urlID == "" {
		fmt.Println("URL not found")
		http.Error(w, "The URL not found", http.StatusNotFound)
		return
	}

	if val, err := linkRepository.GetLink(urlID); err == nil {
		w.Header().Set("Location", val)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		fmt.Println("URL not found")
		http.Error(w, "The URL not found", http.StatusNotFound)
		return
	}
}

func postURL(w http.ResponseWriter, r *http.Request) {
	body, err := middleware.ReadBody(r)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "plain/text")
	w.WriteHeader(http.StatusCreated)

	lastURLID, _ := linkRepository.AddLink(string(body))

	shortLink := generateShortLink(lastURLID)
	_, _ = w.Write([]byte(shortLink))
}

func generateShortLink(lastURLID int) string {
	return fmt.Sprintf("%s/%d", cfg.BaseURL, lastURLID)
}
