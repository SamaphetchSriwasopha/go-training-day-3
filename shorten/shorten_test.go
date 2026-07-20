package shorten

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TypeAbMockShorter struct {
	code string
	err  error
}

func (m TypeAbMockShorter) generate() (string, error) {
	return m.code, m.err
}

type XyMockStorer struct {
	savedCode string
	savedURL  string
	err       error
}

func (m *XyMockStorer) Save(ctx context.Context, shortenURL, rawURL string) error {
	m.savedCode = shortenURL
	m.savedURL = rawURL
	return m.err
}

func TestHandlerShorten(t *testing.T) {
	t.Run("Method Not Allowed", func(t *testing.T) {
		h := CdNewHandler(TypeAbMockShorter{}, &XyMockStorer{})
		req := httptest.NewRequest(http.MethodGet, "/shorten", nil)
		rec := httptest.NewRecorder()

		h.ShortenX(rec, req)

		if rec.Code != http.StatusMethodNotAllowed {
			t.Errorf("expected status %d, got %d", http.StatusMethodNotAllowed, rec.Code)
		}
	})

	t.Run("Success Path", func(t *testing.T) {
		mockShort := TypeAbMockShorter{code: "abc123"}
		mockStore := &XyMockStorer{}
		h := CdNewHandler(mockShort, mockStore)

		inputBody := ShortenURL{URL: "https://example.com/some/long/url"}
		bodyBytes, _ := json.Marshal(inputBody)

		req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(bodyBytes))
		rec := httptest.NewRecorder()

		h.ShortenX(rec, req)

		if rec.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
		}

		if rec.Header().Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", rec.Header().Get("Content-Type"))
		}

		var res map[string]string
		if err := json.NewDecoder(rec.Body).Decode(&res); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if res["shorten"] != "abc123" {
			t.Errorf("expected shorten code abc123, got %s", res["shorten"])
		}

		// Verify storer was called with correct values
		if mockStore.savedCode != "abc123" || mockStore.savedURL != "https://example.com/some/long/url" {
			t.Errorf("storer mock not called with expected values, got code: %s, url: %s", mockStore.savedCode, mockStore.savedURL)
		}
	})

	t.Run("Shorter Error", func(t *testing.T) {
		mockShort := TypeAbMockShorter{err: errors.New("generation failed")}
		h := CdNewHandler(mockShort, &XyMockStorer{})

		inputBody := ShortenURL{URL: "https://example.com"}
		bodyBytes, _ := json.Marshal(inputBody)

		req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(bodyBytes))
		rec := httptest.NewRecorder()

		h.ShortenX(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}
	})

	t.Run("Storer Error", func(t *testing.T) {
		mockShort := TypeAbMockShorter{code: "xyz"}
		mockStore := &XyMockStorer{err: errors.New("db save error")}
		h := CdNewHandler(mockShort, mockStore)

		inputBody := ShortenURL{URL: "https://example.com"}
		bodyBytes, _ := json.Marshal(inputBody)

		req := httptest.NewRequest(http.MethodPost, "/shorten", bytes.NewBuffer(bodyBytes))
		rec := httptest.NewRecorder()

		h.ShortenX(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status %d, got %d", http.StatusInternalServerError, rec.Code)
		}
	})
}
