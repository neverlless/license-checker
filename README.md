license-checker
======

A Go-based tool for Node.js & PHP projects to ensure OSI-approved license compliance. It automates dependency scans, verifies licenses quickly and accurately, generates a report of all dependencies and their licenses in an HTML file, and optionally sends this report to a specified endpoint.

About
-----

license-checker is a streamlined utility tool written in Go, designed to help developers and organizations ensure open source compliance within their Node.js and PHP projects. This efficient scanner automates the process of detecting and verifying licenses against the Open Source Initiative (OSI) approved license list, generating detailed reports, and optionally sending these reports to a remote server for further processing or archiving. Whether you're maintaining a small project or managing large-scale enterprise software, license-checker provides the transparency and peace of mind needed to safeguard your project's legal integrity.

Usage
-----

To use license-checker, specify the path to your project's directory, the name of the report file (optional), and the endpoint to send the report to (optional).

`license-checker -project-dir <path to project> -report-name <report file name> -api-send-endpoint <URL> -ignore-tls true`

Example:
`license-checker -project-dir . -report-name custom_name -api-send-endpoint http://localhost:8080/api/filehandler/`

This will scan the project in the current directory, generate a report named `custom_name.html`, and send it to the specified endpoint.

For a ready-to-use API server example to handle report submissions, see [Web API File Handler](https://github.com/neverlless/web-api-filehandler).

Then, open the generated HTML file in your browser to view the results, or check the standard output for immediate results.

HTML Report Example:

![image](https://github.com/neverlless/license-checker/assets/104908866/c2e7453e-b946-4f60-a82a-31975128a8a3)

Features
--------

- **Automated Dependency Scans**: Quickly scan your project dependencies to identify the licenses used.
- **OSI-Approved License Verification**: Check for OSI compliance against the approved license list.
- **HTML Report Generation**: Generate a detailed report of all dependencies and their licenses in an HTML file. Customize the report name as needed.
- **Remote Report Submission**: Optionally send the generated report to a remote server via POST request.
- **User-Friendly CLI**: Simple and intuitive command-line interface for easy usage.
- **Cross-Platform Support**: Works on Windows, macOS, and Linux operating systems.
- **Open Source**: License-checker is an open-source project released under the MIT License.
- **Lightweight and Efficient**: Fast and efficient scanning process with minimal system resource usage.

Installation
------------

To install license-checker, you can download the latest release from the [releases page](https://github.com/neverlless/license-checker/releases) or use the following Go command:

`go get github.com/neverlless/license-checker`

Alternatively, you can clone the repository and build the binary yourself:

```shell
git clone https://github.com/neverlless/license-checker.git
cd license-checker
go build
```

Once you have the binary, you can run the `license-checker` command from the terminal.

TODO

- [ ] Add support for other languages
- [ ] Add support for custom license lists
- [x] Add support for custom output formats
- [ ] Add support for custom configuration
- [x] Add support for push html report to remote server
