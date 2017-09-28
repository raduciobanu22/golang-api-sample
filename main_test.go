package main

import (
	"os"
	"testing"
	"net/http"
	"net/http/httptest"
)

var service FXService

func TestMain(m *testing.M) {
	service = FXService{}
	service.Init("0f3040b808284cdda36fc571698335c5")
	os.Exit(m.Run())
}

func TestGetAllCurrentRates(t *testing.T) {
	request, _ := http.NewRequest("GET", "/current_rates", nil)
	response := httptest.NewRecorder()
	service.Router.ServeHTTP(response, request)
	
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code 200. Got %d\n", response.Code)
	}
}

func TestGetSingleRate(t *testing.T) {
	request, _ := http.NewRequest("GET", "/current_rates?currency=SGD", nil)
	response := httptest.NewRecorder()
	service.Router.ServeHTTP(response, request)
	
	if response.Code != http.StatusOK {
		t.Errorf("Expected status code 200. Got %d\n", response.Code)
	}
}