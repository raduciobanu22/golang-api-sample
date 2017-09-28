package main

import (
	"os"
)

const OpenExLatestEndpoint = "https://openexchangerates.org/api/latest.json"
const AppID = ""

func main() {
	service := new(FXService)
	service.Init(os.Getenv("APP_ID"), OpenExLatestEndpoint)
	service.Run(":8000")
}