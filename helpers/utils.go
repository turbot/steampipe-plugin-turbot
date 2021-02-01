package helpers

import (
	"fmt"
	"github.com/go-yaml/yaml"
)

// parse given string in YAML format
func ParseYamlString(value string) (interface{}, error) {
	var result interface{}

	err := yaml.Unmarshal([]byte(value), &result)
	// returns value when unmarshal fails
	if err != nil {
		return value, err
	}

	if result == "" {
		return value, nil
	}

	return result, nil
}

// convert value to a string.
func InterfaceToString(value interface{}) string {
	if value == nil {
		return ""
	}
	return fmt.Sprintf("%v", value)
}

// if the value is already a string return it, otherwise convert to the YAML representation
func InterfaceToStringOrYaml(value interface{}) (string, error) {
	if value == nil {
		return "", nil
	}
	if res, ok := value.(string); ok {
		return res, nil
	}

	data, err := yaml.Marshal(value)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// implements a equal operation on 2 YAML strings, ignoring formatting differences
func YamlStringsAreEqual(yaml1, yaml2 string) (bool, error) {
	var yaml1intermediate, yaml2intermediate interface{}

	if err := yaml.Unmarshal([]byte(yaml1), &yaml1intermediate); err != nil {
		return false, fmt.Errorf("Error unmarshaling yaml string: %s", err)
	}

	if err := yaml.Unmarshal([]byte(yaml2), &yaml2intermediate); err != nil {
		return false, fmt.Errorf("Error unmarshaling yaml string: %s", err)
	}

	s1, err := yaml.Marshal(yaml1intermediate)
	if err != nil {
		return false, fmt.Errorf("Error marshaling yaml string: %s", err)
	}

	s2, err := yaml.Marshal(yaml2intermediate)
	if err != nil {
		return false, fmt.Errorf("Error marshaling yaml string: %s", err)
	}

	if string(s1) == string(s2) {
		return true, nil
	}
	return false, nil
}
