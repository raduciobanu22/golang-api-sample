package main

func main() {
	service := new(FXService)
	service.Init("0f3040b808284cdda36fc571698335c5")
	service.Run(":8000")
}