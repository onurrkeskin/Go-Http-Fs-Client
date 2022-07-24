package mid

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"gitlab.com/onurkeskin/go-http-fs-client/foundation/web"
)

func Panics() web.Middleware {

	m := func(handler web.Handler) web.Handler {
		h := func(ctx context.Context, w http.ResponseWriter, r *http.Request) (err error) {
			defer func() {
				if rec := recover(); rec != nil {

					// Stack trace will be provided.
					trace := debug.Stack()
					err = fmt.Errorf("PANIC [%v] TRACE[%s]", rec, string(trace))
				}
			}()
			return handler(ctx, w, r)
		}

		return h
	}

	return m
}
