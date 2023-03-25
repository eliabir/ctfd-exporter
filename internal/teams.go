package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type TeamReturn struct {
	Meta    TeamMeta `json:"meta"`
	Success bool     `json:"success"`
	Data    []Team   `json:"data"`
}

type TeamMeta struct {
	Pagination TeamPagination `json:"pagination"`
}

type TeamPagination struct {
	Page    bool   `json:"page"`
	Next    string `json:"next"`
	Prev    string `json:"prev"`
	Pages   int    `json:"pages"`
	PerPage int    `json:"per_page"`
	Total   int    `json:"total"`
}

type Team struct {
	ID          int      `json:"id"`
	Hidden      bool     `json:"hidden"`
	Country     string   `json:"country"`
	Email       string   `json:"email"`
	CaptainID   int      `json:"captain_id"`
	Secret      string   `json:"secret"`
	Created     string   `json:"created"`
	Affiliation string   `json:"affiliation"`
	Name        string   `json:"name"`
	Bracket     string   `json:"bracket"`
	Fields      []string `json:"fields"`
	Banned      string   `json:"banned"`
	OauthID     int      `json:"oauth_id"`
	Website     string   `json:"website"`
}

type TeamSingleReturn struct {
	Success bool       `json:"success"`
	Data    TeamSingle `json:"data"`
}

type TeamSingle struct {
	CaptainID   int      `json:"captain_id"`
	Hidden      bool     `json:"hidden"`
	Country     string   `json:"country"`
	Created     string   `json:"created"`
	Affiliation string   `json:"affiliation"`
	Name        string   `json:"name"`
	Bracket     string   `json:"bracket"`
	Members     []int    `json:"members"`
	Website     string   `json:"website"`
	ID          int      `json:"id"`
	Email       string   `json:"email"`
	Secret      string   `json:"secret"`
	Fields      []string `json:"fields"`
	Banned      bool     `json:"banned"`
	OauthID     int      `json:"oauth_id"`
	Place       string   `json:"place"`
	Score       int      `json:"score"`
}

func getTeams(apiKey string, apiEndpoint string) []Team {
	// Create a new HTTP request with the Authorization header
	req, err := http.NewRequest("GET", apiEndpoint+"/teams", nil)
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

	var teams TeamReturn
	err = json.NewDecoder(resp.Body).Decode(&teams)

	return teams.Data
}

func getTeam(apiKey string, apiEndpoint string, teamID int) TeamSingle {
	// Create a new HTTP request with the Authorization header
	req, err := http.NewRequest("GET", apiEndpoint+"/teams/"+strconv.Itoa(teamID), nil)
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

	var team TeamSingleReturn
	err = json.NewDecoder(resp.Body).Decode(&team)

	return team.Data
}

func countTeams(apiKey string, apiEndpoint string) {
	go func() {
		for range Ticker.C {
			teams := getTeams(apiKey, apiEndpoint)

			teamsCount := len(teams)
			teamsTotal.Set(float64(teamsCount))
		}
	}()
}

var (
	teamsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ctfd_teams_total",
		Help: "The total number of registered teams",
	})
)

func scoreTeams(apiKey string, apiEndpoint string) {
	go func() {
		for range Ticker.C {
			teams := getScoreboard(apiKey, apiEndpoint)

			for _, team := range teams.Data {
				teamScore.With(prometheus.Labels{"name": team.Name}).Set(float64(team.Score))
			}
		}
	}()
}

var (
	teamScore = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ctfd_team_score",
		Help: "Score per team",
	}, []string{"name"})
)
