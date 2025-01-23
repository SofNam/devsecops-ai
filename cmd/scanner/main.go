package main

import (
	"flag"
	"log"
	"time"

	"github.com/SofNam/devsecops-ai/pkg/ai"
	"github.com/SofNam/devsecops-ai/pkg/reporter"
	"github.com/SofNam/devsecops-ai/pkg/scanner"
	"github.com/SofNam/devsecops-ai/pkg/version"
)

func main() {
	// Command line flags
	targetPath := flag.String("path", ".", "Path to scan")
	modelPath := flag.String("model", "", "Path to AI model")
	outputFormat := flag.String("output", "json", "Output format (json/html)")
	outputPath := flag.String("output-path", "security-report", "Output file path")
	showVersion := flag.Bool("version", false, "Show version information")

	flag.Parse()

	// Show version if requested
	if *showVersion {
		vInfo := version.GetVersion()
		log.Printf("\nScanner Version Information:\n%s\n", vInfo.String())
		return
	}

	// Initialize scanner
	s := scanner.New(&scanner.Config{
		TargetPath: *targetPath,
		ModelPath:  *modelPath,
	})

	// Initialize AI detector
	detector := ai.NewDetector(*modelPath)

	// Run security scan
	findings, err := s.Scan()
	if err != nil {
		log.Fatalf("Scan failed: %v", err)
	}

	// Analyze with AI
	aiResults, err := detector.Analyze(findings)
	if err != nil {
		log.Fatalf("AI analysis failed: %v", err)
	}

	// Get version information
	vInfo := version.GetVersion()

	// Create report configuration
	config := reporter.Config{
		Version:     vInfo.Version,
		RulesUsed:   []string{"SEC-001", "SEC-002"},
		ScanType:    "Security Scan",
		AIEnabled:   true,
		TimeoutSecs: 30,
	}

	// Record start time for report
	startTime := time.Now()

	// Initialize reporter and generate report
	r := reporter.New(*outputFormat, *outputPath+"."+*outputFormat)
	if err := r.Generate(aiResults, config, *targetPath, startTime); err != nil {
		log.Fatalf("Report generation failed: %v", err)
	}

	log.Printf("Report generated successfully at: %s.%s", *outputPath, *outputFormat)
}
