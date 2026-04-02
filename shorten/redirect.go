package shorten

import (
	"net/http"

	"github.com/Techit-Kakaew/urlshorten/httprouter"
)

type redirectHandler struct {
	storage OriginalURLGetter
}

type OriginalURLGetter interface {
	OriginalURL(shortenURL string) (string, error)
}

func NewRedirectHandler(storage OriginalURLGetter) *redirectHandler {
	return &redirectHandler{storage: storage}
}

func (handler *redirectHandler) Handler(c httprouter.Context) {
	shortenURL := c.PathValue("shorturl")

	originalURL, err := handler.storage.OriginalURL(shortenURL)
	if err != nil {
		c.Status(http.StatusNotFound)
		return
	}

	c.Redirect(http.StatusPermanentRedirect, originalURL)
}
