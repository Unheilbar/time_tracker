package entities

import (
	"fmt"
	"time"
)

// /  active -> ptr entries list active
// /
// /
type ListTitle string

type EntriesLists struct {
	// The title of currently active lists
	CurrentActive    ListTitle
	LastActive       ListTitle
	EntriesListsView map[ListTitle]*List
}

type entryStatus uint8

const (
	StatusStop entryStatus = iota
	StatusActive
)

func (es entryStatus) String() string {
	var reset = "\033[0m"
	var green = "\033[32m"
	var red = "\033[31m"

	switch es {
	case StatusActive:
		return fmt.Sprint(green, "ACTIVE", reset)
	case StatusStop:
		return fmt.Sprint(red, "STOPPED", reset)
	}
	return ""
}

// This list state keep snapshots on the moment new entry appeared
type List struct {
	Title   ListTitle
	Created time.Time
	Tags    []string
	States  []*ListState
}

// t.AppendHeader(table.Row{"Title", "Created", "Stopped", "Started", "Total Duration", "Session Duration", "Status"})
func (l *List) AggregateAllRows() []interface{} {
	last := l.States[len(l.States)-1]
	var stopped, started, currentSession string
	var prev *ListState
	started = last.Timestamp.Format(time.DateTime)
	if len(l.States) < 2 || last.Status == StatusActive {
		currentSession = time.Since(last.Timestamp).Truncate(time.Second).String()
	} else {
		prev = l.States[len(l.States)-2]
		started = prev.Timestamp.Format(time.DateTime)
		stopped = last.Timestamp.Format(time.DateTime)
	}

	return []interface{}{
		l.Title,
		l.Created.Format(time.DateTime),
		stopped,
		started,
		last.TotalDuration.Truncate(time.Second).String(),
		currentSession,
		last.Status,
	}
}

type ListState struct {
	Timestamp     time.Time
	TotalDuration time.Duration

	Status entryStatus
}

func (l *List) safeAppend(status entryStatus) {
	// safe check
	if len(l.States)%2 == 0 && status != StatusActive {
		panic("wrong stop append")
	}

	if len(l.States)%2 == 1 && status != StatusStop {
		panic("wrong start append")
	}

	var ns = &ListState{}
	ns.Status = status
	ns.Timestamp = time.Now()
	ns.TotalDuration = l.last().TotalDuration

	var delta time.Duration
	if status == StatusStop {
		delta += time.Since(l.last().Timestamp)
	}

	ns.TotalDuration += delta

	l.States = append(l.States, ns)
}

func (elist *EntriesLists) InsertEntry(title ListTitle, status entryStatus) {
	// if we stop active task, we should remove current active and add stop entry
	if title == elist.CurrentActive && status == StatusStop {
		elist.stopActive(title, status)
	}

	// if we start active task we should skip
	if title == elist.CurrentActive && status == StatusActive {
		return
	}

	// if we start another task we should  stop the current one
	if title != elist.CurrentActive && status == StatusActive {
		elist.switchActive(title)
	}
}

func (elist *EntriesLists) stopActive(title ListTitle, status entryStatus) {
	currentActive, ok := elist.EntriesListsView[title]
	if !ok {
		return
	}
	currentActive.safeAppend(status)
	elist.LastActive = elist.CurrentActive
	elist.CurrentActive = ""

	return
}

func (elist *EntriesLists) RemoveByTitle(title ListTitle) {
	delete(elist.EntriesListsView, title)
	if elist.CurrentActive == title {
		elist.CurrentActive = ""
		elist.LastActive = title
	}
}

func (elist *EntriesLists) RemoveAll() {
	elist.EntriesListsView = make(map[ListTitle]*List)
	elist.CurrentActive = ""
}

var emptyTitle = ListTitle("")

func (elist *EntriesLists) switchActive(title ListTitle) {
	currentActive, ok := elist.EntriesListsView[elist.CurrentActive]
	if ok {
		currentActive.safeAppend(StatusStop)
	}

	elist.LastActive = elist.CurrentActive
	elist.CurrentActive = title
	l, ok := elist.EntriesListsView[title]
	if !ok {
		l = &List{
			Title:   title,
			Created: time.Now(),
		}
		elist.EntriesListsView[title] = l
	}
	l.safeAppend(StatusActive)
}

func (l *List) last() *ListState {
	if len(l.States) == 0 {
		return &ListState{}
	}

	return l.States[len(l.States)-1]
}

// start taskname # creates new one or starts old one
// stop (stops current running) # adds entry to currently runnning task
//
