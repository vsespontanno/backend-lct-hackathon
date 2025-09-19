package main

import (
	"black-pearl/backend-hackathon/internal/app"
	"log"
)

func main() {
	a := app.NewApp()
	log.Println("Server started")
	if err := a.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
