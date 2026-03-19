package model

import "github.com/google/uuid"

type Status string

const (
	StatusPending Status = "PENDING"
	StatusDone    Status = "DONE"
)

type Task struct {
	Id      string
	Status  Status
	Payload string
}

func NewTask(payload string) *Task {
	return &Task{Id: uuid.New().String(), Status: StatusPending, Payload: payload}
}

func (t *Task) ToJsonModel() TaskGetResponseJson {
	return TaskGetResponseJson{
		Id:     t.Id,
		Status: t.Status,
	}
}

func TasksToJsonArrayModel(tasks []Task) []TaskGetResponseJson {
	var result []TaskGetResponseJson

	for _, task := range tasks {
		result = append(result, task.ToJsonModel())
	}

	return result
}

type TasksStats struct {
	Submitted  int
	Completed  int
	InProgress int
}

func (t *TasksStats) ToJsonModel() TasksStatsResponseJson {
	return TasksStatsResponseJson{
		Submitted:  t.Submitted,
		Completed:  t.Completed,
		InProgress: t.InProgress,
	}
}
