#!/bin/bash

# Function to build the application for Linux or macOS
build_for_os() {
    case "$1" in
        linux)
            echo "Building for Linux..."
            GOOS=linux GOARCH=amd64 go build -o dockermi main.go
            echo "Build completed! The executable is dockermi-linux."
            ;;
        darwin)
            echo "Building for macOS..."
            GOOS=darwin GOARCH=amd64 go build -o dockermi main.go
            echo "Build completed! The executable is dockermi-macos."
            ;;
        *)
            echo "Unsupported operating system: $1"
            exit 1
            ;;
    esac
}

# Detect the operating system
OS="$(uname -s)"
case "$OS" in
    Linux)
        build_for_os linux
        ;;
    Darwin)
        build_for_os darwin
        ;;
    *)
        echo "This script only supports Linux and macOS."
        exit 1
        ;;
esac
