package kratix

import (
	"fmt"
	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

type ResourceAccessor interface {
	// GetValue queries the resource and returns the value at the specified path e.g. spec.dbConfig.size
	GetValue(string) (any, error)
	// GetStatus queries the resource and returns the resource.status
	GetStatus(string) (StatusModifier, error)
	// GetName queries the resource and returns the name
	GetName() string
	// GetStatus queries the resource and returns the namespace
	GetNamespace() string
	// GroupVersionKind queries the resource and returns the GroupVersionKind
	GetGroupVersionKind() schema.GroupVersionKind
	// GetLabels queries the resource and returns the labels
	GetLabels() map[string]string
	// GetAnnotations queries the resource and returns the annotations
	GetAnnotations() map[string]string
}

// Resource implements contract.Resource backed by an unstructured object.
type Resource struct {
	obj unstructured.Unstructured
}

var _ ResourceAccessor = (*Resource)(nil)

// GetValue returns the value at the provided path.
func (r *Resource) GetValue(path string) (any, error) {
	val, found, err := unstructured.NestedFieldNoCopy(r.obj.Object, strings.Split(path, ".")...)
	if err != nil {
		return nil, err
	}
	if !found {
		return nil, fmt.Errorf("path %s not found", path)
	}
	return val, nil
}

// GetStatus returns the Status at the provided path.
func (r *Resource) GetStatus(path string) (StatusModifier, error) {
	parts := []string{"status"}
	if path != "" {
		parts = append(parts, strings.Split(path, ".")...)
	}
	val, _, err := unstructured.NestedFieldNoCopy(r.obj.Object, parts...)
	if err != nil {
		return nil, err
	}
	m, ok := val.(map[string]any)
	if !ok {
		m = map[string]any{}
	}
	return &Status{data: m}, nil
}

// GetName returns the resource name.
func (r *Resource) GetName() string { return r.obj.GetName() }

// GetNamespace returns the resource namespace.
func (r *Resource) GetNamespace() string { return r.obj.GetNamespace() }

// GetGroupVersionKind returns the GVK of the resource.
func (r *Resource) GetGroupVersionKind() schema.GroupVersionKind { return r.obj.GroupVersionKind() }

// GetLabels returns the labels of the resource.
func (r *Resource) GetLabels() map[string]string { return r.obj.GetLabels() }

// GetAnnotations returns the annotations of the resource.
func (r *Resource) GetAnnotations() map[string]string { return r.obj.GetAnnotations() }
