package main

import (
	"fmt"
	"net/http"
	"strconv"
)

type mainRequest struct {
	url string
}

var lastUrlId int
var storage map[string]string

func indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		if r.Method != http.MethodGet {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		id := r.URL.Path[1:]

		if val, ok := storage[id]; ok {
			w.Header().Set("Location", val)
			w.WriteHeader(http.StatusTemporaryRedirect)
			w.Write(nil)
			return
		} else {
			http.Error(w, "The URL not found", http.StatusNotFound)
			return
		}

	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		// если не заполнена, возвращаем код ошибки
		http.Error(w, "Bad request", 401)
		return
	}

	w.Header().Set("content-type", "plain/text")
	w.WriteHeader(http.StatusCreated)

	lastUrlId++
	storage[strconv.Itoa(lastUrlId)] = r.FormValue("url")

	w.Write([]byte(fmt.Sprintf("localhost:8080/%d", lastUrlId)))
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
