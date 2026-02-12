# Konversion

```
 ██╗  ██╗ ██████╗ ███╗   ██╗██╗   ██╗███████╗██████╗ ███████╗██╗ ██████╗ ███╗   ██╗
 ██║ ██╔╝██╔═══██╗████╗  ██║██║   ██║██╔════╝██╔══██╗██╔════╝██║██╔═══██╗████╗  ██║
 █████╔╝ ██║   ██║██╔██╗ ██║██║   ██║█████╗  ██████╔╝███████╗██║██║   ██║██╔██╗ ██║
 ██╔═██╗ ██║   ██║██║╚██╗██║╚██╗ ██╔╝██╔══╝  ██╔══██╗╚════██║██║██║   ██║██║╚██╗██║
 ██║  ██╗╚██████╔╝██║ ╚████║ ╚████╔╝ ███████╗██║  ██║███████║██║╚██████╔╝██║ ╚████║
 ╚═╝  ╚═╝ ╚═════╝ ╚═╝  ╚═══╝  ╚═══╝  ╚══════╝╚═╝  ╚═╝╚══════╝╚═╝ ╚═════╝ ╚═╝  ╚═╝
```

**YouTube to MP3 converter for the terminal.** Paste a link, get an MP3. That's it.

Built for producers and artists who just want to rip audio from YouTube without dealing with sketchy websites or bloated apps.

---

## What It Does

- Paste any YouTube URL and it downloads + converts to high quality MP3 (320kbps)
- Shows you the track title, channel, and duration before downloading
- Progress bar so you can see the download happening
- Automatically embeds metadata (title, artist, thumbnail) into the MP3
- Supports full playlists — paste a playlist URL and it grabs every track
- MP3s are saved to your Downloads folder
- Shows you the exact file path and size when it's done

## Quick Install (macOS)

You need 3 things: **Homebrew**, **Go**, and two free tools (**yt-dlp** + **ffmpeg**). Here's how to get everything set up from scratch.

### Step 1: Install Homebrew (if you don't have it)

Open **Terminal** (search "Terminal" in Spotlight) and paste this:

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Follow the instructions on screen. When it's done, close Terminal and reopen it.

### Step 2: Install Go, yt-dlp, and ffmpeg

```bash
brew install go yt-dlp ffmpeg
```

This might take a minute. Wait for it to finish.

### Step 3: Download and build Konversion

```bash
git clone https://github.com/Stevenalvarado025/Konversion.git ~/Konversion
cd ~/Konversion
go build -o konversion .
```

### Step 4: Make it available everywhere (optional)

```bash
sudo cp ~/Konversion/konversion /usr/local/bin/konversion
```

It'll ask for your password. After that you can type `konversion` from anywhere.

---

## Quick Install (Linux)

### Ubuntu / Debian

```bash
sudo apt update
sudo apt install golang ffmpeg
pip install yt-dlp
```

### Then build Konversion

```bash
git clone https://github.com/Stevenalvarado025/Konversion.git ~/Konversion
cd ~/Konversion
go build -o konversion .
sudo cp ~/Konversion/konversion /usr/local/bin/konversion
```

---

## Quick Install (Windows)

1. Install [Go](https://go.dev/dl/) — download the Windows installer and run it
2. Install [ffmpeg](https://www.gyan.dev/ffmpeg/builds/) — download "ffmpeg-release-essentials.zip", extract it, and add the `bin` folder to your PATH
3. Install [yt-dlp](https://github.com/yt-dlp/yt-dlp/releases) — download `yt-dlp.exe` and put it in a folder that's in your PATH
4. Open Command Prompt or PowerShell:

```bash
git clone https://github.com/Stevenalvarado025/Konversion.git %USERPROFILE%\Konversion
cd %USERPROFILE%\Konversion
go build -o konversion.exe .
```

---

## How to Use

Just run it:

```bash
konversion
```

You'll see the Konversion banner and a prompt. That's where you paste your YouTube URL:

```
Enter YouTube URL (or 'quit' to exit): https://www.youtube.com/watch?v=dQw4w9WgXcQ

  Title:    Rick Astley - Never Gonna Give You Up
  Channel:  Rick Astley
  Duration: 3:33

  ℹ Downloading and converting to MP3 (320 kbps)...
  ██████████████████████████████ 100.0%
Saved: /Users/you/Downloads/Rick Astley - Never Gonna Give You Up.mp3 (8.2 MB)

Enter YouTube URL (or 'quit' to exit):
```

It loops — so you can keep pasting URLs back to back. Type `quit` when you're done.

### Download a Playlist

Paste a playlist URL and it grabs every track:

```
Enter YouTube URL (or 'quit' to exit): https://www.youtube.com/playlist?list=PLxxxxxxx

  ℹ Found 12 tracks in playlist

  [1/12] Song One
  [2/12] Song Two
  ...
All 12 tracks downloaded to: /Users/you/Downloads
```

### Where Do My MP3s Go?

Your **Downloads** folder. After every download it prints the full path so you know exactly where to find it.

---

## Troubleshooting

### "yt-dlp is not installed"

```bash
brew install yt-dlp        # macOS
pip install yt-dlp          # any system with Python
sudo apt install yt-dlp     # Ubuntu/Debian
```

### "ffmpeg is not installed"

```bash
brew install ffmpeg         # macOS
sudo apt install ffmpeg     # Ubuntu/Debian
```

### "Invalid YouTube URL format"

Make sure you're pasting the full URL. These formats work:

- `https://www.youtube.com/watch?v=VIDEO_ID`
- `https://youtu.be/VIDEO_ID`
- `https://www.youtube.com/embed/VIDEO_ID`
- `https://www.youtube.com/playlist?list=PLAYLIST_ID`

### "go: command not found"

Install Go first:

```bash
brew install go             # macOS
sudo apt install golang     # Ubuntu/Debian
```

Or download from [go.dev/dl](https://go.dev/dl/)

### How do I update Konversion?

```bash
cd ~/Konversion
git pull
go build -o konversion .
sudo cp konversion /usr/local/bin/konversion
```

---

## How It Works

Konversion is a Go wrapper around two industry-standard open source tools:

- **[yt-dlp](https://github.com/yt-dlp/yt-dlp)** — downloads video/audio from YouTube
- **[ffmpeg](https://ffmpeg.org/)** — converts audio to MP3 format

Konversion handles the terminal UI, progress bar, metadata display, and makes the whole process as simple as paste-and-go. No Python scripts, no browser extensions, no sketchy websites with popups.

## License

MIT — do whatever you want with it.
