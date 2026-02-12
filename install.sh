#!/bin/bash
set -e

REPO="Stvn444/Konversion"
INSTALL_DIR="/usr/local/bin"
BINARY_NAME="konversion"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
CYAN='\033[0;36m'
BOLD='\033[1m'
NC='\033[0m'

info() { echo -e "${CYAN}$1${NC}"; }
success() { echo -e "${GREEN}$1${NC}"; }
warn() { echo -e "${YELLOW}$1${NC}"; }
error() { echo -e "${RED}$1${NC}"; exit 1; }

# Detect OS
OS="$(uname -s)"
case "$OS" in
    Darwin) OS="darwin" ;;
    Linux)  OS="linux" ;;
    *)      error "Unsupported OS: $OS. Use macOS or Linux." ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    amd64)   ARCH="amd64" ;;
    arm64)   ARCH="arm64" ;;
    aarch64) ARCH="arm64" ;;
    *)       error "Unsupported architecture: $ARCH" ;;
esac

FILENAME="konversion-${OS}-${ARCH}"
DOWNLOAD_URL="https://github.com/${REPO}/releases/latest/download/${FILENAME}"

echo ""
echo -e "${CYAN}${BOLD}Installing Konversion...${NC}"
echo -e "  OS:   ${BOLD}${OS}${NC}"
echo -e "  Arch: ${BOLD}${ARCH}${NC}"
echo ""

# Download binary
info "Downloading ${FILENAME}..."
if command -v curl &> /dev/null; then
    curl -fSL "$DOWNLOAD_URL" -o "/tmp/${BINARY_NAME}" || error "Download failed. Check your internet connection and try again."
elif command -v wget &> /dev/null; then
    wget -q "$DOWNLOAD_URL" -O "/tmp/${BINARY_NAME}" || error "Download failed. Check your internet connection and try again."
else
    error "curl or wget is required to download Konversion."
fi

chmod +x "/tmp/${BINARY_NAME}"

# Install binary
info "Installing to ${INSTALL_DIR}..."
if [ ! -d "$INSTALL_DIR" ]; then
    sudo mkdir -p "$INSTALL_DIR"
fi
if [ -w "$INSTALL_DIR" ]; then
    mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    sudo mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

success "Konversion installed successfully!"
echo ""

# Check and install dependencies
MISSING_YTDLP=false
MISSING_FFMPEG=false

if ! command -v yt-dlp &> /dev/null; then
    MISSING_YTDLP=true
fi

if ! command -v ffmpeg &> /dev/null; then
    MISSING_FFMPEG=true
fi

if [ "$MISSING_YTDLP" = true ] || [ "$MISSING_FFMPEG" = true ]; then
    info "Installing missing dependencies..."
    echo ""

    if [ "$OS" = "darwin" ]; then
        # macOS — use Homebrew
        if command -v brew &> /dev/null; then
            BREW_PKGS=""
            [ "$MISSING_YTDLP" = true ] && BREW_PKGS="${BREW_PKGS} yt-dlp"
            [ "$MISSING_FFMPEG" = true ] && BREW_PKGS="${BREW_PKGS} ffmpeg"
            echo -e "  Running: brew install${BREW_PKGS}"
            brew install $BREW_PKGS
        else
            warn "Homebrew not found. Install these manually:"
            [ "$MISSING_YTDLP" = true ] && echo "  brew install yt-dlp"
            [ "$MISSING_FFMPEG" = true ] && echo "  brew install ffmpeg"
            echo ""
            echo "  Install Homebrew first:"
            echo '  /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"'
        fi
    else
        # Linux — use apt + pip/pipx
        if [ "$MISSING_FFMPEG" = true ]; then
            echo -e "  Installing ffmpeg..."
            sudo apt-get update -qq
            sudo apt-get install -y -qq ffmpeg
        fi

        if [ "$MISSING_YTDLP" = true ]; then
            echo -e "  Installing yt-dlp..."
            if command -v pipx &> /dev/null; then
                pipx install yt-dlp
            elif command -v pip3 &> /dev/null; then
                pip3 install --user yt-dlp 2>/dev/null || pip3 install --user --break-system-packages yt-dlp 2>/dev/null || {
                    # If pip fails, try installing pipx
                    sudo apt-get install -y -qq pipx 2>/dev/null && pipx install yt-dlp || {
                        warn "Could not auto-install yt-dlp. Install it manually:"
                        echo "  sudo apt install pipx && pipx install yt-dlp"
                    }
                }
            else
                # No pip3, install it
                sudo apt-get install -y -qq python3-pip 2>/dev/null
                if command -v pip3 &> /dev/null; then
                    pip3 install --user yt-dlp 2>/dev/null || pip3 install --user --break-system-packages yt-dlp
                else
                    warn "Could not auto-install yt-dlp. Install it manually:"
                    echo "  sudo apt install pipx && pipx install yt-dlp"
                fi
            fi

            # Check if yt-dlp ended up in ~/.local/bin and warn if not in PATH
            if ! command -v yt-dlp &> /dev/null; then
                if [ -f "$HOME/.local/bin/yt-dlp" ]; then
                    warn "yt-dlp was installed to ~/.local/bin which is not in your PATH."
                    echo ""
                    echo -e "  ${BOLD}Run this, then restart your terminal:${NC}"
                    echo '  echo '\''export PATH="$HOME/.local/bin:$PATH"'\'' >> ~/.bashrc'
                    echo ""
                fi
            fi
        fi
    fi
    echo ""
fi

# Final check
ALL_GOOD=true
if ! command -v yt-dlp &> /dev/null && [ ! -f "$HOME/.local/bin/yt-dlp" ]; then
    warn "yt-dlp is still missing"
    ALL_GOOD=false
fi
if ! command -v ffmpeg &> /dev/null; then
    warn "ffmpeg is still missing"
    ALL_GOOD=false
fi

if [ "$ALL_GOOD" = true ]; then
    success "All dependencies found!"
fi

echo ""
echo -e "${BOLD}Run it:${NC}"
echo "  konversion"
echo ""
