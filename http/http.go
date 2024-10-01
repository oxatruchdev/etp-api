package http

import (
	"log/slog"
	"net/http"

	"github.com/Evalua-Tu-Profe/etp-api"
	"github.com/a-h/templ"
	"github.com/labstack/echo/v4"
)

func Error(c echo.Context, err error) error {
	code, message := etp.ErrorCode(err), etp.ErrorMessage(err)

	slog.Error(message, "error", err, "code", code)

	return c.JSON(ErrorStatusCode(code), echo.Map{"code": code, "message": message})
}

// lookup of application error codes to HTTP status codes.
var codes = map[string]int{
	etp.ECONFLICT:       http.StatusConflict,
	etp.EINVALID:        http.StatusBadRequest,
	etp.ENOTFOUND:       http.StatusNotFound,
	etp.ENOTIMPLEMENTED: http.StatusNotImplemented,
	etp.EUNAUTHORIZED:   http.StatusUnauthorized,
	etp.EINTERNAL:       http.StatusInternalServerError,
}

// ErrorStatusCode returns the associated HTTP status code for a etp error code.
func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

// FromErrorStatusCode returns the associated etp code for an HTTP status code.
func FromErrorStatusCode(code int) string {
	for k, v := range codes {
		if v == code {
			return k
		}
	}
	return etp.EINTERNAL
}

func Render(ctx echo.Context, statusCode int, t templ.Component) error {
	buf := templ.GetBuffer()
	defer templ.ReleaseBuffer(buf)

	if err := t.Render(ctx.Request().Context(), buf); err != nil {
		return err
	}

	return ctx.HTML(statusCode, buf.String())
}
