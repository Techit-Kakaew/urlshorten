package shorten

import (
	"errors"
	"net/http"
	"testing"

	"github.com/Techit-Kakaew/urlshorten/httprouter"
)

type contextTest struct {
	httprouter.Context
	code int
}

func (c *contextTest) Bind(any) error {
	return nil
}

func (c *contextTest) Status(code int) {
	c.code = code
}

func (c *contextTest) JSON(code int, v any) {
	c.code = code
}

type storageStub struct{}

func (storageStub) Save(string, string) error {
	return errors.New("error test")
}

func TestHandler(t *testing.T) {
	handler := &shortenHandler{
		storage: storageStub{},
	}

	ctx := &contextTest{}

	handler.Handler(ctx)

	if ctx.code != http.StatusInternalServerError {
		t.Errorf("Expected status code %d, got %d", http.StatusInternalServerError, ctx.code)
	}
}

type storageSpy struct {
	shortenKey  string
	originalURL string
}

func (s *storageSpy) Save(shortenKey, originalURL string) error {
	s.shortenKey = shortenKey
	s.originalURL = originalURL
	return errors.New("end test")
}

func TestStorageSaveKeyValueCorrectly(t *testing.T) {
	storage := &storageSpy{}
	ctx := &contextTest{}
	handler := &shortenHandler{
		storage: storage,
	}

	handler.Handler(ctx)

	if len(storage.shortenKey) == 0 {
		t.Errorf("empty shorten key")
	}

	if len(storage.originalURL) > 0 {
		t.Errorf("too long")
	}
}
