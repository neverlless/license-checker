package detector

import (
	"os"
)

// ProjectType represents the project type.
type ProjectType string

const (
	NodeJS ProjectType = "nodejs"
	PHP    ProjectType = "php"
	None   ProjectType = "none"
)

// DetectProjectType determines the type of project in the specified directory.
func DetectProjectType(dir string) ProjectType {
	if _, err := os.Stat(dir + "/package.json"); err == nil {
		return NodeJS
	}
	if _, err := os.Stat(dir + "/composer.json"); err == nil {
		return PHP
	}
	return None
}
