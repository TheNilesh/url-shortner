// BEGIN: 5f7f1d5d7b8a
package rest_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thenilesh/url-shortner/rest"
	"github.com/thenilesh/url-shortner/svc"
)

type mockURLShortner struct {
	mock.Mock
}

func (m *mockURLShortner) CreateShortPath(ctx context.Context, shortPath string, targetURL string) (string, error) {
	args := m.Called(ctx, shortPath, targetURL)
	return args.String(0), args.Error(1)
}

func (m *mockURLShortner) GetTargetURL(ctx context.Context, shortPath string) (string, error) {
	args := m.Called(ctx, shortPath)
	return args.String(0), args.Error(1)
}

func TestShortURLHandler_Create(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	mockSvc := new(mockURLShortner)
	handler := rest.NewShortURLHandler(log, mockSvc)

	router := mux.NewRouter()
	router.HandleFunc("/shorturl", handler.Create).Methods(http.MethodPost)

	t.Run("success", func(t *testing.T) {
		shortURL := rest.ShortURL{
			ShortPath: "test",
			TargetURL: "http://example.com",
		}
		body, _ := json.Marshal(shortURL)
		req, _ := http.NewRequest(http.MethodPost, "/shorturl", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		expectedShortPath := "test"
		mockSvc.On("CreateShortPath", mock.Anything, shortURL.ShortPath, shortURL.TargetURL).Return(expectedShortPath, nil)

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, fmt.Sprintf("/%s", expectedShortPath), rr.Header().Get("Location"))

		var resp rest.Response
		json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Equal(t, fmt.Sprintf("Created short URL: /%s", expectedShortPath), resp.Message)
	})

	t.Run("bad request", func(t *testing.T) {
		shortURL := rest.ShortURL{
			ShortPath: "test",
			TargetURL: "http://example.com",
		}
		body, _ := json.Marshal(shortURL)
		req, _ := http.NewRequest(http.MethodPost, "/shorturl", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		mockSvc.On("CreateShortPath", mock.Anything, shortURL.ShortPath, shortURL.TargetURL).Return("", &svc.ErrValidation{})

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)
	})

	t.Run("conflict", func(t *testing.T) {
		shortURL := rest.ShortURL{
			ShortPath: "test",
			TargetURL: "http://example.com",
		}
		body, _ := json.Marshal(shortURL)
		req, _ := http.NewRequest(http.MethodPost, "/shorturl", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		mockSvc.On("CreateShortPath", mock.Anything, shortURL.ShortPath, shortURL.TargetURL).Return("", &svc.ErrConflict{})

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusConflict, rr.Code)
	})

	t.Run("server error", func(t *testing.T) {
		shortURL := rest.ShortURL{
			ShortPath: "test",
			TargetURL: "http://example.com",
		}
		body, _ := json.Marshal(shortURL)
		req, _ := http.NewRequest(http.MethodPost, "/shorturl", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		mockSvc.On("CreateShortPath", mock.Anything, shortURL.ShortPath, shortURL.TargetURL).Return("", &svc.ErrServerError{})

		rr := httptest.NewRecorder()
		router.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)

		var resp rest.Response
		json.Unmarshal(rr.Body.Bytes(), &resp)
		assert.Equal(t, "Something went wrong", resp.Message)
	})
}
