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
if [ -w "$INSTALL_DIR" ]; then
    mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    sudo mv "/tmp/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
fi

success "Konversion installed successfully!"
echo ""

# Check dependencies
MISSING=""

if ! command -v yt-dlp &> /dev/null; then
    MISSING="${MISSING}yt-dlp "
fi

if ! command -v ffmpeg &> /dev/null; then
    MISSING="${MISSING}ffmpeg "
fi

if [ -n "$MISSING" ]; then
    warn "Missing dependencies: ${MISSING}"
    echo ""
    echo -e "${BOLD}Install them:${NC}"

    if [ "$OS" = "darwin" ]; then
        echo "  brew install ${MISSING}"
    else
        echo "  sudo apt install ${MISSING}    (Debian/Ubuntu)"
        echo ""
        echo "  If yt-dlp isn't available via apt:"
        echo "  pip install yt-dlp"
    fi
    echo ""
else
    success "All dependencies found (yt-dlp, ffmpeg)"
    echo ""
fi

echo -e "${BOLD}Run it:${NC}"
echo "  konversion"
echo ""
