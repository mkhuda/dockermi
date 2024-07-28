# Dockermi

**Dockermi** is a command-line tool for managing Docker services defined in `docker-compose.yml` files. It simplifies the process of starting and stopping multiple Docker services with a single command, generating a shell script (`dockermi.sh`) that can be executed to perform the desired actions.

![Screenshot 2024-07-28 164852](https://github.com/user-attachments/assets/edb4e6b6-e788-49e2-a4db-148896a416c7)

## Features

- Automatically discovers `docker-compose.yml` files in the current directory and its subdirectories.
- Generates a shell script (`dockermi.sh`) for managing services.
- Supports starting and stopping services with simple commands.
- Provides colored logs for better readability and user experience.
- Easy to use and integrate into existing workflows.

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Help](#help)
- [Contributing](#contributing)
- [License](#license)

## Installation

### Prerequisites

- [Go](https://golang.org/dl/) (version 1.18 or later)
- Docker and Docker Compose installed on your system.

### Installing Dockermi

To install the Dockermi application, follow the instructions below based on your operating system.

### For Linux and macOS

1. Open your terminal and clone this repo `git clone https://github.com/mkhuda/dockermi.git`.
2. Navigate to the directory where you have the Dockermi project:
   ```bash
   cd path/to/your/dockermi
   ```
3. Make the build and installation script executable:
   ```bash
   chmod +x install.sh build.sh
   ```
4. Run the build script:
   ```bash
   ./build.sh
   ```
5. Run the installation script with **sudo**:
   ```bash
   sudo ./install.sh
   ```

This will build the application for your OS and install it to `/usr/local/bin`, making it available for execution from anywhere in your terminal.

### For Windows

1. Open Command Prompt and clone this repo `git clone https://github.com/mkhuda/dockermi.git`.
2. Navigate to the directory where you have the Dockermi project:
   ```cmd
   cd path\to\your\dockermi
   ```
3. Run the build script:
   ```cmd
   build.bat
   ```
4. Run the installation script:
   ```cmd
   install.bat
   ```

This will build the application for Windows and install it to `C:\Program Files\dockermi`, making it available for execution from anywhere in your command prompt.

## Running Dockermi

Once installed, you can run the Dockermi application by typing:

```bash
dockermi
```

in your terminal (Linux and macOS) or command prompt (Windows).

### Uninstallation

To uninstall Dockermi, you will need to manually remove the installed binary:

- **Linux and macOS**:
  ```bash
  sudo rm /usr/local/bin/dockermi
  ```

- **Windows**:
  ```cmd
  del "C:\Program Files\dockermi\dockermi.exe"
  ```

### Build the Executable

1. Clone the repository to your local machine:

    ```bash
    git clone https://github.com/mkhuda/dockermi.git
    cd dockermi
    ```

2. Build the `dockermi` executable:

    ```bash
    ./build.sh
    ```

### Make the Executable Available

You can move the built executable to a directory in your `PATH` for easier access:

```bash
sudo mv dockermi /usr/local/bin/
```


## Usage

To generate the `dockermi.sh` script, run:

```bash
dockermi
```

This command creates a `dockermi.sh` script in the current directory, which contains functions for starting and stopping your Docker services.

### Annotations in docker-compose.yml

#### 1. `dockermi.order`

- **Description**: This annotation specifies the order in which the Docker services should be started or stopped. Services with lower order values are started before those with higher values. This is particularly useful when certain services depend on others being up and running first.

- **Type**: String (represents a numeric value)

- **Example**:
    ```yaml
    services:
      web:
        image: nginx:latest
        ports:
          - "80:80"
        labels:
          dockermi.order: "1"  # This service will start first
    ```

- **Usage**: You can set this label in your `docker-compose.yml` to control the startup order of your services, ensuring that dependencies are handled appropriately. For example, a database service might have an order of `1`, while a web service that depends on it could have an order of `2`.

- **Multiple Services**: If multiple services are defined in the `docker-compose.yml` file with the same `dockermi.order`, only the first service that appears in the file will be used for execution. Be mindful of the order in which services are listed.

#### 2. `dockermi.active`

- **Description**: This annotation indicates whether the service is currently active or should be started when the `dockermi.sh` script is executed. If set to `"true"`, the service will be included in the startup process. If set to `"false"`, the service will be skipped during startup.

- **Type**: String (boolean value, "true" or "false")

- **Example**:
    ```yaml
    services:
      web:
        image: nginx:latest
        ports:
          - "80:80"
        labels:
          dockermi.active: "true"  # This service is active and will be started
    ```

- **Usage**: Use this label to manage which services should be actively started or stopped. For instance, if you have a service that is temporarily not needed, you can set `dockermi.active: "false"` to prevent it from starting.

Here is how you might define a service in your `docker-compose.yml` file using both annotations:

```yaml
version: '3.8'

services:
  web:
    image: nginx:latest
    ports:
      - "80:80"
    labels:
      dockermi.order: "1"    # Service start order
      dockermi.active: "true" # This service is active

  db:
    image: mysql:latest
    environment:
      MYSQL_ROOT_PASSWORD: root
    labels:
      dockermi.order: "2"    # This service will start after 'web'
      dockermi.active: "true" # This service is also active
```

#### Summary

- The `dockermi.order` annotation controls the startup order of services.
- The `dockermi.active` annotation determines whether a service should be active during the execution of the `dockermi.sh` script.
- When multiple services have the same `dockermi.order`, only the first service that appears in the `docker-compose.yml` file will be used for execution.
- Using these annotations helps to manage complex service dependencies effectively, ensuring that the right services are up and running when needed.

#### Further Considerations

- **Multiple Services**: When defining multiple services, ensure that their `dockermi.order` values are unique (still ok, if there is multiple version) and reflect the intended startup sequence. If services share the same order value, only the first one listed will be activated, which can lead to unexpected behavior if not managed correctly.
- **Dynamic Activation**: You can dynamically set the `dockermi.active` label based on environment variables or configuration settings to enable/disable services as needed.

By incorporating these annotations into your `docker-compose.yml` file, you can leverage the full power of Dockermi to manage your Docker services efficiently. If you have any further questions or need clarification, feel free to ask!


### Running the Generated Script

1. To start the services defined in your `docker-compose.yml` files, run:

    ```bash
    ./dockermi.sh up
    ```

2. To stop the services, run:

    ```bash
    ./dockermi.sh down
    ```

### Help

To display help information for the `dockermi` command, run:

```bash
dockermi --help
```

This will show the usage details and available commands.

## Contributing

We welcome contributions to improve the functionality and usability of Dockermi. Here’s how you can help:

1. **Fork the Repository**: Click the "Fork" button at the top right of the page to create your own copy of the repository.
2. **Create a Branch**: Create a new branch for your feature or bug fix.

    ```bash
    git checkout -b feature/my-feature
    ```

3. **Make Changes**: Implement your feature or fix.
4. **Commit Your Changes**: Commit your changes with a descriptive message.

    ```bash
    git commit -m "Add my feature"
    ```

5. **Push to Your Branch**: Push your changes to GitHub.

    ```bash
    git push origin feature/my-feature
    ```

6. **Create a Pull Request**: Go to the original repository and click on "New Pull Request."

### Issues

If you encounter any bugs or have suggestions for improvements, please feel free to open an issue on the GitHub repository.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Thanks to the contributors and the open-source community for their support and inspiration.
- Special thanks to the maintainers of Go, Docker, and Docker Compose for providing powerful tools that make this project possible.
````