package tracker

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var noTaskMsg = "No tasks yet. Add your first one with command start"

// Root prints active task or provides usage info
func Root(cmd *cobra.Command, args []string) {
	task, err := getActiveTask()
	if err != nil {
		log.Fatalf("Can't recieve active task %s", err)
	}
	if task.Status == statusNaN {
		log.Println(noTaskMsg)
		return
	}

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Title", "Created", "Stopped", "Started", "Total Duration", "Session Duration", "Status"})
	t.AppendSeparator()
	t.AppendRow([]interface{}{task.Title, task.Created.Truncate(time.Second).Format(time.DateTime), task.stopped(),
		task.Started.Truncate(time.Second).Format(time.DateTime), task.TotalDuration.Truncate(time.Second), task.currentSession(), task.Status})
	t.Render()
}

func Start(cmd *cobra.Command, args []string) {
	taskTitle := cmd.Flag("task").Value.String()

	activeTitle, err := getActiveTaskTitle()
	if err != nil {
		log.Fatalf("Can't get current active task title %s", err)
	}

	if taskTitle == activeTitle {
		return
	}
	err = stopTask(activeTitle)
	if err != nil {
		log.Fatalf("Can't stop active task %s", err)
	}

	startTask(taskTitle)
	if err != nil {
		log.Fatalf("Can't start task %s", err)
	}

}

func Remove(cmd *cobra.Command, args []string) {
	taskTitle := cmd.Flag("task").Value.String()

	if err := removeTask(taskTitle); err != nil {
		log.Fatal(err)
	}
}

func Resume(cmd *cobra.Command, args []string) {
	title, err := getActiveTaskTitle()
	if err != nil {
		log.Fatal(err)
	}
	startTask(title)
}

func Stop(cmd *cobra.Command, args []string) {
	title, err := getActiveTaskTitle()
	if err != nil {
		log.Fatalf("Can't get active task title %s", title)
	}

	err = stopTask(title)
	if err != nil {
		log.Fatalf("Can't stop task %s", err)
	}
}

func List(cmd *cobra.Command, args []string) {
	tasks, err := getAll()
	if len(tasks) == 0 {
		fmt.Println(noTaskMsg)
		return
	}
	if err != nil {
		log.Fatal(err)
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Title", "Created", "Stopped", "Started", "Total Duration", "Session Duration", "Status"})
	t.AppendSeparator()
	for idx, task := range tasks {
		t.AppendRow([]interface{}{idx, task.Title, task.Created.Format(time.DateTime), task.stopped(), task.Started.Format(time.DateTime), task.TotalDuration.Truncate(time.Second), task.currentSession(), task.Status})
		t.AppendSeparator()
	}
	t.Render()
}
func startTask(t string) error {
	err := updateActiveTaskTitle(t)
	if err != nil {
		return err
	}

	task, err := getTask(t)
	if err != nil {
		return err
	}

	if task.Status == statusNaN {
		return saveTask(Task{
			Status:  statusActive,
			Created: time.Now(),
			Started: time.Now(),
			Title:   t,
		})
	}

	task.Status = statusActive
	task.Started = time.Now()

	return saveTask(task)
}

func stopTask(t string) error {
	task, err := getTask(t)
	if err != nil {
		return err
	}
	if task.Status == statusNaN {
		return nil
	}

	task.TotalDuration += task.currentSession()
	task.Stopped = time.Now()
	task.Status = statusStopped

	return saveTask(task)
}
