package debber

import (
	"errors"
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

// Represents the config file that represents the debian/ directory content
type DebianFile struct {
	Source   DebianSource    `toml:"source"`
	Packages []DebianPackage `toml:"packages"`
}

// Represents the source information that is used for the build
type DebianSource struct {
	Name             string                 `toml:"name"`
	Maintainer       string                 `toml:"maintainer"`
	Section          string                 `toml:"section"`
	Priority         string                 `toml:"priority"`
	StandardsVersion string                 `toml:"standards-version"`
	Rules            string                 `toml:"rules"`
	RawRules         string                 `toml:"raw_rules"`
	Changelog        []DebianChangelogEntry `toml:"changelog"`
	Copyright        DebianCopyright        `toml:"copyright"`
}

// Represents the information about a binary package that is build from the source
type DebianPackage struct {
	Name        string `toml:"name"`
	Arch        string `toml:"arch"`
	Section     string `toml:"section"`
	Priority    string `toml:"priority"`
	Description string `toml:"description"`
}

// Represents the data for `debian/changelog`
type DebianChangelogEntry struct {
	Version      string   `toml:"version"`
	Distribution []string `toml:"distribution"`
	Urgency      string   `toml:"urgency"`
	Changes      []string `toml:"changes"`
	ChangedBy    string   `toml:"changed_by"`
	Date         string   `toml:"date"`
}

// Represents the data for `debian/copyright`
type DebianCopyright struct {
	Format          string                 `toml:"format"`
	Source          string                 `toml:"source"`
	UpstreamName    string                 `toml:"upstream_name"`
	UpstreamContact string                 `toml:"upstream_contact"`
	Files           []DebianCopyrightFiles `toml:"files"`
}

// Represents the data for the files stanza in `debian/copyright`
type DebianCopyrightFiles struct {
	Files     string   `toml:"files"`
	Copyright []string `toml:"copyright"`
	License   string   `toml:"license"`
}

// Validate that the required fields are provided
func (d *DebianFile) Validate() error {
	if d.Source.Name == "" {
		return errors.New("source.name cannot be missing or an empty string!")
	}
	if d.Source.Maintainer == "" {
		return errors.New("source.maintainer cannot be missing or an empty string!")
	}
	if d.Source.Section == "" {
		return errors.New("source.section cannot be missing or an empty string!")
	}
	if d.Source.Priority == "" {
		return errors.New("source.priority cannot be missing or an empty string!")
	}
	if d.Source.StandardsVersion == "" {
		return errors.New("source.standards-version cannot be missing or an empty string!")
	}
	if d.Source.Rules != "" && d.Source.RawRules != "" {
		return errors.New("you cannot have both source.rules and source.raw_rules defined!")
	}
	if len(d.Source.Changelog) == 0 {
		return errors.New("At least one changelog entry is needed! None were provided in [[source.changelog]]")
	} else {
		for idx, change := range d.Source.Changelog {
			if change.Version == "" {
				return fmt.Errorf("No `version` was provided for change with index %d", idx)
			}
			if len(change.Distribution) == 0 {
				return fmt.Errorf("No `distribution` was provided for change with index %d and version %s", idx, change.Version)
			}
			if change.Urgency == "" {
				return fmt.Errorf("No `urgency` was provided for change with index %d and version %s", idx, change.Version)
			}
			if len(change.Changes) == 0 {
				return fmt.Errorf("No `changes` was provided for change with index %d and version %s", idx, change.Version)
			}
			if change.Date == "" {
				return fmt.Errorf("No `date` was provided for change with index %d and version %s", idx, change.Version)
			}
		}
	}
	if len(d.Source.Copyright.Files) == 0 {
		return errors.New("At least one entry is required in `source.copyright.files`!")
	} else {
		for idx, fileStanza := range d.Source.Copyright.Files {
			if fileStanza.Files == "" {
				return fmt.Errorf("Expected `files` pattern for `source.copyright.files` with index %d", idx)
			}
			if len(fileStanza.Copyright) == 0 {
				return fmt.Errorf("Expected at least one `copyright` entry for `source.copyright.files` with Files '%s' and index %d", fileStanza.Files, idx)
			}
			if fileStanza.License == "" {
				return fmt.Errorf("Expected `license` entry for `source.copyright.files` with Files '%s' and index %d", fileStanza.Files, idx)
			}
		}
	}

	if len(d.Packages) == 0 {
		return errors.New("At least one binary package needs to be built from source! Non were defined in [[packages]]")
	} else {
		for idx, pkg := range d.Packages {
			if pkg.Name == "" {
				return fmt.Errorf("Package %d is missing 'name' or it is set to an empty string", idx)
			}
			if pkg.Arch == "" {
				return fmt.Errorf("Package %d is missing 'name' or it is set to an empty string", idx)
			}
			if pkg.Section == "" {
				return fmt.Errorf("Package %d is missing 'name' or it is set to an empty string", idx)
			}
			if pkg.Priority == "" {
				return fmt.Errorf("Package %d is missing 'name' or it is set to an empty string", idx)
			}
			if pkg.Description == "" {
				return fmt.Errorf("Package %d is missing 'name' or it is set to an empty string", idx)
			}
		}
	}

	return nil
}

// Parse the contents of the debian config file
func ParseFile(debianFile string) (*DebianFile, error) {
	rawData, err := os.ReadFile(debianFile)
	if err != nil {
		return nil, err
	}

	var content DebianFile
	err = toml.Unmarshal(rawData, &content)
	if err != nil {
		var derr *toml.DecodeError
		if errors.As(err, &derr) {
			return nil, fmt.Errorf("\n%s", derr.String())
		}
		return nil, err
	}

	if err := content.Validate(); err != nil {
		return nil, err
	}

	return &content, nil
}

// Method handles the creation of a new debian package file
func CreateNewDebFile(fileName string) error {
	_, err := os.Stat(fileName)
	if err == nil {
		return fmt.Errorf("%s already exists! Either specify a new name via `-f` or remove existing file", fileName)
	}

	var newData DebianFile
	newData.Packages = append(newData.Packages, DebianPackage{})
	newData.Source.Changelog = append(newData.Source.Changelog, DebianChangelogEntry{})
	newData.Source.Copyright = DebianCopyright{
		Format:          "https://www.debian.org/doc/packaging-manuals/copyright-format/1.0/",
		Source:          "<url://example.com>",
		UpstreamName:    "<project-name>",
		UpstreamContact: "<preferred name and address to reach the upstream project>",
		Files: []DebianCopyrightFiles{
			{
				Files:     "*",
				Copyright: []string{"<years> <author's name here>"},
				License:   "<license name here>",
			},
			{
				Files:     "debian/*",
				Copyright: []string{"<years> <author's name here>"},
				License:   "GPL-2+",
			},
		},
	}
	content, err := toml.Marshal(&newData)
	if err != nil {
		return err
	}

	headerContent := fmt.Sprintf("# Generated by debber\n%s", string(content))
	err = os.WriteFile(fileName, []byte(headerContent), 0644)
	if err != nil {
		return err
	}

	return nil
}
