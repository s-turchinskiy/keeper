package main

import (
	"github.com/joho/godotenv"
	"github.com/s-turchinskiy/keeper/internal/client"
	"log"
)

// go run -ldflags "-X buildinfo.BuildVersion=v1.0.1 -X buildinfo.BuildDate=18.12.2025 -X buildinfo.BuildCommit=Comment"
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {

	err := godotenv.Load("./.env")
	if err != nil {
		_ = godotenv.Load("./cmd/client/.env")
	}

	app, err := client.NewApp()
	if err != nil {
		log.Fatal(err)
	}
	defer app.Close()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}
}
