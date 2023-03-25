package main

//   "id": 51,
//   "name": "Flagg 4",
//   "solves": 16

type TaskReturn struct {
	Success bool   `json:"success"`
	Data    []Task `json:"data"`
}

type Task struct {
}
