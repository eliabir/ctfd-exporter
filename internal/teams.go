package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

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

func countTeams(apiKey string, apiEndpoint string) {
	go func() {
		for {
			teams := getTeams(apiKey, apiEndpoint)

			fmt.Println(teams[1].Name)

			teamsCount := len(teams)
			teamsTotal.Set(float64(teamsCount))

			time.Sleep(5 * time.Second)
		}
	}()
}

var (
	teamsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ctfd_teams_total",
		Help: "The total number of teams",
	})
)
