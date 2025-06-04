package tracker

import (
	"time"
)

type Task struct {
	Name      string        `json:"name"`
	StartTime time.Time     `json:"start_time"`
	EndTime   time.Time     `json:"end_time"`
	Duration  time.Duration `json:"duration"`
}

type Tracker struct {
	ActiveTask *Task  `json:"active_task"`
	History    []Task `json:"history"`
}
