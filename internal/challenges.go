package main

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Challenge struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

func getChallenges(apiKey string, apiEndpoint string) Challenge {

	// Create a new HTTP request with the Authorization header
	req, err := http.NewRequest("GET", apiEndpoint+"/challenges", nil)
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Token "+apiKey)
	req.Header.Set("Content-Type", "application/json")

	// Send the HTTP request and retrieve the response
	client := http.DefaultClient
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	var challenges Challenge
	err = json.NewDecoder(resp.Body).Decode(&challenges)

	return challenges
}

// func countChallenges(apiKey string, apiEndpoint string) {
// 	go func() {

// 		fmt.Printf(challenges["name"])

// 		challengesCount := len(challenges)
// 		challengesTotal.Set(float64(challengesCount))
// 	}()
// }

var (
	challengesTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ctfd_challenges_total",
		Help: "The total number of challenges",
	})
)
