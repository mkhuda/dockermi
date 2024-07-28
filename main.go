package main

import (
	"os"

	"github.com/fatih/color"
	dockermi "github.com/mkhuda/dockermi/pkg" // Import your dockermi package
)

func main() {
	projectDir, err := os.Getwd()
	if err != nil {
		color.Red("Error getting current directory: %v", err)
		os.Exit(1)
	}
	// Execute the RunDockermi function and handle any errors
	if _, err := dockermi.RunDockermi(projectDir); err != nil {
		color.Red("Error: %v", err)
		os.Exit(1)
	}
}
