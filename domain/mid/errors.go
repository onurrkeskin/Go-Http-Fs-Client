package mid

import (
	"context"
	"errors"
	"net/http"

	"gitlab.com/onurkeskin/go-http-fs-client/foundation/web"
	"go.uber.org/zap"
)

type ErrorResponse struct {
	Error  string            `json:"error"`
	Fields map[string]string `json:"fields,omitempty"`
}

type RequestError struct {
	Err    error
	Status int
}

func NewRequestError(err error, status int) error {
	return &RequestError{err, status}
}

func (re *RequestError) Error() string {
	return re.Err.Error()
}

func IsRequestError(err error) bool {
	var re *RequestError
	return errors.As(err, &re)
}

func GetRequestError(err error) *RequestError {
	var re *RequestError
	if !errors.As(err, &re) {
		return nil
	}
	return re
}

func Errors(log *zap.SugaredLogger) web.Middleware {
	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) error {
			if err := handler(ctx, w, r); err != nil {
				var er ErrorResponse
				var status int
				switch {
				case IsRequestError(err):
					reqErr := GetRequestError(err)
					er = ErrorResponse{
						Error: reqErr.Error(),
					}
					status = reqErr.Status

				default:
					er = ErrorResponse{
						Error: http.StatusText(http.StatusInternalServerError),
					}
					status = http.StatusInternalServerError
				}
				if err := web.Respond(ctx, w, er, status); err != nil {
					return err
				}
				if web.IsShutdown(err) {
					return err
				}
			}
			return nil
		}

		return h
	}

	return m
}
