package config

import (
	"github.com/joho/godotenv"
	"log"
)

type envData struct {
}

var Env *envData

func InitializeEnv() {
	//TODO: when setting up in docker have to change env
	err := godotenv.Load("C:\\Users\\Aritr\\GolandProjects\\ProductService\\.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	} else {
		log.Println(".env file loaded")
	}
}
