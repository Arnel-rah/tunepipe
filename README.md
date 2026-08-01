# tunepipe

A minimal, terminal-based music player powered by [`yt-dlp`](https://github.com/yt-dlp/yt-dlp). Search, browse, and play audio streams directly from your terminal — no browser, no bloat.

![style](https://img.shields.io/badge/style-monochrome-000000)
![go version](https://img.shields.io/badge/go-1.22%2B-00ADD8)

## Features

- **TUI interface** built with [Bubbletea](https://github.com/charmbracelet/bubbletea) — clean, minimal, monochrome design
- **Streaming playback** via `yt-dlp`, no local downloads required by default
- **URL cache with prefetch-on-navigate** — reduces playback latency by resolving the next track's stream URL while you browse
- Low-latency playback, cross-platform (Linux, macOS, Windows)

## Requirements

`tunepipe` relies on the following external tools being available in your `PATH`:

| Tool | Purpose | Install |
|------|---------|---------|
| [`yt-dlp`](https://github.com/yt-dlp/yt-dlp) | Resolving and extracting audio streams | `pip install yt-dlp` or see [releases](https://github.com/yt-dlp/yt-dlp/releases) |
| `ffmpeg` | Audio decoding/playback | [ffmpeg.org/download](https://ffmpeg.org/download.html) |

`tunepipe` checks for these on startup and will warn you if they're missing.

## Installation

### Homebrew (macOS/Linux)

```bash
brew install nel/tunepipe/tunepipe
```

### Go install

```bash
go install github.com/Arnel-rah/tunepipe@latest
```

### Prebuilt binaries

Download the latest binary for your platform from the [Releases page](https://github.com/Arnel-rah/tunepipe/releases).

```bash
# Example: Linux amd64
curl -L https://github.com/Arnel-rah/tunepipe/releases/latest/download/tunepipe_linux_amd64.tar.gz | tar xz
sudo mv tunepipe /usr/local/bin/
```

### From source

```bash
git clone https://github.com/Arnel-rah/tunepipe.git
cd tunepipe
go build -o tunepipe
```

## Usage

```bash
tunepipe
```

Launches the TUI. From there:

| Key | Action |
|-----|--------|
| `/` | Search |
| `↑` / `↓` | Navigate list |
| `Enter` | Play selected track |
| `Space` | Pause / resume |
| `n` | Next track |
| `p` | Previous track |
| `q` | Quit |

> Adjust the key bindings above to match your actual implementation.

## Configuration

`tunepipe` reads optional configuration from `~/.config/tunepipe/config.yaml`:

```yaml
cache_dir: ~/.cache/tunepipe
prefetch: true
theme: monochrome
```

## How it works

1. You search or queue a track.
2. `tunepipe` shells out to `yt-dlp` to resolve a direct stream URL.
3. The resolved URL is cached; while you're browsing or listening to the current track, `tunepipe` prefetches the next one in the background.
4. Playback is handed off for decoding/output, minimizing the gap between selecting a track and hearing audio.

## Development

```bash
git clone https://github.com/Arnel-rah/tunepipe.git
cd tunepipe
go mod tidy
go run main.go
```

### Releasing

Releases are automated with [GoReleaser](https://goreleaser.com) via GitHub Actions. Push a semver tag to trigger a release:

```bash
git tag -a v0.1.0 -m "first release"
git push origin v0.1.0
```

## License

MIT — see [LICENSE](LICENSE) for details.