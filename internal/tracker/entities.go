package tracker

import (
	"fmt"
	"time"
)

type Task struct {
	ID            []byte
	Title         string
	Started       time.Time
	Stoped        time.Time
	TotalDuration time.Duration
}

func getKey(title string) []byte {
	return []byte(fmt.Sprintf("%s_%s", tasksPrefix, title))
}
