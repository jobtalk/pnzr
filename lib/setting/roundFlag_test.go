package setting

import (
	"reflect"
	"testing"
)

func TestRoundFlags(t *testing.T) {
	tests := []struct {
		in   []string
		want struct {
			o       []string
			region  string
			profile string
		}
	}{
		{
			[]string{"-foo", "-bar", "-baz"},
			struct {
				o       []string
				region  string
				profile string
			}{
				[]string{"-foo", "-bar", "-baz"},
				"ap-northeast-1",
				"default",
			},
		},
		{
			[]string{"foo", "bar", "baz"},
			struct {
				o       []string
				region  string
				profile string
			}{
				[]string{"foo", "bar", "baz"},
				"ap-northeast-1",
				"default",
			},
		},
		{
			[]string{"-region", "bar", "baz"},
			struct {
				o       []string
				region  string
				profile string
			}{
				[]string{"baz"},
				"bar",
				"default",
			},
		},
		{
			[]string{"-region=bar", "baz"},
			struct {
				o       []string
				region  string
				profile string
			}{
				[]string{"baz"},
				"bar",
				"default",
			},
		},
		{
			[]string{"-profile", "hoge", "baz"},
			struct {
				o       []string
				region  string
				profile string
			}{
				[]string{"baz"},
				"ap-northeast-1",
				"hoge",
			},
		},
		{
			[]string{"-profile=hoge", "baz"},
			struct {
				o       []string
				region  string
				profile string
			}{
				[]string{"baz"},
				"ap-northeast-1",
				"hoge",
			},
		},
		{
			[]string{"-profile", "hoge", "baz", "-region", "ap-northeast-1"},
			struct {
				o       []string
				region  string
				profile string
			}{
				[]string{"baz"},
				"ap-northeast-1",
				"hoge",
			},
		},
		{
			[]string{"-profile", "hoge", "baz", "-region=ap-northeast-1"},
			struct {
				o       []string
				region  string
				profile string
			}{
				[]string{"baz"},
				"ap-northeast-1",
				"hoge",
			},
		},
		{
			[]string{"-profile=hoge", "baz", "-region", "ap-northeast-1"},
			struct {
				o       []string
				region  string
				profile string
			}{
				[]string{"baz"},
				"ap-northeast-1",
				"hoge",
			},
		},
		{
			[]string{"-profile=hoge", "baz", "-region=ap-northeast-1"},
			struct {
				o       []string
				region  string
				profile string
			}{
				[]string{"baz"},
				"ap-northeast-1",
				"hoge",
			},
		},
		{
			[]string{"-profile=", "baz", "-region=ap-northeast-1"},
			struct {
				o       []string
				region  string
				profile string
			}{
				[]string{"baz"},
				"ap-northeast-1",
				"",
			},
		},
		{
			[]string{"-profile=hoge", "baz", "-region="},
			struct {
				o       []string
				region  string
				profile string
			}{
				[]string{"baz"},
				"",
				"hoge",
			},
		},
		{
			[]string{"-profile=", "baz", "-region="},
			struct {
				o       []string
				region  string
				profile string
			}{
				[]string{"baz"},
				"",
				"",
			},
		},
	}

	for i, test := range tests {
		o, r, p := roundFlags(test.in)
		if !reflect.DeepEqual(o, test.want.o) ||
			!reflect.DeepEqual(r, test.want.region) ||
			!reflect.DeepEqual(p, test.want.profile) {
			t.Fatalf("%d: want: %v, but: %v", i, test.want, struct {
				o       []string
				region  string
				profile string
			}{
				o,
				r,
				p,
			})
		}
	}
}