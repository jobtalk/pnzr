package edit

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

			got := getEditor()

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

			got := getEditor()

			if got != test.want {
				t.Fatalf("want %q, but %q:", test.want, got)
			}
		}()
	}
}

func TestGetEditorNoEnv(t *testing.T) {
	os.Unsetenv("PNZR_EDITOR")
	os.Unsetenv("EDITOR")

	got := getEditor()

	if got != "nano" {
		t.Fatalf("want %q, but %q:", "nano", got)
	}

}
