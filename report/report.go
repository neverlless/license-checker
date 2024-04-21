package report

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/neverlless/license-checker/licenses"
	"github.com/neverlless/license-checker/scanner"
)

const reportTemplate = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>License Report</title>
    <style>
        body {
            font-family: 'Segoe UI', 'Arial', sans-serif;
            color: #333;
            background-color: #f4f4f4;
            margin: 0;
            padding: 20px;
        }
        h1 {
            color: #007acc;
        }
        ul {
            list-style-type: none;
            padding: 0;
        }
        li {
            background-color: #fff;
            margin-bottom: 10px;
            padding: 10px;
            border-radius: 5px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
        }
        li:hover {
            background-color: #f9f9f9;
        }
        .osi-approved {
            color: #28a745;
            font-weight: bold;
        }
        .not-osi-approved, .license-unknown {
            color: #dc3545;
            font-weight: bold;
        }
    </style>
</head>
<body>
    <h1>License Report</h1>
    <ul>
        {{range .}}
        <li>
            {{cleanName .Name}} - {{.Version}} - 
            {{if eq .License "Unknown"}}
                <span class="license-unknown">License Unknown</span>
            {{else if isOSIApproved .License}}
                <span class="osi-approved">{{.License}} (OSI Approved)</span>
            {{else}}
                <span class="not-osi-approved">{{.License}} (Not OSI Approved)</span>
            {{end}}
        </li>
        {{end}}
    </ul>
</body>
</html>
`

func GenerateHTMLReport(dependencies []scanner.Dependency, reportFilePath string) error {
	funcMap := template.FuncMap{
		"cleanName": func(name string) string {
			return strings.Replace(name, "@", "", -1)
		},
		"isOSIApproved": func(license string) bool {
			return licenses.IsOSIApproved(license)
		},
	}

	tmpl, err := template.New("report").Funcs(funcMap).Parse(reportTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(reportFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := tmpl.Execute(file, dependencies); err != nil {
		return err
	}

	fmt.Printf("License report generated: %s\n", reportFilePath)
	return nil
}
