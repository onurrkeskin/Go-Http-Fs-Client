package web

import (
	"encoding/json"
	"net/http"

	"github.com/dimfeld/httptreemux/v5"
)

func Param(r *http.Request, key string) string {
	m := httptreemux.ContextParams(r.Context())
	return m[key]
}

func Decode(r *http.Request, val any) error {
	decoder := json.NewDecoder(r.Body)

	if err := decoder.Decode(val); err != nil {
		return err
	}

	return nil
}
