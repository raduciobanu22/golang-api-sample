package main

import (
	"net/http"
	"log"
	"time"
	"encoding/json"
	"strings"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
)

const CacheExpiration = 60 //minutes

// Structure used to map the response
// from OpenExchange
type OpenExchangeLatest struct {
	Disclaimer string
	License string
	Timestamp int64
	Base string
	Rates map[string]float32
}

type FXService struct {
	Router *mux.Router
	AppID string
	OpenExchangeUrl string
	CacheService *cache.Cache
}

// Initialize the FXService by instantiating the Router and CacheService components,
// registering the API route and setting the OpenExchangeRates AppID
func (service *FXService) Init(openExchangeAppID string, openExchangeUrl string) {
	service.Router = mux.NewRouter()
	service.CacheService = cache.New(CacheExpiration * time.Minute, 1 * time.Minute)
	service.AppID = openExchangeAppID
	service.OpenExchangeUrl = openExchangeUrl

	service.Router.HandleFunc("/current_rates", service.GetCurrentRates).Methods("GET")
}

// Start the HTTP Service
func (service *FXService) Run(port string) {
	log.Println("Starting service")
	log.Fatal(http.ListenAndServe(port, service.Router))
}

// /current_rates API endpoint handler
func (service *FXService) GetCurrentRates(w http.ResponseWriter, r *http.Request) {
	params, ok := r.URL.Query()["currency"]
	var currency string

	if !ok || len(params) < 1 {
		currency = ""
	} else {
		currency = params[0]
	}

	rates, err := service.FetchRates(strings.ToUpper(currency))

	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	var response []byte
	response, err = json.MarshalIndent(rates, "", "    ")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

func (service *FXService) FetchRates(currency string) (map[string]float32, error) {
	//Check for rates in cache
	rates, found := service.CacheService.Get("rates")

	//Fetch new rates if there's nothing cached
	var err error
	if !found {
		log.Println("Fetching new rates")
		rates, err = service.GetNewRates()

		if err != nil || len(rates.(map[string]float32)) == 0 {
			return map[string]float32{}, err
		}

		//Cache new rates
		log.Println("Caching new rates")
		service.CacheService.Set("rates", rates, cache.DefaultExpiration)
	} else {
		log.Println("Retrieved rates from cache")
	}

	//Handle currency filter
	if currency != "" {
		var mappedRates = rates.(map[string]float32)
		rate := mappedRates[currency]
		return map[string]float32{currency: rate}, nil
	}
	return rates.(map[string]float32), nil
}

// Retrieve the latest fx rates from openexchange.org
func (service *FXService) GetNewRates() (map[string]float32, error) {
	var client = &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(service.OpenExchangeUrl + "?app_id=" + service.AppID)
	if err != nil || res.StatusCode != 200 {
		return nil, err
	}

	defer res.Body.Close()

	var dataDecoded OpenExchangeLatest
	decoder := json.NewDecoder(res.Body)
	err = decoder.Decode(&dataDecoded)

	if err != nil {
		return nil, err
	}

	return dataDecoded.Rates, nil
}