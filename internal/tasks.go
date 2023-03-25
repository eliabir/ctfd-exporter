package main

type TaskReturn struct {
	Success bool   `json:"success"`
	Data    []Task `json:"data"`
}

type Task struct {
}
