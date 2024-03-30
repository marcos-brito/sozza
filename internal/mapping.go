package internal

import (
	"fmt"
	"io"
	"os"

	"gopkg.in/yaml.v3"
)

type Mapping map[string][]map[string]string

func ReadMappingFromFile(path string) (*Mapping, error) {
	file, err := os.Open(path)

	if err != nil {
		return nil, fmt.Errorf("Could not read the mapping file: %s", err)
	}

	content, err := io.ReadAll(file)

	if err != nil {
		return nil, err
	}

	mapping, err := ReadMapping(content)
	if err != nil {
		return nil, err
	}

	return mapping, nil
}

func ReadMapping(content []byte) (*Mapping, error) {
	m := &Mapping{}
	err := yaml.Unmarshal(content, &m)

	if err != nil {
		return nil, fmt.Errorf("Could not unmarshal the mapping: %s", err)
	}

	return m, nil
}
