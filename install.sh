#!/usr/bin/env bash

# datFlux installer script

set -e

# needed for the initial gum installation
BLUE='\033[1;34m'
RED='\033[1;31m'
YELLOW='\033[1;33m'
NC='\033[0m'

clear

command_exists() {
    command -v "$1" &> /dev/null
}

# install gum based on OS
install_gum() {
    echo -e "${BLUE}Installing gum for the datFlux installer UI...${NC}"

    if [[ "$OSTYPE" == "darwin"* ]]; then
        # macOS
        if command_exists brew; then
            brew install gum
        else
            echo -e "${YELLOW}Homebrew not found. Please install Homebrew first:${NC}"
            echo -e "https://brew.sh"
            exit 1
        fi
    elif command_exists apt-get; then
        # Debian/Ubuntu
        sudo mkdir -p /etc/apt/keyrings
        curl -fsSL https://repo.charm.sh/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/keyrings/charm.gpg
        echo "deb [signed-by=/etc/apt/keyrings/charm.gpg] https://repo.charm.sh/apt/ * *" | sudo tee /etc/apt/sources.list.d/charm.list
        sudo apt update && sudo apt install -y gum
    elif command_exists yum; then
        # RHEL/Fedora
        echo '[charm]
name=Charm
baseurl=https://repo.charm.sh/yum/
enabled=1
gpgcheck=1
gpgkey=https://repo.charm.sh/yum/gpg.key' | sudo tee /etc/yum.repos.d/charm.repo
        sudo rpm --import https://repo.charm.sh/yum/gpg.key
        sudo yum install -y gum
    elif command_exists zypper; then
        # OpenSUSE
        echo '[charm]
name=Charm
baseurl=https://repo.charm.sh/yum/
enabled=1
gpgcheck=1
gpgkey=https://repo.charm.sh/yum/gpg.key' | sudo tee /etc/yum.repos.d/charm.repo
        sudo rpm --import https://repo.charm.sh/yum/gpg.key
        sudo zypper refresh
        sudo zypper install -y gum
    elif command_exists pacman; then
        # Arch
        sudo pacman -S gum
    elif command_exists pkg; then
        # FreeBSD
        sudo pkg install gum
    else
        echo -e "${RED}Cannot detect package manager. Please install gum manually:${NC}"
        echo -e "https://github.com/charmbracelet/gum#installation"
        exit 1
    fi

    # gum verification
    if ! command_exists gum; then
        echo -e "${RED}Failed to install gum. Please install it manually and run this script again.${NC}"
        exit 1
    fi
}

# install gum, if not already
if ! command_exists gum; then
    echo -e "${YELLOW}This installer requires gum for its UI. Would you like to install it? (y/n)${NC}"
    read -r INSTALL_GUM
    if [[ "$INSTALL_GUM" =~ ^[Yy]$ ]]; then
        install_gum
        clear
    else
        echo -e "${RED}Gum is required for this installer. Exiting.${NC}"
        exit 1
    fi
fi

# now that gum is installed
# Tokyo Night to match datFlux (approximate matches, sadly)
TN_BG="60"        # background
TN_BLUE="75"      # info
TN_GREEN="114"    # success
TN_PURPLE="141"   # accent
TN_RED="204"      # error
TN_YELLOW="222"   # warning

gum style --foreground "$TN_PURPLE" --border-foreground "$TN_PURPLE" --border double \
    --align center --width 50 --margin "1 0" --padding "1 2" \
    "datFlux Installer" "Entropy-Borne Password Generator"

gum spin --spinner dot --title "Checking Go installation..." -- sleep 0.5

if ! command_exists go; then
    gum style --foreground "$TN_RED" "Error: Go is not installed or not in your PATH."
    gum style --foreground "$TN_BLUE" "Please install Go from https://golang.org/dl/"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
gum style --foreground "$TN_GREEN" "✓ Go is installed: $GO_VERSION"

# installation directory
DEFAULT_INSTALL_DIR="/usr/local/bin"
if [[ "$OSTYPE" == "darwin"* ]]; then
    DEFAULT_INSTALL_DIR="/usr/local/bin"
elif [[ -d "$HOME/.local/bin" && "$PATH" == *"$HOME/.local/bin"* ]]; then
    DEFAULT_INSTALL_DIR="$HOME/.local/bin"
fi

# request installation directory
echo ""
gum style --foreground "$TN_BLUE" "Where would you like to install datFlux?"
gum style --foreground "$TN_BLUE" "Press Enter for default: $DEFAULT_INSTALL_DIR"
echo ""
INSTALL_DIR=$(gum input --placeholder "$DEFAULT_INSTALL_DIR")
INSTALL_DIR=${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}

# create directory
if [[ ! -d "$INSTALL_DIR" ]]; then
    gum style --foreground "$TN_YELLOW" "$INSTALL_DIR does not exist."
    CREATE_DIR=$(gum confirm "Create it?" && echo "yes" || echo "no")
    if [[ "$CREATE_DIR" == "yes" ]]; then
        gum spin --spinner dot --title "Creating $INSTALL_DIR..." -- mkdir -p "$INSTALL_DIR" || {
            gum style --foreground "$TN_RED" "Failed to create directory. Try running with sudo?";
            exit 1;
        }
    else
        gum style --foreground "$TN_RED" "Installation cancelled."
        exit 1
    fi
fi

# installation directory not in PATH warning !!!
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    gum style --foreground "$TN_YELLOW" "Warning: $INSTALL_DIR is not in your PATH."
    gum style --foreground "$TN_BLUE" "You might want to add it with:"
    gum style --foreground "$TN_BLUE" "export PATH=\"\$PATH:$INSTALL_DIR\""
fi

# build datFlux
gum spin --spinner dot --title "Building datFlux..." -- go build -o datflux ./cmd/datflux || {
    gum style --foreground "$TN_RED" "Build failed.";
    exit 1;
}

CURRENT_DIR=$(pwd)

# install datFlux
gum spin --spinner dot --title "Installing to $INSTALL_DIR/datflux..." -- cp "$CURRENT_DIR/datflux" "$INSTALL_DIR/" || {
    gum style --foreground "$TN_RED" "Failed to copy binary. Try running with sudo?"
    gum style --foreground "$TN_BLUE" "Try: sudo cp \"$CURRENT_DIR/datflux\" \"$INSTALL_DIR/\""
    exit 1
}

# remove local binary after installing
gum spin --spinner dot --title "Cleaning up..." -- rm -f "$CURRENT_DIR/datflux"

chmod +x "$INSTALL_DIR/datflux" || {
    gum style --foreground "$TN_RED" "Failed to make binary executable."
    exit 1;
}

clear

# success
gum style --foreground "$TN_GREEN" --border-foreground "$TN_GREEN" --border double \
    --align center --width 50 --margin "1 0" --padding "1 2" \
    "datFlux installed successfully!"

# gum style --foreground "$TN_GREEN" "Thank you for installing!"

# is installation directory in PATH?
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    SHELL_NAME=$(basename "$SHELL")

    gum style --foreground "$TN_YELLOW" "Note: $INSTALL_DIR is not in your PATH."

    case $SHELL_NAME in
        bash)
            gum style "Add to your PATH by running:"
            gum style --foreground "$TN_PURPLE" "echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.bashrc"
            gum style --foreground "$TN_PURPLE" "source ~/.bashrc"
            ;;
        zsh)
            gum style "Add to your PATH by running:"
            gum style --foreground "$TN_PURPLE" "echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.zshrc"
            gum style --foreground "$TN_PURPLE" "source ~/.zshrc"
            ;;
        *)
            gum style "Add $INSTALL_DIR to your PATH to run datflux from anywhere."
            ;;
    esac
fi

echo ""
gum style " ⥤ Run $(gum style --foreground "$TN_PURPLE" "datflux") to start the application."

exit 0
