package main

import (
	"bufio"
	"bytes"
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

	var stderrBuf bytes.Buffer
	cmd.Stderr = &stderrBuf

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
		errMsg := strings.TrimSpace(stderrBuf.String())
		if errMsg != "" {
			return fmt.Errorf("yt-dlp failed: %s", errMsg)
		}
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

