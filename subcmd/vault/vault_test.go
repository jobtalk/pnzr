package vault

import "testing"

func b(b bool) *bool {
	return &b
}

func TestMode_CheckMultiFlagSet(t *testing.T) {
	tests := []struct {
		in   mode
		want bool
	}{
		{
			mode{
				b(true),
				b(true),
				b(true),
				b(true)},
			true,
		},
		{
			mode{
				b(true),
				b(true),
				b(true),
				b(false)},
			true,
		},
		{
			mode{
				b(true),
				b(true),
				b(false),
				b(true)},
			true,
		},
		{
			mode{
				b(true),
				b(true),
				b(false),
				b(false)},
			true,
		},
		{
			mode{
				b(true),
				b(false),
				b(true),
				b(true)},
			true,
		},
		{
			mode{
				b(true),
				b(false),
				b(true),
				b(false)},
			true,
		},
		{
			mode{
				b(true),
				b(false),
				b(false),
				b(true)},
			true,
		},
		{
			mode{
				b(true),
				b(false),
				b(false),
				b(false)},
			false,
		},
		{
			mode{
				b(false),
				b(true),
				b(true),
				b(true)},
			true,
		},
		{
			mode{
				b(false),
				b(true),
				b(true),
				b(false)},
			true,
		},
		{
			mode{
				b(false),
				b(true),
				b(false),
				b(true)},
			true,
		},
		{
			mode{
				b(false),
				b(true),
				b(false),
				b(false)},
			false,
		},
		{
			mode{
				b(false),
				b(false),
				b(true),
				b(true)},
			true,
		},
		{
			mode{
				b(false),
				b(false),
				b(true),
				b(false)},
			false,
		},
		{
			mode{
				b(false),
				b(false),
				b(false),
				b(true)},
			false,
		},
		{
			mode{
				b(false),
				b(false),
				b(false),
				b(false)},
			false,
		},
	}

	for _, test := range tests {
		got := test.in.checkMultiFlagSet()

		if test.want != got {
			t.Fatalf("want %q, but %q:", test.want, got)
		}
	}
}
