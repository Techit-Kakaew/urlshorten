package httprouter

import (
	"encoding/json"
	"net/http"
)

type Context interface {
	PathValue(key string) string
	Bind(v any) error
	Status(code int)
	JSON(code int, v any)
	Redirect(code int, url string)
}

func NewHandler(handler func(Context)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := NewContext(w, r)
		handler(ctx)
	}
}

type context struct {
	w http.ResponseWriter
	r *http.Request
}

func NewContext(w http.ResponseWriter, r *http.Request) *context {
	return &context{w: w, r: r}
}

func (c *context) PathValue(key string) string {
	return c.r.PathValue(key)
}

func (c *context) Bind(v any) error {
	return json.NewDecoder(c.r.Body).Decode(v)
}

func (c *context) Status(code int) {
	c.w.WriteHeader(code)
}

func (c *context) JSON(code int, v any) {
	c.w.WriteHeader(code)
	json.NewEncoder(c.w).Encode(v)
}

func (c *context) Redirect(code int, url string) {
	http.Redirect(c.w, c.r, url, code)
}
