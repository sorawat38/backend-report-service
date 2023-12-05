# Ice Cream Shop - Report service

[![License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

Welcome to the Report Service, a versatile tool designed to facilitate report generation and management for your business. This service provides an efficient and customizable solution for creating, storing, and accessing reports in various formats.

## Features

- **Generate Analysis Report:** Generate report for business analysis.

## Table of Contents

- [Ice Cream Shop - Report service](#ice-cream-shop---report-service)
  - [Features](#features)
  - [Table of Contents](#table-of-contents)
  - [Getting Started](#getting-started)
    - [Prerequisites](#prerequisites)
    - [Installation](#installation)
    - [Configuration](#configuration)
    - [API Documentation](#api-documentation)
    - [License](#license)

## Getting Started

Follow these instructions to get a copy of the project up and running on your local machine for development and testing purposes.

### Prerequisites

- Go 1.16+ installed
- A MySQL database

### Installation

1. Clone the repository:

   ```bash
   git clone https://github.com/CLCM3102-Ice-Cream-Shop/backend-report-service.git
   
   ```
2. Go to project directory
   ```bash
   cd backend-report-service
   ```
3. Install dependencies:
   ```bash
   go mod tidy
   ```
4. Run proejct
   ```bash
   make run
   ```
### Configuration
- Copy the example configuration file and customize it according to your environment:

    ```bash
    cp config.example.yaml config.yaml
    ```
- Update the config.yaml file with your database and other settings.

### API Documentation
The API documentation is available [here](https://satrawo38.atlassian.net/wiki/spaces/CP/pages/4555062/API+Specification).

### License
This project is licensed under the MIT License - see the [LICENSE](LICENSE.md) file for details.

