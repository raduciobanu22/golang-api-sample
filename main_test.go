package main

import (
	"os"
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/stretchr/testify/assert"
	"log"
	"encoding/json"
)

var service FXService

func TestMain(m *testing.M) {
	service = FXService{}
	service.Init("")
	os.Exit(m.Run())
}

func TestGetAllCurrentRates(t *testing.T) {
	request, _ := http.NewRequest("GET", "/current_rates", nil)
	response := httptest.NewRecorder()
	service.Router.ServeHTTP(response, request)
	assert.Equal(t, response.Code, http.StatusOK, "Expected status code is 200")

	var data map[string]float32
	decoder := json.NewDecoder(response.Body)
	decoder.Decode(&data)

	assert.True(t, len(data) > 1, "Expected more than one currency rate")
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
	assert.NotEqual(t, data[currency], 0,  "Currency rate shouldn't be 0")
	log.Println(data[currency])
}

func TestGetNewRates(t *testing.T) {
	//TODO Mock/Hijack the http request to openexchangerates
}