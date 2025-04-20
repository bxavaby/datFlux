#!/usr/bin/env bash

# datFlux installer script

set -e

BLUE='\033[1;34m'
GREEN='\033[1;32m'
YELLOW='\033[1;33m'
RED='\033[1;31m'
CYAN='\033[1;36m'
NC='\033[0m'

echo -e "${BLUE}"
echo "┌─────────────────────────────────────────┐"
echo "│                                         │"
echo "│           datFlux Installer             │"
echo "│     Entropy-Borne Password Generator    │"
echo "│                                         │"
echo "└─────────────────────────────────────────┘"
echo -e "${NC}"

echo -e "${CYAN}Checking Go installation...${NC}"
if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in your PATH.${NC}"
    echo -e "${YELLOW}Please install Go from https://golang.org/dl/${NC}"
    exit 1
fi

GO_VERSION=$(go version | awk '{print $3}')
echo -e "${GREEN}✓ Go is installed: ${GO_VERSION}${NC}"

DEFAULT_INSTALL_DIR="/usr/local/bin"
if [[ "$OSTYPE" == "darwin"* ]]; then
    DEFAULT_INSTALL_DIR="/usr/local/bin"
elif [[ -d "$HOME/.local/bin" && "$PATH" == *"$HOME/.local/bin"* ]]; then
    DEFAULT_INSTALL_DIR="$HOME/.local/bin"
fi

echo -e "${CYAN}Where would you like to install datFlux?${NC}"
echo -e "${YELLOW}Press Enter for default: ${DEFAULT_INSTALL_DIR}${NC}"
read -p "Installation path: " INSTALL_DIR
INSTALL_DIR=${INSTALL_DIR:-$DEFAULT_INSTALL_DIR}

if [[ ! -d "$INSTALL_DIR" ]]; then
    echo -e "${YELLOW}Directory ${INSTALL_DIR} does not exist. Create it?${NC} (y/n)"
    read -r CREATE_DIR
    if [[ "$CREATE_DIR" =~ ^[Yy]$ ]]; then
        echo -e "${CYAN}Creating directory ${INSTALL_DIR}...${NC}"
        mkdir -p "$INSTALL_DIR" || { echo -e "${RED}Failed to create directory. Try running with sudo?${NC}"; exit 1; }
    else
        echo -e "${RED}Installation cancelled.${NC}"
        exit 1
    fi
fi

if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${YELLOW}Warning: ${INSTALL_DIR} is not in your PATH.${NC}"
    echo -e "You might want to add it with: ${CYAN}export PATH=\"\$PATH:${INSTALL_DIR}\"${NC}"
fi

# build datFlux
echo -e "${CYAN}Building datFlux...${NC}"
go build -o datflux ./cmd/datflux || { echo -e "${RED}Build failed.${NC}"; exit 1; }

CURRENT_DIR=$(pwd)

echo -e "${CYAN}Installing to ${INSTALL_DIR}/datflux...${NC}"
if cp "$CURRENT_DIR/datflux" "$INSTALL_DIR/"; then
    echo -e "${GREEN}✓ Installation successful!${NC}"
else
    echo -e "${RED}Failed to copy binary. Try running with sudo?${NC}"
    echo -e "${YELLOW}Try: ${CYAN}sudo cp \"$CURRENT_DIR/datflux\" \"$INSTALL_DIR/\"${NC}"
    exit 1
fi

chmod +x "$INSTALL_DIR/datflux" || { echo -e "${RED}Failed to make binary executable.${NC}"; exit 1; }

echo -e "${GREEN}"
echo "┌─────────────────────────────────────────┐"
echo "│                                         │"
echo "│     datFlux installed successfully!     │"
echo "│                                         │"
echo "└─────────────────────────────────────────┘"
echo -e "${NC}"
echo -e "Run ${CYAN}datflux${NC} to start the application."

# is installation directory in PATH?
if [[ ":$PATH:" != *":$INSTALL_DIR:"* ]]; then
    echo -e "${YELLOW}Note: ${INSTALL_DIR} is not in your PATH.${NC}"

    SHELL_NAME=$(basename "$SHELL")
    case $SHELL_NAME in
        bash)
            echo -e "Add to your PATH by running:"
            echo -e "${CYAN}echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.bashrc${NC}"
            echo -e "${CYAN}source ~/.bashrc${NC}"
            ;;
        zsh)
            echo -e "Add to your PATH by running:"
            echo -e "${CYAN}echo 'export PATH=\"\$PATH:$INSTALL_DIR\"' >> ~/.zshrc${NC}"
            echo -e "${CYAN}source ~/.zshrc${NC}"
            ;;
        *)
            echo -e "Add ${INSTALL_DIR} to your PATH to run datflux from anywhere."
            ;;
    esac
fi

echo -e "\n${GREEN}Thank you for installing datFlux!${NC}"
exit 0
