package main

import (
	"compress/gzip"
	"errors"
	"io"
	"net/http"
)

func readBody(r *http.Request) ([]byte, error) {
	var reader io.Reader

	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {

			return []byte{}, errors.New("error when reading the gzipped body")
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return []byte{}, errors.New("error when reading the gzipped body")
	}

	return body, nil
}
