package main

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type ScoreboardReturn struct {
	Success bool             `json:"success"`
	Data    []ScoreboardTeam `json:"data"`
}

type ScoreboardTeam struct {
	Pos         int    `json:"pos"`
	AccountID   int    `json:"account_id"`
	AccountURL  string `json:"account_url"`
	AccountType string `json:"account_type"`
	OauthID     int    `json:"oauth_id"`
	Name        string `json:"name"`
	Score       int    `json:"score"`
	Members     []ScoreboardTeamMember
}

type ScoreboardTeamMember struct {
	ID      int    `json:"id"`
	OauthID int    `json:"oauth_id"`
	Name    string `json:"name"`
	Score   int    `json:"score"`
}

func getScoreboard(apiKey string, apiEndpoint string) ScoreboardReturn {
	// Create a new HTTP request with the Authorization header
	req, err := http.NewRequest("GET", apiEndpoint+"/scoreboard", nil)
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

	// bodyBytes, err := io.ReadAll(resp.Body)
	// bodyString := string(bodyBytes)
	// fmt.Println(bodyString)

	var scoreboard ScoreboardReturn
	err = json.NewDecoder(resp.Body).Decode(&scoreboard)

	return scoreboard
}

func countScoreboardTeams(scoreboardC chan ScoreboardReturn) {
	go func() {
		scoreboard := <-scoreboardC

		scoreboardTeamsCount := len(scoreboard.Data)
		scoreboardTeams.Set(float64(scoreboardTeamsCount))
	}()
}

var (
	scoreboardTeams = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ctfd_scoreboard_teams_total",
		Help: "Number of teams on scoreboard",
	})
)
