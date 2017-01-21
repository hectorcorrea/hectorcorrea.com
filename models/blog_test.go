package models

import "testing"

func TestSlug(t *testing.T) {
	testA := []string{"", ""}
	testB := []string{"abc 345 DEF", "abc-345-def"}
	testC := []string{"hello c#", "hello-c-sharp"}
	testD := []string{"a<b", "a-b"}
	testE := []string{"a <  b", "a-b"}
	testF := []string{"a b<", "a-b"}
	testG := []string{"a b<<", "a-b"}
	testH := []string{"<", ""}
	tests := [][]string{testA, testB, testC, testD, testE, testF,
		testG, testH}
	for _, test := range tests {
		value := test[0]
		slug := getSlug(value)
		expected := test[1]
		if slug != expected {
			t.Errorf("Unexpected slug (%s) for (%s)", slug, value)
		}
	}
}

// func TestPostedOnRSS(t *testing.T) {
// 	dbDate := "2015-09-17 02:06:31 +0000 UTC"
// 	// layout := "2006-01-02 15:04:06 -0700 MST"
// 	layout := "2006-01-02 15:04:05 -0700 MST"
// 	goDate, _ := time.Parse(layout, dbDate)
// 	t.Error(goDate)
// }
