package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"go-time-tracker/tracker"
	"os"
	"strings"
	"time"

	"github.com/olekukonko/tablewriter"
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
		name := strings.Join(os.Args[2:], " ")
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
		// for _, task := range t.History {
		// 	fmt.Printf("%s - %s\n", task.Name, task.Duration)
		// }
		if len(t.History) == 0 {
			fmt.Println("No tasks tracked yet.")
			return
		}

		table := tablewriter.NewWriter(os.Stdout)
		table.SetHeader([]string{"Task", "Start", "End", "Duration"})

		for _, task := range t.History {
			table.Append([]string{
				task.Name,
				task.StartTime.Format("01-02-2006 03:04:05 PM"),
				task.EndTime.Format("01-02-2006 03:04:05 PM"),
				task.Duration.Truncate(time.Second).String(),
			})
		}

		table.Render()

	case "export":
		exportAll := false
		condense := false
		for _, arg := range os.Args[2:] {
			if arg == "-all" {
				exportAll = true
			}
			if arg == "-condense" {
				condense = true
			}
		}

		file, err := os.Create("time-tracker-export.csv")
		if err != nil {
			fmt.Println("Failed to create time-tracker-export.csv")
			return
		}
		defer file.Close()

		writer := csv.NewWriter(file)
		defer writer.Flush()

		if condense {
			totals := make(map[string]time.Duration)
			for _, task := range t.History {
				totals[task.Name] += task.Duration
			}

			writer.Write([]string{"Task Name,", "Total Duration"})
			for name, dur := range totals {
				writer.Write([]string{name, dur.String()})
			}
			fmt.Println("Exported condensed data to time-tracker-export.csv")
			return
		}

		if exportAll {
			writer.Write([]string{"Task Name", "Start Time", "End Time", "Duration"})
			for _, task := range t.History {
				writer.Write([]string{
					task.Name,
					task.StartTime.Format(time.RFC3339),
					task.EndTime.Format(time.RFC3339),
					task.Duration.String(),
				})
			}
			fmt.Println("Exported full data to time-tracker-export.csv")
			return
		}

		fmt.Println("Please provide one of: -all or -condense")

	case "clear":
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Are you sure you want to clear all tracking data? (y/N): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input != "y" && input != "yes" {
			fmt.Println("Aborted.")
			return
		}

		t.ActiveTask = nil
		t.History = nil

		err := tracker.SaveTracker(t)
		if err != nil {
			fmt.Println("Failed to save tracker:", err)
		} else {
			fmt.Println("All tracking has been cleared")
		}

	default:
		fmt.Println("Unko	wn command:", cmd)

	}

	tracker.SaveTracker(t)
}
