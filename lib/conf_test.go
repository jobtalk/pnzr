package lib

import (
	"bytes"
	"encoding/json"
	"testing"
)

func isJSON(s string) bool {
	var v interface{}
	return nil == json.Unmarshal([]byte(s), &v)
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
