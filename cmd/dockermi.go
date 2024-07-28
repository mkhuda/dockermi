package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"dockermi/internal/dockercompose"
	"dockermi/internal/script"

	"github.com/fatih/color"
)

const version = "0.0.1"

func displayHelp() {
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

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("Dockermi version:", version)
		os.Exit(0)
	}

	help := flag.Bool("help", false, "Display help information")
	flag.Parse()

	if *help {
		displayHelp()
		return
	}

	projectDir, err := os.Getwd()
	if err != nil {
		color.Red("Error getting current directory: %v", err)
		return
	}

	// Find docker-compose.yml files
	services, foundDockerCompose := dockercompose.FindServices(projectDir)

	if !foundDockerCompose {
		color.Yellow("No docker-compose.yml found within this folder")
		return
	}

	// Create the dockermi.sh script
	scriptPath := filepath.Join(projectDir, "dockermi.sh")
	if err := script.CreateDockermiScript(scriptPath, services); err != nil {
		color.Red("Error creating dockermi.sh file: %v", err)
		return
	}

	color.Green("Generated script: %s", scriptPath)
	color.Blue("You can now run ./dockermi.sh up or ./dockermi.sh down")
}
