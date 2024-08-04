package dockercompose

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/fatih/color"
	DockermiTypes "github.com/mkhuda/dockermi/types"

	"gopkg.in/yaml.v2"
	// "github.com/goccy/go-yaml"
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
func FindServices(root string, force bool) (DockermiTypes.ServiceScriptReturn, error) {
	var services DockermiTypes.ServiceScriptReturn

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".yml") {
			return nil
		}

		composedFiles, err := ParseComposeFile(path, false, force)

		if err != nil {
			return err
		}

		for serviceName, service := range composedFiles {
			order, active := "", ""
			orderExists, activeExists := false, false

			// Access the labels directly
			if val, exists := service.Labels["dockermi.order"]; exists {
				order, orderExists = val, true
			}
			if val, exists := service.Labels["dockermi.active"]; exists {
				active, activeExists = val, true
			}

			// Determine if the service should be included
			var includeService bool

			if force {
				// If force is true, always include the service
				includeService = true
			} else {
				// Otherwise, check the order and active labels
				includeService = (orderExists && activeExists && active == "true")
			}

			if includeService {
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

	return services, err
}

// [Proposed Feature]
// FindServicesWithKey searches for docker-compose.yml files in the specified directory
// and groups the services by their 'dockermi.key' label. It parses each file to extract
// services that are active and have an associated order. If the 'dockermi.key' label is
// missing, a default key will be generated. The function returns a map where the keys are
// the values of 'dockermi.key' and the values are slices of ServiceScript structures
// containing the order, service name, and the compose file path. In case of an error during
// the file traversal or parsing, it returns the error encountered.
//
// Parameters:
//   - root: the root directory to start the search from
//
// Returns:
//   - map[string][]DockermiTypes.ServiceScript: a map of services grouped by their
//     'dockermi.key' labels
//   - error: if any errors occur during the execution, they are returned
func FindServicesWithKey(root string) (map[string][]DockermiTypes.ServiceScript, error) {
	groups := make(map[string][]DockermiTypes.ServiceScript)

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(info.Name(), ".yml") {
			return nil
		}

		services, err := ParseComposeFile(path, true, false)

		if err != nil {
			return err
		}

		for serviceName, service := range services {
			order, orderExists := service.Labels["dockermi.order"]
			active, activeExists := service.Labels["dockermi.active"]
			key := service.Labels["dockermi.key"] // Check for dockermi.key

			if key != "" && orderExists && activeExists && active == "true" {
				if _, exists := groups[key]; !exists {
					groups[key] = []DockermiTypes.ServiceScript{}
				}
				groups[key] = append(groups[key], DockermiTypes.ServiceScript{
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

	return groups, err
}

// ParseComposeFile reads and parses a docker-compose.yml file located at the specified path.
// It extracts the services defined in the file and returns a map of these services. If
// the 'withKey' parameter is true, it assigns a default 'dockermi.key' label to services
// that do not have one defined. If a service is inactive (as indicated by the 'dockermi.active'
// label set to "false"), it will not be included in the returned map.
//
// Parameters:
//   - path: the path to the docker-compose.yml file
//   - withKey: a boolean indicating whether to assign a default 'dockermi.key' to services
//     that lack this label
//
// Returns:
//   - map[string]DockermiTypes.Service: a map where the keys are service names and the values
//     are the corresponding Service structures
//   - error: if any errors occur during reading or parsing the file, they are returned
func ParseComposeFile(path string, withKey bool, force bool) (map[string]DockermiTypes.Service, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var composeFile map[string]interface{}
	err = yaml.Unmarshal(file, &composeFile)
	if err != nil {
		return nil, nil
	}

	services := make(map[string]DockermiTypes.Service)

	servicesData, exists := composeFile["services"]
	if !exists {
		return nil, nil
	}

	if servicesData, ok := servicesData.(map[interface{}]interface{}); ok {
		for name, data := range servicesData {
			if serviceData, ok := data.(map[interface{}]interface{}); ok {
				service, err := unmarshalService(serviceData)
				if err != nil {
					return nil, err
				}
				service.Name = name.(string) // Set the service name
				services[name.(string)] = service
			}
		}
	}

	return services, nil
}

// Custom unmarshal function to handle labels
func unmarshalService(data map[interface{}]interface{}) (DockermiTypes.Service, error) {
	service := DockermiTypes.Service{
		Labels: make(map[string]string),
	}

	if val, ok := data["image"].(string); ok {
		service.Image = val
	}
	if ports, ok := data["ports"].([]interface{}); ok {
		for _, port := range ports {
			if p, ok := port.(string); ok {
				service.Ports = append(service.Ports, p)
			}
		}
	}

	// Handle labels
	if labels, ok := data["labels"]; ok {
		switch labels := labels.(type) {
		case []interface{}:
			for _, label := range labels {
				if strLabel, ok := label.(string); ok {
					parts := strings.SplitN(strLabel, "=", 2)
					if len(parts) == 2 {
						service.Labels[strings.TrimSpace(parts[0])] = strings.TrimSpace(parts[1])
					}
				}
			}
		case map[interface{}]interface{}:
			for k, v := range labels {
				if key, ok := k.(string); ok {
					if value, ok := v.(string); ok {
						service.Labels[key] = value
					}
				}
			}
		}
	}

	return service, nil
}
