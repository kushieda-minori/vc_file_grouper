package wiki

import (
	"testing"
)

func TestGetTime(t *testing.T) {
	s := "2020-12-16 12:10:02 +0900 JST"
	tm := getTime(s)
	if tm == nil {
		t.Errorf("Time was nil")
		return
	}
	tms := tm.Format("Mon Jan 2 15:04:05 -0700 MST 2006")
	expected := "Wed Dec 16 12:10:02 +0900 JST 2020"
	if tms != expected {
		t.Errorf("Time was expected to be %v but was %v", expected, tms)
	}
}
