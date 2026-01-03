package model

type TasksPostRequestJson struct {
	Payload string `json:"payload"`
}

type TaskGetResponseJson struct {
	Id     string `json:"id"`
	Status Status `json:"status"`
}

type TasksStatsResponseJson struct {
	Submitted  int `json:"submitted"`
	Completed  int `json:"completed"`
	InProgress int `json:"in_progress"`
}
