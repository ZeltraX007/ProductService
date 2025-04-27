package utils

import (
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

func GetHttpClient() (*http.Client, error) {
	// reading timeout value from env
	httpTimeout, err := strconv.Atoi(os.Getenv("HTTP_CLIENT_TIMEOUT"))
	if err != nil {
		log.Println("Error while reading http timeout from the env file: ", err)
		return nil, err
	}

	httpClient := &http.Client{Timeout: time.Duration(httpTimeout) * time.Second}
	return httpClient, nil
}
