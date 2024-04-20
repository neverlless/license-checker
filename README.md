license-checker
======

A Go-based tool for Node.js & PHP projects to ensure OSI-approved license compliance. It automates dependency scans and verifies licenses quickly and accurately and generates a report of all dependencies and their licenses in html file.

About
-----

license-checker is a streamlined utility tool written in Go, designed to help developers and organizations ensure open source compliance within their Node.js and PHP projects. This efficient scanner automates the process of detecting and verifying licenses against the Open Source Initiative (OSI) approved license list. With license-checker, you can quickly scan your project dependencies, identify the licenses used, and check for OSI compliance, all with a simple and user-friendly command-line interface. Whether you're maintaining a small project or managing large-scale enterprise software, license-checker provides the transparency and peace of mind needed to safeguard your project's legal integrity.

Usage
-----

`license-checker <path to project>`

Example:
`license-checker .`

Then open the generated `license-report.html` file in your browser to view the results or just check stdout for the results.

HTML Report Example:

![image](https://github.com/neverlless/license-checker/assets/104908866/c2e7453e-b946-4f60-a82a-31975128a8a3)

Features
--------

- **Automated Dependency Scans**: Quickly scan your project dependencies to identify the licenses used.
- **OSI-Approved License Verification**: Check for OSI compliance against the approved license list.
- **HTML Report Generation**: Generate a detailed report of all dependencies and their licenses in an HTML file.
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
git clone
cd license-checker
go build
```

Once you have the binary, you can run the `license-checker` command from the terminal.

TODO

- [ ] Add support for other languages
- [ ] Add support for custom license lists
- [ ] Add support for custom output formats
- [ ] Add support for custom configuration
- [ ] Add support for push html report to remote server
