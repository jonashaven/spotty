package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)
const playerURL = "https://api.spotify.com/v1/me/player/currently-playing"

type nowPlayingResponse struct {
	IsPlaying bool `json:"is_playing"`
	Item      struct {
		Name    string `json:"name"`
		Artists []struct {
			Name string `json:"name"`
		} `json:"artists"`
		Album struct {
			Name string `json:"name"`
		} `json:"album"`
	} `json:"item"`
}

func GetNowPlaying(accessToken string, short bool) error {
	req, err := http.NewRequest("GET", playerURL, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == 204 || resp.StatusCode == 202 {
		SaveCache(&Cache{})
		if !short {
			fmt.Println("Nothing playing right now.")
		}
		return nil
	}

	if resp.StatusCode == 401 {
		return fmt.Errorf("token expired")
	}

	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API error (%d): %s", resp.StatusCode, body)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var np nowPlayingResponse
	if err := json.Unmarshal(body, &np); err != nil {
		return err
	}

	if np.Item.Name == "" {
		fmt.Println("Nothing playing right now.")
		return nil
	}

	artists := make([]string, len(np.Item.Artists))
	for i, a := range np.Item.Artists {
		artists[i] = a.Name
	}

	status := "▶"
	if !np.IsPlaying {
		status = "⏸"
	}

	text := fmt.Sprintf("%s — %s", np.Item.Name, strings.Join(artists, ", "))

	SaveCache(&Cache{
		Text:      np.Item.Name,
		Artists:   strings.Join(artists, ", "),
		Album:     np.Item.Album.Name,
		IsPlaying: np.IsPlaying,
		FetchedAt: time.Now(),
	})

	if short {
		printShort(np.Item.Name, np.IsPlaying)
	} else {
		fmt.Printf("%s %s\n", status, text)
		fmt.Printf("  %s\n", np.Item.Album.Name)
	}
	return nil
}

func printFull(c *Cache) {
	if c.Text == "" {
		fmt.Println("Nothing playing right now.")
		return
	}
	status := "▶"
	if !c.IsPlaying {
		status = "⏸"
	}
	fmt.Printf("%s %s — %s\n", status, c.Text, c.Artists)
	fmt.Printf("  %s\n", c.Album)
}

func printShort(text string, isPlaying bool) {
	status := "▶"
	if !isPlaying {
		status = "⏸"
	}
	fmt.Printf("%s %s\n", status, text)
}
