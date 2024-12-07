# <img src="https://bruce.tools/images/logo.png" alt="bruce logo" width="32"/> Bruce

Basic Runtime for Uniform Compute Environments

---
Bruce is a lightweight, single-binary tool designed to configure and install operating system packages and settings in a reproducible and uniform way. It operates without the heavy dependencies required by tools like Ansible or Chef, making it ideal for quickly setting up fleets of servers or bootstrapping instances in environments like AWS EC2.

[![Latest Release](https://img.shields.io/github/release/brucehq/bruce.svg)](https://github.com/brucehq/bruce/releases/latest)
[![License](https://img.shields.io/github/license/brucehq/bruce.svg)](LICENSE)
[![Go Report Card](https://goreportcard.com/badge/github.com/brucehq/bruce)](https://goreportcard.com/report/github.com/brucehq/bruce)
![Linux](https://img.shields.io/badge/Linux-amd64%20%7C%20arm64-772953?logo=linux&logoColor=white)
![macOS](https://img.shields.io/badge/macOS-amd64%20%7C%20arm64-0c77e3?logo=apple&logoColor=white)
![Windows](https://img.shields.io/badge/Windows-amd64%20%7C%20arm64-0078D6?logo=windows&logoColor=white)

The repository is also supported and backed by Runtime Dynamics LLC. Which is primarily focused on enterprise level support and services for Bruce.  The Bruce project is a community driven project that is supported by the community and the enterprise level support is provided by Runtime Dynamics LLC.

## Table of Contents
- [Features](#features)
- [Getting Started](#getting-started)
    - [Documentation](#documentation)
    - [Installation](#installation)
- [Quick Start](#quick-start)
    - [Prepare a Configuration File](#prepare-a-configuration-file)
    - [Run Bruce](#run-bruce)
    - [Load Configuration from Different Sources](#load-configuration-from-different-sources)
- [Usage](#usage)
- [Configuration](#configuration)
- [Operators](#operators)
- [Contributing](#contributing)
- [Contact](#contact)
- [Principles and Context](#principles-and-context)
- [Credits](#credits)c

## Features
* **Zero Dependencies**: No additional OS dependencies; works even on minimal installations.
* **Single Binary**: Easy distribution and deployment as a standalone binary.
* **Multi-Platform Support**: Compatible with Linux, macOS, and basic Windows support.
* **Fast Execution**: Configures entire systems rapidly, outperforming similar tools on resource-constrained instances.
* **Flexible Configuration**: Supports native commands with OS limiters, services management, package installations, file ownership settings, and more.
* **Template Engine**: Inject variables and environment variables into templates, enabling dynamic configurations.
* **Secure Loaders**: Load configuration files from multiple sources like local files, S3 buckets, or HTTP URLs.
* **Conditional Execution**: Execute or exclude commands based on conditions.
* **Service Management**: Restart services only on change detection to minimize downtime.

## Getting Started

### Documentation
For the full documentation on Bruce, visit the [Bruce Documentation](https://docs.bruce.tools/).

### Installation

Download the [latest release](https://github.com/brucehq/bruce/releases) for your operating system from the [Releases Page](https://github.com/brucehq/bruce/releases).

#### One-liner Installation for Linux AMD64
Note: This command downloads the latest release for Linux AMD64, and uses sudo to extract the tarball to /usr/local/bin. Ensure you have the necessary permissions, you can skip the last part and manually move the binary to your desired location.
```bash
wget -qO- $(curl -s https://api.github.com/repos/brucehq/bruce/releases/latest | grep "linux_amd64" | grep https | cut -d : -f 2,3 | tr -d \" | awk '{$1=$1};1') | sudo tar -xz -C /usr/local/bin
```

## Quick Start
### Prepare a Configuration File

Create an `install.yml` file. You can use the example configuration as a starting point. Here's a basic example:

```yaml
---
variables:
  Person: "Steven"
steps:
- cmd: echo "{{.Person}} is using Bruce"
  setEnv: Person
- template: ./output2.txt
  source: https://raw.githubusercontent.com/brucehq/bruce/refs/heads/main/template-example.txt
- api: https://postman-echo.com/get?foo1=bar1&foo2=bar2
  jsonKey: headers.host
  setEnv: apiResponse
- cmd: echo {{.apiResponse}}
```
This is an extremely rudimentary example, that uses a "global variable Person" to echo a message, then uses a template to create a file, then uses an API call to get a response and echo it.

### Run Bruce

```bash
./bruce ./install.yml
```
The output on running that would look like this:
```shell
11:30AM INF cmd: echo "Steven is using Bruce"
11:30AM INF template: https://raw.githubusercontent.com/brucehq/bruce/refs/heads/main/template-example.txt => ./output2.txt
11:30AM INF template written: ./output2.txt
11:30AM INF no backup file for ./output2.txt
11:30AM INF API request: GET https://postman-echo.com/get?foo1=bar1&foo2=bar2
11:30AM INF cmd: echo postman-echo.com
```
And the output of the output2.txt file would be:
```shell
There be a person named: Steven is using Bruce here.
(string) (len=21) "Steven is using Bruce"
```
### Load Configuration from Different Sources
Alternatively, load the configuration from an S3 bucket or an HTTP URL:

#### From S3

```bash
./bruce s3://your-bucket/install.yml
```
Note: This requires that you have set up your AWS credentials appropriately first.

#### From HTTP

```bash
./bruce https://your.domain/$(hostname -f).yml
```

## Usage
Bruce operates based on a configuration file that defines the desired state of your system. It supports various operators to perform tasks like running command line, managing services, and templating files.

## Configuration
The configuration file is written in YAML format. Below is a basic example:

```yaml
---
variables:
# note do not set Option here as it will be overwritten
steps:
  - cmd: echo "Secondary test command received"
```
The yaml file is divided into 2 primary sections, first and foremost we start with the variables section which is optional, but allows you to set variables that can be used in the steps section.  The steps section is where the magic happens, it is a list of commands that Bruce will execute in order.  Each command is a dictionary with a key of the operator you want to use, and the value is the configuration for that operator.

## Operators
Bruce supports a variety of operators within the configuration file:

| Operator            | Description                                                                                                                                                                      |
|---------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| **API**             | Make HTTP(s) requests and set the response as an environment variable, or dump it to a file to be used later.  Has built in capabilities to parse json for dot notation of keys. |
| **Command**         | Execute native commands with arguments.  Optional capabilities for running `onlyIF` or `notIf` conditions. Output can also be set as envars for use later in the configuration.  |
| **Copy**            |  Copy provides a means to copy files from one location to another, sources could include http(s) or s3 or local files on the same box.                                           |
| **Cron**            | Enables the ability to add or remove cron jobs from the system, specifically for *nix environmetns.|
| **Git**             | Clone or pull a git repository to a specified location, does not require git to be installed on the system.|
| **Loop**            | Allows for looping over a list of commands, useful for iterating over a list of items. Example use case would be installing kafka and providing a different ID for each instance.|
| **RecursiveCopy**   | Copy files from one location to another, recursively.  Sources could include http(s) or s3 or local files on the same box, this uses concurrency to allow for multiple files to be processed at the same time.|
| **RemoteExecution** | Execute commands on remote hosts via ssh, this allows you to execute a remote command and also set the output as an environment variable.|
| **Signals**         | Send signals like SIGINT or SIGHUP to running processes.|
| **Tarball**         | Extract a tarball to a specified location, removing the tarball after extraction, without having to have tar installed on the system.|
| **Template**       | Render templates with injected variables and environment variables, which allows for stable system configurations across multiple environments.|


## Contributing
We welcome contributions from the community to enhance Bruce's functionality, especially in extending Windows support.

### To contribute:

* Fork the repository.
* Create a new branch for your feature or bugfix.
* Commit your changes with clear messages.
* Submit a pull request to the main branch.

## Contact
For any questions or suggestions, feel free to open an issue or contact the maintainer:

Website: [bruce.tools](https://bruce.tools)
Email: [support@bruce.tools](mailto:support@bruce.tools)

Expand your capabilities by using Bruce with an event driven backend and build advanced automation workflows, on [bruce.tools](https://bruce.tools)

## Principles and Context
Bruce is built with the principles of simplicity and efficiency, aiming to make system configuration as straightforward as possible. Originally designed for machine learning applications, Bruce facilitates the transition from ML training to hosting in environments lacking dedicated and advanced operations capabilities. It consequently expanded to handle distributed execution across multiple agents via its advanced backend features on https://bruce.tools

## Credits
https://postman-echo.com/ - For providing a free API endpoint for testing purposes.