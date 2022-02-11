package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
)

type mainRequest struct {
	url string
}

var lastUrlId int
var storage map[string]string

func indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		fmt.Printf("GET request to %s \n", r.URL)
		if r.Method != http.MethodGet {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		split := strings.Split(r.URL.Path, "/")
		id := split[1]

		if val, ok := storage[id]; ok {
			w.WriteHeader(http.StatusTemporaryRedirect)
			w.Header().Set("Location", val)
			w.Write([]byte(val))
			return
		} else {
			fmt.Println("URL not found \n")
			http.Error(w, "The URL not found", http.StatusNotFound)
			return
		}

	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	fmt.Printf("POST request to %s \n", r.URL)

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

	w.Write([]byte(fmt.Sprintf("http://localhost:8080/%d", lastUrlId)))
}

func main() {
	storage = make(map[string]string)
	// маршрутизация запросов обработчику
	http.HandleFunc("/", indexHandler)
	//http.HandleFunc()
	// запуск сервера с адресом localhost, порт 8080

	fmt.Println("starting the server")
	http.ListenAndServe("localhost:8080", nil)
}
