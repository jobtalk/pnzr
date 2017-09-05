package config

import (
	"bytes"
	"encoding/json"
	"testing"
)

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

func TestEmbedde(t *testing.T) {
	dst := new(bytes.Buffer)

	baseConfStr := `
	{
		"foo": "bar",
		"any_val": $val
	}
	`
	embeddedVal := `{"val": {"hoge":"huga"}}`

	resultJSON := `{"foo":"bar","any_val":{"hoge":"huga"}}`

	result, err := Embedde(baseConfStr, embeddedVal)
	if err != nil {
		t.Fatal("error: %v", err.Error())
	}
	if !isJSON(result) {
		t.Fatal("%v is not json", result)
	}

	src := []byte(result)
	if err := json.Compact(dst, src); err != nil {
		t.Fatal("error: %v", err.Error())
	}

	if resultJSON != dst.String() {
		errString := ""
		errString += "result: \n"
		errString += dst.String()
		errString += "\n\n"

		errString += "answer: \n"
		errString += resultJSON
		errString += "\n\n"

		t.Fatal(errString)
	}
}
