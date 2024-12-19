package tracker

import (
	"fmt"
	"time"
)

type taskStatus uint8

func (ts taskStatus) String() string {
	var reset = "\033[0m"
	var green = "\033[32m"
	var red = "\033[31m"

	switch ts {
	case statusActive:
		return fmt.Sprint(green, "ACTIVE", reset)
	case statusStopped:
		return fmt.Sprint(red, "STOPPED", reset)
	case statusNaN:
		return "NaN"
	}
	return ""
}

const (
	statusNaN taskStatus = iota
	statusActive
	statusStopped
)

var statusView = map[taskStatus]string{
	statusNaN:     "NaN",
	statusActive:  "ACTIVE",
	statusStopped: "STOPPED",
}

type Task struct {
	ID            []byte
	Title         string
	Status        taskStatus
	Created       time.Time
	Stopped       time.Time
	Started       time.Time
	TotalDuration time.Duration
	Tags          []string
}

func (t Task) currentSession() time.Duration {
	if t.Status != statusActive {
		return 0
	}

	return time.Since(t.Started).Truncate(time.Second)
}

func (t Task) stopped() string {
	if t.Stopped.Unix() < 0 {
		return ""
	}

	return t.Stopped.Format(time.DateTime)
}

func getKey(title string) []byte {
	return []byte(fmt.Sprintf("%s_%s", tasksPrefix, title))
}
