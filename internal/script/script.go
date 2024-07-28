package script

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/fatih/color"
	"github.com/schollz/progressbar/v3"
)

// CreateDockermiScript generates the dockermi.sh script based on the provided services.
func CreateDockermiScript(scriptPath string, services []struct {
	Order       string
	ServiceName string
	ComposeFile string
}) error {
	dockermiScript, err := os.Create(scriptPath)
	if err != nil {
		return err
	}
	defer dockermiScript.Close()

	dockermiScript.WriteString("#!/bin/bash\n\n")
	dockermiScript.WriteString("# Usage: ./dockermi.sh [up|down]\n\n")

	// Sort services for starting (ascending order)
	sort.Slice(services, func(i, j int) bool {
		return services[i].Order < services[j].Order
	})

	// Generate start_services function
	dockermiScript.WriteString("start_services() {\n")

	// Create a progress bar
	bar := progressbar.New(len(services))

	for _, service := range services {
		dockermiScript.WriteString(fmt.Sprintf("    echo \"Starting %s...\"\n", service.ServiceName))
		dockermiScript.WriteString(fmt.Sprintf("    docker-compose -f \"%s\" up -d\n", service.ComposeFile))

		bar.Add(1)
		time.Sleep(500 * time.Millisecond) // Simulate delay for demonstration
	}
	dockermiScript.WriteString("}\n\n")

	// Generate stop_services function (descending order)
	dockermiScript.WriteString("stop_services() {\n")
	sort.Slice(services, func(i, j int) bool {
		return services[i].Order > services[j].Order
	})
	for _, service := range services {
		dockermiScript.WriteString(fmt.Sprintf("    echo \"Stopping %s...\"\n", service.ServiceName))
		dockermiScript.WriteString(fmt.Sprintf("    docker-compose -f \"%s\" down\n", service.ComposeFile))

		bar.Add(1)
		time.Sleep(500 * time.Millisecond) // Simulate delay for demonstration
	}
	dockermiScript.WriteString("}\n\n")

	// Add main logic to call the appropriate function based on the argument
	dockermiScript.WriteString(`if [ "$#" -ne 1 ]; then
    echo "Invalid argument!"
    echo "Usage: $0 [up|down]"
    exit 1
fi

ACTION=$1

case "$ACTION" in
    up)
        start_services
        ;;
    down)
        stop_services
        ;;
    *)
        echo "Invalid argument: $ACTION"
        echo "Usage: $0 [up|down]"
        exit 1
        ;;
esac
`)

	// Make the dockermi.sh script executable (Unix systems)
	if err := os.Chmod(scriptPath, 0755); err != nil {
		color.Red("Error making the script executable: %v", err)
		return err
	}

	return nil
}
