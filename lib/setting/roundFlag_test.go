package setting

import (
	"os"
	"reflect"
	"testing"
)

func setenv(m map[string]string) {
	for k, v := range m {
		err := os.Setenv(k, v)
		if err != nil {
			panic(err)
		}
	}
}

func unsetenv(m map[string]string) {
	for k, _ := range m {
		err := os.Unsetenv(k)
		if err != nil {
			panic(err)
		}
	}
}

func TestRoundFlags(t *testing.T) {
	type want struct {
		o       []string
		region  string
		profile string
	}
	tests := []struct {
		in  []string
		env map[string]string
		w   want
	}{
		{
			[]string{"-foo", "-bar", "-baz"},
			nil,
			want{
				[]string{"-foo", "-bar", "-baz"},
				"ap-northeast-1",
				"default",
			},
		},
		{
			[]string{"foo", "bar", "baz"},
			nil,
			want{
				[]string{"foo", "bar", "baz"},
				"ap-northeast-1",
				"default",
			},
		},
		{
			[]string{"-region", "bar", "baz"},
			nil,
			want{
				[]string{"baz"},
				"bar",
				"default",
			},
		},
		{
			[]string{"-region=bar", "baz"},
			nil,
			want{
				[]string{"baz"},
				"bar",
				"default",
			},
		},
		{
			[]string{"-profile", "hoge", "baz"},
			nil,
			want{
				[]string{"baz"},
				"ap-northeast-1",
				"hoge",
			},
		},
		{
			[]string{"-profile=hoge", "baz"},
			nil,
			want{
				[]string{"baz"},
				"ap-northeast-1",
				"hoge",
			},
		},
		{
			[]string{"-profile", "hoge", "baz", "-region", "ap-northeast-1"},
			nil,
			want{
				[]string{"baz"},
				"ap-northeast-1",
				"hoge",
			},
		},
		{
			[]string{"-profile", "hoge", "baz", "-region=ap-northeast-1"},
			nil,
			want{
				[]string{"baz"},
				"ap-northeast-1",
				"hoge",
			},
		},
		{
			[]string{"-profile=hoge", "baz", "-region", "ap-northeast-1"},
			nil,
			want{
				[]string{"baz"},
				"ap-northeast-1",
				"hoge",
			},
		},
		{
			[]string{"-profile=hoge", "baz", "-region=ap-northeast-1"},
			nil,
			want{
				[]string{"baz"},
				"ap-northeast-1",
				"hoge",
			},
		},
		{
			[]string{"-profile=", "baz", "-region=ap-northeast-1"},
			nil,
			want{
				[]string{"baz"},
				"ap-northeast-1",
				"",
			},
		},
		{
			[]string{"-profile=hoge", "baz", "-region="},
			nil,
			want{
				[]string{"baz"},
				"",
				"hoge",
			},
		},
		{
			[]string{"-profile=", "baz", "-region="},
			nil,
			want{
				[]string{"baz"},
				"",
				"",
			},
		},
		{
			[]string{},
			nil,
			want{
				[]string{},
				"ap-northeast-1",
				"default",
			},
		},
	}

	for i, test := range tests {
		setenv(test.env)
		o, r, p := roundFlags(test.in)
		if (!reflect.DeepEqual(o, test.w.o) && (len(o) != 0 || len(test.w.o) != 0)) ||
			!reflect.DeepEqual(r, test.w.region) ||
			!reflect.DeepEqual(p, test.w.profile) {
			t.Fatalf("%d: want: %v, but: %v", i, test.w, want{
				o,
				r,
				p,
			})
		}
		unsetenv(test.env)
	}
}
