package main

import (
	"bytes"
	"github.com/caarlos0/env/v6"
	"github.com/naneri/shortener/cmd/shortener/router"
	"github.com/naneri/shortener/internal/app/link"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type testStruct struct {
	name    string
	request *http.Request
	want    want
}
type want struct {
	code        int
	response    string
	contentType string
}

func Test_indexHandler(t *testing.T) {
	err := env.Parse(&cfg)

	if err != nil {
		log.Fatalf("unable to parse env vars: %v", err)
	}

	linkRepo := link.InitMemoryRepo()

	appRouter := router.Router{
		Repository: linkRepo,
		Config:     cfg,
	}

	postRequest := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(`{"url" : "https://google.com"}`)))
	jsonPostRequest := httptest.NewRequest("POST", "/api/shorten", bytes.NewReader([]byte(`{"url" : "https://google.com"}`)))

	tests := []testStruct{
		{
			name:    "test posting url",
			request: postRequest,
			want: want{
				code:        201,
				response:    "http://localhost:8080/1",
				contentType: "plain/text",
			},
		},
		{
			name:    "test '/api/shorten' ",
			request: jsonPostRequest,
			want: want{
				code:        201,
				response:    `{"result":"http://localhost:8080/2"}`,
				contentType: "application/json",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			// определяем хендлер
			r := appRouter.GetHandler()
			h := http.HandlerFunc(r.ServeHTTP)
			// запускаем сервер
			h.ServeHTTP(w, test.request)
			res := w.Result()

			if res.StatusCode != test.want.code {
				t.Errorf("Expected status code %d, got %d", test.want.code, w.Code)
			}

			defer res.Body.Close()

			if strings.TrimSpace(w.Body.String()) != strings.TrimSpace(test.want.response) {
				t.Errorf("Expected body %s, got %s", test.want.response, w.Body.String())
			}
		})
	}
}
