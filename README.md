# CTFd Exporter

Prometheus exporter for CTFd metrics written in Go.


## Prerequisites

* [Docker and Docker-Compose](https://docs.docker.com/get-docker/)

## Usage

Copy the `example.env`, rename it to `.env` and fill out the appropriate values.

``` ini
# API key for the CTFd instance
CTFD_API=<key>

# URL of the CTFd instance
CTFD_URL=<url>

# The polling rate of the exporter in seconds
POLL_RATE=<seconds>
```

Start the exporter with Docker Compose.

``` bash
docker compose up -d
```

## Metrics

The exporter provides the following metrics:

* `ctfd_challenge_solves`: The amount of solves per challenge
  * Labels:
    * `category`: Category of the challenge
    * `id`: Challenge ID
    * `name`: Name of the challenge
    * `value`: Amount of points the challenge is worth
  * Type: Gauge

* `ctfd_challenges_total`: The total number of challenges
  * Type: Gauge

* `ctfd_challenges_total_points`: The total amount of available points
  * Type: Gauge

* `ctfd_scoreboard_teams_total`: Number of teams on scoreboard
  * Type: Gauge

* `ctfd_submission_fails`: Number of incorrect submissions per task
  * Labels:
    * `name`: Name of tasbk
  * Type: Gauge

* `ctfd_submission_solves`: Number of incorrect submissions per task
  * Labels:
    * `name`: Name of task
  * Type: Gauge

* `ctfd_submissions_fails_total`: Total number of incorrect submissions
  * Type: Gauge

* `ctfd_submissions_solves_total`: Total number of correct submissions
  * Type: Gauge

* `ctfd_submissions_total`: Total number of submissions
  * Type: Gauge

* `ctfd_teams_total`: The total number of registered teams
  * Type: Gauge

* `ctfd_unique_ips`: Number of unique IPs that have submitted flags
    * Type: Gauge

* `ctfd_user_score`: Score per user
  * Labels:
    * `name`: Name of user
  * Type: Gauge

* `ctfd_users_total`: Total number of registered users
  * Type: Gauge
