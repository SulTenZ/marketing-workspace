package httpx

import (
	"encoding/json"
	"net/http"
)

type ErrorCode string

const (
	NotFound           ErrorCode = "NOT_FOUND"
	MethodNotAllowed   ErrorCode = "METHOD_NOT_ALLOWED"
	InvalidOrigin      ErrorCode = "INVALID_ORIGIN"
	InvalidJSON        ErrorCode = "INVALID_JSON"
	RequestTooLarge    ErrorCode = "REQUEST_TOO_LARGE"
	UnsupportedMedia   ErrorCode = "UNSUPPORTED_MEDIA_TYPE"
	InternalError      ErrorCode = "INTERNAL_ERROR"
	RequestTimeout     ErrorCode = "REQUEST_TIMEOUT"
	RequestCancelled   ErrorCode = "REQUEST_CANCELLED"
	ServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
)

type ErrorEnvelope struct {
	Error ErrorDetail `json:"error"`
}

type ErrorDetail struct {
	Code      ErrorCode         `json:"code"`
	Message   string            `json:"message"`
	Fields    map[string]string `json:"fields"`
	RequestID string            `json:"request_id"`
}

func errorInfo(code ErrorCode) (ErrorCode, int, string) {
	switch code {
	case NotFound:
		return code, http.StatusNotFound, "Resource tidak ditemukan."
	case MethodNotAllowed:
		return code, http.StatusMethodNotAllowed, "Metode HTTP tidak diizinkan."
	case InvalidOrigin:
		return code, http.StatusForbidden, "Origin permintaan tidak diizinkan."
	case InvalidJSON:
		return code, http.StatusBadRequest, "Body JSON tidak valid."
	case RequestTooLarge:
		return code, http.StatusRequestEntityTooLarge, "Ukuran permintaan melebihi batas."
	case UnsupportedMedia:
		return code, http.StatusUnsupportedMediaType, "Gunakan Content-Type application/json."
	case RequestTimeout:
		return code, http.StatusServiceUnavailable, "Waktu pemrosesan permintaan habis."
	case RequestCancelled:
		return code, http.StatusServiceUnavailable, "Permintaan dibatalkan."
	case ServiceUnavailable:
		return code, http.StatusServiceUnavailable, "Layanan belum tersedia."
	default:
		return InternalError, http.StatusInternalServerError, "Terjadi kesalahan internal."
	}
}

func errorBody(r *http.Request, code ErrorCode) []byte {
	code, _, message := errorInfo(code)
	// Hanya katalog tetap; arbitrary error/panic/config tidak keluar ke klien.
	body, _ := json.Marshal(ErrorEnvelope{Error: ErrorDetail{
		Code: code, Message: message, Fields: map[string]string{}, RequestID: RequestID(r),
	}})
	return append(body, '\n')
}

func WriteError(w http.ResponseWriter, r *http.Request, code ErrorCode) {
	_, status, _ := errorInfo(code)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("Cache-Control", "no-store")
	w.WriteHeader(status)
	_, _ = w.Write(errorBody(r, code))
}
