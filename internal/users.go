package main

import (
	"encoding/json"
	"net/http"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type UserReturn struct {
	Meta    UserMeta `json:"meta"`
	Success bool     `json:"success"`
	Data    []User   `json:"data"`
}

type UserMeta struct {
	Pagination UserPagination `json:"pagination"`
}

type UserPagination struct {
	Page    int    `json:"page"`
	Next    string `json:"next"`
	Prev    string `json:"prev"`
	Pages   int    `json:"pages"`
	PerPage int    `json:"per_page"`
	Total   int    `json:"total"`
}

type User struct {
	ID          int      `json:"id"`
	Country     string   `json:"country"`
	TeamID      int      `json:"team_id"`
	Affiliation string   `json:"affiliation"`
	Bracket     string   `json:"bracket"`
	Name        string   `json:"name"`
	Fields      []string `json:"fields"`
	OauthID     int      `json:"oauth_id"`
	Website     string   `json:"website"`
}

func getUsers(apiKey string, apiEndpoint string) UserReturn {
	// Create a new HTTP request with the Authorization header
	req, err := http.NewRequest("GET", apiEndpoint+"/users", nil)
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

	var users UserReturn
	err = json.NewDecoder(resp.Body).Decode(&users)

	return users
}

func countUsers(usersC chan UserReturn) {
	go func() {
		for {
			users := <-usersC

			usersCount := users.Meta.Pagination.Total
			usersTotal.Set(float64(usersCount))
		}
	}()
}

var (
	usersTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ctfd_users_total",
		Help: "The total number of registered users",
	})
)

func scoreUser(scoreboardC chan ScoreboardReturn) {
	go func() {
		for {
			teams := <-scoreboardC

			for _, team := range teams.Data {
				for _, user := range team.Members {
					userScore.With(prometheus.Labels{"name": user.Name}).Set(float64(user.Score))
				}
			}
		}
	}()
}

var (
	userScore = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ctfd_user_score",
		Help: "Score per user",
	}, []string{"name"})
)
