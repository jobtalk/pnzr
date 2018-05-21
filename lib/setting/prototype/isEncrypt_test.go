package prototype

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
		{
			`{"hoge":"huga"}`,
			false,
		},
		{
			`{"hoge":"huga", "cipher": "hoge"}`,
			true,
		},
	}

	for _, test := range tests {
		got := NewLoader(nil, nil).isEncrypted([]byte(test.input))
		if got != test.want {
			t.Fatalf("want %v, but %v:", test.want, got)
		}
	}
}
