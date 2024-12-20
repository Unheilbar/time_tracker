package entities

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	firstTitle  = ListTitle("first")
	secondTitle = ListTitle("second")
	thirdTitle  = ListTitle("third")
)

var elist *EntriesLists

func Test_EntriesList(t *testing.T) {
	t.Run(`add 1 active 2 active check 1 stop
	then add 3 check 3 is current 2 is last`,
		func(t *testing.T) {
			elist = initEmptyElist()

			elist.InsertEntry(firstTitle, StatusActive)

			assert.Equal(t, firstTitle, elist.CurrentActive)

			firstState := getListLastState(elist, firstTitle)
			assert.Equal(t, StatusActive, firstState.Status)

			elist.InsertEntry(secondTitle, StatusActive)

			secondState := getListLastState(elist, secondTitle)
			assert.Equal(t, StatusActive, secondState.Status)

			firstState = getListLastState(elist, firstTitle)
			assert.Equal(t, StatusStop, firstState.Status)

			elist.InsertEntry(thirdTitle, StatusActive)
			assert.Equal(t, thirdTitle, elist.CurrentActive)
			assert.Equal(t, secondTitle, elist.LastActive)

		})

	t.Run(`add 1 active wait 1 second add 2 active check 1 duration 1<d<3`, func(t *testing.T) {
		elist = initEmptyElist()

		elist.InsertEntry(firstTitle, StatusActive)

		time.Sleep(time.Second)
		elist.InsertEntry(secondTitle, StatusActive)

		firstState := getListLastState(elist, firstTitle)

		checkLess0 := time.Second < firstState.TotalDuration
		checkLess1 := firstState.TotalDuration < time.Second*2

		assert.Equal(t, true, checkLess0)
		assert.Equal(t, true, checkLess1)
	})

	t.Run("start 1 -> sleep -> start 2 -> sleep -> start 1 -> -> sleep -> stop 1 -> check duration", func(t *testing.T) {

		elist = initEmptyElist()

		elist.InsertEntry(firstTitle, StatusActive)

		time.Sleep(time.Second)
		elist.InsertEntry(secondTitle, StatusActive)

		time.Sleep(time.Second)
		elist.InsertEntry(firstTitle, StatusActive)
		time.Sleep(time.Second)

		elist.InsertEntry(firstTitle, StatusStop)

		firstState := getListLastState(elist, firstTitle)

		checkLess0 := time.Second*2 < firstState.TotalDuration
		checkLess1 := firstState.TotalDuration < time.Second*3

		assert.Equal(t, true, checkLess0)
		assert.Equal(t, true, checkLess1)
	})
}

func initEmptyElist() *EntriesLists {
	return &EntriesLists{
		EntriesListsView: make(map[ListTitle]*List),
	}
}

func getListLastState(l *EntriesLists, t ListTitle) *ListState {
	length := len(l.EntriesListsView[t].States)
	last := l.EntriesListsView[t].States[length-1]
	return last
}
