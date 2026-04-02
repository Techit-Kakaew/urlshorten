package shorten

import (
	"fmt"
	"math/rand"
	"net/http"

	"github.com/Techit-Kakaew/urlshorten/httprouter"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

type URLBody struct {
	URL string `json:"url" binding:"required"`
}

type shortenHandler struct {
	// db      *sql.DB
	storage ShortenURLSaver
}

type ShortenURLSaver interface {
	Save(shortenURL, originalURL string) error
}

func NewHandler(storage ShortenURLSaver) *shortenHandler {
	return &shortenHandler{storage: storage}
}

func (handler *shortenHandler) Handler(c httprouter.Context) {
	var payload URLBody

	if err := c.Bind(&payload); err != nil {
		c.Status(http.StatusBadRequest)
	}

	const size = 6
	var shortenKey [size]byte
	for i := range size {
		shortenKey[i] = charset[rand.Intn(len(charset))]
	}

	if err := handler.storage.Save(string(shortenKey[:]), payload.URL); err != nil {
		c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, map[string]string{
		"short_url": fmt.Sprintf("http://localhost:8080/%s", shortenKey),
	})
}
