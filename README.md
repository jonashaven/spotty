# spotty

Spotify CLI for your terminal & tmux. Shows what's playing, controls playback.

## Install

```bash
curl -L https://github.com/jonashaven/spotty/releases/latest/download/spotty-$(uname -s | tr '[:upper:]' '[:lower:]')-$(uname -m | sed 's/x86_64/amd64/' | sed 's/aarch64/arm64/') -o spotty
chmod +x spotty
sudo mv spotty /usr/local/bin/spotty
```

## Setup

1. Create/reuse a Spotify app at https://developer.spotify.com/dashboard (Web API)
2. Add `http://127.0.0.1:8888/callback` as redirect URI
3. Run `spotty login`

## Usage

```
spotty              # show current track
spotty play         # resume
spotty pause        # pause
spotty next         # skip
spotty prev         # previous
spotty update       # self-update
```

## Tmux

Add to `.tmux.conf`:

```tmux
set -g status-interval 30
set -g status-right "#[fg=#a6e3a1]#(command -v spotty >/dev/null && spotty --short 2>/dev/null) #[fg=#cdd6f4,bg=#313244] %H:%M "
```

API is cached locally — only called every 30s regardless of session count.

## License

MIT
