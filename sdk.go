package kratix

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/syntasso/kratix/api/v1alpha1"
	"sigs.k8s.io/yaml"
)

// The SDK interface implements the Kratix SDK core library function
type SDKInvoker interface {
	// ReadResourceInput reads the file in /kratix/input/object.yaml and returns a Resource
	ReadResourceInput() (Resource, error)
	// ReadPromiseInput reads the file in /kratix/input/object.yaml and returns a Resource
	ReadPromiseInput() (Promise, error)
	// ReadDestinationSelectors
	ReadDestinationSelectors() ([]DestinationSelector, error)
	// ReadStatus reads the /kratix/metadata/status.yaml
	ReadStatus() (StatusModifier, error)
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
	PublishStatus(Resource, StatusModifier) error
}

// ensure SDKInvoker implemented
var _ SDKInvoker = (*KratixSDK)(nil)

// KratixSDK implements the SDKInvoker interface for reading and writing
// Kratix workflow data.
type KratixSDK struct {
	outputDir   string
	inputDir    string
	metadataDir string

	inputObject string
}

// Option configures KratixSDK.
type Option func(*KratixSDK)

// WithInputDir overrides the path to the input directory.
func WithInputDir(p string) Option {
	return func(k *KratixSDK) { k.inputDir = p }
}

// WithInputObject overrides the name of the input object file.
func WithInputObject(p string) Option {
	return func(k *KratixSDK) { k.inputObject = p }
}

func WithMetadataDir(p string) Option {
	return func(k *KratixSDK) { k.metadataDir = p }
}

// WithOutputDir overrides the output directory path.
func WithOutputDir(p string) Option {
	return func(k *KratixSDK) { k.outputDir = p }
}

// New creates a KratixSDK with optional configuration overrides.
func New(opts ...Option) *KratixSDK {
	sdk := &KratixSDK{
		inputDir:    "/kratix/input",
		metadataDir: "/kratix/metadata",
		outputDir:   "/kratix/output",
		inputObject: "object.yaml",
	}
	for _, opt := range opts {
		opt(sdk)
	}
	return sdk
}

// ReadResourceInput reads the object YAML and returns a Resource.
func (k *KratixSDK) ReadResourceInput() (Resource, error) {
	data, err := os.ReadFile(filepath.Join(k.inputDir, k.inputObject))
	if err != nil {
		return nil, fmt.Errorf("read object input: %w", err)
	}
	r := &ResourceImpl{}
	if err := yaml.Unmarshal(data, &r.obj.Object); err != nil {
		return nil, fmt.Errorf("unmarshal object: %w", err)
	}
	return r, nil
}

// ReadPromiseInput reads the object YAML and returns it as a Promise.
func (k *KratixSDK) ReadPromiseInput() (Promise, error) {
	data, err := os.ReadFile(filepath.Join(k.inputDir, k.inputObject))
	if err != nil {
		return nil, fmt.Errorf("read promise input: %w", err)
	}
	p := &v1alpha1.Promise{}
	if err := yaml.Unmarshal(data, &p); err != nil {
		return nil, fmt.Errorf("unmarshal promise: %w", err)
	}
	obj, err := p.ToUnstructured()
	if err != nil {
		return nil, fmt.Errorf("unmarshal promise: %w", err)
	}
	return &PromiseImpl{ResourceImpl: ResourceImpl{obj: *obj}, promise: p}, nil
}

// ReadDestinationSelectors reads destination selectors from file.
func (k *KratixSDK) ReadDestinationSelectors() ([]DestinationSelector, error) {
	data, err := os.ReadFile(filepath.Join(k.metadataDir, "destination_selectors.yaml"))
	if err != nil {
		return nil, fmt.Errorf("read destination selectors: %w", err)
	}
	var selectors []DestinationSelector
	if err := yaml.Unmarshal(data, &selectors); err != nil {
		return nil, fmt.Errorf("unmarshal destination selectors: %w", err)
	}
	return selectors, nil
}

// ReadStatus reads the status.yaml from the output directory.
func (k *KratixSDK) ReadStatus() (StatusModifier, error) {
	// TODO: fix this; status.yaml should be read from the /kratix/metadata directory
	data, err := os.ReadFile(filepath.Join(k.metadataDir, "status.yaml"))
	if err != nil {
		return nil, fmt.Errorf("read status: %w", err)
	}
	var m map[string]any
	if err := yaml.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("unmarshal status: %w", err)
	}
	return &Status{data: m}, nil
}

func (k *KratixSDK) write(dir, relPath string, content []byte) error {
	full := filepath.Join(dir, relPath)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return fmt.Errorf("create output dir: %w", err)
	}
	if err := os.WriteFile(full, content, 0o644); err != nil {
		return fmt.Errorf("write output file: %w", err)
	}
	return nil
}

// WriteOutput writes content to the named file under the output directory.
func (k *KratixSDK) WriteOutput(relPath string, content []byte) error {
	return k.write(k.outputDir, relPath, content)
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
	return k.write(k.metadataDir, "status.yaml", data)
}

// WriteDestinationSelectors writes the selectors to destination_selectors.yaml.
func (k *KratixSDK) WriteDestinationSelectors(ds []DestinationSelector) error {
	// TODO: make sure this merges the existing destination selectors with the new ones
	data, err := yaml.Marshal(ds)
	if err != nil {
		return fmt.Errorf("marshal destination selectors: %w", err)
	}
	// TODO: fix this; destination_selectors.yaml should be written to the /kratix/metadata directory
	return k.write(k.metadataDir, "destination_selectors.yaml", data)
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
func (k *KratixSDK) PublishStatus(res Resource, s StatusModifier) error {
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
