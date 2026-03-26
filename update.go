package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
)

const repoAPI = "https://api.github.com/repos/jonashaven/spotty/releases/latest"

var version = "dev"

type ghRelease struct {
	TagName string    `json:"tag_name"`
	Assets  []ghAsset `json:"assets"`
}

type ghAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

func runUpdate() {
	resp, err := http.Get(repoAPI)
	if err != nil {
		fatal(err)
	}
	defer resp.Body.Close()

	var rel ghRelease
	if err := json.NewDecoder(resp.Body).Decode(&rel); err != nil {
		fatal(err)
	}

	if rel.TagName == version {
		fmt.Printf("Already on latest (%s)\n", version)
		return
	}

	fmt.Printf("Update available: %s → %s\n", version, rel.TagName)
	fmt.Print("Install? [y/N] ")
	reader := bufio.NewReader(os.Stdin)
	answer, _ := reader.ReadString('\n')
	if strings.TrimSpace(strings.ToLower(answer)) != "y" {
		fmt.Println("Skipped.")
		return
	}

	fmt.Println("Downloading...")

	assetName := fmt.Sprintf("spotty-%s-%s", runtime.GOOS, runtime.GOARCH)
	if runtime.GOOS == "windows" {
		assetName += ".exe"
	}

	var downloadURL string
	for _, a := range rel.Assets {
		if a.Name == assetName {
			downloadURL = a.BrowserDownloadURL
			break
		}
	}
	if downloadURL == "" {
		fatal(fmt.Errorf("no binary for %s/%s", runtime.GOOS, runtime.GOARCH))
	}

	binResp, err := http.Get(downloadURL)
	if err != nil {
		fatal(err)
	}
	defer binResp.Body.Close()

	exe, err := os.Executable()
	if err != nil {
		fatal(err)
	}
	// Resolve symlinks
	if resolved, err := resolveSymlink(exe); err == nil {
		exe = resolved
	}

	tmp := exe + ".tmp"
	f, err := os.OpenFile(tmp, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
	if err != nil {
		if strings.Contains(err.Error(), "permission denied") {
			fmt.Fprintln(os.Stderr, "Permission denied. Try: sudo spotty update")
			os.Exit(1)
		}
		fatal(err)
	}

	if _, err := io.Copy(f, binResp.Body); err != nil {
		f.Close()
		os.Remove(tmp)
		fatal(err)
	}
	f.Close()

	if err := os.Rename(tmp, exe); err != nil {
		os.Remove(tmp)
		fatal(err)
	}

	fmt.Printf("Updated to %s\n", rel.TagName)
}

func resolveSymlink(path string) (string, error) {
	info, err := os.Lstat(path)
	if err != nil {
		return "", err
	}
	if info.Mode()&os.ModeSymlink != 0 {
		return os.Readlink(path)
	}
	return path, nil
}
