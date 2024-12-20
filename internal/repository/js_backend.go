package repository

import (
	"encoding/json"
	"os"

	"github.com/Unheilbar/time_tracker/internal/entities"
)

type FileBackend struct {
	path string
}

func NewFileBackend(path string) *FileBackend {
	return &FileBackend{
		path: path,
	}
}

func (js *FileBackend) LoadList() (*entities.EntriesLists, error) {
	fileBytes, _ := os.ReadFile(js.path)

	elist := &entities.EntriesLists{
		EntriesListsView: make(map[entities.ListTitle]*entities.List),
	}

	if len(fileBytes) == 0 {
		return elist, nil
	}

	err := json.Unmarshal(fileBytes, elist)

	if err != nil {
		return nil, err
	}

	return elist, nil
}

func (fb *FileBackend) DumpList(list *entities.EntriesLists) error {
	// Marshal the data into JSON format
	jsonData, err := json.MarshalIndent(list, "", "  ")
	if err != nil {
		return err
	}

	// Write the JSON data to a file
	err = os.WriteFile(fb.path, jsonData, 0644)
	if err != nil {
		return err
	}

	return nil
}
