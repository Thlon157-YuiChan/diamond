package io

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chihaya/diamond/internal/model"
)

// SaveProject writes a project to path as pretty JSON.
func SaveProject(path string, p *model.Project) error {
	if p == nil {
		return fmt.Errorf("nil project")
	}
	if p.Version == 0 {
		p.Version = model.CurrentSchemaVersion
	}
	data, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	return os.WriteFile(path, data, 0644)
}

// LoadProject reads a project from path.
func LoadProject(path string) (*model.Project, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var p model.Project
	if err := json.Unmarshal(data, &p); err != nil {
		return nil, err
	}
	if p.Version == 0 {
		p.Version = model.CurrentSchemaVersion
	}
	return &p, nil
}
