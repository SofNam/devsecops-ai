#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print with color
echo_color() {
    color=$1
    message=$2
    echo -e "${color}${message}${NC}"
}

# Check if Docker is installed and running
check_docker() {
    if ! command -v docker &> /dev/null; then
        echo_color $RED "Docker is not installed. Please install Docker first."
        exit 1
    fi

    if ! docker info &> /dev/null; then
        echo_color $RED "Docker daemon is not running. Please start Docker."
        exit 1
    fi
}

# Build the Docker image
build_image() {
    echo_color $YELLOW "Building Docker image..."
    docker build -t securityai-scanner . || {
        echo_color $RED "Failed to build Docker image"
        exit 1
    }
}

# Run tests
run_tests() {
    local test_type=$1
    echo_color $YELLOW "Running tests..."

    case $test_type in
        "all")
            echo_color $YELLOW "Running all tests..."
            go test ./... -v
            ;;
        "coverage")
            echo_color $YELLOW "Running tests with coverage..."
            go test ./... -v -cover -coverprofile=coverage.out
            go tool cover -html=coverage.out
            ;;
        "integration")
            echo_color $YELLOW "Running integration tests..."
            go test ./test -v
            ;;
        "version")
            echo_color $YELLOW "Running version package tests..."
            go test ./pkg/version -v
            ;;
        "reporter")
            echo_color $YELLOW "Running reporter package tests..."
            go test ./pkg/reporter -v
            ;;
        "benchmark")
            echo_color $YELLOW "Running benchmark tests..."
            go test ./... -bench=.
            ;;
        *)
            echo_color $RED "Unknown test type: $test_type"
            echo "Available test types:"
            echo "  all       - Run all tests"
            echo "  coverage  - Run tests with coverage"
            echo "  integration - Run integration tests"
            echo "  version   - Run version package tests"
            echo "  reporter  - Run reporter package tests"
            echo "  benchmark - Run benchmark tests"
            exit 1
            ;;
    esac
}

# Start the application
start_app() {
    echo_color $YELLOW "Starting Security Scanner..."

    # Create logs directory if it doesn't exist
    mkdir -p logs

    # Run the container with security options
    docker run \
        --name security-scanner \
        --security-opt seccomp=seccomp-profile.json \
        --security-opt no-new-privileges \
        --cap-drop ALL \
        --cap-add NET_BIND_SERVICE \
        --read-only \
        --tmpfs /tmp \
        -v "$(pwd)/target:/scan:ro" \
        -v "$(pwd)/logs:/logs" \
        -e SCANNER_LOG_LEVEL=info \
        -e SCANNER_OUTPUT_FORMAT=json \
        securityai-scanner "$@" || {
            echo_color $RED "Failed to start container"
            exit 1
        }

    echo_color $GREEN "Scanner started successfully!"
}

# Clean up function
cleanup() {
    echo_color $YELLOW "Cleaning up..."
    docker rm -f security-scanner &> /dev/null
}

# Show help message
show_help() {
    echo "Usage: ./start.sh [command] [options]"
    echo ""
    echo "Commands:"
    echo "  start     - Start the security scanner"
    echo "  test      - Run tests"
    echo "  help      - Show this help message"
    echo ""
    echo "Test Options:"
    echo "  all       - Run all tests"
    echo "  coverage  - Run tests with coverage"
    echo "  integration - Run integration tests"
    echo "  version   - Run version package tests"
    echo "  reporter  - Run reporter package tests"
    echo "  benchmark - Run benchmark tests"
}

# Main execution
main() {
    command=$1
    shift

    case $command in
        "test")
            test_type=${1:-"all"}  # Default to "all" if no test type specified
            run_tests $test_type
            ;;
        "start")
            check_docker
            cleanup
            build_image
            start_app "$@"
            ;;
        "help")
            show_help
            ;;
        "")
            check_docker
            cleanup
            build_image
            start_app
            ;;
        *)
            echo_color $RED "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# Set up trap for cleanup on script exit
trap cleanup EXIT

# Run main function with all arguments passed to script
main "$@"
