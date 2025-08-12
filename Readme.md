# Kratix Go SDK

The Kratix Go SDK provides a Go implementation of the [Kratix SDK Contract](https://github.com/syntasso/kratix/blob/main/sdk/contract.md), enabling you to build Kratix workflows in Go. This SDK simplifies the development of Kratix promises and resources by providing a clean, idiomatic Go interface for reading inputs, writing outputs, managing status, and handling destination selectors.

## Features

- **Resource Management**: Read and write Kratix resources and promises
- **Status Handling**: Update resource status and write status files
- **Output Generation**: Write workflow outputs to files
- **Destination Selectors**: Configure where resources should be deployed
- **Environment Variables**: Access workflow context (action, type, promise name, pipeline name)
- **Testing Support**: Built-in testing utilities and mocks

## Installation

Add the Kratix Go SDK to your Go module:

```bash
go get github.com/syntasso/kratix-go
```

## Documentation

For detailed API documentation, visit the [GoDoc page](https://pkg.go.dev/github.com/syntasso/kratix-go).

## Usage

### Basic Workflow Example

Here's a simple example of how to use the SDK in a Kratix workflow:

```go
package main

import (
	"log"

	kratix "github.com/syntasso/kratix-go"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"
)

func main() {
	// Initialize the SDK
	sdk := kratix.New()

	// Access workflow context
	log.Printf("Workflow action: %s", sdk.WorkflowAction())
	log.Printf("Workflow type: %s", sdk.WorkflowType())
	log.Printf("Promise name: %s", sdk.PromiseName())
	log.Printf("Pipeline name: %s", sdk.PipelineName())

	// Read the resource input
	resource, err := sdk.ReadResourceInput()
	if err != nil {
		log.Fatalf("failed to read resource input: %v", err)
	}

	// Extract values from the resource
	fields, err := resource.GetValue("spec.fields")
	if err != nil {
		log.Fatalf("failed to get fields: %v", err)
	}

	// Create a Kubernetes resource (e.g., ConfigMap)
	cm := &corev1.ConfigMap{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "ConfigMap",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.GetName() + "-config",
			Namespace: resource.GetNamespace(),
		},
		Data: map[string]string{
			"workflowAction": sdk.WorkflowAction(),
			"workflowType":   sdk.WorkflowType(),
			"promiseName":    sdk.PromiseName(),
			"pipelineName":   sdk.PipelineName(),
		},
	}

	// Marshall it to YAML
	cmContent, err := yaml.Marshal(cm)
	if err != nil {
		log.Fatalf("failed to marshal config map: %v", err)
	}

	// Write to the output directory
	err = sdk.WriteOutput("config.yaml", cmContent)
	if err != nil {
		log.Fatalf("failed to write output: %v", err)
	}

	// Persist Status to the Resource immediately
	newStatus := kratix.NewStatus()
	newStatus.Set("viaPublishStatus", true)
	if err := sdk.PublishStatus(resource, newStatus); err != nil {
		log.Fatalf("failed to publish status: %v", err)
	}

	// Write Status file (to be persisted by Kratix at the end of the Workflow)
	status := kratix.NewStatus()
	status.Set("viaStatusFile", true)
	if err := sdk.WriteStatus(status); err != nil {
		log.Fatalf("failed to write status: %v", err)
	}

	// Configure and write Destination Selectors
	destinationSelectors := []kratix.DestinationSelector{
		{MatchLabels: map[string]string{"environment": "dev"}},
	}
	if err := sdk.WriteDestinationSelectors(destinationSelectors); err != nil {
		log.Fatalf("failed to write destination selectors: %v", err)
	}
}
```

### Key SDK Methods

- **`ReadResourceInput()`**: Read the resource from `/kratix/input/object.yaml` (for Resource Workflows)
- **`ReadPromiseInput()`**: Read the promise from `/kratix/input/object.yaml` (for Promise Workflows)
- **`WriteOutput(filename, content)`**: Write content to `/kratix/output/`
- **`WriteStatus(status)`**: Write status to `/kratix/metadata/status.yaml`
- **`WriteDestinationSelectors(selectors)`**: Write destination selectors to `/kratix/metadata/destination_selectors.yaml`
- **`PublishStatus(resource, status)`**: Update the resource status in Kubernetes

## Development

### Prerequisites

- Go 1.24.5 or later
- Docker (for end-to-end testing)
- Kind (for local Kubernetes testing)

### Getting Started

1. Clone the repository:
   ```bash
   git clone https://github.com/syntasso/kratix-go.git
   cd kratix-go
   ```

2. Install dependencies:
   ```bash
   go mod download
   ```

3. Run tests:
   ```bash
   make test
   ```

### Testing

The SDK includes comprehensive testing support:

- **Unit Tests**: Run with `make test`
- **End-to-End Tests**: Run with `make e2e-test`
- **Test Coverage**: Generate coverage reports with `make test-coverage`
- **Watch Mode**: Run tests in watch mode with `make test-watch`

### Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## License

This project is licensed under the same license as the Kratix project. See the [LICENSE](LICENSE) file for details.
