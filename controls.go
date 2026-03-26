package main

import (
	"fmt"
	"net/http"
)

const playerBase = "https://api.spotify.com/v1/me/player"

func runControl(cfg *Config, action string) {
	var method, endpoint string
	switch action {
	case "play":
		method, endpoint = "PUT", playerBase+"/play"
	case "pause":
		method, endpoint = "PUT", playerBase+"/pause"
	case "next":
		method, endpoint = "POST", playerBase+"/next"
	case "prev":
		method, endpoint = "POST", playerBase+"/previous"
	}

	req, err := http.NewRequest(method, endpoint, nil)
	if err != nil {
		fatal(err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.AccessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 || resp.StatusCode == 200 {
		return
	}
	if resp.StatusCode == 401 {
		fatal(fmt.Errorf("token expired. Run: spotty login"))
	}
	if resp.StatusCode == 404 {
		fatal(fmt.Errorf("no active device found"))
	}
	if resp.StatusCode == 403 {
		fatal(fmt.Errorf("premium required for playback control"))
	}
	fatal(fmt.Errorf("API error (%d)", resp.StatusCode))
}
