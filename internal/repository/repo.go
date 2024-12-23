package repository

import (
	"encoding/json"
	"log"

	"github.com/Unheilbar/time_tracker/internal/entities"
	"github.com/dgraph-io/badger"
)

// DB defines an embedded key/value store database interface.
type DB interface {
	Get(namespace, key []byte) (value []byte, err error)
	Remove(namespace, key []byte) error
	Set(namespace, key, value []byte) error
	Has(namespace, key []byte) (bool, error)
	All(namespace, prefix []byte) (vals [][]byte, err error)
	Close() error
}

type Repository struct {
	db DB
}

var Repo *Repository

func NewRepo(db DB) *Repository {
	return &Repository{db}
}

var listPrefix = []byte("my_list")
var ns = []byte("ns")

func (repo *Repository) LoadList() (*entities.EntriesLists, error) {
	enc, err := repo.db.Get(ns, listPrefix)
	if err != nil && err != badger.ErrKeyNotFound {
		return nil, err
	}

	res := &entities.EntriesLists{
		EntriesListsView: make(map[entities.ListTitle]*entities.List),
		Tags:             entities.TagsView{View: make(map[entities.Tag][]entities.ListTitle)},
	}

	if err == badger.ErrKeyNotFound {
		log.Print("Create task list at db")
		return res, nil
	}

	err = json.Unmarshal(enc, &res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (repo *Repository) DumpList(l *entities.EntriesLists) error {
	enc, err := json.Marshal(l)
	if err != nil {
		return err
	}

	return repo.db.Set(ns, listPrefix, enc)
}
