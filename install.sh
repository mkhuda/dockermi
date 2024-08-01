#!/bin/bash

# Function to build the application for Linux or macOS
build_app() {
    echo "Building for $1..."
    GOOS="$1" GOARCH=amd64 go build -o dockermi main.go
    echo "Build completed! The executable is dockermi."
}

# Function to install the application
install_app() {
    echo "Installing dockermi on $1..."
    cp dockermi /usr/local/bin/dockermi
    chmod +x /usr/local/bin/dockermi
    echo "dockermi installed! You can run it by typing 'dockermi' in your terminal."
}

# Detect the operating system and perform actions
OS="$(uname -s)"
case "$OS" in
    Linux)
        build_app linux
        install_app linux
        ;;
    Darwin)
        build_app darwin
        install_app darwin
        ;;
    *)
        echo "This script only supports Linux and macOS."
        exit 1
        ;;
esac