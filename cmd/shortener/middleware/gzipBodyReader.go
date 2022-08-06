package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
)

// ReadBody reads the body and checks if it is gzipped or not
func ReadBody(r *http.Request) ([]byte, error) {
	var reader io.Reader

	// had problems with replacing the default request.Body with the gz and because of that moved all the reading logic to this function
	if r.Header.Get(`Content-Encoding`) == `gzip` {
		gz, err := gzip.NewReader(r.Body)
		if err != nil {
			return []byte{}, err
		}
		reader = gz
		defer gz.Close()
	} else {
		reader = r.Body
	}

	body, err := io.ReadAll(reader)
	if err != nil {
		return []byte{}, err
	}

	return body, nil
}
