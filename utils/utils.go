package utils

import "fmt"

// DisplayHelp prints the usage information for the dockermi command to the console.
func DisplayHelp() {
	fmt.Println(`Usage: dockermi [--help]

This command generates a dockermi.sh script to manage Docker services defined in docker-compose.yml files.

Options:
    --help  Display this help message and exit.
    --version Display current installed version

Examples:
    dockermi
    ./dockermi.sh up
    ./dockermi.sh down`)
}
