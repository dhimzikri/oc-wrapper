#!/bin/bash
# OpenCode Wrapper Installer for Linux
# Usage: curl -sSL https://raw.githubusercontent.com/dhimzikri/oc-wrapper/master/install.sh | bash

set -e

INSTALL_DIR="$HOME/.config/opencode/wrapper"
REPO_URL="https://github.com/dhimzikri/oc-wrapper"

echo -e "\033[36mInstalling OpenCode Wrapper...\033[0m"

# Create install directory
mkdir -p "$INSTALL_DIR"

# Check if Go is installed
if command -v go &> /dev/null; then
    echo -e "\033[33mBuilding from source...\033[0m"
    
    TEMP_DIR=$(mktemp -d)
    trap "rm -rf $TEMP_DIR" EXIT
    
    git clone --depth 1 "$REPO_URL" "$TEMP_DIR"
    cd "$TEMP_DIR"
    GOOS=linux GOARCH=amd64 go build -o "$INSTALL_DIR/opencode" main_linux.go
else
    echo -e "\033[33mGo not found. Downloading pre-built binary...\033[0m"
    BINARY_URL="$REPO_URL/releases/latest/download/opencode-wrapper-linux-amd64"
    curl -sSL "$BINARY_URL" -o "$INSTALL_DIR/opencode"
fi

# Make executable
chmod +x "$INSTALL_DIR/opencode"

# Shell detection and alias setup
ALIAS_LINE="alias opencode='$INSTALL_DIR/opencode'"
MARKER="# OpenCode Wrapper"

add_to_shell_config() {
    local config_file="$1"
    
    if [ -f "$config_file" ]; then
        if ! grep -q "opencode-wrapper\|oc-wrapper" "$config_file" 2>/dev/null; then
            echo -e "\033[33mAdding to $config_file...\033[0m"
            echo "" >> "$config_file"
            echo "$MARKER" >> "$config_file"
            echo "$ALIAS_LINE" >> "$config_file"
        else
            echo -e "\033[33m$config_file already configured, skipping...\033[0m"
        fi
    fi
}

# Add to bash
if [ -f "$HOME/.bashrc" ]; then
    add_to_shell_config "$HOME/.bashrc"
elif [ -f "$HOME/.bash_profile" ]; then
    add_to_shell_config "$HOME/.bash_profile"
fi

# Add to zsh
if [ -f "$HOME/.zshrc" ]; then
    add_to_shell_config "$HOME/.zshrc"
fi

echo ""
echo -e "\033[32mInstallation complete!\033[0m"
echo -e "\033[36mRestart your terminal or run:\033[0m"
echo "  source ~/.bashrc   # for bash"
echo "  source ~/.zshrc    # for zsh"
echo ""
echo -e "\033[36mThen use: opencode\033[0m"
