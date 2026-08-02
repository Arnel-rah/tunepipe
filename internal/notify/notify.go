package notify

import (
	"github.com/gen2brain/beeep"
	"github.com/Arnel-rah/tunepipe/internal/ytdlp"
)

// NotifyTrack shows a simple OS notification with track title and uploader.
// Best-effort: errors are returned but callers may ignore them.
func NotifyTrack(t ytdlp.Track) error {
	title := t.Title
	msg := t.Uploader
	return beeep.Notify(title, msg, "")
}
