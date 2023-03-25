package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Global ticker
var Ticker *time.Ticker

func main() {
	// Get CTFd api token
	apiKey := os.Getenv("CTFD_API")
	if apiKey == "" {
		log.Fatalln("Could not retrieve API token from environment variable")
	}

	apiEndpoint := "https://ctf.uia.no/api/v1"

	// Create ticker for every second
	Ticker = time.NewTicker(1 * time.Second)

	countChallenges(apiKey, apiEndpoint)
	countTeams(apiKey, apiEndpoint)
	countScoreboardTeams(apiKey, apiEndpoint)
	scoreTeams(apiKey, apiEndpoint)
	countUsers(apiKey, apiEndpoint)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
