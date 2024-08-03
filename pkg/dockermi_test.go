package dockermi_test

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/mkhuda/dockermi/internal/dockercompose"
	dockermi "github.com/mkhuda/dockermi/pkg"
)

// Helper function to verify the existence of a docker-compose.yml file
func verifyComposeFileExists(t *testing.T, relativePath string) {
	t.Helper()
	// Construct the absolute path to the docker-compose.yml file
	composeFile := filepath.Join("../test", relativePath)
	if _, err := os.Stat(composeFile); os.IsNotExist(err) {
		t.Fatalf("Expected docker-compose file does not exist: %v", composeFile)
	}
}

// Test for valid docker-compose.yml in the nginx and postgres folders
func TestCreatedDockermi(t *testing.T) {
	// Verify the existence of the required docker-compose.yml files
	verifyComposeFileExists(t, "nginx/docker-compose.yml")
	verifyComposeFileExists(t, "postgres/docker-compose.yml")

	// Change to the directory of the compose files
	if err := os.Chdir("../test/"); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	projectDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}

	// Call the main function
	if _, err := dockermi.RunDockermi(projectDir); err != nil { // Assuming you refactored to RunDockermi
		fmt.Printf("project dir %v", projectDir)
		t.Errorf("RunDockermi failed: %v", err)
	}

	// Get the current working directory
	currentDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}

	// Check if the dockermi.sh file is created
	scriptPath := filepath.Join(currentDir, "dockermi.sh")
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		t.Errorf("Expected dockermi.sh to be created, but it was not.")
	}

	// Clean up after test
	// if err := os.Remove(scriptPath); err != nil {
	// 	t.Errorf("Failed to remove dockermi.sh: %v", err)
	// }
}

func TestFindService(t *testing.T) {
	// Change to the directory of the compose files
	if err := os.Chdir("../test/"); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	projectDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}

	services, err := dockercompose.FindServices(projectDir, false)
	if err != nil {
		t.Fatalf("Error finding services: %v", err)
	}

	servicesLength := len(services)

	expectedLength := 6
	if servicesLength != expectedLength {
		t.Fatalf("Expected %v keys to be created. Created keys are: %v", expectedLength, servicesLength)
	}
}

func TestFindServiceForce(t *testing.T) {
	// Change to the directory of the compose files
	if err := os.Chdir("../test/"); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	projectDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}

	services, err := dockercompose.FindServices(projectDir, true)
	if err != nil {
		t.Fatalf("Error finding services: %v", err)
	}
	servicesLength := len(services)

	expectedLength := 9
	if servicesLength != expectedLength {
		t.Fatalf("Expected %v keys to be created. Created keys are: %v", expectedLength, servicesLength)
	}
}

func TestFindServiceByKey(t *testing.T) {
	// Change to the directory of the compose files
	if err := os.Chdir("../test/"); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	projectDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}

	services, err := dockercompose.FindServicesWithKey(projectDir)
	if err != nil {
		t.Fatalf("Error finding services: %v", err)
	}
	servicesLength := len(services)
	if servicesLength != 3 {
		t.Fatalf("Expected 3 keys to be created. Create keys are: %v", servicesLength)
	}
}

func TestDirectoriesWithoutCompose(t *testing.T) {
	// Change to the directory without compose files
	if err := os.Chdir("../test/empty"); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	projectDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Error getting current working directory: %v", err)
	}

	services, err := dockercompose.FindServices(projectDir, false)
	if err != nil {
		t.Fatalf("Error finding services: %v", err)
	}
	servicesLength := len(services)
	if servicesLength != 0 {
		t.Fatalf("Expected 0 keys to be created. Create keys are: %v", servicesLength)
	}
}
