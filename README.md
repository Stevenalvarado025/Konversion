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

---

## Install (macOS)

Open **Terminal** (search "Terminal" in Spotlight) and follow these steps. **Run each command one at a time** — wait for each one to finish before pasting the next.

### Step 1: Install Homebrew (if you don't have it)

```bash
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
```

Follow the instructions on screen. When it's done, **close Terminal and reopen it**.

### Step 2: Install Go, yt-dlp, and ffmpeg

```bash
brew install go yt-dlp ffmpeg
```

Wait for it to finish. This might take a few minutes.

### Step 3: Download Konversion

```bash
git clone https://github.com/Stvn444/Konversion.git ~/Konversion
```

> **Already have a ~/Konversion folder?** Delete it first: `rm -rf ~/Konversion` then run the clone command again.

### Step 4: Build it

```bash
cd ~/Konversion
```

```bash
go build -o konversion .
```

### Step 5: Install it

```bash
sudo cp konversion /usr/local/bin/konversion
```

It will ask for your Mac password. Type it in and press Enter. You won't see the characters as you type — that's normal.

### Step 6: Run it

Open a **new Terminal window** and type:

```bash
konversion
```

> **IMPORTANT:** The command is all lowercase: `konversion` (not `Konversion`)

You should see the Konversion banner and a prompt to paste a URL. You're good to go!

---

## Install (Linux)

Run each command one at a time:

```bash
sudo apt update
```

```bash
sudo apt install golang ffmpeg
```

```bash
pip install yt-dlp
```

```bash
git clone https://github.com/Stvn444/Konversion.git ~/Konversion
```

```bash
cd ~/Konversion
```

```bash
go build -o konversion .
```

```bash
sudo cp konversion /usr/local/bin/konversion
```

Then run:

```bash
konversion
```

---

## Install (Windows)

1. Install [Go](https://go.dev/dl/) — download the Windows installer and run it
2. Install [Git](https://git-scm.com/downloads/win) — download and run the installer
3. Install [ffmpeg](https://www.gyan.dev/ffmpeg/builds/) — download "ffmpeg-release-essentials.zip", extract it, and add the `bin` folder to your PATH
4. Install [yt-dlp](https://github.com/yt-dlp/yt-dlp/releases) — download `yt-dlp.exe` and put it in a folder that's in your PATH

Then open Command Prompt or PowerShell and run each command one at a time:

```bash
git clone https://github.com/Stvn444/Konversion.git %USERPROFILE%\Konversion
```

```bash
cd %USERPROFILE%\Konversion
```

```bash
go build -o konversion.exe .
```

Then to run it:

```bash
%USERPROFILE%\Konversion\konversion.exe
```

---

## How to Use

Just type `konversion` in Terminal (all lowercase):

```bash
konversion
```

You'll see the Konversion banner and a prompt. Paste your YouTube URL and press Enter:

```
Enter YouTube URL (or 'quit' to exit): https://www.youtube.com/watch?v=dQw4w9WgXcQ

  Title:    Rick Astley - Never Gonna Give You Up
  Channel:  Rick Astley
  Duration: 3:33

  Downloading and converting to MP3 (320 kbps)...
  ██████████████████████████████ 100.0%
Saved: /Users/you/Downloads/Rick Astley - Never Gonna Give You Up.mp3 (8.2 MB)

Enter YouTube URL (or 'quit' to exit):
```

It loops — keep pasting URLs back to back. Type `quit` when you're done.

### Download a Playlist

Paste a playlist URL and it grabs every track:

```
Enter YouTube URL (or 'quit' to exit): https://www.youtube.com/playlist?list=PLxxxxxxx

  Found 12 tracks in playlist

  [1/12] Song One
  [2/12] Song Two
  ...
All 12 tracks downloaded to: /Users/you/Downloads
```

### Where Do My MP3s Go?

Your **Downloads** folder. After every download it prints the full path so you know exactly where to find it.

---

## Troubleshooting

### "command not found: konversion"

This means the install didn't finish. Run these commands:

```bash
cd ~/Konversion
```

```bash
go build -o konversion .
```

```bash
sudo cp konversion /usr/local/bin/konversion
```

Then try `konversion` again. Make sure you're typing it **all lowercase**.

### "destination path already exists"

You already have a ~/Konversion folder. Delete it and try again:

```bash
rm -rf ~/Konversion
git clone https://github.com/Stvn444/Konversion.git ~/Konversion
```

### "go: command not found"

Go isn't installed. Install it:

```bash
brew install go             # macOS
sudo apt install golang     # Ubuntu/Debian
```

Or download from [go.dev/dl](https://go.dev/dl/)

### "yt-dlp is not installed"

```bash
brew install yt-dlp         # macOS
pip install yt-dlp           # any system with Python
sudo apt install yt-dlp      # Ubuntu/Debian
```

### "ffmpeg is not installed"

```bash
brew install ffmpeg          # macOS
sudo apt install ffmpeg      # Ubuntu/Debian
```

### "Invalid YouTube URL format"

Make sure you're pasting the full URL. These formats work:

- `https://www.youtube.com/watch?v=VIDEO_ID`
- `https://youtu.be/VIDEO_ID`
- `https://www.youtube.com/embed/VIDEO_ID`
- `https://www.youtube.com/playlist?list=PLAYLIST_ID`

### How do I update Konversion?

```bash
cd ~/Konversion
```

```bash
git pull
```

```bash
go build -o konversion .
```

```bash
sudo cp konversion /usr/local/bin/konversion
```

---

## How It Works

Konversion is a Go wrapper around two industry-standard open source tools:

- **[yt-dlp](https://github.com/yt-dlp/yt-dlp)** — downloads video/audio from YouTube
- **[ffmpeg](https://ffmpeg.org/)** — converts audio to MP3 format

Konversion handles the terminal UI, progress bar, metadata display, and makes the whole process as simple as paste-and-go. No Python scripts, no browser extensions, no sketchy websites with popups.

## Disclaimer

Konversion is provided for **personal and educational use only**. This tool is intended for downloading content that you have the right to download, such as:

- Videos you own or have created
- Content licensed under Creative Commons or similar open licenses
- Content where the creator has explicitly allowed downloading
- Fair use purposes (commentary, education, research)

**You are solely responsible for how you use this tool.** Downloading copyrighted material without permission may violate YouTube's Terms of Service and copyright laws in your country. The author(s) of Konversion are not responsible for any misuse of this software or any violations of third-party rights, terms of service, or applicable laws.

By using Konversion, you agree that you will comply with all applicable laws and respect the rights of content creators.

## License

MIT — do whatever you want with it.
