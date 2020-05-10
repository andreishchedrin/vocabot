package main

import (	
	"os"
	"github.com/joho/godotenv"
	"log"
)

func init() {
    // loads values from .env into the system
    if err := godotenv.Load(); err != nil {
        log.Print("No .env file found")
    }
}

func main() {
	var token string 
	token, _ = os.LookupEnv("TG_TOKEN")	
	telegramBot(token)
}
