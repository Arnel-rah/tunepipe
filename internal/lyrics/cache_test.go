package lyrics

import (
	"context"
	"testing"
)

func TestClientCachesLyrics(t *testing.T) {
	client := NewClient(nil)
	lines := []Line{{Text: "hello"}}
	client.store("artist\x00title", lines)
	got, err := client.Get(context.Background(), "artist", "title")
	if err != nil || len(got) != 1 || got[0].Text != "hello" {
		t.Fatalf("expected cached lyrics, got lines=%v err=%v", got, err)
	}
}
