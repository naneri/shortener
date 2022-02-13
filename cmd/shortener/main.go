package main

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"io"
	"net/http"
	"strconv"
)

var lastUrlId int
var storage map[string]string

func main() {
	r := chi.NewRouter()

	r.Post("/", postUrl)
	r.Get("/{url}", getUrl)

	_ = http.ListenAndServe(":8080", r)
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
	if storage == nil {
		storage = make(map[string]string)
	}
	storage[strconv.Itoa(lastUrlId)] = string(body)

	_, _ = w.Write([]byte(fmt.Sprintf("http://localhost:8080/%d", lastUrlId)))
}
