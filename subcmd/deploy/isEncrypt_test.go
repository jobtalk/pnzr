package deploy

import "testing"

func TestIsEncrypted(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{
			``,
			false,
		},
		{
			`{}`,
			false,
		},
		{
			`{"cipher"}`,
			false,
		},
		{
			`{"cipher": "hoge"}`,
			true,
		},
		{
			`{"cipher":{}}`,
			false,
		},
	}

	for _, test := range tests {
		got := isEncrypted([]byte(test.input))
		if got != test.want {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}
