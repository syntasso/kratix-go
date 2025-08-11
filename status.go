package kratix

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/itchyny/gojq"
)

// TODO: Do we want a GetConditions() in the Status interface object?

type StatusModifier interface {
	// Get queries the Status and retrieves the value at the specified path e.g. healthStatus.state
	Get(string) any
	// Set updates the value at the specified path e.g. healthStatus.state
	Set(string, any) error
	// Set removes the value at the specified path e.g. healthStatus.state
	Remove(string) error
}

// Status implements StatusModifier using a generic map.
type Status struct {
	data map[string]any
}

var _ StatusModifier = (*Status)(nil)

// Get retrieves the value at the provided path.
// It can be used to execute a jq-like query on the Status data and returns the results
// Examples:
//   - ".pods[].name" -> returns all pod names
//   - ".pods[] | select(.status == \"Running\")" -> returns all running pods
//   - ".pods[].containers[] | select(.ready == true)" -> returns all ready containers
//   - ".pods | length" -> returns the number of pods
func (s *Status) Get(path string) any {
	results, err := s.query("get", path, nil)
	if err != nil || len(results) == 0 {
		return nil
	}
	return results[0]
}

// Set updates the value at the provided path.
// It accepts jq-like paths, like ".pods[].name" or ".pods[] | select(.status == \"Running\")"
func (s *Status) Set(path string, val any) error {
	_, err := s.query("set", path, val)
	return err
}

// Remove deletes the value at the provided path.
// It accepts jq-like paths, like ".pods[].name" or ".pods[] | select(.status == \"Running\")"
func (s *Status) Remove(path string) error {
	_, err := s.query("remove", path, nil)
	return err
}

func normalisePath(path string) (string, error) {
	if path == "" {
		return "", errors.New("path cannot be empty")
	}

	if !strings.HasPrefix(path, ".") {
		path = "." + path
	}

	return path, nil
}

func buildQuery(op, path string, val any) (string, bool, error) {
	var query string
	persist := true
	switch op {
	case "get":
		query = fmt.Sprintf(`%s`, path)
		persist = false
	case "set":
		jsonObj, err := json.Marshal(val)
		if err != nil {
			return "", false, err
		}
		query = fmt.Sprintf(`%s = %s`, path, string(jsonObj))
	case "remove":
		query = fmt.Sprintf(`del(%s)`, path)
	default:
		return "", false, fmt.Errorf("invalid operation: %s", op)
	}

	return query, persist, nil
}

func (s *Status) query(op, path string, val any) ([]any, error) {
	var err error
	if path, err = normalisePath(path); err != nil {
		return nil, err
	}

	query, persist, err := buildQuery(op, path, val)
	if err != nil {
		return nil, err
	}

	jqQuery, err := gojq.Parse(query)

	if err != nil {
		return nil, err
	}

	var results []any
	iter := jqQuery.Run(s.data)
	for {
		v, ok := iter.Next()
		if !ok {
			break
		}
		if err, ok := v.(error); ok {
			return nil, err
		}
		results = append(results, v)
	}

	if persist {
		s.data = results[0].(map[string]any)
	}

	return results, nil
}
