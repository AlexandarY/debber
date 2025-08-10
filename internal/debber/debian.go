package debber

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*
var debianTemplates embed.FS

// Represents the debian directory
type DebDir struct {
	path string
	data *DebianFile
}

// Initialize a new debian directory
func CreateDebianDirectory(path string, data *DebianFile) (*DebDir, error) {
	debDirPath := filepath.Join(path, "debian")
	err := os.Mkdir(debDirPath, 0755)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return nil, fmt.Errorf("The directory '%s' already exists!", debDirPath)
		} else {
			return nil, err
		}
	}

	return &DebDir{path: debDirPath, data: data}, nil
}

// Generate the `debian/control` file
func (d *DebDir) CreateControl() error {
	control, err := template.ParseFS(debianTemplates, "templates/control.tmpl")
	if err != nil {
		return err
	}

	controlFile, err := os.Create(fmt.Sprintf("%s/control", d.path))
	if err != nil {
		return err
	}

	err = control.Execute(controlFile, d.data)
	if err != nil {
		return err
	}

	return nil
}
