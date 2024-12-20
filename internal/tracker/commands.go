package tracker

import (
	"log"
	"os"
	"strings"

	"github.com/Unheilbar/time_tracker/internal/entities"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
)

var noTaskMsg = "No tasks yet. Add your first one with command start"

type Repository interface {
	LoadList() (*entities.EntriesLists, error)
	DumpList(*entities.EntriesLists) error
}

type App struct {
	repo Repository
}

func NewApp(repo Repository) *App {
	return &App{
		repo: repo,
	}
}

// Root prints active task or provides usage info
func (a *App) Root(cmd *cobra.Command, args []string) {
	list, err := a.repo.LoadList()
	if err != nil {
		log.Fatal("failed to upload list from db", err)
	}

	activeTitle := list.CurrentActive
	lastActive := list.LastActive
	if activeTitle == "" {
		log.Println("No tasks are running. Run a task with start [taskname] command")
		if lastActive != "" {
			// display last active
		}
		return
	}

	// yet to implement
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Title", "Created", "Stopped", "Started", "Total Duration", "Session Duration", "Status"})
	t.AppendSeparator()
	t.AppendRow(list.EntriesListsView[activeTitle].AggregateAllRows())
	t.Render()
}

func (a *App) Start(cmd *cobra.Command, args []string) {
	if len(args) == 0 {
		log.Fatal("Provide task title")
	}

	title := getTitleByArgs(args)

	list, err := a.repo.LoadList()
	if err != nil {
		log.Fatal("failed to upload list from db", err)
	}

	list.InsertEntry(title, entities.StatusActive)

	err = a.repo.DumpList(list)
	if err != nil {
		log.Fatal("failed to save list to db ", err)
	}

}

func (a *App) Stop(cmd *cobra.Command, args []string) {
	list, err := a.repo.LoadList()
	if err != nil {
		log.Fatal("failed to upload list from db")
	}

	title := list.CurrentActive

	list.InsertEntry(title, entities.StatusStop)

	err = a.repo.DumpList(list)
	if err != nil {
		log.Fatal("failed to save list to db")
	}
}

func (a *App) Remove(cmd *cobra.Command, args []string) {
	title := getTitleByArgs(args)

	list, err := a.repo.LoadList()
	if err != nil {
		log.Fatal("failed to upload list from db")
	}

	if title != "" {
		list.RemoveByTitle(title)
	}

	isAll := cmd.Flags().Lookup("all").Changed
	if isAll {
		list.RemoveAll()
	}

	a.repo.DumpList(list)

}

func (a *App) Resume(cmd *cobra.Command, args []string) {
	list, err := a.repo.LoadList()
	if err != nil {
		log.Fatal("failed to upload list from db")
	}

	title := list.CurrentActive
	if title == "" {
		title = list.LastActive
	}

	list.InsertEntry(title, entities.StatusActive)

	err = a.repo.DumpList(list)
	if err != nil {
		log.Fatal("failed to save list to db")
	}
}

func (a *App) List(cmd *cobra.Command, args []string) {
	list, err := a.repo.LoadList()
	if err != nil {
		log.Fatal("failed to upload list from db")
	}

	renderAggregatedAll(list)
}

func getTitleByArgs(args []string) entities.ListTitle {
	return entities.ListTitle(strings.Join(args, " "))
}

func renderAggregatedAll(list *entities.EntriesLists) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Title", "Created", "Stopped", "Started", "Total Duration", "Session Duration", "Status"})
	t.AppendSeparator()
	for title, entries := range list.EntriesListsView {
		if title != list.CurrentActive {
			t.AppendRow(entries.AggregateAllRows())
			t.AppendSeparator()
		}
	}
	if list.CurrentActive != "" {
		t.AppendRow(list.EntriesListsView[list.CurrentActive].AggregateAllRows())
	}

	t.Render()
}
