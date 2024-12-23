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

	firstTag  = Tag("#tag1")
	secondTag = Tag("#tag2")
	thirdTag  = Tag("#tag3")
)

type tester struct {
	elist *EntriesLists
}

func (tester *tester) reset() {
	tester.elist = InitEmptyElist()
}

func (tester *tester) start1() {
	tester.elist.InsertEntry(firstTitle, StatusActive)
}

func (tester *tester) start2() {
	tester.elist.InsertEntry(secondTitle, StatusActive)
}

func (tester *tester) start3() {
	tester.elist.InsertEntry(thirdTitle, StatusActive)
}

func (tester *tester) stop1() {
	tester.elist.InsertEntry(firstTitle, StatusStop)
}

func (tester *tester) stop2() {
	tester.elist.InsertEntry(secondTitle, StatusStop)
}

func (tester *tester) stop3() {
	tester.elist.InsertEntry(thirdTitle, StatusStop)
}

func (tester *tester) tag1(tag Tag) {
	tester.elist.AddTag(tag, firstTitle)
}

func (tester *tester) tag2(tag Tag) {
	tester.elist.AddTag(tag, secondTitle)
}

func (tester *tester) tag3(tag Tag) {
	tester.elist.AddTag(tag, thirdTitle)
}

func (tester *tester) checkStatus(t *testing.T, title ListTitle, expstatus entryStatus) {
	state := getListLastState(tester.elist, title)
	assert.Equal(t, expstatus, state.Status)
}

func Test_EntriesList(t *testing.T) {
	tester := &tester{}

	t.Run(`add 1 active 2 active check 1 stop
	then add 3 check 3 is current 2 is last`,
		func(t *testing.T) {
			tester.reset()

			tester.start1()
			tester.start2()

			tester.checkStatus(t, secondTitle, StatusActive)
			tester.checkStatus(t, firstTitle, StatusStop)

			tester.start3()
			assert.Equal(t, thirdTitle, tester.elist.CurrentActive)
			assert.Equal(t, secondTitle, tester.elist.LastActive)

		})

	t.Run(`add 1 active wait 1 second add 2 active check 1 duration 1<d<3`,
		func(t *testing.T) {
			tester.reset()

			tester.start1()

			time.Sleep(time.Millisecond)

			tester.start2()

			firstState := getListLastState(tester.elist, firstTitle)

			checkLess0 := time.Millisecond < firstState.TotalDuration
			checkLess1 := firstState.TotalDuration < time.Millisecond*2

			assert.Equal(t, true, checkLess0)
			assert.Equal(t, true, checkLess1)
		})

	t.Run("start 1 -> sleep -> start 2 -> sleep -> start 1 -> -> sleep -> stop 1 -> check duration",
		func(t *testing.T) {
			tester.reset()

			tester.start1()
			time.Sleep(time.Millisecond)

			tester.start2()
			time.Sleep(time.Millisecond)

			tester.start1()
			time.Sleep(time.Millisecond)

			tester.stop1()

			firstState := getListLastState(tester.elist, firstTitle)

			checkLess0 := time.Millisecond*2 < firstState.TotalDuration
			checkLess1 := firstState.TotalDuration < time.Millisecond*3

			assert.Equal(t, true, checkLess0)
			assert.Equal(t, true, checkLess1)
		})

	t.Run(` -> start task1 
		-> add tag1 
		-> start task2
		-> add tag1
		-> add tag2
		->start3 
		->get by any tag1,tag2 
		-> check task1,task2 ->
		-> get by all tag1, tag2
		-> check only task1
		-> remove tag1 
		-> get by all tag1,tag2
		-> check empty
		check empty`,
		func(t *testing.T) {
			tester.reset()
			check := map[ListTitle]bool{
				firstTitle:  true,
				secondTitle: true,
			}

			tester.start1()
			tester.start2()

			tester.tag1(firstTag)
			tester.tag2(firstTag)
			tester.tag1(secondTag)

			res := tester.elist.Filter([]Tag{firstTag}, ContainsAny)

			for _, l := range res {
				_, ok := check[l.Title]
				assert.Equal(t, true, ok)
			}

			res = tester.elist.Filter([]Tag{secondTag}, ContainsAll)

			assert.Equal(t, 1, len(res))
			assert.Equal(t, firstTitle, res[0].Title)

			tester.elist.RemoveTag(firstTag)

			res = tester.elist.Filter([]Tag{firstTag}, ContainsAny)

			assert.Equal(t, 0, len(res))

			res = tester.elist.Filter([]Tag{firstTag, secondTag}, ContainsAny)

			assert.Equal(t, 1, len(res))

			res = tester.elist.Filter([]Tag{firstTag, secondTag}, ContainsAll)
			assert.Equal(t, 0, len(res))

		})
	t.Run("start t1, tag1 t1, rm t1, check tag removed",
		func(t *testing.T) {
			tester.reset()

			tester.start1()
			tester.tag1(firstTag)

			tester.elist.RemoveByTitle(firstTitle)

			_, ok := tester.elist.Tags.View[firstTag]
			assert.Equal(t, false, ok)

			res := tester.elist.Filter([]Tag{firstTag}, ContainsAny)
			assert.Equal(t, 0, len(res))

		})
}

func getListLastState(l *EntriesLists, t ListTitle) *ListState {
	length := len(l.EntriesListsView[t].States)
	last := l.EntriesListsView[t].States[length-1]
	return last
}
