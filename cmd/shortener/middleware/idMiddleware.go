package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"net/http"
	"time"
)

var secretKey = []byte("secret key")
var userId uint32

func IdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			data   []byte
			err    error
			idSign []byte
		)

		// parse cookie
		cookie, err := r.Cookie("user")

		if err != nil {
			fmt.Println("error getting cookie: " + err.Error())
			httpCookie := generateUserCookie()
			http.SetCookie(w, &httpCookie)
		} else {
			data, err = hex.DecodeString(cookie.Value)

			if err != nil {
				fmt.Println("error decoding cookie: " + err.Error())
				httpCookie := generateUserCookie()
				http.SetCookie(w, &httpCookie)
			} else {
				userId = binary.BigEndian.Uint32(data[:4])
				h := hmac.New(sha256.New, secretKey)
				h.Write(data[:4])
				idSign = h.Sum(nil)

				// if parse correctly, add the cookie to context
				if !hmac.Equal(idSign, data[4:]) {
					fmt.Println("wrong sign")
					httpCookie := generateUserCookie()
					http.SetCookie(w, &httpCookie)
				}
			}
		}

		ctx := r.Context()
		req := r.WithContext(context.WithValue(ctx, "userId", userId))
		*r = *req
		// else grant user the signed cookie with Unique identifier
		next.ServeHTTP(w, r)
	})
}

func generateUserCookie() http.Cookie {
	userId++
	uint32userIdBuf := make([]byte, 4)
	binary.BigEndian.PutUint32(uint32userIdBuf[0:], uint32(userId))

	hash := hmac.New(sha256.New, secretKey)
	hash.Write(uint32userIdBuf)
	sign := hash.Sum(uint32userIdBuf)
	userCookie := hex.EncodeToString(sign)

	expire := time.Now().Add(10 * time.Minute)
	httpCookie := http.Cookie{Name: "user", Value: userCookie, Path: "/", Expires: expire, MaxAge: 90000}

	return httpCookie
}