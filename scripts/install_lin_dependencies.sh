#!/bin/bash

if command -v apt > /dev/null 2>&1; then
    PACKAGE_MANAGER="apt"
elif command -v dnf > /dev/null 2>&1; then
    PACKAGE_MANAGER="dnf"
elif command -v pacman > /dev/null 2>&1; then
    PACKAGE_MANAGER="pacman"
else
    echo "No supported package manager found. Exiting..."
    exit 1
fi

if [[ "$PACKAGE_MANAGER" == "apt" ]]; then
    sudo apt update
    sudo apt install -y librsvg2-bin python3 python3-pip
elif [[ "$PACKAGE_MANAGER" == "dnf" ]]; then
    sudo dnf install -y librsvg2 python3 python3-pip
elif [[ "$PACKAGE_MANAGER" == "pacman" ]]; then
    sudo pacman -Syu --noconfirm librsvg python python-pip
fi

python3 -m pip install --upgrade pip
python3 -m pip install svg.path

read -p "Do you want to install Go? (Y/N): " user_input
if [[ "$user_input" =~ ^[Yy]$ ]]; then
    if [[ "$PACKAGE_MANAGER" == "apt" ]]; then
        sudo apt install -y golang
    elif [[ "$PACKAGE_MANAGER" == "dnf" ]]; then
        sudo dnf install -y golang
    elif [[ "$PACKAGE_MANAGER" == "pacman" ]]; then
        sudo pacman -S --noconfirm go
    fi
    echo "Go installed successfully!"
else
    echo "Skipping Go installation."
fi

echo "All installations completed successfully!"
