package rest_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/thenilesh/url-shortner/mocks"
	"github.com/thenilesh/url-shortner/rest"
	"github.com/thenilesh/url-shortner/svc"
)

func TestShortURLHandler_Create_Success(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	mockSvc := new(mocks.URLShortner)
	handler := rest.NewShortURLHandler(log, mockSvc)

	router := mux.NewRouter()
	router.HandleFunc("/shorturl", handler.Create).Methods(http.MethodPost)
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
}

func TestShortURLHandler_Create_BadRequest(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	mockSvc := new(mocks.URLShortner)
	handler := rest.NewShortURLHandler(log, mockSvc)

	router := mux.NewRouter()
	router.HandleFunc("/shorturl", handler.Create).Methods(http.MethodPost)
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
}

func TestShortURLHandler_Create_conflict(t *testing.T) {
	log := logrus.New()
	log.SetLevel(logrus.ErrorLevel)

	mockSvc := new(mocks.URLShortner)
	handler := rest.NewShortURLHandler(log, mockSvc)

	router := mux.NewRouter()
	router.HandleFunc("/shorturl", handler.Create).Methods(http.MethodPost)
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
}
