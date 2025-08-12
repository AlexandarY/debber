package debber

import (
	"embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

//go:embed templates/*
var debianTemplates embed.FS

// define custom functions to be used in template parsing
func CustomFunctions() template.FuncMap {
	return template.FuncMap{
		"joinStr": strings.Join,
		"add":     func(a, b int) int { return a + b },
	}
}

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

// Generate the `debian/rules` file
func (d *DebDir) CreateRules() error {
	rules, err := template.ParseFS(debianTemplates, "templates/rules.tmpl")
	if err != nil {
		return err
	}

	rulesFile, err := os.Create(fmt.Sprintf("%s/rules", d.path))
	if err != nil {
		return err
	}
	if d.data.Source.Rules == "" && d.data.Source.RawRules == "" {
		err = rules.Execute(rulesFile, struct{ Default bool }{Default: true})
		return err
	} else if d.data.Source.Rules == "" && d.data.Source.RawRules != "" {
		err = rules.Execute(rulesFile, struct {
			Default  bool
			RawRules string
		}{Default: false, RawRules: d.data.Source.RawRules})
		return err
	}
	return nil
}

// Genreate the `debian/changelog` file
func (d *DebDir) CreateChangelog() error {
	changelog, err := template.New("changelog.tmpl").Funcs(CustomFunctions()).ParseFS(debianTemplates, "templates/changelog.tmpl")
	if err != nil {
		return err
	}

	changelogFile, err := os.Create(fmt.Sprintf("%s/changelog", d.path))
	if err != nil {
		return err
	}

	err = changelog.Execute(changelogFile, struct {
		Package string
		Changes []DebianChangelogEntry
	}{
		Package: d.data.Source.Name,
		Changes: d.data.Source.Changelog,
	})
	return err
}

// Genreate the `debian/copyright` file
// TODO: Add License text to the copyright file
func (d *DebDir) CreateCopyright() error {
	copyright, err := template.New("copyright.tmpl").Funcs(CustomFunctions()).ParseFS(debianTemplates, "templates/copyright.tmpl")
	if err != nil {
		return err
	}

	copyrightFile, err := os.Create(fmt.Sprintf("%s/copyright", d.path))
	if err != nil {
		return err
	}

	err = copyright.Execute(copyrightFile, d.data.Source.Copyright)
	return err
}
