package httpx

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"mime"
	"net/http"
	"unicode/utf8"
)

const MaxBodyBytes int64 = 256 << 10

func LimitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.ContentLength > MaxBodyBytes {
			WriteError(w, r, RequestTooLarge)
			return
		}
		r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
		next.ServeHTTP(w, r)
	})
}

// DecodeJSON dipakai handler endpoint berikutnya. Validasi field bisnis tetap
// menjadi tanggung jawab domain; detail error decoder tidak dikirim ke klien.
func DecodeJSON(w http.ResponseWriter, r *http.Request, destination any) ErrorCode {
	media, _, err := mime.ParseMediaType(r.Header.Get("Content-Type"))
	if err != nil || media != "application/json" {
		return UnsupportedMedia
	}
	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		var size *http.MaxBytesError
		if errors.As(err, &size) {
			return RequestTooLarge
		}
		return InvalidJSON
	}
	if !utf8.Valid(payload) {
		return InvalidJSON
	}
	decoder := json.NewDecoder(bytes.NewReader(payload))
	decoder.DisallowUnknownFields()
	if decoder.Decode(destination) != nil {
		return InvalidJSON
	}
	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return InvalidJSON
	}
	return ""
}
