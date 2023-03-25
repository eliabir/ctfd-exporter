package main

import (
	"encoding/json"
	"net/http"

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

func getChallenges(apiKey string, apiEndpoint string) ChallengeReturn {

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

	return challenges
}

func countChallenges(challengesC chan ChallengeReturn) {
	go func() {
		for {
			challenges := <-challengesC

			challengesCount := len(challenges.Data)
			challengesTotal.Set(float64(challengesCount))

			// for _, challenge := range challenges.Data {
			// 	solvesChallenges.With(prometheus.Labels{"name": challenge.Name}).Set(float64(challenge.Solves))
			// }
		}
	}()
}

var (
	challengesTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ctfd_challenges_total",
		Help: "The total number of challenges",
	})
)

func getSolvesChallenges(challengesC chan ChallengeReturn) {
	go func() {
		for {
			challenges := <-challengesC

			for _, challenge := range challenges.Data {
				solvesChallenges.With(prometheus.Labels{"name": challenge.Name}).Set(float64(challenge.Solves))
			}
		}
	}()
}

var (
	solvesChallenges = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ctfd_challenge_solves",
		Help: "The amount of solves per challenge",
	}, []string{"name"})
)
