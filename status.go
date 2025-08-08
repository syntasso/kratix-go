package kratix

import (
	"errors"
	"strings"
)

// TODO: Do we want a GetConditions() in the Status interface object?

type StatusModifier interface {
	// Get queries the Status and retrieves the value at the specified path e.g. healthStatus.state
	Get(string) any
	// Set updates the value at the specified path e.g. healthStatus.state
	Set(string, any) error
	// Set removes the value at the specified path e.g. healthStatus.state
	Remove(string) bool
}

// Status implements StatusModifier using a generic map.
type Status struct {
	data map[string]any
}

var _ StatusModifier = (*Status)(nil)

// Get retrieves the value at the provided path.
func (s *Status) Get(path string) any {
	parts := strings.Split(path, ".")
	var current any = s.data
	for _, p := range parts {
		m, ok := current.(map[string]any)
		if !ok {
			return nil
		}
		current = m[p]
	}
	return current
}

// Set updates the value at the provided path.
func (s *Status) Set(path string, val any) error {
	if path == "" {
		return errors.New("path cannot be empty")
	}
	parts := strings.Split(path, ".")
	if s.data == nil {
		s.data = map[string]any{}
	}
	m := s.data
	for i, p := range parts {
		if i == len(parts)-1 {
			m[p] = val
			return nil
		}
		next, ok := m[p].(map[string]any)
		if !ok {
			next = map[string]any{}
			m[p] = next
		}
		m = next
	}
	return nil
}

// Remove deletes the value at the provided path.
func (s *Status) Remove(path string) bool {
	parts := strings.Split(path, ".")
	m := s.data
	for i, p := range parts {
		if i == len(parts)-1 {
			if _, ok := m[p]; ok {
				delete(m, p)
				return true
			}
			return false
		}
		next, ok := m[p].(map[string]any)
		if !ok {
			return false
		}
		m = next
	}
	return false
}
