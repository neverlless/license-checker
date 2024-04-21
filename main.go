package main

import (
	"bytes"
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

	flag.Parse()

	licenses.LoadLicenses()

	projectDir := *projectDirFlag
	projectType := detector.DetectProjectType(projectDir)
	if projectType == detector.None {
		fmt.Println("Failed to detect project type. Make sure the project directory path is correct.")
		return
	}
	fmt.Printf("Defining the Project Type: %s\n", projectType)

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

	reportFileName := filepath.Join(filepath.Dir(projectDir), *reportNameFlag+".html")
	err = report.GenerateHTMLReport(dependencies, reportFileName)
	if err != nil {
		fmt.Printf("Error generating license report: %s\n", err)
		return
	}
	fmt.Println("The license report has been successfully generated.")

	if *apiSendEndpointFlag != "" {
		file, err := os.Open(reportFileName)
		if err != nil {
			fmt.Printf("Error opening license report for sending: %s\n", err)
			return
		}
		defer file.Close()

		var requestBody bytes.Buffer
		multipartWriter := multipart.NewWriter(&requestBody)

		fileWriter, err := multipartWriter.CreateFormFile("file", filepath.Base(reportFileName))
		if err != nil {
			fmt.Printf("Error adding file to multipart message: %s\n", err)
			return
		}
		_, err = io.Copy(fileWriter, file)
		if err != nil {
			fmt.Printf("Error writing file to multipart message: %s\n", err)
			return
		}

		multipartWriter.Close()

		request, err := http.NewRequest("POST", *apiSendEndpointFlag, &requestBody)
		if err != nil {
			fmt.Printf("Error creating request: %s\n", err)
			return
		}
		request.Header.Set("Content-Type", multipartWriter.FormDataContentType())

		response, err := http.DefaultClient.Do(request)
		if err != nil {
			fmt.Printf("Error sending license report: %s\n", err)
			return
		}
		defer response.Body.Close()

		fmt.Println("The license report has been successfully sent.")
	}
}
