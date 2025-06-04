package main

import (
	"fmt"
	"go-time-tracker/tracker"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: go-timer-tracker [start|stop|log|export]")
	}

	t, _ := tracker.LoadTracker()
	cmd := os.Args[1]

	switch cmd {
	case "start":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task name")
			return
		}
		name := os.Args[2]
		t.ActiveTask = &tracker.Task{
			Name:      name,
			StartTime: time.Now(),
		}
		fmt.Println("Started:", name)

		err := tracker.SaveTracker(t)
		if err != nil {
			fmt.Println("Failed to save tracker:", err)
		}

	case "stop":
		if t.ActiveTask == nil {
			fmt.Println("No task is currently running")
			return
		}
		t.ActiveTask.EndTime = time.Now()
		t.ActiveTask.Duration = t.ActiveTask.EndTime.Sub(t.ActiveTask.StartTime)
		fmt.Printf("Stopped %s (Duration: %s)\n", t.ActiveTask.Name, t.ActiveTask.Duration)
		t.History = append(t.History, *t.ActiveTask)
		t.ActiveTask = nil

	case "log":
		for _, task := range t.History {
			fmt.Printf("%s - %s\n", task.Name, task.Duration)
		}

	default:
		fmt.Println("Unkown command:", cmd)

	}

	tracker.SaveTracker(t)
}
