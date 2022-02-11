package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func indexHandler(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
		return
	}

	b, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	fmt.Println(b)
	os.Exit(0)
	//urlToShorten := r.URL.Query().Get("url")
	//if urlToShorten == "" {
	//    http.Error(w, "The url parameter is missing", http.StatusBadRequest)
	//    return
	//}
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(http.StatusCreated)

	responseJson := map[string]string{
		"id":           "1",
		"shortenedUrl": "localhost:8080/1",
	}
	fmt.Println(responseJson)
}

func main() {

	// маршрутизация запросов обработчику
	http.HandleFunc("/", indexHandler)
	// запуск сервера с адресом localhost, порт 8080

	fmt.Println("starting the server")
	http.ListenAndServe(":8080", nil)
}
