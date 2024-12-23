package entities

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type EntriesLists struct {
	// The title of currently active lists
	CurrentActive    ListTitle
	LastActive       ListTitle
	EntriesListsView map[ListTitle]*List
	Tags             TagsView
}

func (elist *EntriesLists) AddTag(tag Tag, title ListTitle) error {
	list, ok := elist.EntriesListsView[title]
	if !ok {
		return errors.New("title doesn't exist")
	}

	list.Tags = append(list.Tags, tag)
	elist.Tags.View[tag] = append(elist.Tags.View[tag], title)

	return nil

}

func (elist *EntriesLists) RemoveTag(tag Tag) {
	for _, title := range elist.Tags.View[tag] {
		elist.EntriesListsView[title].RemoveTag(tag)
	}

	delete(elist.Tags.View, tag)
}

func (elist *EntriesLists) Lists() []*List {
	var lists []*List

	for _, list := range elist.EntriesListsView {
		lists = append(lists, list)
	}

	return lists
}

type Filter func(*List, []Tag) bool

func (elist *EntriesLists) Filter(tags []Tag, ok Filter) []*List {
	if len(tags) == 0 {
		return elist.Lists()
	}
	var lists []*List

	for _, tag := range tags {
		for _, title := range elist.Tags.View[tag] {
			list := elist.EntriesListsView[title]
			if ok(list, tags) {
				lists = append(lists, list)
			}
		}
	}

	return lists
}

func ContainsAll(l *List, tags []Tag) bool {
	if len(l.Tags) == 0 {
		return false
	}

	filter := toTagFilter(l.Tags)

	for _, tag := range tags {
		if _, ok := filter[tag]; !ok {
			return false
		}
	}

	return true
}

func ContainsAny(l *List, tags []Tag) bool {
	// repeating code
	if len(l.Tags) == 0 {
		return false
	}

	filter := toTagFilter(tags)

	for _, tag := range l.Tags {
		if _, ok := filter[tag]; ok {
			return true
		}
	}

	return false
}

func toTagFilter(tags []Tag) map[Tag]bool {
	filter := make(map[Tag]bool)
	for _, tag := range tags {
		filter[tag] = true
	}

	return filter
}

type Tag string

type TagsView struct {
	View map[Tag][]ListTitle
}

// /  active -> ptr entries list active
// /
// /
type ListTitle string

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
	Id      uint64
	Title   ListTitle
	Created time.Time
	Tags    []Tag
	States  []*ListState
}

// t.AppendHeader(table.Row{"#","Title", "Created",  "Started","Stopped", "Total Duration", "Session Duration", "Status"})
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
		//l.Id,
		titleAggregate(l.Title),
		l.Created.Format(time.DateTime),
		started,
		stopped,
		last.TotalDuration.Truncate(time.Second).String(),
		currentSession,
		last.Status,
	}
}

func titleAggregate(title ListTitle) string {
	words := strings.Split(string(title), " ")
	var res string
	trigger := 10
	current := 0
	for _, word := range words {
		if current+len(word) > trigger {
			current = 0
			res += "\n"
		}
		res += word
		res += " "
		current += len(word)
	}
	return res
}

func (l *List) RemoveTag(t Tag) {
	var tags []Tag

	for _, tag := range l.Tags {
		if tag != t {
			tags = append(tags, tag)
		}
	}

	l.Tags = tags
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
	for _, tag := range elist.EntriesListsView[title].Tags {
		delete(elist.Tags.View, tag)
	}

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

func InitEmptyElist() *EntriesLists {
	return &EntriesLists{
		EntriesListsView: make(map[ListTitle]*List),
		Tags: TagsView{
			View: make(map[Tag][]ListTitle),
		},
	}
}
