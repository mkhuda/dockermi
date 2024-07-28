package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "os"
    "path/filepath"
    "sort"
    "time"

    "github.com/fatih/color"
    "github.com/schollz/progressbar/v3"
    "gopkg.in/yaml.v2"
)

// Service represents a service in the docker-compose.yml file.
type Service struct {
    Image  string            `yaml:"image"`
    Ports  []string          `yaml:"ports"`
    Labels map[string]string `yaml:"labels"`
}

// DockerCompose represents the structure of the docker-compose.yml file.
type DockerCompose struct {
    Services map[string]Service `yaml:"services"`
}

// displayHelp prints the help message.
func displayHelp() {
    fmt.Println(`Usage: dockermi [--help]

This command generates a dockermi.sh script to manage Docker services defined in docker-compose.yml files.

Options:
    --help  Display this help message and exit.

Examples:
    dockermi
    ./dockermi.sh up
    ./dockermi.sh down`)
}

func main() {
    // Define the command-line flags
    help := flag.Bool("help", false, "Display help information")
    flag.Parse()

    // Show help if requested
    if *help {
        displayHelp()
        return
    }

    // Get current directory
    projectDir, err := os.Getwd()
    if err != nil {
        color.Red("Error getting current directory: %v", err)
        return
    }

    dockermiScriptPath := filepath.Join(projectDir, "dockermi.sh")
    var services []struct {
        Order       string
        ServiceName string
        ComposeFile string
    }

    // Flag to check if any docker-compose.yml file is found
    foundDockerCompose := false

    // Traverse directories up to a depth of 2
    err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        if info.IsDir() || filepath.Base(path) != "docker-compose.yml" {
            return nil
        }

        // Set the flag to true when a docker-compose.yml is found
        foundDockerCompose = true

        // Read and parse the docker-compose.yml file
        data, err := ioutil.ReadFile(path)
        if err != nil {
            return err
        }

        var dockerCompose DockerCompose
        if err := yaml.Unmarshal(data, &dockerCompose); err != nil {
            return err
        }

        // Process each service
        for serviceName, service := range dockerCompose.Services {
            order, orderExists := service.Labels["dockermi.order"]
            active, activeExists := service.Labels["dockermi.active"]

            if orderExists && activeExists && active == "true" { // Check if active is true
                color.Green("Service: %s", serviceName)
                color.Blue("Order: %s", order)
                color.Yellow("Active: %s", active)

                // Store service info
                services = append(services, struct {
                    Order       string
                    ServiceName string
                    ComposeFile string
                }{
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
        return
    }

    // Check if any docker-compose.yml files were found
    if !foundDockerCompose {
        color.Yellow("No docker-compose.yml found within this folder")
        return // Exit without creating the dockermi.sh script
    }

    // Create or clear the dockermi.sh file
    dockermiScript, err := os.Create(dockermiScriptPath)
    if err != nil {
        color.Red("Error creating dockermi.sh file: %v", err)
        return
    }
    defer dockermiScript.Close()

    // Write the script header
    dockermiScript.WriteString("#!/bin/bash\n\n")
    dockermiScript.WriteString("Usage: ./dockermi.sh [up|down]\n\n")

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

        // Simulate some progress (optional)
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

        // Simulate some progress (optional)
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
    if err := os.Chmod(dockermiScriptPath, 0755); err != nil {
        color.Red("Error making the script executable: %v", err)
        return
    }

    color.Green("Generated script: %s", dockermiScriptPath)
    color.Blue("You can now run ./dockermi.sh up or ./dockermi.sh down")
}
