package tracker

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

// Root prints active task or provides usage info
func Root(cmd *cobra.Command, args []string) {
}

func Create(cmd *cobra.Command, args []string) {
	taskTitle := cmd.Flag("task").Value.String()

	task := Task{
		Title:   taskTitle,
		Started: time.Now(),
	}

	if err := saveTask(task); err != nil {
		log.Fatal(err)
	}
}

func Resume(cmd *cobra.Command, args []string) {
	fmt.Println("start triggered ", args)
}

func AFK(cmd *cobra.Command, args []string) {
	fmt.Println("start triggered ", args)
}

func Stop(cmd *cobra.Command, args []string) {
	taskTitle := cmd.Flag("task").Value.String()
	task, err := getTask(taskTitle)
	if err != nil {
		log.Fatal(err)
	}
	task.Stoped = time.Now()
	task.TotalDuration += time.Since(task.Started)

	if err := saveTask(task); err != nil {
		log.Fatal(err)
	}
}

func List(cmd *cobra.Command, args []string) {
	tasks, err := getAll()
	if err != nil {
		log.Fatal(err)
	}
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"#", "Title", "Started", "Stoped", "Total Duration"})
	t.AppendSeparator()
	for idx, task := range tasks {
		t.AppendRow([]interface{}{idx, task.Title, task.Started, task.Stoped, task.TotalDuration})
		t.AppendSeparator()
	}
	t.Render()
}
