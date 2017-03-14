package conf

import "testing"

func TestIsJSON(t *testing.T) {
	jsonStr := `{"hoge":"huga", "flag": true, "num": 0}`
	nonJSON := "This is a pen."

	if !isJSON(jsonStr) {
		t.Errorf("The expected value is %v, but actually it is %v.", true, !isJSON(jsonStr))
	}

	if isJSON(nonJSON) {
		t.Errorf("The expected value is %v, but actually it is %v.", false, isJSON(nonJSON))
	}
}
