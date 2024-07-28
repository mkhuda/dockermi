#!/bin/bash

# Function to install the application
install_app() {
    case "$1" in
        linux)
            echo "Installing dockermi on Linux..."
            cp dockermi /usr/local/bin/dockermi
            chmod +x /usr/local/bin/dockermi
            echo "dockermi installed! You can run it by typing 'dockermi' in your terminal."
            ;;
        darwin)
            echo "Installing dockermi on macOS..."
            cp dockermi /usr/local/bin/dockermi
            chmod +x /usr/local/bin/dockermi
            echo "dockermi installed! You can run it by typing 'dockermi' in your terminal."
            ;;
        *)
            echo "Unsupported operating system for installation: $1"
            exit 1
            ;;
    esac
}

# Detect the operating system
OS="$(uname -s)"
case "$OS" in
    Linux)
        install_app linux
        ;;
    Darwin)
        install_app darwin
        ;;
    *)
        echo "This script only supports Linux and macOS."
        exit 1
        ;;
esac
