package middleware

import (
	"fmt"
	"net/http"
)

func IdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, _ := r.Cookie("user")
		fmt.Println(cookie)
		next.ServeHTTP(w, r)
	})
}
