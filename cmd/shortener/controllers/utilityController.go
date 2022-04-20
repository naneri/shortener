package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

type UtilityController struct {
	DbConnection *sql.DB
}

func (cont *UtilityController) PingDb(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := cont.DbConnection.PingContext(ctx); err != nil {
		http.Error(w, "error when pinging database", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
	return
}
