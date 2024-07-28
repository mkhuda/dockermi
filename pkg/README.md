# Dockermi (Package)

The `dockermi` package provides functionality to generate a `dockermi.sh` script for managing Docker services defined in `docker-compose.yml` files.
For **cli** version, please refer to [this](https://github.com/mkhuda/dockermi/README.md) readme.

## Installation

To use the `dockermi` package in your Go project, you can import it directly:

```go
import "github.com/mkhuda/dockermi/pkg"
```

## Usage

### RunDockermi Function

You can use the `RunDockermi` function to generate the Docker management script programmatically. Here’s an example:

```go
package main

import (
    "github.com/mkhuda/dockermi/pkg"
    "log"
)

func main() {
    projectDir := "/path/to/your/project"
    
    // Run the dockermi function
    if err := dockermi.RunDockermi(projectDir); err != nil {
        log.Fatalf("Error running Dockermi: %v", err)
    }
}
```

### Testing

The `dockermi_test.go` file contains tests for the `dockermi` package. You can run the tests using:

```bash
go test ./pkg
```

## License

This package is licensed under the MIT License. [LICENSE](../LICENSE)