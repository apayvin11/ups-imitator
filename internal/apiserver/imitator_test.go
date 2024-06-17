package apiserver

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/alex11prog/ups-imitator/internal/app/imitator"
	"github.com/alex11prog/ups-imitator/internal/app/model"
	"github.com/stretchr/testify/assert"
)

func TestServer_handlerGetMode(t *testing.T) {
	imitator := imitator.New(nil, model.TestConfig(t))
	s := newServer(imitator)
	testCase := struct {
		name         string
		uri          string
		expectedCode int
	}{"valid", "/imitator/mode", http.StatusOK}
	t.Run(testCase.name, func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, testCase.uri, nil)
		s.router.ServeHTTP(rec, req)
		assert.Equal(t, testCase.expectedCode, rec.Code)
	})
}

func TestServer_handlerUpdateMode(t *testing.T) {
	imitator := imitator.New(nil, model.TestConfig(t))
	s := newServer(imitator)
	testCases := []struct {
		name         string
		payload      interface{}
		expectedCode int
	}{
		{
			"invalid payload",
			"invalid",
			http.StatusBadRequest,
		},
		{
			"invalid param",
			map[string]string{
				"mode": "invalid",
			},
			http.StatusBadRequest,
		},
		{
			"valid",
			map[string]interface{}{
				"mode": true,
			},
			http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPut, "/imitator/mode", b)
			s.router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_handlerGetAllUps(t *testing.T) {
	imitator := imitator.New(nil, model.TestConfig(t))
	s := newServer(imitator)
	testCase := struct {
		name         string
		uri          string
		expectedCode int
	}{"valid", "/imitator/ups", http.StatusOK}
	t.Run(testCase.name, func(t *testing.T) {
		rec := httptest.NewRecorder()
		req, _ := http.NewRequest(http.MethodGet, testCase.uri, nil)
		s.router.ServeHTTP(rec, req)
		assert.Equal(t, testCase.expectedCode, rec.Code)
	})
}