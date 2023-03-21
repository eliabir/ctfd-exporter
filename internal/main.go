package main

import (
	"fmt"
	"log"
	"os"
)

func main() {
	// Get CTFd api token
	apiKey := os.Getenv("CTFD_API")
	if apiKey == "" {
		log.Fatalln("Could not retrieve API token from environment variable")
	}

	apiEndpoint := "http://localhost:8000/api/v1"

	challenges := getChallenges(apiKey, apiEndpoint)

	fmt.Printf(challenges.Category)

	// countChallenges(apiKey, apiEndpoint)

	// http.Handle("/metrics", promhttp.Handler())
	// http.ListenAndServe(":2112", nil)
}
