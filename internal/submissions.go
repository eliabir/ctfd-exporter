package main

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type SubmissionReturn struct {
	Meta    SubmissionMeta `json:"meta"`
	Success bool           `json:"success"`
	Data    []Submission   `json:"data"`
}

type SubmissionMeta struct {
	Pagination SubmissionPagination `json:"pagination"`
}

type SubmissionPagination struct {
	Page    int `json:"page"`
	Next    int `json:"next"`
	Prev    int `json:"prev"`
	Pages   int `json:"pages"`
	PerPage int `json:"per_page"`
	Total   int `json:"total"`
}

type Submission struct {
	ID          int            `json:"id"`
	UserID      int            `json:"user_id"`
	IP          string         `json:"ip"`
	Date        string         `json:"date"`
	Type        string         `json:"type"`
	TeamID      int            `json:"team_id"`
	Team        SubmissionTeam `json:"team"`
	Challenge   SubmissionChallenge
	Provided    string         `json:"provided"`
	User        SubmissionUser `json:"user"`
	ChallengeID int            `json:"challenge_id"`
}

type SubmissionTeam struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type SubmissionChallenge struct {
	ID       int    `json:"id"`
	Value    int    `json:"value"`
	Category string `json:"category"`
	Name     string `json:"name"`
}

type SubmissionUser struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// func getSubmissions(apiKey string, apiEndpoint string) SubmissionReturn {

// 	// Create a new HTTP request with the Authorization header
// 	req, err := http.NewRequest("GET", apiEndpoint+"/submissions", nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	req.Header.Set("Authorization", "Token "+apiKey)
// 	req.Header.Set("Content-Type", "application/json")

// 	// Send the HTTP request and retrieve the response
// 	client := http.DefaultClient
// 	resp, err := client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	var submissions SubmissionReturn
// 	err = json.NewDecoder(resp.Body).Decode(&submissions)
// 	if err != nil {
// 		panic(err)
// 	}

// 	pages_total := submissions.Meta.Pagination.Pages
// 	per_page := submissions.Meta.Pagination.PerPage
// 	submissions_total := per_page * pages_total

// 	fmt.Println(submissions_total)

// 	// Create a new HTTP request with the Authorization header
// 	req, err = http.NewRequest("GET", apiEndpoint+"/submissions?per_page="+strconv.Itoa(submissions_total), nil)
// 	if err != nil {
// 		panic(err)
// 	}
// 	req.Header.Set("Authorization", "Token "+apiKey)
// 	req.Header.Set("Content-Type", "application/json")

// 	// Send the HTTP request and retrieve the response
// 	client = http.DefaultClient
// 	resp, err = client.Do(req)
// 	if err != nil {
// 		panic(err)
// 	}
// 	defer resp.Body.Close()

// 	err = json.NewDecoder(resp.Body).Decode(&submissions)
// 	if err != nil {
// 		panic(err)
// 	}

// 	return submissions
// }

func getSubmissionsAll(apiKey string, apiEndpoint string) []SubmissionReturn {
	var allSubmissions []SubmissionReturn
	per_page := 100
	page := 1

	for {
		// Create a new HTTP request with the Authorization header
		req, err := http.NewRequest("GET", apiEndpoint+"/submissions"+"?per_page="+strconv.Itoa(per_page)+"&page="+strconv.Itoa(page), nil)
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

		var submissions SubmissionReturn
		err = json.NewDecoder(resp.Body).Decode(&submissions)
		if err != nil {
			panic(err)
		}

		allSubmissions = append(allSubmissions, submissions)

		if submissions.Meta.Pagination.Next == 0 {
			break
		} else {
			page++
		}
	}

	return allSubmissions
}

func countSubmissions(submissionC chan []SubmissionReturn) {
	go func() {
		for {
			submissions := <-submissionC

			type Submission struct {
				Category string
				Solves   int
				Fails    int
			}

			submissionsMap := make(map[string]Submission)

			var uniqueIPsSlice []string

			var submissionSolvesCount float64
			var submissionFailsCount float64

			for _, submissions := range submissions {
				for _, submission := range submissions.Data {
					challengeName := submission.Challenge.Name

					if submission.Type == "correct" {
						submissionSolvesCount++
						submissionsMap[challengeName] = Submission{
							Category: submission.Challenge.Category,
							Solves:   submissionsMap[challengeName].Solves + 1,
							Fails:    submissionsMap[challengeName].Fails,
						}
					} else if submission.Type == "incorrect" {
						submissionFailsCount++
						submissionsMap[challengeName] = Submission{
							Category: submission.Challenge.Category,
							Solves:   submissionsMap[challengeName].Solves,
							Fails:    submissionsMap[challengeName].Fails + 1,
						}
					}

					if _, ok := 
				}
			}

			submissionsTotal.Set(float64(submissions[0].Meta.Pagination.Total))
			submissionsSolves.Set(float64(submissionSolvesCount))
			submissionsFails.Set(float64(submissionFailsCount))

			for name, submission := range submissionsMap {
				submissionSolves.With(prometheus.Labels{"name": name}).Set(float64(submission.Solves))
				submissionFails.With(prometheus.Labels{"name": name}).Set(float64(submission.Fails))
			}

		}
	}()
}

// func countSubmissionsOld(submissionsC chan []SubmissionReturn) {
// 	go func() {
// 		for {
// 			submissionsAll := <-submissionsC

// 			submissionsTotal.Set(float64(submissionsAll[0].Meta.Pagination.Total))

// 			var submissionsSolvesCount int = 0
// 			var submissionsFailsCount int = 0

// 			for _, submissions := range submissionsAll {
// 				for _, submission := range submissions.Data {
// 					if submission.Type == "correct" {
// 						submissionsSolvesCount++
// 					} else {
// 						submissionsFailsCount++
// 					}

// 				}
// 			}

// 			submissionsSolves.Set(float64(submissionsSolvesCount))
// 			submissionsFails.Set(float64(submissionsFailsCount))
// 		}
// 	}()
// }

var (
	submissionsTotal = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ctfd_submissions_total",
		Help: "Total number of submissions",
	})
)

var (
	submissionsSolves = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ctfd_submissions_solves_total",
		Help: "Amount of correct submissions",
	})
)

var (
	submissionsFails = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ctfd_submissions_fails_total",
		Help: "Amount of incorrect submissions",
	})
)

var (
	submissionSolves = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ctfd_submission_solves",
		Help: "Amount of correct submissions per task",
	}, []string{"name"})
)

var (
	submissionFails = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "ctfd_submission_fails",
		Help: "Amount of incorrect submissions per task",
	}, []string{"name"})
)

var (
	uniqueIPs = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "ctfd_unique_ips",
		Help: "Amount of unique IPs that have submitted flags",
	})
)
