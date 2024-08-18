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
			"valid, true",
			map[string]interface{}{
				"mode": true,
			},
			http.StatusOK,
		},
		{
			"valid, false",
			map[string]interface{}{
				"mode": false,
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

func TestServer_handlerGetAllUpsParams(t *testing.T) {
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

func TestServer_handlerUpdateUpsParams(t *testing.T) {
	imitator := imitator.New(nil, model.TestConfig(t))
	s := newServer(imitator)
	testCases := []struct {
		name         string
		prapare      func()
		payload      any
		expectedCode int
	}{
		{
			"invalid payload",
			nil,
			"invalid",
			http.StatusBadRequest,
		},
		{
			"invalid param",
			nil,
			map[string]string{
				"input_ac_voltage": "invalid",
			},
			http.StatusBadRequest,
		},
		{
			"invalid, auto mode",
			nil,
			map[string]interface{}{
				"input_ac_voltage": 230,
			},
			http.StatusForbidden,
		},
		{
			"valid, InputAcVoltage",
			func() {
				imitator.SetMode(false)
			},
			map[string]interface{}{
				"input_ac_voltage": 230,
			},
			http.StatusOK,
		},
		{
			"valid, InputAcCurrent",
			nil,
			map[string]interface{}{
				"input_ac_current": 10,
			},
			http.StatusOK,
		},
		{
			"valid, BatGroupVoltage",
			nil,
			map[string]interface{}{
				"bat_group_voltage": 55,
			},
			http.StatusOK,
		},
		{
			"valid, BatGroupCurrent",
			nil,
			map[string]interface{}{
				"bat_group_current": 20,
			},
			http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.prapare != nil {
				tc.prapare()
			}
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPatch, "/imitator/ups/params", b)
			s.router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_handlerUpdateBattery(t *testing.T) {
	imitator := imitator.New(nil, model.TestConfig(t))
	s := newServer(imitator)
	testCases := []struct {
		name         string
		prapare      func()
		payload      any
		batId        string
		expectedCode int
	}{
		{
			"invalid payload",
			nil,
			"invalid",
			"0",
			http.StatusBadRequest,
		},
		{
			"invalid bat id",
			nil,
			map[string]any{
				"voltage": 12,
			},
			"invalid",
			http.StatusBadRequest,
		},
		{
			"invalid param",
			nil,
			map[string]any{
				"voltage": "invalid",
			},
			"0",
			http.StatusBadRequest,
		},
		{
			"invalid, auto mode",
			nil,
			map[string]interface{}{
				"voltage": 12,
			},
			"0",
			http.StatusForbidden,
		},
		{
			"valid, voltage",
			func() {
				imitator.SetMode(false)
			},
			map[string]interface{}{
				"voltage": 12,
			},
			"0",
			http.StatusOK,
		},
		{
			"invalid bat id, range over",
			nil,
			map[string]any{
				"voltage": 12,
			},
			"4",
			http.StatusUnprocessableEntity,
		},

		{
			"valid, voltage",
			nil,
			map[string]interface{}{
				"voltage": 12,
			},
			"0",
			http.StatusOK,
		},
		{
			"valid, resist",
			nil,
			map[string]interface{}{
				"resist": 6,
			},
			"0",
			http.StatusOK,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.prapare != nil {
				tc.prapare()
			}
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPatch, "/imitator/ups/"+tc.batId, b)
			s.router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_handlerUpdateAlarms(t *testing.T) {
	imitator := imitator.New(nil, model.TestConfig(t))
	s := newServer(imitator)
	testCases := []struct {
		name         string
		prapare      func()
		payload      interface{}
		expectedCode int
	}{
		{
			"invalid payload",
			nil,
			"invalid",
			http.StatusBadRequest,
		},
		{
			"invalid param",
			nil,
			map[string]string{
				"upc_in_battery_mode": "invalid",
			},
			http.StatusBadRequest,
		},
		{
			"invalid, auto mode",
			nil,
			map[string]interface{}{
				"upc_in_battery_mode": true,
			},
			http.StatusForbidden,
		},
		{
			"valid, UpcInBatteryMode",
			func() {
				imitator.SetMode(false)
			},
			map[string]interface{}{
				"upc_in_battery_mode": true,
			},
			http.StatusOK,
		},
		{
			"valid, LowBattery",
			nil,
			map[string]interface{}{
				"low_battery": true,
			},
			http.StatusOK,
		},
		{
			"valid, Overload",
			nil,
			map[string]interface{}{
				"overload": true,
			},
			http.StatusOK,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.prapare != nil {
				tc.prapare()
			}
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPatch, "/imitator/ups/alarms", b)
			s.router.ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}
