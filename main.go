package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)

func main() {
	PrintBanner()

	// Check dependencies
	fmt.Println("Checking dependencies...")
	if err := CheckDependencies(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Dependencies OK!")
	fmt.Println()

	home, _ := os.UserHomeDir()
	downloadDir := filepath.Join(home, "Downloads")

	opts := Options{
		Quality:   320,
		OutputDir: downloadDir,
	}

	for {
		fmt.Print("Enter YouTube URL (or 'quit' to exit): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)

		if input == "quit" || input == "q" {
			fmt.Println("Goodbye!")
			break
		}

		if input == "" {
			continue
		}

		// Validate URL
		if !IsValidYouTubeURL(input) {
			fmt.Println("Error: Invalid YouTube URL format")
			fmt.Println("Supported formats:")
			fmt.Println("  - https://www.youtube.com/watch?v=VIDEO_ID")
			fmt.Println("  - https://youtu.be/VIDEO_ID")
			fmt.Println("  - https://www.youtube.com/embed/VIDEO_ID")
			fmt.Println("  - https://www.youtube.com/playlist?list=PLAYLIST_ID")
			fmt.Println()
			continue
		}

		// Route to playlist or single download
		if IsPlaylistURL(input) {
			if err := HandlePlaylist(input, opts); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		} else {
			if err := DownloadAndConvert(input, opts); err != nil {
				fmt.Printf("Error: %v\n", err)
			}
		}
		fmt.Println()
	}
}
