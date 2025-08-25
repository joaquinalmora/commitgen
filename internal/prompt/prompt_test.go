package prompt

import (
	"strings"
	"testing"
)

func TestMakePromptSingleLine(t *testing.T) {
	cases := []struct {
		files  []string
		patch  string
		expect string
	}{
		{[]string{"README.md"}, "", "Update documentation"},
		{[]string{"main.go"}, "--- a/main.go\n+++ b/main.go\n@@ -1 +1 @@\n+hello", "Update main.go"},
		{[]string{"test/foo_test.go"}, "", "Update tests"},
	}

	for _, c := range cases {
		msg := MakePrompt(c.files, c.patch)
		if strings.Contains(msg, "\n") {
			t.Fatalf("expected single-line prompt, got newline in: %q", msg)
		}
		if msg == "" {
			t.Fatalf("expected non-empty prompt for files %v", c.files)
		}
	}
}
