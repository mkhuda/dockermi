// Package dockermi provides the core functionality to generate a dockermi.sh script
// to manage Docker services defined in docker-compose.yml files.
package dockermi

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mkhuda/dockermi/internal/dockercompose"
	"github.com/mkhuda/dockermi/internal/script"
	dockermiUtils "github.com/mkhuda/dockermi/utils" // Import the utils package

	"github.com/fatih/color"
)

const version = "0.0.9"

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

	return generateScripts(projectDir)
}

// generateScripts finds docker-compose.yml files and generates corresponding scripts.
func generateScripts(projectDir string) (string, error) {
	services, foundDockerCompose := dockercompose.FindServices(projectDir)

	if !foundDockerCompose {
		color.Yellow("No docker-compose.yml found within this folder")
		return "No docker-compose.yml found within this folder", nil
	}

	// Create the dockermi.sh script
	scriptPath := filepath.Join(projectDir, "dockermi.sh")
	if err := script.CreateDockermiScript(scriptPath, services); err != nil {
		color.Red("Error creating dockermi.sh file: %v", err)
		return "", err
	}

	color.Green("Generated script: %s", scriptPath)
	color.Blue("You can now run [./dockermi.sh up] or [./dockermi.sh down]")

	return scriptPath, nil
}
