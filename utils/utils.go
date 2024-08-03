package utils

import (
	"fmt"
)

// DisplayHelp prints the usage information for the dockermi command to the console.
func DisplayHelp(version string) {
	fmt.Printf(`
Dockermi version: %s  
Usage: dockermi [command] [options]

This command generates a dockermi.sh script to manage Docker services defined in docker-compose.yml files.  
The 'dockermi run | down' command can be run within a folder that contains the dockermi.sh script, 
and the 'dockermi.sh' script is created in the current directory where the 'dockermi' command is executed.

Commands:
    create <service-key>   Generate a dockermi.sh script for the specified service key.
    up [options]           Start the Docker services defined in the dockermi.sh file in current directory.
    down [options]         Stop the Docker services defined in the dockermi.sh file in current directory.

Options:
    --help                 Display this help message and exit.
    --version              Display current installed version.
    
Examples:
    dockermi                        # Generates a dockermi.sh script in the current directory.
    dockermi create myservicekey    # [Experimental] Create a script for the specified service key.
    dockermi up -d --build              # Start services with the --build option.
    dockermi down --remove-orphans   # Stop services and remove orphan containers.`, version)
}
