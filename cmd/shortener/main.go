package main

import (
	"fmt"
	"net/http"
)

type mainRequest struct {
	url string
}

var storage map[int]string

func indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path != "/" {
		if r.Method != http.MethodGet {
			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
			return
		}

		w.WriteHeader(http.StatusTemporaryRedirect)
		w.Write([]byte("http://yandex.ru"))
		return
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

	responseJson := map[string]string{
		"id":           "1",
		"shortenedUrl": "localhost:8080/1",
	}
	w.Write([]byte(responseJson["shortenedUrl"]))
}

func main() {

	// маршрутизация запросов обработчику
	http.HandleFunc("/", indexHandler)
	//http.HandleFunc()
	// запуск сервера с адресом localhost, порт 8080

	fmt.Println("starting the server")
	http.ListenAndServe(":8080", nil)
}
