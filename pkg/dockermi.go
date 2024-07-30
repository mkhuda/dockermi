// Package dockermi provides the core functionality to generate a dockermi.sh script
// to manage Docker services defined in docker-compose.yml files.
package dockermi

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/mkhuda/dockermi/internal/dockercompose"
	"github.com/mkhuda/dockermi/internal/script"
	dockermiUtils "github.com/mkhuda/dockermi/utils" // Import the utils package

	"github.com/fatih/color"
)

const version = "0.1.0"

// RunDockermi executes the main logic of the dockermi command. It takes a
// projectDir parameter, which specifies the directory where the function
// will look for docker-compose.yml files and create the dockermi.sh script.
//
// Parameters:
//   - projectDir: the directory where docker-compose.yml files may located
//
// Returns:
//   - string: Path location of created dockermi.sh
//   - error: if any errors occur during the execution, they are returned
func RunDockermi(projectDir string) (string, error) {
	if len(os.Args) > 1 && os.Args[1] == "--version" {
		fmt.Println("Dockermi version:", version)
		os.Exit(0)
	}

	help := flag.Bool("help", false, "Display help information")
	flag.Parse()

	if *help {
		dockermiUtils.DisplayHelp()
		return "", nil
	}

	// Check if the command is provided
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "up":
			return handleUpDownCommand(projectDir, "up", os.Args[2:])
		case "down":
			return handleUpDownCommand(projectDir, "down", os.Args[2:]) // Handle the down command
		case "create":
			if len(os.Args) < 3 {
				return "", fmt.Errorf("missing key for create command")
			}
			return createDockermiScript(projectDir, os.Args[2])
		default:
			return generateScripts(projectDir)
		}
	}
	// If no specific command is provided, generate the scripts
	return generateScripts(projectDir)
}

// handleUpDownCommand handles the 'up' command logic.
func handleUpDownCommand(projectDir string, command string, args []string) (string, error) {
	color.Green("Executing 'up' command...")

	return runDockermiScript(projectDir, command, args)
}

// generateScripts finds docker-compose.yml files and generates corresponding scripts.
func generateScripts(projectDir string) (string, error) {
	services, err := dockercompose.FindServices(projectDir)
	servicesLength := len(services)

	if servicesLength == 0 {
		color.Yellow("No docker-compose.yml found within this folder")
		return "No docker-compose.yml found within this folder", nil
	}

	if err != nil {
		return "", err
	}

	// Create the dockermi.sh script
	scriptPath := filepath.Join(projectDir, "dockermi.sh")
	if err := script.CreateDockermiScript(scriptPath, services); err != nil {
		color.Red("Error creating dockermi.sh file: %v", err)
		return "", err
	}

	fmt.Println()
	color.Green("Generated script: %s", scriptPath)
	fmt.Println()
	color.Blue("You can now run [dockermi up] or [dockermi down]")

	return scriptPath, nil
}

// runDockermiScript executes the dockermi.sh script located in the current directory
// with the specified subcommand (e.g., "up" or "down").
func runDockermiScript(currentDir, subcommand string, options []string) (string, error) {
	scriptPath := filepath.Join(currentDir, "dockermi.sh")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return "", fmt.Errorf("dockermi.sh script not found in current directory")
	}

	color.Green("Running script: %s with subcommand: %s and options: %v", scriptPath, subcommand, options)

	// Prepare command with subcommand and options
	cmd := exec.Command("bash", append([]string{scriptPath, subcommand}, options...)...) // Pass the subcommand and options to the script
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run dockermi.sh: %w", err)
	}

	return scriptPath, nil
}

// createDockermiScript creates a dockermi-{key}.sh script in the user's home directory.
func createDockermiScript(projectDir string, key string) (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	// Ensure the .dockermi directory exists
	dockermiDir := filepath.Join(homeDir, ".dockermi")
	err = os.MkdirAll(dockermiDir, os.ModePerm)
	if err != nil {
		return "", err
	}

	// Use FindServicesWithKey to get the services associated with the key
	services, err := dockercompose.FindServicesWithKey(projectDir)
	if err != nil {
		return "", err
	}

	color.Red("find %v", services)

	// Check if the key exists in the services map
	groupedServices, exists := services[key]
	if !exists {
		return "", fmt.Errorf("no services found for key: %s", key)
	}

	// Create the dockermi-{key}.sh script
	scriptPath := filepath.Join(dockermiDir, fmt.Sprintf("dockermi-%s.sh", key))
	if err := script.CreateDockermiScript(scriptPath, groupedServices); err != nil {
		color.Red("Error creating dockermi-%s.sh file: %v", key, err)
		return "", err
	}

	color.Green("Generated script: %s", scriptPath)
	return scriptPath, nil
}
