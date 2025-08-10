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
	Name             string `toml:"name"`
	Maintainer       string `toml:"maintainer"`
	Section          string `toml:"section"`
	Priority         string `toml:"priority"`
	StandardsVersion string `toml:"standards-version"`
	Rules            string `toml:"rules"`
}

// Represents the information about a binary package that is build from the source
type DebianPackage struct {
	Name        string `toml:"name"`
	Arch        string `toml:"arch"`
	Section     string `toml:"section"`
	Priority    string `toml:"priority"`
	Description string `toml:"description"`
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
	if d.Source.Rules == "" {
		return errors.New("source.rules cannot be missing or an empty string!")
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
