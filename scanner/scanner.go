package scanner

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// Dependency represents a project dependency.
type Dependency struct {
	Name    string
	Version string
	License string
}

type packageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

type composerJSON struct {
	Require map[string]string `json:"require"`
}

// ScanDependencies scans project dependencies and returns license information.
func ScanDependencies(projectType string, dir string) ([]Dependency, error) {
	switch projectType {
	case "nodejs":
		return scanNodeJSDependencies(dir)
	case "php":
		return scanPHPDependencies(dir)
	default:
		return nil, nil
	}
}

func fetchLicenseFromNPM(packageName string) (string, error) {
	url := fmt.Sprintf("https://registry.npmjs.org/%s/latest", packageName)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var data struct {
		License string `json:"license"`
	}
	if err := json.Unmarshal(body, &data); err != nil {
		return "", err
	}

	return data.License, nil
}

func scanNodeJSDependencies(dir string) ([]Dependency, error) {
	var dependencies []Dependency
	packageJSONPath := filepath.Join(dir, "package.json")

	data, err := os.ReadFile(packageJSONPath)
	if err != nil {
		return nil, err
	}

	var pkg packageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}

	allDependencies := make(map[string]string)
	for name, version := range pkg.Dependencies {
		allDependencies[name] = version
	}
	for name, version := range pkg.DevDependencies {
		allDependencies[name] = version
	}

	for name, version := range allDependencies {
		license, err := fetchLicenseFromNPM(name)
		if err != nil {
			fmt.Printf("Failed to obtain a license for %s: %s\n", name, err)
			license = "Unknown"
		}
		dependencies = append(dependencies, Dependency{Name: name, Version: version, License: license})
	}

	return dependencies, nil
}

func fetchLicenseFromPackagist(packageName string) (string, error) {

	url := fmt.Sprintf("https://repo.packagist.org/p/%s.json", packageName)
	resp, err := http.Get(url)
	if err != nil {
		return "", fmt.Errorf("error when requesting Packagist for %s: %w", packageName, err)
	}
	defer resp.Body.Close()

	// Checking the status of the HTTP response
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Packagist returned status %d for %s", resp.StatusCode, packageName)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error when requesting Packagist for %s: %w", packageName, err)
	}

	var data struct {
		Packages map[string]map[string]struct {
			License []string `json:"license"`
		} `json:"packages"`
	}

	if err := json.Unmarshal(body, &data); err != nil {
		return "", fmt.Errorf("error when deserializing data from Packagist for %s: %w", packageName, err)
	}

	for _, versions := range data.Packages {
		for _, version := range versions {
			if len(version.License) > 0 {
				return version.License[0], nil // Returning the first license found
			}
		}
	}

	return "License not found", nil // License not found
}

func scanPHPDependencies(dir string) ([]Dependency, error) {
	var dependencies []Dependency
	composerJSONPath := filepath.Join(dir, "composer.json")

	data, err := os.ReadFile(composerJSONPath)
	if err != nil {
		return nil, err
	}

	var comp composerJSON
	if err := json.Unmarshal(data, &comp); err != nil {
		return nil, err
	}

	for name, version := range comp.Require {
		license, err := fetchLicenseFromPackagist(name)
		if err != nil {
			fmt.Printf("Failed to obtain a license for %s: %s\n", name, err)
			license = "Unknown"
		}
		dependencies = append(dependencies, Dependency{Name: name, Version: version, License: license})
	}

	return dependencies, nil
}
