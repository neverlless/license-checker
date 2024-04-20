package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/neverlless/license-checker/detector"
	"github.com/neverlless/license-checker/licenses"
	"github.com/neverlless/license-checker/report"
	"github.com/neverlless/license-checker/scanner"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: license-checker <path to project>")
		return
	}

	projectDir := os.Args[1]

	// Downloading license data
	licenses.LoadLicenses()

	// Определение типа проекта
	projectType := detector.DetectProjectType(projectDir)
	if projectType == detector.None {
		fmt.Println("Failed to detect project type. Make sure the project directory path is correct.")
		return
	}
	fmt.Printf("Defining the Project Type: %s\n", projectType)

	// Dependency scanning
	dependencies, err := scanner.ScanDependencies(string(projectType), projectDir)
	if err != nil {
		fmt.Printf("Error while scanning dependencies: %s\n", err)
		return
	}

	for _, dep := range dependencies {
		if !licenses.IsOSIApproved(dep.License) {
			fmt.Printf("Warning: The %s license for dependency %s is not OSI approved.\n", dep.License, dep.Name)
		}
	}

	// Generating a license report
	err = report.GenerateHTMLReport(dependencies, filepath.Dir(projectDir))
	if err != nil {
		fmt.Printf("Error generating license report: %s\n", err)
		return
	}

	fmt.Println("The license report has been successfully generated.")
}
