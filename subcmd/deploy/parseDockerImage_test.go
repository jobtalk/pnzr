package deploy

import "testing"

func TestParseDockerImage(t *testing.T) {
	type wantSt struct {
		url string
		tag string
	}
	tests := []struct {
		in         string
		want       wantSt
		expectsErr bool
	}{
		{
			"foo.bar.baz",
			wantSt{
				url: "foo.bar.baz",
			},
			false,
		},
		{
			"foo.bar.baz:hoge",
			wantSt{
				url: "foo.bar.baz",
				tag: "hoge",
			},
			false,
		},
		{
			in:         "foo:bar:baz:hoge",
			expectsErr: true,
		},
	}

	for _, test := range tests {
		url, tag, err := parseDockerImage(test.in)

		if !test.expectsErr && err != nil {
			t.Fatalf("should not be error for %v but: %v", test.in, err)
		}
		if test.expectsErr && err == nil {
			t.Fatalf("should be error for %v but not:", test.in)
		}
		if test.want.tag != tag || test.want.url != url {
			t.Fatalf("want tag: %q, want url: %q, but tag: %q, url: %q", test.want.tag, test.want.url, tag, url)
		}
	}
}
