package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

func main() {
	cmd := "now"
	short := false
	refresh := false
	for _, arg := range os.Args[1:] {
		switch arg {
		case "--short":
			short = true
		case "--refresh":
			refresh = true
			short = true
		default:
			cmd = arg
		}
	}

	switch cmd {
	case "login":
		runLogin()
	case "now":
		runNow(short, refresh)
	case "update":
		runUpdate()
	case "version":
		fmt.Println(version)
	default:
		fmt.Fprintf(os.Stderr, "Usage: spotty [login|now|update|version] [--short]\n")
		os.Exit(1)
	}
}

func runLogin() {
	cfg, _ := LoadConfig()

	reader := bufio.NewReader(os.Stdin)

	if cfg.ClientID == "" {
		fmt.Print("Client ID: ")
		cfg.ClientID, _ = reader.ReadString('\n')
		cfg.ClientID = strings.TrimSpace(cfg.ClientID)
	}
	if cfg.ClientSecret == "" {
		fmt.Print("Client Secret: ")
		cfg.ClientSecret, _ = reader.ReadString('\n')
		cfg.ClientSecret = strings.TrimSpace(cfg.ClientSecret)
	}

	if err := SaveConfig(cfg); err != nil {
		fatal(err)
	}

	if err := Login(cfg); err != nil {
		fatal(err)
	}
	fmt.Println("Logged in!")
}

func runNow(short bool, refresh bool) {
	// In short mode, use cache to avoid hitting API too often
	if short && !refresh {
		if cache, err := LoadCache(); err == nil && time.Since(cache.FetchedAt) < 30*time.Second {
			printShort(cache.Text, cache.IsPlaying)
			return
		}
	}

	cfg, err := LoadConfig()
	if err != nil || cfg.AccessToken == "" {
		if short {
			return // silent in tmux
		}
		fmt.Fprintln(os.Stderr, "Not logged in. Run: spotty login")
		os.Exit(1)
	}

	if time.Now().After(cfg.Expiry) {
		if err := RefreshAccessToken(cfg); err != nil {
			if short {
				return
			}
			fmt.Fprintln(os.Stderr, "Token refresh failed. Run: spotty login")
			os.Exit(1)
		}
	}

	err = GetNowPlaying(cfg.AccessToken, short)
	if err != nil && err.Error() == "token expired" {
		if err := RefreshAccessToken(cfg); err != nil {
			if short {
				return
			}
			fmt.Fprintln(os.Stderr, "Token refresh failed. Run: spotty login")
			os.Exit(1)
		}
		err = GetNowPlaying(cfg.AccessToken, short)
	}
	if err != nil && !short {
		fatal(err)
	}
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
