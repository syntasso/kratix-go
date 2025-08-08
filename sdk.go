package kratix

import (
	"fmt"
	"os"
	"path/filepath"

	"sigs.k8s.io/yaml"
)

// TODO: add some logging to the functions in the KratixSDK implementation; make the log levels configurable

// The SDK interface implements the Kratix SDK core library function
type SDKInvoker interface {
	// ReadResourceInput reads the file in /kratix/input/object.yaml and returns a Resource
	ReadResourceInput() (ResourceAccessor, error)
	// ReadPromiseInput reads the file in /kratix/input/object.yaml and returns a Resource
	ReadPromiseInput() (PromiseAccessor, error)
	// ReadDestinationSelectors
	ReadDestinationSelectors() ([]DestinationSelector, error)
	// WriteOutput writes the content to the specifies file at the path /kratix/output/filepath
	WriteOutput(string, []byte) error
	// WriteStatus writes the specified status to the /kratix/output/status.yaml
	WriteStatus(StatusModifier) error
	// WriteDestinationSelectors writes the specified Destination Selectors to the /kratix/output/destination_selectors.yaml
	WriteDestinationSelectors([]DestinationSelector) error
	// WorkflowAction returns the value of KRATIX_WORKFLOW_ACTION environment variable
	WorkflowAction() string
	// WorkflowType returns the value of KRATIX_WORKFLOW_TYPE environment variable
	WorkflowType() string
	// PromiseName returns the value of the KRATIX_PROMISE_NAME environment variable
	PromiseName() string
	// PipelineName returns the value of the KRATIX_PIPELINE_NAME environment variable
	PipelineName() string
	// PublishStatus updates the status of the provided resource with the provided status
	PublishStatus(ResourceAccessor, StatusModifier) error
	// ReadStatus reads the /kratix/output/status.yaml
	ReadStatus() (StatusModifier, error)
}

// ensure SDKInvoker implemented
var _ SDKInvoker = (*KratixSDK)(nil)

// KratixSDK implements the SDKInvoker interface for reading and writing
// Kratix workflow data.
type KratixSDK struct {
	objectPath               string
	destinationSelectorsPath string
	outputDir                string
}

// Option configures KratixSDK.
type Option func(*KratixSDK)

// WithObjectPath overrides the path to the object input file.
func WithObjectPath(p string) Option {
	return func(k *KratixSDK) { k.objectPath = p }
}

// WithDestinationSelectorsPath overrides the path to the destination selectors input file.
func WithDestinationSelectorsPath(p string) Option {
	return func(k *KratixSDK) { k.destinationSelectorsPath = p }
}

// WithOutputDir overrides the output directory path.
func WithOutputDir(p string) Option {
	return func(k *KratixSDK) { k.outputDir = p }
}

// New creates a KratixSDK with optional configuration overrides.
func New(opts ...Option) *KratixSDK {
	sdk := &KratixSDK{
		objectPath:               "/kratix/input/object.yaml",
		destinationSelectorsPath: "/kratix/metadata/destination_selectors.yaml",
		outputDir:                "/kratix/output",
	}
	for _, opt := range opts {
		opt(sdk)
	}
	return sdk
}

// ReadResourceInput reads the object YAML and returns a Resource.
func (k *KratixSDK) ReadResourceInput() (ResourceAccessor, error) {
	// TODO: change the name of the interface to Resource, update the struct to something else
	data, err := os.ReadFile(k.objectPath)
	if err != nil {
		return nil, fmt.Errorf("read object input: %w", err)
	}
	r := &Resource{}
	if err := yaml.Unmarshal(data, &r.obj.Object); err != nil {
		return nil, fmt.Errorf("unmarshal object: %w", err)
	}
	return r, nil
}

// ReadPromiseInput reads the object YAML and returns it as a Promise.
func (k *KratixSDK) ReadPromiseInput() (PromiseAccessor, error) {
	// TODO: change the name of the interface to Promise, update the struct to something else
	data, err := os.ReadFile(k.objectPath)
	if err != nil {
		return nil, fmt.Errorf("read promise input: %w", err)
	}

	// TODO: change this to return an unstructured object, similar to the ReadResourceInput function
	//       make sure to update the PromiseAccessor interface to encapsulate the unstructured object
	var out map[string]any
	if err := yaml.Unmarshal(data, &out); err != nil {
		return nil, fmt.Errorf("unmarshal promise: %w", err)
	}
	return out, nil
}

// ReadDestinationSelectors reads destination selectors from file.
func (k *KratixSDK) ReadDestinationSelectors() ([]DestinationSelector, error) {
	data, err := os.ReadFile(k.destinationSelectorsPath)
	if err != nil {
		return nil, fmt.Errorf("read destination selectors: %w", err)
	}
	var selectors []DestinationSelector
	if err := yaml.Unmarshal(data, &selectors); err != nil {
		return nil, fmt.Errorf("unmarshal destination selectors: %w", err)
	}
	return selectors, nil
}

// WriteOutput writes content to the named file under the output directory.
func (k *KratixSDK) WriteOutput(relPath string, content []byte) error {
	full := filepath.Join(k.outputDir, relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	if err := os.WriteFile(full, content, 0o644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	return nil
}

// WriteStatus writes the provided Status to status.yaml.
func (k *KratixSDK) WriteStatus(s StatusModifier) error {
	// TODO: do we need to passa StatusModifier or is the Status object enough?
	// TODO: make sure this merges the existing status from the file with the new status
	sts, ok := s.(*Status)
	if !ok {
		return fmt.Errorf("unsupported status type %T", s)
	}
	data, err := yaml.Marshal(sts.data)
	if err != nil {
		return fmt.Errorf("marshal status: %w", err)
	}
	// TODO: fix this; status.yaml should be written to the /kratix/metadata directory
	return k.WriteOutput("status.yaml", data)
}

// WriteDestinationSelectors writes the selectors to destination_selectors.yaml.
func (k *KratixSDK) WriteDestinationSelectors(ds []DestinationSelector) error {
	// TODO: make sure this merges the existing destination selectors with the new ones
	data, err := yaml.Marshal(ds)
	if err != nil {
		return fmt.Errorf("marshal destination selectors: %w", err)
	}
	// TODO: fix this; destination_selectors.yaml should be written to the /kratix/metadata directory
	return k.WriteOutput("destination_selectors.yaml", data)
}

// WorkflowAction returns the workflow action environment variable.
func (k *KratixSDK) WorkflowAction() string {
	return os.Getenv("KRATIX_WORKFLOW_ACTION")
}

// WorkflowType returns the workflow type environment variable.
func (k *KratixSDK) WorkflowType() string {
	return os.Getenv("KRATIX_WORKFLOW_TYPE")
}

// PromiseName returns the promise name environment variable.
func (k *KratixSDK) PromiseName() string {
	return os.Getenv("KRATIX_PROMISE_NAME")
}

// PipelineName returns the pipeline name environment variable.
func (k *KratixSDK) PipelineName() string {
	return os.Getenv("KRATIX_PIPELINE_NAME")
}

// PublishStatus merges the provided status into the resource and persists it.
func (k *KratixSDK) PublishStatus(res ResourceAccessor, s StatusModifier) error {
	panic("not implemented")
	// r, ok := res.(*Resource)
	// if !ok {
	// 	return fmt.Errorf("unsupported resource type %T", res)
	// }
	// newStatus, ok := s.(*Status)
	// if !ok {
	// 	return fmt.Errorf("unsupported status type %T", s)
	// }
	// existing, ok := r.obj.Object["status"].(map[string]any)
	// if !ok {
	// 	existing = map[string]any{}
	// }
	// mergeMaps(existing, newStatus.data)
	// r.obj.Object["status"] = existing
	// data, err := yaml.Marshal(existing)
	// if err != nil {
	// 	return fmt.Errorf("marshal status: %w", err)
	// }
	// return k.WriteOutput("status.yaml", data)
}

// ReadStatus reads the status.yaml from the output directory.
func (k *KratixSDK) ReadStatus() (StatusModifier, error) {
	// TODO: fix this; status.yaml should be read from the /kratix/metadata directory
	data, err := os.ReadFile(filepath.Join(k.outputDir, "status.yaml"))
	if err != nil {
		return nil, fmt.Errorf("read status: %w", err)
	}
	var m map[string]any
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("unmarshal status: %w", err)
	}
	return &Status{data: m}, nil
}

// mergeMaps recursively merges src into dst.
func mergeMaps(dst, src map[string]any) {
	for k, v := range src {
		if mv, ok := v.(map[string]any); ok {
			if existing, ok := dst[k].(map[string]any); ok {
				mergeMaps(existing, mv)
				continue
			}
		}
		dst[k] = v
	}
}
