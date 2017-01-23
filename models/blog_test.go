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

// func TestLegacySlug(t *testing.T) {
// 	legacySlug := "Something-ABC.aspx"
// 	slug := strings.ToLower(legacySlug)
//
// 	if strings.HasSuffix(slug, ".aspx") {
// 		slug = slug[0 : len(slug)-5]
// 		t.Errorf("(1) %s", slug)
// 	} else {
// 		t.Errorf("(2) %s", slug)
// 	}
// }
