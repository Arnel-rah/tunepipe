package lyrics

import (
	"testing"
	"time"
)

func TestParseSynced(t *testing.T) {
	lines, err := parseSynced("[00:01.20]First line\n[00:03.50]Second line")
	if err != nil {
		t.Fatal(err)
	}
	if len(lines) != 2 || lines[0].At != 1200*time.Millisecond || lines[1].Text != "Second line" {
		t.Fatalf("unexpected parsed lyrics: %+v", lines)
	}
}
