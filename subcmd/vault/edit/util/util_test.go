package util

import (
	"os"
	"testing"
)

func TestGetEditorEnvPnzrEditor(t *testing.T) {
	os.Unsetenv("PNZR_EDITOR")
	os.Unsetenv("EDITOR")
	tests := []struct {
		in   string
		want string
	}{
		{"vi", "vi"},
		{"", "nano"},
	}

	for _, test := range tests {
		func() {
			defer func() {
				os.Unsetenv("PNZR_EDITOR")
				os.Unsetenv("EDITOR")
			}()

			os.Setenv("PNZR_EDITOR", test.in)

			got := GetEditor()

			if got != test.want {
				t.Fatalf("want %q, but %q:", test.want, got)
			}
		}()
	}
}

func TestGetEditorEnvEditor(t *testing.T) {
	os.Unsetenv("PNZR_EDITOR")
	os.Unsetenv("EDITOR")
	tests := []struct {
		in   string
		want string
	}{
		{"vi", "vi"},
		{"", "nano"},
	}

	for _, test := range tests {
		func() {
			defer func() {
				os.Unsetenv("PNZR_EDITOR")
				os.Unsetenv("EDITOR")
			}()

			os.Setenv("EDITOR", test.in)

			got := GetEditor()

			if got != test.want {
				t.Fatalf("want %q, but %q:", test.want, got)
			}
		}()
	}
}

func TestGetEditorNoEnv(t *testing.T) {
	os.Unsetenv("PNZR_EDITOR")
	os.Unsetenv("EDITOR")

	got := GetEditor()

	if got != "nano" {
		t.Fatalf("want %q, but %q:", "nano", got)
	}

}
