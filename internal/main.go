package main

import (
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Global ticker
var Ticker *time.Ticker

func main() {
	// Get CTFd API token
	apiKey := os.Getenv("CTFD_API")
	if apiKey == "" {
		log.Fatalln("Could not retrieve API token from environment variable")
	}

	// Get CTFd URL
	ctfdUrl := os.Getenv("CTFD_URL")
	if ctfdUrl == "" {
		log.Fatalln("Could not retrieve CTFd URL from environment variable")
	}

	// Remove trailing /
	ctfdUrl = strings.TrimSuffix(ctfdUrl, "/")

	// Build API URL
	apiEndpoint := ctfdUrl + "/api/v1"

	// Create ticker for every second
	Ticker = time.NewTicker(1 * time.Second)

	usersC := make(chan UserReturn)
	teamsC := make(chan TeamReturn)
	scoreboardC := make(chan ScoreboardReturn)
	challengesC := make(chan ChallengeReturn)
	submissionsC := make(chan []SubmissionReturn)

	go func() {
		for range Ticker.C {
			usersC <- getUsers(apiKey, apiEndpoint)
			teamsC <- getTeams(apiKey, apiEndpoint)
			scoreboardC <- getScoreboard(apiKey, apiEndpoint)
			challengesC <- getChallenges(apiKey, apiEndpoint)
			submissionsC <- getSubmissionsAll(apiKey, apiEndpoint)
		}
	}()

	countUsers(usersC)
	countTeams(teamsC)
	countChallenges(challengesC)
	getSolvesChallenges(challengesC)
	getTotalPoints(challengesC)
	countScoreboardTeams(scoreboardC)
	scoreTeams(scoreboardC)
	scoreUser(scoreboardC)
	countSubmissions(submissionsC)

	http.Handle("/metrics", promhttp.Handler())
	http.ListenAndServe(":2112", nil)
}
