package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/naneri/shortener/cmd/shortener/config"
	"github.com/naneri/shortener/cmd/shortener/dto"
	"github.com/naneri/shortener/cmd/shortener/middleware"
	"github.com/naneri/shortener/internal/app/link"
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

	body, err := middleware.ReadBody(r)

	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &requestBody)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, ok := r.Context().Value("userId").(uint32)

	if !ok {
		http.Error(w, "wrong user ID", http.StatusInternalServerError)
		return
	}

	lastURLID, err := controller.DB.AddLink(requestBody.URL, userId)

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

	val, err := controller.DB.GetLink(urlID)

	if err != nil {
		fmt.Println("URL not found")
		http.Error(w, "The URL not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Location", val)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

func (controller *MainController) PostURL(w http.ResponseWriter, r *http.Request) {
	body, err := middleware.ReadBody(r)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userId, ok := r.Context().Value("userId").(uint32)

	if !ok {
		http.Error(w, "wrong user ID", http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "plain/text")
	w.WriteHeader(http.StatusCreated)

	lastURLID, err := controller.DB.AddLink(string(body), userId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	shortLink := generateShortLink(lastURLID, controller.Config.BaseURL)
	_, _ = w.Write([]byte(shortLink))
}

func (controller *MainController) UserUrls(w http.ResponseWriter, r *http.Request) {
	var userLinks []*link.Link
	userId, ok := r.Context().Value("userId").(uint32)

	if !ok {
		http.Error(w, "wrong user ID", http.StatusInternalServerError)
		return
	}

	links := controller.DB.GetAllLinks()

	for _, userLink := range links {
		if userLink.UserId == userId {
			userLinks = append(userLinks, userLink)
		}
	}

	fmt.Println(userLinks)
}

func generateShortLink(lastURLID int, baseURL string) string {
	return fmt.Sprintf("%s/%d", baseURL, lastURLID)
}
