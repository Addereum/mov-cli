#!/usr/bin/env bash
set -e

echo -e "\033[1;35m=========================================\033[0m"
echo -e "\033[1;35m  mov-cli (Addereum Fork) Installer      \033[0m"
echo -e "\033[1;35m=========================================\033[0m"

REPO_URL="git+https://github.com/Addereum/mov-cli.git"

if command -v pipx &> /dev/null; then
    echo -e "\033[1;32m[✓] pipx found! Installing via pipx...\033[0m"
    pipx install "$REPO_URL" --force
    echo -e "\033[1;32m[✓] Installation complete! You can now run 'mov-cli' from anywhere.\033[0m"
elif command -v pip &> /dev/null; then
    echo -e "\033[1;33m[!] pipx not found. Falling back to pip...\033[0m"
    pip install --user --upgrade "$REPO_URL"
    echo -e "\033[1;32m[✓] Installation complete! Make sure ~/.local/bin is in your PATH.\033[0m"
else
    echo -e "\033[1;31m[x] Error: Neither pipx nor pip is installed. Please install Python first.\033[0m"
    exit 1
fi
