package main

import (
	"bytes"
	"crypto/tls"
	"flag"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"

	"github.com/neverlless/license-checker/detector"
	"github.com/neverlless/license-checker/licenses"
	"github.com/neverlless/license-checker/report"
	"github.com/neverlless/license-checker/scanner"
)

func main() {
	projectDirFlag := flag.String("project-dir", ".", "Path to the project directory")
	reportNameFlag := flag.String("report-name", "license_report", "Name of the report file without extension")
	apiSendEndpointFlag := flag.String("api-send-endpoint", "", "URL to send the report file via POST request (optional)")
	ignoreTLSFlag := flag.Bool("ignore-tls", false, "Ignore TLS verification for API endpoint (optional)")

	flag.Parse()

	licenses.LoadLicenses()

	projectDir := *projectDirFlag
	projectType := detector.DetectProjectType(projectDir)
	if projectType == detector.None {
		fmt.Println("Failed to detect project type. Make sure the project directory path is correct.")
		return
	}
	fmt.Printf("Defining the Project Type: %s\n", projectType)

	if err := processDependencies(projectType, projectDir, reportNameFlag, apiSendEndpointFlag, *ignoreTLSFlag); err != nil {
		fmt.Printf("Error: %s\n", err)
	}
}

func processDependencies(projectType detector.ProjectType, projectDir string, reportNameFlag, apiSendEndpointFlag *string, ignoreTLS bool) error {
	dependencies, err := scanner.ScanDependencies(string(projectType), projectDir)
	if err != nil {
		return fmt.Errorf("error while scanning dependencies: %w", err)
	}

	warnNonOSIApproved(dependencies)

	reportFileName := filepath.Join(filepath.Dir(projectDir), *reportNameFlag+".html")
	if err := generateAndSendReport(dependencies, reportFileName, apiSendEndpointFlag, ignoreTLS); err != nil {
		return err
	}

	return nil
}

func warnNonOSIApproved(dependencies []scanner.Dependency) {
	for _, dep := range dependencies {
		if !licenses.IsOSIApproved(dep.License) {
			fmt.Printf("Warning: The %s license for dependency %s is not OSI approved.\n", dep.License, dep.Name)
		}
	}
}

func generateAndSendReport(dependencies []scanner.Dependency, reportFileName string, apiSendEndpointFlag *string, ignoreTLS bool) error {
	if err := report.GenerateHTMLReport(dependencies, reportFileName); err != nil {
		return fmt.Errorf("error generating license report: %w", err)
	}
	fmt.Println("The license report has been successfully generated.")

	if *apiSendEndpointFlag != "" {
		if err := sendReport(reportFileName, *apiSendEndpointFlag, ignoreTLS); err != nil {
			return err
		}
	}

	return nil
}

func sendReport(reportFileName, apiSendEndpoint string, ignoreTLS bool) error {
	file, err := os.Open(reportFileName)
	if err != nil {
		return fmt.Errorf("error opening license report for sending: %w", err)
	}
	defer file.Close()

	var requestBody bytes.Buffer
	multipartWriter := multipart.NewWriter(&requestBody)

	fileWriter, err := multipartWriter.CreateFormFile("file", filepath.Base(reportFileName))
	if err != nil {
		return fmt.Errorf("error adding file to multipart message: %w", err)
	}
	if _, err = io.Copy(fileWriter, file); err != nil {
		return fmt.Errorf("error writing file to multipart message: %w", err)
	}

	multipartWriter.Close()

	client := &http.Client{}
	if ignoreTLS {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client = &http.Client{Transport: tr}
	}

	request, err := http.NewRequest("POST", apiSendEndpoint, &requestBody)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

	response, err := client.Do(request)
	if err != nil {
		return fmt.Errorf("error sending license report: %w", err)
	}
	defer response.Body.Close()

	fmt.Println("The license report has been successfully sent.")
	return nil
}
