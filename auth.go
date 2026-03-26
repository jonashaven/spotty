package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

const (
	authURL  = "https://accounts.spotify.com/authorize"
	tokenURL = "https://accounts.spotify.com/api/token"
	redirect = "http://127.0.0.1:8888/callback"
	scopes   = "user-read-currently-playing"
)

type tokenResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"`
}

func Login(cfg *Config) error {
	state := randomState()

	params := url.Values{
		"client_id":     {cfg.ClientID},
		"response_type": {"code"},
		"redirect_uri":  {redirect},
		"scope":         {scopes},
		"state":         {state},
	}

	authLink := authURL + "?" + params.Encode()
	fmt.Println("Open this URL to log in:")
	fmt.Println(authLink)
	fmt.Println()
	fmt.Print("Paste the redirect URL here: ")

	reader := bufio.NewReader(os.Stdin)
	rawURL, _ := reader.ReadString('\n')
	rawURL = strings.TrimSpace(rawURL)

	parsed, err := url.Parse(rawURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	if s := parsed.Query().Get("state"); s != state {
		return fmt.Errorf("state mismatch")
	}
	if e := parsed.Query().Get("error"); e != "" {
		return fmt.Errorf("auth error: %s", e)
	}

	code := parsed.Query().Get("code")
	if code == "" {
		return fmt.Errorf("no code in URL")
	}

	tok, err := exchangeCode(cfg, code)
	if err != nil {
		return err
	}

	cfg.AccessToken = tok.AccessToken
	cfg.RefreshToken = tok.RefreshToken
	cfg.Expiry = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	return SaveConfig(cfg)
}

func exchangeCode(cfg *Config, code string) (*tokenResponse, error) {
	data := url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {code},
		"redirect_uri": {redirect},
	}
	return doTokenRequest(cfg, data)
}

func RefreshAccessToken(cfg *Config) error {
	data := url.Values{
		"grant_type":    {"refresh_token"},
		"refresh_token": {cfg.RefreshToken},
	}
	tok, err := doTokenRequest(cfg, data)
	if err != nil {
		return err
	}
	cfg.AccessToken = tok.AccessToken
	if tok.RefreshToken != "" {
		cfg.RefreshToken = tok.RefreshToken
	}
	cfg.Expiry = time.Now().Add(time.Duration(tok.ExpiresIn) * time.Second)
	return SaveConfig(cfg)
}

func doTokenRequest(cfg *Config, data url.Values) (*tokenResponse, error) {
	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.SetBasicAuth(cfg.ClientID, cfg.ClientSecret)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("token error (%d): %s", resp.StatusCode, body)
	}

	var tok tokenResponse
	if err := json.Unmarshal(body, &tok); err != nil {
		return nil, err
	}
	return &tok, nil
}

func randomState() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
