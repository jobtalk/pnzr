package lib

import "testing"

func TestIsJson(t *testing.T) {
	notJSONstr := "this is not json"
	jsonStr := `
	{
		"foo":"bar",
		"bool": true
	}
	`

	if isJSON(notJSONstr) != false {
		t.Fatalf("\"%v\"であるべきだが\"%v\"だった", false, !false)
	}

	if isJSON(jsonStr) != true {
		t.Fatalf("\"%v\"であるべきだが\"%v\"だった", true, !true)
	}
}
