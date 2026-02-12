package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Options holds all CLI-derived settings.
type Options struct {
	Quality   int
	OutputDir string
	KeepVideo bool
	NoMeta    bool
	Verbose   bool
}

// CheckDependencies verifies that yt-dlp and ffmpeg are installed.
func CheckDependencies() error {
	missing := []string{}

	if _, err := exec.LookPath("yt-dlp"); err != nil {
		missing = append(missing, "yt-dlp")
	}
	if _, err := exec.LookPath("ffmpeg"); err != nil {
		missing = append(missing, "ffmpeg")
	}

	if len(missing) > 0 {
		msg := "Missing required dependencies: " + strings.Join(missing, ", ") + "\n\n"
		msg += "  Install with:\n"
		for _, dep := range missing {
			switch dep {
			case "yt-dlp":
				msg += "    brew install yt-dlp       (macOS)\n"
				msg += "    pip install yt-dlp         (pip)\n"
				msg += "    sudo apt install yt-dlp    (Debian/Ubuntu)\n"
			case "ffmpeg":
				msg += "    brew install ffmpeg        (macOS)\n"
				msg += "    sudo apt install ffmpeg    (Debian/Ubuntu)\n"
			}
		}
		return fmt.Errorf("%s", msg)
	}
	return nil
}

// DownloadAndConvert downloads a single URL and converts it to MP3.
func DownloadAndConvert(url string, opts Options) error {
	info, err := FetchMetadata(url, opts.Verbose)
	if err != nil {
		return err
	}
	PrintTrackInfo(info)

	Info("Downloading and converting to MP3 (%d kbps)…", opts.Quality)

	args := []string{
		url,
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", strconv.Itoa(opts.Quality) + "K",
		"--newline", // one progress line per update
		"-o", filepath.Join(opts.OutputDir, "%(title)s.%(ext)s"),
	}

	if !opts.NoMeta {
		args = append(args, "--embed-thumbnail", "--embed-metadata")
		// yt-dlp needs atomicparsley or mutagen for thumbnails on mp3;
		// --convert-thumbnails jpg ensures compatibility
		args = append(args, "--convert-thumbnails", "jpg")
	}

	if opts.KeepVideo {
		args = append(args, "-k")
	}

	Verbose(opts.Verbose, "Running: yt-dlp %s", strings.Join(args, " "))

	cmd := exec.Command("yt-dlp", args...)
	cmd.Stderr = nil

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to capture output: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start yt-dlp: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	progressRe := regexp.MustCompile(`\[download\]\s+([\d.]+)%\s+of`)
	for scanner.Scan() {
		line := scanner.Text()
		Verbose(opts.Verbose, "%s", line)

		if matches := progressRe.FindStringSubmatch(line); matches != nil {
			pct, _ := strconv.ParseFloat(matches[1], 64)
			ProgressBar(pct, 30, "")
		}
	}

	if err := cmd.Wait(); err != nil {
		ClearLine()
		return fmt.Errorf("yt-dlp failed: %w", err)
	}

	ClearLine()

	// Show full file path and size
	savedPath := filepath.Join(opts.OutputDir, info.Title+".mp3")
	if fi, err := os.Stat(savedPath); err == nil {
		sizeMB := float64(fi.Size()) / 1024 / 1024
		fmt.Printf("Saved: %s (%.1f MB)\n", savedPath, sizeMB)
	} else {
		fmt.Printf("Saved to: %s\n", opts.OutputDir)
	}
	return nil
}

// HandlePlaylist downloads all items in a playlist.
func HandlePlaylist(url string, opts Options) error {
	items, err := FetchPlaylistItems(url, opts.Verbose)
	if err != nil {
		return err
	}
	if len(items) == 0 {
		return fmt.Errorf("no items found in playlist")
	}

	Info("Found %d tracks in playlist", len(items))
	fmt.Println()

	var errors []string
	for i, item := range items {
		fmt.Printf("  %s[%d/%d]%s %s\n", BoldCyan, i+1, len(items), Reset, item.Title)
		dlURL := item.WebpageURL
		if dlURL == "" {
			Warn("Skipping item %d: no URL", i+1)
			continue
		}
		if err := downloadSingle(dlURL, opts); err != nil {
			Error("Failed: %v", err)
			errors = append(errors, fmt.Sprintf("%s: %v", item.Title, err))
			continue
		}
	}

	fmt.Println()
	if len(errors) > 0 {
		Warn("%d of %d tracks failed", len(errors), len(items))
		for _, e := range errors {
			Error("  %s", e)
		}
		return fmt.Errorf("%d downloads failed", len(errors))
	}
	fmt.Printf("All %d tracks downloaded to: %s\n", len(items), opts.OutputDir)
	return nil
}

// HandleBatch reads URLs from a file and processes each one.
func HandleBatch(path string, opts Options) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open batch file: %w", err)
	}
	defer f.Close()

	var urls []string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" && !strings.HasPrefix(line, "#") {
			urls = append(urls, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading batch file: %w", err)
	}
	if len(urls) == 0 {
		return fmt.Errorf("no URLs found in %s", path)
	}

	Info("Batch mode: %d URLs loaded from %s", len(urls), path)
	fmt.Println()

	var errors []string
	for i, url := range urls {
		fmt.Printf("  %s[%d/%d]%s %s\n", BoldCyan, i+1, len(urls), Reset, url)
		if err := DownloadAndConvert(url, opts); err != nil {
			Error("Failed: %v", err)
			errors = append(errors, fmt.Sprintf("%s: %v", url, err))
		}
		fmt.Println()
	}

	if len(errors) > 0 {
		Warn("%d of %d downloads failed", len(errors), len(urls))
		return fmt.Errorf("%d downloads failed", len(errors))
	}
	Success("All %d downloads completed!", len(urls))
	return nil
}

// HandleSearch performs a YouTube search and lets the user pick a result.
func HandleSearch(query string, opts Options) error {
	results, err := SearchYouTube(query, 5, opts.Verbose)
	if err != nil {
		return err
	}
	if len(results) == 0 {
		return fmt.Errorf("no results found for: %s", query)
	}

	PrintSearchResults(results)

	idx, err := PromptSelection(len(results))
	if err != nil {
		return err
	}

	selected := results[idx]
	url := selected.WebpageURL
	if url == "" {
		url = "https://www.youtube.com/watch?v=" + selected.ID
	}

	fmt.Println()
	return DownloadAndConvert(url, opts)
}

// downloadSingle is a lightweight download without re-fetching metadata display
// (used in playlist context where we already printed the track title).
func downloadSingle(url string, opts Options) error {
	args := []string{
		url,
		"--extract-audio",
		"--audio-format", "mp3",
		"--audio-quality", strconv.Itoa(opts.Quality) + "K",
		"--newline",
		"-o", filepath.Join(opts.OutputDir, "%(title)s.%(ext)s"),
	}

	if !opts.NoMeta {
		args = append(args, "--embed-thumbnail", "--embed-metadata")
		args = append(args, "--convert-thumbnails", "jpg")
	}

	if opts.KeepVideo {
		args = append(args, "-k")
	}

	Verbose(opts.Verbose, "Running: yt-dlp %s", strings.Join(args, " "))

	cmd := exec.Command("yt-dlp", args...)
	cmd.Stderr = nil

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to capture output: %w", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start yt-dlp: %w", err)
	}

	scanner := bufio.NewScanner(stdout)
	progressRe := regexp.MustCompile(`\[download\]\s+([\d.]+)%\s+of`)
	for scanner.Scan() {
		line := scanner.Text()
		Verbose(opts.Verbose, "%s", line)

		if matches := progressRe.FindStringSubmatch(line); matches != nil {
			pct, _ := strconv.ParseFloat(matches[1], 64)
			ProgressBar(pct, 30, "")
		}
	}

	if err := cmd.Wait(); err != nil {
		ClearLine()
		return fmt.Errorf("yt-dlp failed: %w", err)
	}

	ClearLine()
	fmt.Println("Done")
	return nil
}
