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

//func indexHandler(w http.ResponseWriter, r *http.Request) {
//	if r.URL.Path != "/" {
//		fmt.Printf("GET request to %s \n", r.URL)
//		if r.Method != http.MethodGet {
//			http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
//			return
//		}
//
//		split := strings.Split(r.URL.Path, "/")
//		id := split[1]
//
//		if val, ok := storage[id]; ok {
//			w.Header().Set("Location", val)
//			w.WriteHeader(http.StatusTemporaryRedirect)
//			return
//		} else {
//			fmt.Println("URL not found")
//			http.Error(w, "The URL not found", http.StatusNotFound)
//			return
//		}
//
//	}

//if r.Method != http.MethodPost {
//	http.Error(w, "Only POST requests are allowed!", http.StatusMethodNotAllowed)
//	return
//}
//
//fmt.Printf("POST request to %s \n", r.URL)
//
//body, err := io.ReadAll(r.Body)
//// обрабатываем ошибку
//if err != nil {
//	http.Error(w, err.Error(), 500)
//	return
//}
//
//w.Header().Set("content-type", "plain/text")
//w.WriteHeader(http.StatusCreated)
//
//lastUrlId++
//if storage == nil {
//	storage = make(map[string]string)
//}
//storage[strconv.Itoa(lastUrlId)] = string(body)
//
//w.Write([]byte(fmt.Sprintf("http://localhost:8080/%d", lastUrlId)))
//}

//func getUrl(c *gin.Context) {
//
//}
//
//
//func postUrl(c *gin.Context) {
//	body, err := io.ReadAll(c.Request.Body)
//
//	if err != nil {
//		c.String(http.StatusInternalServerError, err.Error())
//		return
//	}
//
//	lastUrlId++
//	if storage == nil {
//		storage = make(map[string]string)
//	}
//	storage[strconv.Itoa(lastUrlId)] = string(body)
//
//	c.String(http.StatusCreated, fmt.Sprintf("http://localhost:8080/%d", lastUrlId))
//}
func getUrl(w http.ResponseWriter, r *http.Request) {
	urlId := chi.URLParam(r, "urlId")

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

func main() {
	r := chi.NewRouter()

	r.Post("/", postUrl)
	r.Get("{url}", getUrl)

	_ = http.ListenAndServe(":8080", r)
	// маршрутизация запросов обработчику
	//http.HandleFunc("/", indexHandler)
	////http.HandleFunc()
	//// запуск сервера с адресом localhost, порт 8080
	//
	//fmt.Println("starting the server")
	//http.ListenAndServe("localhost:8080", nil)
}
