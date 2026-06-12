# app-common-go

Simple client access library for the configuration of the [Clowder operator](https://github.com/RedHatInsights/clowder). This library provides a Go interface for applications running in Clowder-managed environments to access their runtime configuration.

## Installation

Add this library to your Go project:

```sh
go get github.com/redhatinsights/app-common-go
```

### Prerequisites

- Go 1.18 or later
- Applications must be deployed in a Clowder-managed environment to use runtime configuration features

## Usage

Import the library and check if Clowder is enabled before accessing configuration:

```go
import (
    clowder "github.com/redhatinsights/app-common-go/pkg/api/v1"
)

func main() {
    if clowder.IsClowderEnabled() {
        fmt.Printf("Public Port: %s", clowder.LoadedConfig.PublicPort)
    }
}
```

The library automatically loads configuration from the `ACG_CONFIG` environment variable when running in a Clowder environment.

## API Overview

The library provides several helper functions and global variables for accessing common configuration elements:

### Configuration Globals

- `clowder.LoadedConfig` - The parsed application configuration
- `clowder.KafkaTopics` - Map of Kafka topics keyed by requested name
- `clowder.KafkaServers` - List of Kafka broker URLs
- `clowder.ObjectBuckets` - Map of object storage buckets keyed by requested name
- `clowder.DependencyEndpoints` - Nested map `[appName][deploymentName]` for public service endpoints
- `clowder.PrivateDependencyEndpoints` - Nested map `[appName][deploymentName]` for private service endpoints

### Helper Methods

- `clowder.IsClowderEnabled()` - Returns true if the `ACG_CONFIG` environment variable is set
- `clowder.LoadedConfig.RdsCa()` - Creates a temporary file with the RDS CA certificate and returns the filename
- `clowder.LoadedConfig.KafkaCa(<BrokerConfig>)` - Creates a temporary file with the Kafka CA certificate and returns the filename (if broker not given, first is chosen)

### Example: Accessing Kafka Configuration

```go
if clowder.IsClowderEnabled() {
    // Access a specific topic
    if topic, ok := clowder.KafkaTopics["my-topic"]; ok {
        fmt.Printf("Topic name: %s\n", topic.Name)
    }
    
    // Get all Kafka brokers
    brokers := clowder.KafkaServers
    fmt.Printf("Kafka brokers: %v\n", brokers)
}
```

## Development

### Running Tests

Set the `ACG_CONFIG` environment variable to point to a test configuration file:

```sh
ACG_CONFIG="../../../tests/test.json" go test -v ./pkg/api/v1/...
```

### Project Structure

- `pkg/api/v1/` - Main API package containing configuration types and loading logic
- `tests/` - Test fixtures and sample configuration files
- `sync_config.sh` - Script to synchronize configuration schema

### Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines on commit messages, signing commits, and opening pull requests.

## License

This project does not currently include a license file. Please contact the maintainers for licensing information.
