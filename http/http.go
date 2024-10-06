package http

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/a-h/templ"
)

type HTTPError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// Error handles errors by writing the appropriate HTTP status and JSON response
func Error(w http.ResponseWriter, r *http.Request, err error) {
	code, message := etp.ErrorCode(err), etp.ErrorMessage(err)

	slog.Error(message, "error", err, "code", code)

	w.WriteHeader(ErrorStatusCode(code))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HTTPError{Code: code, Message: message})
}

// lookup of application error codes to HTTP status codes
var codes = map[string]int{
	etp.ECONFLICT:       http.StatusConflict,
	etp.EINVALID:        http.StatusBadRequest,
	etp.ENOTFOUND:       http.StatusNotFound,
	etp.ENOTIMPLEMENTED: http.StatusNotImplemented,
	etp.EUNAUTHORIZED:   http.StatusUnauthorized,
	etp.EINTERNAL:       http.StatusInternalServerError,
}

// ErrorStatusCode returns the associated HTTP status code for a given `etp` error code
func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

// FromErrorStatusCode returns the associated etp code for an HTTP status code
func FromErrorStatusCode(code int) string {
	for k, v := range codes {
		if v == code {
			return k
		}
	}
	return etp.EINTERNAL
}

// Render renders a templated HTML response with a status code
func Render(w http.ResponseWriter, r *http.Request, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	// Render the component into the buffer
	if err := t.Render(r.Context(), buf); err != nil {
		return err
	}

	// Write the HTML response
	w.Header().Set("Content-Type", "text/html")
	w.WriteHeader(statusCode)
	_, err := w.Write([]byte(buf.Bytes()))
	return err
}

type Session struct{}
