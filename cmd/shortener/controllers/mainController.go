package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/naneri/shortener/cmd/shortener/config"
	"github.com/naneri/shortener/cmd/shortener/dto"
	"github.com/naneri/shortener/internal/app/link"
	"io"
	"log"
	"net/http"
)

type MainController struct {
	DB     link.Repository // <--
	Config config.Config
}

func (controller *MainController) ShortenURL(w http.ResponseWriter, r *http.Request) {
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

	lastURLID, err := controller.DB.AddLink(requestBody.URL)

	if err != nil {
		log.Print("error when adding a link:" + err.Error())
		http.Error(w, "error shortening the link.", http.StatusInternalServerError)
	}

	responseStruct := response{Result: generateShortLink(lastURLID, controller.Config.BaseURL)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(responseStruct)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (controller *MainController) GetURL(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "url")

	if urlID == "" {
		fmt.Println("URL not found")
		http.Error(w, "The URL not found", http.StatusNotFound)
		return
	}

	if val, err := controller.DB.GetLink(urlID); err == nil {
		w.Header().Set("Location", val)
		w.WriteHeader(http.StatusTemporaryRedirect)
		return
	} else {
		fmt.Println("URL not found")
		http.Error(w, "The URL not found", http.StatusNotFound)
		return
	}
}

func (controller *MainController) PostURL(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "plain/text")
	w.WriteHeader(http.StatusCreated)

	lastURLID, _ := controller.DB.AddLink(string(body))

	shortLink := generateShortLink(lastURLID, controller.Config.BaseURL)
	_, _ = w.Write([]byte(shortLink))
}

func generateShortLink(lastURLID int, baseURL string) string {
	return fmt.Sprintf("%s/%d", baseURL, lastURLID)
}
