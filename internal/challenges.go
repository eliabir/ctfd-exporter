package main

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ChallengeReturn struct {
	Success bool        `json:"success"`
	Data    []Challenge `json:"data"`
}

type Challenge struct {
	ID         int      `json:"id"`
	Type       string   `json:"type"`
	Name       string   `json:"name"`
	Value      int      `json:"value"`
	Solves     int      `json:"solves"`
	SolvedByMe string   `json:"solved_by_me"`
	Category   string   `json:"category"`
	Tags       []string `json:"tags"`
	Template   string   `json:"template"`
	Script     string   `json:"script"`
}

func getChallenges(apiKey string, apiEndpoint string) []Challenge {

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

	var challenges ChallengeReturn
	err = json.NewDecoder(resp.Body).Decode(&challenges)

	return challenges.Data
}

func countChallenges(apiKey string, apiEndpoint string) {
	go func() {
		for range Ticker.C {
			challenges := getChallenges(apiKey, apiEndpoint)

			challengesCount := len(challenges)
			challengesTotal.Set(float64(challengesCount))

			time.Sleep(5 * time.Second)
		}
	}()
}

var (
	challengesTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ctfd_challenges_total",
		Help: "The total number of challenges",
	})
)
