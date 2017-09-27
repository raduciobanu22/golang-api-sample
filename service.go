package main

import (
	"net/http"
	"log"
	"github.com/gorilla/mux"
	"github.com/patrickmn/go-cache"
	"time"		
	"encoding/json"
	"strings"
)

const OpenExLatestEndpoint = "https://openexchangerates.org/api/latest.json"
const CacheExpiration = 60 //minutes

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
	CacheService *cache.Cache
}

func (service *FXService) Init(openExchangeAppID string) {
	service.Router = mux.NewRouter()
	service.Router.HandleFunc("/current_rates", service.GetCurrentRates).Methods("GET")
	service.CacheService = cache.New(CacheExpiration * time.Minute, 1 * time.Minute)
	service.AppID = openExchangeAppID
}

func (service *FXService) Run(port string) {
	log.Fatal(http.ListenAndServe(port, service.Router))
}

func (service *FXService) GetCurrentRates(w http.ResponseWriter, r *http.Request) {
	currency, ok := r.URL.Query()["currency"]
	rates, err := service.FetchRates()
	
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	if !ok || len(currency) < 1 {		
		ratesJson, err := json.MarshalIndent(rates, "", "    ")
		
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(ratesJson)
	} else {					
		currency := strings.ToUpper(currency[0])
		value := rates[currency]
		if value == 0 {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)			
			return
		}		

		var rate = make(map[string]float32)
		rate[currency] = value
		rateJson, err := json.MarshalIndent(rate, "", "    ")
		
		if (err != nil) {			
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}		

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(rateJson)
	}
}

func (service *FXService) FetchRates() (map[string]float32, error) {
	//Check for rates in cache
	rates, found := service.CacheService.Get("rates")	
	if found {
		log.Println("Rates retrieved from cache")
		return rates.(map[string]float32), nil
	}
	
	//Fetch new rates if there's nothing cached
	log.Println("Fetching new rates")	

	var client = &http.Client{Timeout: 10 * time.Second}
	res, err := client.Get(OpenExLatestEndpoint + "?app_id=" + service.AppID)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	var dataDecoded OpenExchangeLatest
	decoder := json.NewDecoder(res.Body)	
	err = decoder.Decode(&dataDecoded)
	 
	if err != nil {
		return nil, err
	}

	//Cache new rates
	service.CacheService.Set("rates", dataDecoded.Rates, cache.DefaultExpiration)
	log.Println("Cached new rates")

	return dataDecoded.Rates, nil
}