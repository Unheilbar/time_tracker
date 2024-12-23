package tracker

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/Unheilbar/time_tracker/internal/entities"
	"github.com/Unheilbar/time_tracker/internal/flags"
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

	tags := getTags(cmd)
	for _, tag := range tags {
		list.AddTag(tag, title)
	}

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

	tags := getTags(cmd)

	renderAggregatedAll(list, tags)
}

func getTags(cmd *cobra.Command) []entities.Tag {
	tagsStr := cmd.Flags().Lookup(flags.Tag.Name).Value.String()
	if len(tagsStr) == 0 {
		return nil
	}
	tags := strings.Split(tagsStr, "#")
	if len(tags) == 0 {
		log.Fatal("No tags found. Make sure your tags start with #")
	}
	fmt.Println(tags)

	var res []entities.Tag
	for _, tag := range tags {
		tag = strings.Trim(tag, " ")
		tag = strings.ToLower(tag)

		if strings.Contains(tag, " ") {
			log.Fatalf("Wrong tag format %s. Make sure your tags start with # lie in #youidiot.", tag)
		}

		res = append(res, entities.Tag(fmt.Sprint("#", tag)))
	}

	return res
}

func getTitleByArgs(args []string) entities.ListTitle {
	return entities.ListTitle(strings.Join(args, " "))
}

func renderAggregatedAll(list *entities.EntriesLists, tags []entities.Tag) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"Title", "Created", "Started", "Stopped", "Total Duration", "Session Duration", "Status"})
	t.AppendSeparator()
	var hasActive bool
	for _, entries := range list.Filter(tags, entities.ContainsAll) {
		if entries.Title != list.CurrentActive {
			t.AppendRow(entries.AggregateAllRows())
			t.AppendSeparator()
		} else {
			hasActive = true
		}
	}

	if list.CurrentActive != "" && hasActive {
		t.AppendRow(list.EntriesListsView[list.CurrentActive].AggregateAllRows())
	}

	t.Render()
}
