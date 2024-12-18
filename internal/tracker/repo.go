package tracker

import (
	"encoding/json"
	"log"
	"os"
)

var dataDir = ".time_tracker/badger"
var logDir = ".time_tracker/logs"

const (
	envDataDir = "GO_TIME_TRACKER_DATA_PATH"
	envLogDir  = "GO_TIME_TRACKER_LOG_PATH"
)

var db DB

func init() {
	if val, ok := os.LookupEnv(envDataDir); ok {
		dataDir = val
	}

	if val, ok := os.LookupEnv(envLogDir); ok {
		logDir = val
	}

	db = getDB()
}

// DB defines an embedded key/value store database interface.
type DB interface {
	Get(namespace, key []byte) (value []byte, err error)
	Set(namespace, key, value []byte) error
	Has(namespace, key []byte) (bool, error)
	All(namespace, prefix []byte) (vals [][]byte, err error)
	Close() error
}

func getDB() DB {
	db, err := NewBadgerDB(dataDir)
	if err != nil {
		log.Fatal(err)
	}

	return db
}

var (
	ns          = []byte("tasks")
	tasksPrefix = []byte("task_pfx")
)

func saveTask(t Task) error {
	enc, _ := json.Marshal(t)

	return db.Set(ns, getKey(t.Title), enc)
}

func getTask(title string) (Task, error) {
	var t Task
	enc, err := db.Get(ns, getKey(title))
	if err != nil {
		return Task{}, err
	}

	json.Unmarshal(enc, &t)
	return t, nil
}

func getAll() ([]Task, error) {
	encs, err := db.All(ns, tasksPrefix)
	if err != nil {
		return nil, err
	}

	var res []Task
	for _, enc := range encs {
		var t Task
		json.Unmarshal(enc, &t)
		res = append(res, t)
	}

	return res, nil
}
