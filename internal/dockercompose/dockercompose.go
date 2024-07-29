package dockercompose

import (
	"os"
	"path/filepath"

	"github.com/fatih/color"
	DockermiTypes "github.com/mkhuda/dockermi/types"
	"gopkg.in/yaml.v2"
)

// FindServices searches for docker-compose.yml files in the specified directory.
// It scans the directory and its subdirectories for docker-compose.yml files,
// parses them to extract services with specific labels, and returns a list of
// these services along with a boolean indicating if any docker-compose.yml files were found.
//
// Parameters:
//   - root: the root directory to start the search from
//
// Returns:
//   - []struct{Order, ServiceName, ComposeFile string}: a slice of structs containing
//     the order, service name, and path to the docker-compose file for each relevant service
//   - bool: a boolean indicating if any docker-compose.yml files were found
func FindServices(root string) (DockermiTypes.ServiceScriptReturn, bool) {
	var services DockermiTypes.ServiceScriptReturn

	foundDockerCompose := false

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || filepath.Base(path) != "docker-compose.yml" {
			return nil
		}

		foundDockerCompose = true

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		var dockerCompose DockermiTypes.DockerCompose
		if err := yaml.Unmarshal(data, &dockerCompose); err != nil {
			return err
		}

		for serviceName, service := range dockerCompose.Services {
			order, orderExists := service.Labels["dockermi.order"]
			active, activeExists := service.Labels["dockermi.active"]

			if orderExists && activeExists && active == "true" {
				color.Green("Service: %s", serviceName)
				color.Blue("Order: %s", order)
				color.Yellow("Active: %s", active)

				services = append(services, DockermiTypes.ServiceScript{
					Order:       order,
					ServiceName: serviceName,
					ComposeFile: path,
				})
			} else if activeExists {
				color.Yellow("Service '%s' is inactive (dockermi.active=false). Skipping...", serviceName)
			} else {
				color.Red("Service '%s' is missing 'dockermi.order' or 'dockermi.active' labels. Skipping...", serviceName)
			}
		}
		return nil
	})

	if err != nil {
		color.Red("Error walking the path: %v", err)
	}

	return services, foundDockerCompose
}
