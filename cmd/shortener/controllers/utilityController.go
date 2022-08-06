package controllers

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

type UtilityController struct {
	DBConnection *sql.DB
}

// PingDB - pings the database to see if it has connected
func (cont *UtilityController) PingDB(w http.ResponseWriter, r *http.Request) {

	ctx, cancel := context.WithTimeout(r.Context(), 1*time.Second)
	defer cancel()
	if err := cont.DBConnection.PingContext(ctx); err != nil {
		http.Error(w, "error when pinging database", http.StatusInternalServerError)
	}

	w.WriteHeader(http.StatusOK)
}
