---
name: spotty
description: Install and configure spotty — a Spotify CLI for now-playing display and playback controls, with tmux status bar integration
---

# Spotty

Spotify CLI built in Go. Shows currently playing track and controls playback via Spotify Web API.

## Install

Download the latest binary for the user's platform:

```bash
curl -L https://github.com/jonashaven/spotty/releases/latest/download/spotty-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/') -o spotty
chmod +x spotty
sudo mv spotty /usr/local/bin/spotty
```

Or with Go: `go install github.com/jonashaven/spotty@latest`

## Setup

Requires a Spotify Developer app:

1. Go to https://developer.spotify.com/dashboard
2. Create or reuse an app with **Web API** enabled
3. Add `http://127.0.0.1:8888/callback` as a redirect URI
4. Run `spotty login` — enter Client ID and Secret when prompted
5. Open the printed URL in a browser, log in to Spotify
6. Copy the redirect URL from the browser (page won't load — that's expected) and paste it back into the terminal

Credentials are stored at `~/.config/spotty/config.json` with restricted permissions (0600).

## Commands

- `spotty` or `spotty now` — show currently playing track (artist + album)
- `spotty now --short` — title only, single line (for tmux/scripts)
- `spotty play` — resume playback
- `spotty pause` — pause playback
- `spotty next` — skip to next track
- `spotty prev` — go to previous track
- `spotty login` — authenticate with Spotify
- `spotty update` — check for updates and self-update
- `spotty version` — print version

## Tmux Integration

Add to `.tmux.conf`:

```tmux
set -g status-interval 30
set -g status-right "#[fg=#a6e3a1]#(command -v spotty >/dev/null && spotty --short 2>/dev/null) #[fg=#cdd6f4,bg=#313244] %H:%M "
set -g status-right-length 100

# Optional: click status bar to refresh
bind -n MouseDown1StatusRight run-shell "command -v spotty >/dev/null && spotty --refresh >/dev/null 2>&1" \; refresh-client -S
```

The `--short` flag uses a local cache (`~/.config/spotty/cache.json`) so the Spotify API is only called every 30 seconds, regardless of how many tmux sessions are running.

## Notes

- Requires Spotify Premium for playback controls (play/pause/next/prev)
- Now-playing works on any Spotify tier
- Token auto-refreshes; re-login only needed if refresh token is revoked
