package main

import (
	"os"
	"testing"
	"net/http"
	"net/http/httptest"
	"encoding/json"
	"github.com/stretchr/testify/assert"
)

var service FXService

func TestMain(m *testing.M) {
	service = FXService{}
	service.Init("123456", "http://example.com")
	os.Exit(m.Run())
}

func TestGetAllCurrentRates(t *testing.T) {
	testHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"rates": {"AUD": 1.27, "SGD": 1.36, "EUR": 0.85}}`))
	}))
	defer testHttpServer.Close()

	service.OpenExchangeUrl = testHttpServer.URL

	request, _ := http.NewRequest("GET", "/current_rates", nil)
	response := httptest.NewRecorder()
	service.Router.ServeHTTP(response, request)
	assert.Equal(t, response.Code, http.StatusOK, "Expected status code is 200")

	var data map[string]float32
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&data)

	assert.True(t, len(data) > 1, "Expected more than one currency rate")
	assert.True(t, data["SGD"] == 1.36, "SGD rate should be 1.36")
}

func TestGetSingleRate(t *testing.T) {
	currency := "SGD"
	request, _ := http.NewRequest("GET", "/current_rates?currency=" + currency, nil)
	response := httptest.NewRecorder()
	service.Router.ServeHTTP(response, request)
	assert.Equal(t, response.Code, http.StatusOK, "Expected status code is 200")

	var data map[string]float32
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&data)

	assert.Equal(t, len(data), 1, "Response should contain one currency rate")
	assert.True(t, data[currency] != 0,  "Currency rate shouldn't be 0")
	assert.True(t, data[currency] == 1.36, "SGD rate should be 1.36")
}

func TestGetNewRates(t *testing.T) {
	testHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"rates": {"AUD": 1.27, "SGD": 1.36, "EUR": 0.85}}`))
	}))
	defer testHttpServer.Close()

	service.OpenExchangeUrl = testHttpServer.URL
	rates, _ := service.GetNewRates()

	assert.True(t, rates["SGD"] == 1.36, "SGD rate should be 1.36")
	assert.True(t, rates["AUD"] == 1.27, "AUD rate should be 1.27")
	assert.True(t, rates["EUR"] == 0.85, "AUD rate should be 0.85")
}

func TestOpenExchangeIsUnavailable(t *testing.T) {
	testHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	}))
	defer testHttpServer.Close()

	service.OpenExchangeUrl = testHttpServer.URL
	rates, _ := service.GetNewRates()

	assert.True(t, len(rates) == 0, "Rates map should be empty")
}

func TestOpenExchangeDoesNotReturnJSON(t *testing.T) {
	testHttpServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Not a json"))
	}))
	defer testHttpServer.Close()

	service.OpenExchangeUrl = testHttpServer.URL
	rates, _ := service.GetNewRates()

	assert.True(t, len(rates) == 0, "Rates map should be empty")
}