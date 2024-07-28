package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mkhuda/dockermi/internal/dockercompose"
	"github.com/mkhuda/dockermi/internal/script"

	"github.com/fatih/color"
)

const version = "0.0.5"

// displayHelp prints the usage information for the dockermi command to the console.
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

// RunDockermi executes the main logic of the dockermi command.
func RunDockermi() error {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("Dockermi version:", version)
		os.Exit(0)
	}

	help := flag.Bool("help", false, "Display help information")
	flag.Parse()

	if *help {
		displayHelp()
		return nil
	}

	projectDir, err := os.Getwd()
	if err != nil {
		color.Red("Error getting current directory: %v", err)
		return err
	}

	// Find docker-compose.yml files
	services, foundDockerCompose := dockercompose.FindServices(projectDir)

	if !foundDockerCompose {
		color.Yellow("No docker-compose.yml found within this folder")
		return err
	}

	// Create the dockermi.sh script
	scriptPath := filepath.Join(projectDir, "dockermi.sh")
	if err := script.CreateDockermiScript(scriptPath, services); err != nil {
		color.Red("Error creating dockermi.sh file: %v", err)
		return err
	}

	color.Green("Generated script: %s", scriptPath)
	color.Blue("You can now run ./dockermi.sh up or ./dockermi.sh down")

	return nil

}

// main is the entry point of the application.
func main() {
	// Execute the RunDockermi function and handle any errors
	if err := RunDockermi(); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}
