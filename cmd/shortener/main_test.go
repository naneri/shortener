package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func Test_indexHandler(t *testing.T) {
	type want struct {
		code        int
		response    string
		contentType string
	}

	postRequest := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(`{"url" : "https://google.com"}`)))

	tests := []struct {
		name    string
		request *http.Request
		want    want
	}{
		{
			name:    "test posting url",
			request: postRequest,
			want: want{
				code:        201,
				response:    "http://localhost:8080/1",
				contentType: "plain/text",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			// определяем хендлер
			h := http.HandlerFunc(postUrl)
			// запускаем сервер
			h.ServeHTTP(w, tt.request)
			res := w.Result()

			if res.StatusCode != tt.want.code {
				t.Errorf("Expected status code %d, got %d", tt.want.code, w.Code)
			}

			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)
			if err != nil {
				t.Fatal(err)
			}
			if string(resBody) != tt.want.response {
				t.Errorf("Expected body %s, got %s", tt.want.response, w.Body.String())
			}
		})
	}
}
