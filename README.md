<div align="center">

# tunepipe

**A minimal, terminal-based music player powered by [`yt-dlp`](https://github.com/yt-dlp/yt-dlp).**
Search, browse, and play audio streams directly from your terminal — no browser, no bloat.

![style](https://img.shields.io/badge/style-monochrome-000000)
![go version](https://img.shields.io/badge/go-1.26%2B-00ADD8)
![license](https://img.shields.io/badge/license-MIT-green)
![platform](https://img.shields.io/badge/platform-linux%20%7C%20macos%20%7C%20windows-lightgrey)
[![release](https://img.shields.io/github/v/release/Arnel-rah/tunepipe)](https://github.com/Arnel-rah/tunepipe/releases/latest)

</div>

<!-- Replace with an asciinema recording or GIF — this sells the tool in 5 seconds -->
<!-- ![demo](docs/demo.gif) -->

---

## Table of contents

- [Features](#features)
- [Requirements](#requirements)
- [Installation](#installation)
- [Usage](#usage)
- [Configuration](#configuration)
- [How it works](#how-it-works)
- [Troubleshooting](#troubleshooting)
- [Development](#development)
- [Contributing](#contributing)
- [License](#license)

## Features

- 🖥️ **TUI interface** built with [Bubbletea](https://github.com/charmbracelet/bubbletea) — clean, minimal, monochrome design
- 🎧 **Streaming playback** via `yt-dlp`, no local downloads required by default
- ⚡ **URL cache with prefetch-on-navigate** — resolves the next track's stream URL while you browse, minimizing playback latency
- 📉 **Data-conscious by default** — prioritizes low-bitrate audio formats (64–96 kbps) to reduce bandwidth usage, ideal on limited or metered connections
- 🌍 **Cross-platform** — Linux, macOS, Windows

## Requirements

`tunepipe` relies on the following external tools:

| Tool | Purpose | Install |
|------|---------|---------|
| `mpv` | Audio playback engine | `winget install mpv-player.mpv` (Windows), `brew install mpv` (macOS), your package manager (Linux) |
| [`yt-dlp`](https://github.com/yt-dlp/yt-dlp) | Resolving and extracting audio streams | Optional in `PATH`; if missing, `tunepipe` downloads a standalone binary automatically |

`tunepipe` checks these at startup and exits with an install hint if `mpv` is missing.

> **Note:** `yt-dlp` needs a JS challenge solver to reliably resolve YouTube streams. `tunepipe` passes `--remote-components ejs:github` automatically — no extra setup needed on your end.

## Installation

### Homebrew (macOS/Linux)

```bash
brew install Arnel-rah/tunepipe/tunepipe
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
2. `tunepipe` shells out to `yt-dlp` to resolve a direct stream URL, preferring low-bitrate audio-only formats.
3. The resolved URL is cached; while you're browsing or listening to the current track, `tunepipe` prefetches the next one in the background.
4. Playback is handed off to `mpv` for decoding/output, minimizing the gap between selecting a track and hearing audio.

## Troubleshooting

**`Signature solving failed` / `n challenge solving failed` warnings**
`yt-dlp` needs an up-to-date JS challenge solver. `tunepipe` passes `--remote-components ejs:github` by default, but if you're calling `yt-dlp` manually, add the same flag or install [Deno](https://deno.land) as a local runtime.

**`This video is unavailable`**
The video itself may be private, deleted, or region-restricted — check it directly in a browser before filing an issue.

**Stale `yt-dlp` behind an old system install**
If you have multiple `yt-dlp` binaries on your `PATH` (e.g. one from `apt` and one from `pipx`), the wrong one may take priority. Check with:
```bash
which yt-dlp
yt-dlp --version
```
Remove the outdated system package if needed (`sudo apt remove yt-dlp`, `hash -r`).

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
git tag -a v0.2.2 -m "v0.2.2: fix JS challenge solver, reduce data usage"
git push origin v0.2.2
```

## Contributing

Issues and pull requests are welcome. If you're fixing a bug or adding a feature:

1. Fork the repo and create a branch (`fix/...`, `feat/...`)
2. Follow the existing commit convention (`type(scope): message`)
3. Open a PR describing what changed and why

## License

MIT — see [LICENSE](LICENSE) for details.