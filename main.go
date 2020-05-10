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
	token, _ := os.LookupEnv("TG_TOKEN")	
	hook, _ := os.LookupEnv("WEB_HOOK_ADDRESS")
	cert, _ := os.LookupEnv("CERT")
	key, _ := os.LookupEnv("KEY")
	telegramBot(token, hook, cert, key)
}
