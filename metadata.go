package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// VideoInfo holds parsed metadata from yt-dlp --dump-json.
type VideoInfo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Uploader  string `json:"uploader"`
	Duration  int    `json:"duration"`
	Thumbnail string `json:"thumbnail"`
	WebpageURL string `json:"webpage_url"`
	Playlist  string `json:"playlist_title"`
}

// FetchMetadata runs yt-dlp --dump-json for a single URL and returns parsed info.
func FetchMetadata(url string, verbose bool) (*VideoInfo, error) {
	done := make(chan struct{})
	go showSpinner("Fetching metadata", done)

	cmd := exec.Command("yt-dlp", "--dump-json", "--no-warnings", "--no-playlist", url)
	Verbose(verbose, "Running: %s", strings.Join(cmd.Args, " "))

	out, err := cmd.Output()
	close(done)
	ClearLine()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch metadata: %w", err)
	}

	var info VideoInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return nil, fmt.Errorf("failed to parse metadata: %w", err)
	}
	return &info, nil
}

// FetchPlaylistItems returns metadata for every item in a playlist.
func FetchPlaylistItems(url string, verbose bool) ([]VideoInfo, error) {
	done := make(chan struct{})
	go showSpinner("Fetching playlist", done)

	cmd := exec.Command("yt-dlp", "--flat-playlist", "--dump-json", "--no-warnings", url)
	Verbose(verbose, "Running: %s", strings.Join(cmd.Args, " "))

	out, err := cmd.Output()
	close(done)
	ClearLine()

	if err != nil {
		return nil, fmt.Errorf("failed to fetch playlist: %w", err)
	}

	var items []VideoInfo
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		var info VideoInfo
		if err := json.Unmarshal([]byte(line), &info); err != nil {
			continue
		}
		// flat-playlist entries have an "id" but may lack a full webpage_url
		if info.WebpageURL == "" && info.ID != "" {
			info.WebpageURL = "https://www.youtube.com/watch?v=" + info.ID
		}
		items = append(items, info)
	}
	return items, nil
}

// SearchYouTube searches YouTube via yt-dlp and returns up to n results.
func SearchYouTube(query string, n int, verbose bool) ([]VideoInfo, error) {
	done := make(chan struct{})
	go showSpinner("Searching YouTube", done)

	searchTerm := fmt.Sprintf("ytsearch%d:%s", n, query)
	cmd := exec.Command("yt-dlp", searchTerm, "--dump-json", "--no-warnings", "--no-download")
	Verbose(verbose, "Running: %s", strings.Join(cmd.Args, " "))

	out, err := cmd.Output()
	close(done)
	ClearLine()

	if err != nil {
		return nil, fmt.Errorf("search failed: %w", err)
	}

	var results []VideoInfo
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		var info VideoInfo
		if err := json.Unmarshal([]byte(line), &info); err != nil {
			continue
		}
		results = append(results, info)
	}
	return results, nil
}

// PromptSelection asks the user to pick a search result (1-indexed).
func PromptSelection(count int) (int, error) {
	fmt.Printf("  %sSelect a track (1-%d):%s ", Bold, count, Reset)
	var input string
	fmt.Scanln(&input)
	choice, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil || choice < 1 || choice > count {
		return 0, fmt.Errorf("invalid selection: %s", input)
	}
	return choice - 1, nil
}

func showSpinner(label string, done <-chan struct{}) {
	tick := 0
	for {
		select {
		case <-done:
			return
		default:
			ClearLine()
			fmt.Printf("  %s %s", SpinnerFrame(tick), label)
			tick++
			time.Sleep(80 * time.Millisecond)
		}
	}
}

// IsPlaylistURL is a simple heuristic to detect playlist URLs.
func IsPlaylistURL(url string) bool {
	return strings.Contains(url, "playlist?list=") || strings.Contains(url, "&list=")
}

// IsValidYouTubeURL checks if the provided URL is a valid YouTube URL.
func IsValidYouTubeURL(url string) bool {
	patterns := []string{
		`^https?://(www\.)?youtube\.com/watch\?v=[\w-]+`,
		`^https?://youtu\.be/[\w-]+`,
		`^https?://(www\.)?youtube\.com/embed/[\w-]+`,
		`^https?://(www\.)?youtube\.com/playlist\?list=[\w-]+`,
	}
	for _, pattern := range patterns {
		matched, _ := regexp.MatchString(pattern, url)
		if matched {
			return true
		}
	}
	return false
}
