package router

import (
	"database/sql"
	"github.com/go-chi/chi/v5"
	"github.com/naneri/shortener/cmd/shortener/config"
	"github.com/naneri/shortener/cmd/shortener/controllers"
	"github.com/naneri/shortener/cmd/shortener/middleware"
	"github.com/naneri/shortener/internal/app/link"
)

type Router struct {
	Repository link.Repository
	Config     config.Config
	DB         *sql.DB
}

func (router *Router) GetHandler() *chi.Mux {
	r := chi.NewRouter()

	mainController := controllers.MainController{
		LinkRepository: router.Repository,
		Config:         router.Config,
	}

	utilityController := controllers.UtilityController{
		DBConnection: router.DB,
	}

	r.Use(middleware.GzipMiddleware)
	r.Use(middleware.IDMiddleware)
	r.Post("/", mainController.PostURL)
	r.Post("/api/shorten", mainController.ShortenURL)
	r.Get("/{url}", mainController.GetURL)
	r.Get("/api/user/urls", mainController.UserUrls)
	r.Delete("/api/user/urls", mainController.DeleteUserUrls)
	r.Get("/ping", utilityController.PingDB)
	r.Post("/api/shorten/batch", mainController.ShortenBatch)

	return r
}
