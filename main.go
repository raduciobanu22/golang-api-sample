package main

func main() {
	service := new(FXService)
	service.Init("")
	service.Run(":8000")
}