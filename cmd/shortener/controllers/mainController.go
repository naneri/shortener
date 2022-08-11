package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/naneri/shortener/cmd/shortener/config"
	"github.com/naneri/shortener/cmd/shortener/dto"
	"github.com/naneri/shortener/cmd/shortener/middleware"
	"github.com/naneri/shortener/internal/app/link"
	"log"
	"net/http"
)

type ShortenResponse struct {
	Result string `json:"result"`
}

type MainController struct {
	LinkRepository link.Repository // <--
	Config         config.Config
}

// ShortenURL shortens a URL passed by user
// @Summary Shorten URL endpoint
// @Accept json
// @Produce json
// @Success 201 {object} ShortenResponse
// @Failure 409 {object} ShortenResponse
// @Failure 500 {string} string
// @Router /api/shorten [post]
func (controller *MainController) ShortenURL(w http.ResponseWriter, r *http.Request) {
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

	userID, ok := r.Context().Value(middleware.UserID(middleware.UserIDContextKey)).(uint32)

	if !ok {
		http.Error(w, "wrong user ID", http.StatusInternalServerError)
		return
	}

	lastURLID, err := controller.LinkRepository.AddLink(requestBody.URL, userID)

	if err != nil {
		if lastURLID != 0 {
			responseStruct := ShortenResponse{Result: generateShortLink(lastURLID, controller.Config.BaseURL)}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusConflict)

			err = json.NewEncoder(w).Encode(responseStruct)

			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}
		log.Print("error when adding a link:" + err.Error())
		http.Error(w, "error shortening the link.", http.StatusInternalServerError)
	}

	responseStruct := ShortenResponse{Result: generateShortLink(lastURLID, controller.Config.BaseURL)}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(responseStruct)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetURL gets a URL from the storage if the urlID exists
// @Summary Get URL using given ID
// @Accept plain
func (controller *MainController) GetURL(w http.ResponseWriter, r *http.Request) {
	urlID := chi.URLParam(r, "url")

	if urlID == "" {
		fmt.Println("empty URL")
		http.Error(w, "The URL param is empty", http.StatusNotFound)
		return
	}

	val, err := controller.LinkRepository.GetLink(urlID)

	if err != nil {
		var modelDeletedErr *link.ModelDeletedError

		if errors.As(err, &modelDeletedErr) {
			w.Header().Set("content-type", "plain/text")
			w.WriteHeader(http.StatusGone)
		} else {
			fmt.Println("URL not found: " + err.Error())
			http.Error(w, "The URL not found", http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Location", val)
	w.WriteHeader(http.StatusTemporaryRedirect)
}

// PostURL stores a URL provided by use, almost identical to ShortenURL
// @Summary post URL (similar to ShortenURL)
// @Accept plain
// @Produce plain
// @Success 201 {string} string
// @Failure 500 {string} string
// @Failure 409 {string} string
// @Router /	[post]
func (controller *MainController) PostURL(w http.ResponseWriter, r *http.Request) {
	body, err := middleware.ReadBody(r)
	// обрабатываем ошибку
	if err != nil {
		http.Error(w, "error reading body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	userID, ok := r.Context().Value(middleware.UserID(middleware.UserIDContextKey)).(uint32)

	if !ok {
		http.Error(w, "wrong user ID", http.StatusInternalServerError)
		return
	}

	lastURLID, err := controller.LinkRepository.AddLink(string(body), userID)

	if err != nil {
		if lastURLID != 0 {
			shortenedURL := generateShortLink(lastURLID, controller.Config.BaseURL)
			w.WriteHeader(http.StatusConflict)
			_, err = w.Write([]byte(shortenedURL))

			if err != nil {
				http.Error(w, "error generating short link: "+err.Error(), http.StatusInternalServerError)
				return
			}

			return
		}

		http.Error(w, "error with data storge: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "plain/text")
	w.WriteHeader(http.StatusCreated)
	shortLink := generateShortLink(lastURLID, controller.Config.BaseURL)
	_, _ = w.Write([]byte(shortLink))
}

// UserUrls shows a list of user shortened URLs
// @Summary Lists user shortened URLs
// @Accept json
// @Produce json
// @Success 200 {array} dto.ListLink
// @Success 204 {string} string
// @Failure 500 {string} string
// @Router /api/user/urls [get]
func (controller *MainController) UserUrls(w http.ResponseWriter, r *http.Request) {
	var userLinks []dto.ListLink
	userID, ok := r.Context().Value(middleware.UserID(middleware.UserIDContextKey)).(uint32)

	if !ok {
		http.Error(w, "wrong user ID", http.StatusInternalServerError)
		return
	}

	links, dbErr := controller.LinkRepository.GetAllLinks()

	if dbErr != nil {
		log.Println("Error getting data from the storage: " + dbErr.Error())
		http.Error(w, "error in the storage system", http.StatusInternalServerError)
	}

	for _, userLink := range links {
		if userLink.UserID == userID {
			dtoLink := dto.ListLink{
				ShortURL:    generateShortLink(userLink.ID, controller.Config.BaseURL),
				OriginalURL: userLink.URL,
			}
			userLinks = append(userLinks, dtoLink)
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if len(userLinks) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	err := json.NewEncoder(w).Encode(userLinks)

	if err != nil {
		log.Println("error generating shorten response when encoding user links: " + err.Error())
		http.Error(w, "error generating ShortenResponse", http.StatusInternalServerError)
		return
	}

}

// ShortenBatch shortens a batch of user URLs
// @Summary Shorten multiple URLs in batch
// @Accept json
// @Produce json
// @Success 201 {array} dto.ResponseBatchLink
// @Success 204 {string} string
// @Failure 500 {string} string
// @Router /api/shorten/batch [post]
func (controller *MainController) ShortenBatch(w http.ResponseWriter, r *http.Request) {
	var batchLinks []dto.BatchLink
	responseLinks := make([]dto.ResponseBatchLink, 0, 8)

	body, err := middleware.ReadBody(r)

	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &batchLinks)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	userID, ok := r.Context().Value(middleware.UserID(middleware.UserIDContextKey)).(uint32)

	if !ok {
		http.Error(w, "wrong user ID", http.StatusInternalServerError)
		return
	}

	for _, batchLink := range batchLinks {
		lastURLID, addErr := controller.LinkRepository.AddLink(batchLink.OriginalURL, userID)

		if addErr != nil {
			http.Error(w, addErr.Error(), http.StatusInternalServerError)
			return
		}

		responseLinks = append(responseLinks, dto.ResponseBatchLink{
			CorrelationID: batchLink.CorrelationID,
			ShortURL:      generateShortLink(lastURLID, controller.Config.BaseURL),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	if len(responseLinks) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.WriteHeader(http.StatusCreated)
	encodeErr := json.NewEncoder(w).Encode(responseLinks)

	if encodeErr != nil {
		http.Error(w, "error generating ShortenResponse", http.StatusInternalServerError)
		return
	}
}

// DeleteUserUrls deletes URLs of the user, that a taken from a list that the user has provided.
// @Summary deletes all urls of the User
// @Accept json
// @Produce plain
// @Success 202
// @Failure 500 {string} string
// @Router /api/user/urls [delete]
func (controller *MainController) DeleteUserUrls(w http.ResponseWriter, r *http.Request) {
	var urlIds []string

	body, err := middleware.ReadBody(r)

	// обрабатываем ошибку
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	err = json.Unmarshal(body, &urlIds)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("content-type", "plain/text")
	w.WriteHeader(http.StatusAccepted)

	go func() {
		deleteErr := controller.LinkRepository.DeleteLinks(urlIds)
		if deleteErr != nil {
			log.Println("error deleting links: " + deleteErr.Error())
		} else {
			log.Println("links deleted ")
		}
	}()
}

func generateShortLink(lastURLID int, baseURL string) string {
	return fmt.Sprintf("%s/%d", baseURL, lastURLID)
}
