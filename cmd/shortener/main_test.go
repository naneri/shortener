package main

import (
	"bytes"
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
				response:    `{"result":"http://localhost:8080/1"}`,
				contentType: "application/json",
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			// определяем хендлер
			r := mainHandler()
			h := http.HandlerFunc(r.ServeHTTP)
			// запускаем сервер
			h.ServeHTTP(w, test.request)
			res := w.Result()

			if res.StatusCode != test.want.code {
				t.Errorf("Expected status code %d, got %d", test.want.code, w.Code)
			}

			defer res.Body.Close()
			//resBody, err := io.ReadAll(res.Body)

			//if err != nil {
			//	t.Fatal(err)
			//}
			if strings.TrimSpace(w.Body.String()) != strings.TrimSpace(test.want.response) {
				t.Errorf("Expected body %s, got %s", test.want.response, w.Body.String())
			}
		})
	}
}
