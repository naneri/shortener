package controllers

import (
	"context"
	"database/sql"
	"github.com/naneri/shortener/cmd/shortener/config"
	"net"
	"net/http"
	"time"
)

type UtilityController struct {
	DBConnection *sql.DB
	Config       config.Config
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

func (cont *UtilityController) Stats(w http.ResponseWriter, r *http.Request) {
	if cont.Config.TrustedSubnet == "" {
		http.Error(w, "access forbidden", http.StatusForbidden)
		return
	}

	_, subnet, err := net.ParseCIDR(cont.Config.TrustedSubnet)

	if err != nil {
		http.Error(w, "error parsing subnet from Config, you need to restart the app with the right subnet setting", http.StatusInternalServerError)
		return
	}

	realIp := r.Header.Get("X-Real-IP")

	userIp := net.ParseIP(realIp)

	if userIp == nil {
		http.Error(w, "error - user has passed a wrong IP", http.StatusBadRequest)
		return
	}

	if !subnet.Contains(userIp) {
		http.Error(w, "access forbidden", http.StatusForbidden)
		return
	}

	w.WriteHeader(http.StatusOK)
}
