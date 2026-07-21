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
- `clowder.DependencyEndpoints` - Nested map `[appName][deploymentName]` for V1 public service endpoints
- `clowder.PrivateDependencyEndpoints` - Nested map `[appName][deploymentName]` for V1 private service endpoints
- `clowder.DependencyEndpointsV2` - Nested map `[appName][serviceName]` for V2 public service endpoints (URI-based)
- `clowder.PrivateDependencyEndpointsV2` - Nested map `[appName][serviceName]` for V2 private service endpoints (URI-based)

### Helper Methods

- `clowder.IsClowderEnabled()` - Returns true if the `ACG_CONFIG` environment variable is set
- `clowder.GetV2DependencyEndpoint(app, name)` - Retrieves a V2 public endpoint by app and service name; returns `(endpoint, bool)`
- `clowder.GetV2PrivateDependencyEndpoint(app, name)` - Retrieves a V2 private endpoint by app and service name; returns `(endpoint, bool)`
- `clowder.LoadedConfig.RdsCa()` - Creates a temporary file with the RDS CA certificate and returns the filename
- `clowder.LoadedConfig.KafkaCa(<BrokerConfig>)` - Creates a temporary file with the Kafka CA certificate and returns the filename (if broker not given, first is chosen)
- `clowder.LoadedConfig.KafkaFirstCa()` - Convenience method: creates a temporary file with the Kafka CA certificate from the first broker, with nil-safety checks

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

### Example: Accessing V2 Dependency Endpoints

V2 endpoints provide a simplified URI-based connection model with explicit authentication and certificate handling:

```go
if clowder.IsClowderEnabled() {
    // Use the getter function for safe lookups
    if endpoint, ok := clowder.GetV2DependencyEndpoint("rbac-service", "api"); ok {
        fmt.Printf("Service URI: %s\n", endpoint.Uri)
        fmt.Printf("Requires authentication: %v\n", endpoint.Authenticated)
        
        // If TLS is in use, CA certificate path is provided
        if endpoint.CaCertificate != nil {
            fmt.Printf("CA certificate: %s\n", *endpoint.CaCertificate)
        }
    }
    
    // Or directly access the nested map (less safe if app/service might not exist)
    if appEndpoints, ok := clowder.DependencyEndpointsV2["rbac-service"]; ok {
        if endpoint, ok := appEndpoints["api"]; ok {
            fmt.Printf("Found endpoint: %s\n", endpoint.Uri)
        }
    }
}
```

**V2 Endpoint Fields:**

- `Uri` - Complete URI including protocol, hostname, and port (e.g., `http://service.svc:8000` or `https://service:8443`)
- `Authenticated` - Boolean flag indicating if the endpoint requires authentication
  - By default, `true` for cross-cluster dependencies (ClowdAppRef routed through gateways)
  - By default, `false` for in-cluster dependencies (ClowdApp with network isolation)
  - This default can be overridden per-deployment via `webServices.public.authenticated` or `webServices.private.authenticated` on the ClowdApp or ClowdAppRef resource
- `CaCertificate` - Optional path to CA certificate file (only present for HTTPS URIs)

## Development

### Running Tests

Set the `ACG_CONFIG` environment variable to point to a test configuration file:

```sh
ACG_CONFIG="testdata/test.json" go test -v ./...
```

### Project Structure

- `pkg/api/v1/` - Main API package containing configuration types and loading logic
- `pkg/api/v1/testdata/` - Test fixtures and sample configuration files
- `sync_config.sh` - Script to synchronize configuration schema

### Contributing

Contributions are welcome. Please read [CONTRIBUTING.md](./CONTRIBUTING.md) for guidelines on commit messages, signing commits, and opening pull requests.

## Release Notes

### V2 Dependency Endpoints (ENGPROD-10121)

**Added:** Support for Clowder V2 dependency endpoints with simplified URI-based configuration.

**What's New:**

- New globals: `DependencyEndpointsV2` and `PrivateDependencyEndpointsV2` expose V2 endpoints in indexed nested maps `[appName][serviceName]`
- New getter functions: `GetV2DependencyEndpoint()` and `GetV2PrivateDependencyEndpoint()` provide safe lookups with bounds-checking
- New field on V2 endpoints: `Authenticated` boolean flag indicates whether authentication is required
  - By default, `true` for cross-cluster dependencies (ClowdAppRef with gateway routing)
  - By default, `false` for in-cluster dependencies (ClowdApp with network isolation)
  - This default can be overridden per-deployment via `webServices.public.authenticated` or `webServices.private.authenticated` on the ClowdApp or ClowdAppRef resource
- Per-endpoint CA certificates: `CaCertificate` field provides the path to TLS CA for HTTPS endpoints (null for HTTP)

**Backward Compatibility:**

- Existing V1 API unchanged: `DependencyEndpoints`, `PrivateDependencyEndpoints`, and all existing helper methods remain fully supported
- V2 and V1 endpoints can coexist in the same configuration; applications can use either API depending on Clowder version
- If Clowder deployment does not emit V2 endpoints, the V2 globals remain `nil` and getters safely return `(zero, false)` instead of panicking

**Migration Path:**

Applications should prefer the V2 API for new code:
- Simpler connection logic: single URI per endpoint, no protocol/port selection needed
- Explicit authentication semantics: clear signal of whether cross-cluster auth is required
- Per-endpoint CA handling: eliminates shared CA path ambiguity

## License

This project does not currently include a license file. Please contact the maintainers for licensing information.
