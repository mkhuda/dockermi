package types

// ServiceScript is a struct at the ServiceScriptReturn
type ServiceScript struct {
	Order       string
	ServiceName string
	ComposeFile string
}

// ServiceScriptReturn represent the return of some internal methods
type ServiceScriptReturn []ServiceScript

// Service represents a service in the docker-compose.yml file.
type Service struct {
	Name string
	// Image  string            `yaml:"image"`
	// Ports  []string          `yaml:"ports"`
	Labels map[string]string `yaml:"labels"`
}

// DockerCompose represents the structure of the docker-compose.yml file.
type DockerCompose struct {
	Services map[string]Service `yaml:"services"`
}
