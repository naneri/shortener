package middleware

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"net/http"
	"time"
)

var secretKey = []byte("secret key")
var userId uint

func IdMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		var (
			data          []byte
			decodedUserId uint32
			err           error
			idSign        []byte
		)
		cookie, _ := r.Cookie("user")

		data, err = hex.DecodeString(cookie.String())
		if err != nil {
			panic(err)
		}
		decodedUserId = binary.BigEndian.Uint32(data[:4])
		h := hmac.New(sha256.New, secretKey)
		h.Write(data[:4])
		idSign = h.Sum(nil)

		// parse cookie

		// if parse correctly, add the cookie to context
		if hmac.Equal(idSign, data[4:]) {
			ctx := r.Context()
			req := r.WithContext(context.WithValue(ctx, "userId", decodedUserId))
			*r = *req
			next.ServeHTTP(w, r)
		}

		// else grant user the signed cookie with Unique identifier
		userId++
		uint32userIdBuf := make([]byte, 4)
		binary.BigEndian.PutUint32(uint32userIdBuf[0:], uint32(userId))

		hash := hmac.New(sha256.New, secretKey)
		hash.Write(uint32userIdBuf)
		sign := hash.Sum(uint32userIdBuf)
		userCookie := hex.EncodeToString(sign)

		expire := time.Now().Add(10 * time.Minute)
		httpCookie := http.Cookie{Name: "user", Value: userCookie, Path: "/", Expires: expire, MaxAge: 90000}
		http.SetCookie(w, &httpCookie)
		next.ServeHTTP(w, r)
	})
}
