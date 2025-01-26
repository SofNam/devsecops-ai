# DevSecOps AI Scanner

An AI-powered security scanner that combines traditional static analysis with artificial intelligence to detect potential security vulnerabilities in your codebase. This tool is designed to enhance your DevSecOps pipeline by providing intelligent security analysis and detailed reporting.

## Features

- AI-Enhanced Security Analysis
  - Pattern recognition using machine learning
  - Context-aware vulnerability detection
  - Automated severity assessment
  - Smart false positive reduction

- Comprehensive Security Scanning
  - Static code analysis
  - Configuration file analysis
  - Dependency checking
  - Custom rule support

- Advanced Reporting
  - Multiple output formats (JSON/HTML)
  - Detailed vulnerability descriptions
  - Code snippets with context
  - Actionable remediation suggestions
  - Statistical summaries

- Security-First Design
  - Secure container execution
  - Read-only filesystem access
  - Minimal container privileges
  - Seccomp profile integration

## Installation

### Prerequisites

- Docker
- Go 1.x (for development)
- Git

### Quick Start

1. Clone the repository:
   ```bash
   git clone https://github.com/SofNam/devsecops-ai.git
   cd devsecops-ai
   ```

2. Start the scanner:
   ```bash
   ./start.sh start
   ```

## Usage

The scanner can be run using the provided `start.sh` script or directly using the compiled binary.

### Using start.sh

```bash
# Run security scan
./start.sh start --path /path/to/code

# Run tests
./start.sh test all

# Show help
./start.sh help
```

### Command Line Options

```bash
Usage: ./start.sh [command] [options]

Commands:
  start     - Start the security scanner
  test      - Run tests
  help      - Show this help message

Test Options:
  all       - Run all tests
  coverage  - Run tests with coverage
  integration - Run integration tests
  version   - Run version package tests
  reporter  - Run reporter package tests
  benchmark - Run benchmark tests
```

### Direct Binary Usage

```bash
./scanner --path /path/to/scan \
          --model /path/to/model \
          --output json \
          --output-path report
```

## Configuration

### Scanner Configuration

The scanner can be configured using environment variables or configuration files:

```yaml
SCANNER_LOG_LEVEL=info
SCANNER_OUTPUT_FORMAT=json
```

### Security Rules

Custom security rules can be defined in `rules.json`:

```json
{
  "id": "RULE-001",
  "name": "SQL Injection Detection",
  "pattern": "(?i)(SELECT|INSERT|UPDATE|DELETE).*\\$\\{",
  "severity": "HIGH",
  "category": "Injection",
  "description": "Potential SQL injection vulnerability detected"
}
```

### Docker Security Settings

The scanner runs with enhanced security settings:
- Read-only container filesystem
- No privilege escalation
- Minimal capabilities
- Seccomp profile restrictions
- Temporary filesystem for required writes

## Development

### Project Structure

```
├── pkg/
│   ├── ai/         - AI detection logic
│   ├── models/     - Data structures
│   ├── reporter/   - Report generation
│   ├── scanner/    - Core scanning logic
│   └── version/    - Version information
├── docker-compose.yml
├── Dockerfile
├── seccomp-profile.json
└── start.sh
```

### Running Tests

```bash
# Run all tests
./start.sh test all

# Run with coverage
./start.sh test coverage

# Run benchmarks
./start.sh test benchmark
```

## Security Considerations

- The scanner runs in a containerized environment with restricted privileges
- All file operations are performed in read-only mode
- Network access is limited to required services
- System calls are restricted using seccomp profiles
- No persistent storage outside designated volumes
