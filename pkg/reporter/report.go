package reporter

import (
	"encoding/json"
	"fmt"
	"html/template"
	"os"
	"time"

	"github.com/SofNam/devsecops-ai/pkg/models"
)

// Severity represents the severity level of a finding
type Severity string

const (
	Critical models.Severity = "CRITICAL"
	High     models.Severity = "HIGH"
	Medium   models.Severity = "MEDIUM"
	Low      models.Severity = "LOW"
	Info     models.Severity = "INFO"
)

// Finding represents a security finding
type Finding struct {
	ID          string          `json:"id"`
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Severity    models.Severity `json:"severity"`
	Category    string          `json:"category"`
	Location    string          `json:"location"`
	CodeSnippet string          `json:"codeSnippet,omitempty"`
	Timestamp   time.Time       `json:"timestamp"`
	Remediation string          `json:"remediation,omitempty"`
}

// Report represents the complete security scan report
type Report struct {
	ScanID        string           `json:"scanId"`
	Timestamp     time.Time        `json:"timestamp"`
	Target        string           `json:"target"`
	Findings      []models.Finding `json:"findings"`
	SummaryStats  Stats            `json:"summaryStats"`
	ScanDuration  string           `json:"scanDuration"`
	ScannerConfig Config           `json:"scannerConfig"`
}

// Stats represents statistical information about the findings
type Stats struct {
	TotalFindings int `json:"totalFindings"`
	CriticalCount int `json:"criticalCount"`
	HighCount     int `json:"highCount"`
	MediumCount   int `json:"mediumCount"`
	LowCount      int `json:"lowCount"`
	InfoCount     int `json:"infoCount"`
}

// Config represents scanner configuration
type Config struct {
	Version     string   `json:"version"`
	RulesUsed   []string `json:"rulesUsed"`
	ScanType    string   `json:"scanType"`
	AIEnabled   bool     `json:"aiEnabled"`
	TimeoutSecs int      `json:"timeoutSecs"`
}

// Reporter handles report generation
type Reporter struct {
	OutputFormat string
	OutputPath   string
}

// NewReporter creates a new reporter instance
func New(format, path string) *Reporter {
	return &Reporter{
		OutputFormat: format,
		OutputPath:   path,
	}
}

// Generate creates a report in the specified format
func (r *Reporter) Generate(findings []models.Finding, config Config, target string, duration time.Time) error {
	report := r.createReport(findings, config, target, duration)

	switch r.OutputFormat {
	case "json":
		return r.generateJSON(report)
	case "html":
		return r.generateHTML(report)
	default:
		return fmt.Errorf("unsupported format: %s", r.OutputFormat)
	}
}

// createReport assembles the complete report
func (r *Reporter) createReport(findings []models.Finding, config Config, target string, duration time.Time) Report {
	stats := r.calculateStats(findings)

	return Report{
		ScanID:        fmt.Sprintf("SCAN-%d", time.Now().Unix()),
		Timestamp:     time.Now(),
		Target:        target,
		Findings:      findings,
		SummaryStats:  stats,
		ScanDuration:  time.Since(duration).String(),
		ScannerConfig: config,
	}
}

// calculateStats calculates statistics for findings
func (r *Reporter) calculateStats(findings []models.Finding) Stats {
	stats := Stats{}

	for _, finding := range findings {
		stats.TotalFindings++
		switch finding.Severity {
		case Critical:
			stats.CriticalCount++
		case High:
			stats.HighCount++
		case Medium:
			stats.MediumCount++
		case Low:
			stats.LowCount++
		case Info:
			stats.InfoCount++
		}
	}

	return stats
}

// generateJSON creates a JSON report
func (r *Reporter) generateJSON(report Report) error {
	file, err := os.Create(r.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(report); err != nil {
		return fmt.Errorf("failed to encode report: %v", err)
	}

	return nil
}

// generateHTML creates an HTML report
func (r *Reporter) generateHTML(report Report) error {
	tmpl, err := template.New("report").Parse(htmlTemplate)
	if err != nil {
		return fmt.Errorf("failed to parse HTML template: %v", err)
	}

	file, err := os.Create(r.OutputPath)
	if err != nil {
		return fmt.Errorf("failed to create report file: %v", err)
	}
	defer file.Close()

	if err := tmpl.Execute(file, report); err != nil {
		return fmt.Errorf("failed to generate HTML report: %v", err)
	}

	return nil
}

// HTML template for report generation
const htmlTemplate = `
<!DOCTYPE html>
<html>
<head>
    <title>Security Scan Report</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            color: #333;
        }
        .header {
            background-color: #f8f9fa;
            padding: 20px;
            border-radius: 5px;
            margin-bottom: 20px;
        }
        .finding {
            border: 1px solid #ddd;
            padding: 15px;
            margin-bottom: 15px;
            border-radius: 5px;
        }
        .critical { border-left: 5px solid #dc3545; }
        .high { border-left: 5px solid #fd7e14; }
        .medium { border-left: 5px solid #ffc107; }
        .low { border-left: 5px solid #28a745; }
        .info { border-left: 5px solid #17a2b8; }
        .stats {
            display: grid;
            grid-template-columns: repeat(auto-fit, minmax(150px, 1fr));
            gap: 10px;
            margin: 20px 0;
        }
        .stat-item {
            padding: 10px;
            background-color: #f8f9fa;
            border-radius: 5px;
            text-align: center;
        }
        code {
            background-color: #f8f9fa;
            padding: 10px;
            display: block;
            border-radius: 5px;
            margin: 10px 0;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>Security Scan Report</h1>
        <p>Scan ID: {{.ScanID}}</p>
        <p>Target: {{.Target}}</p>
        <p>Timestamp: {{.Timestamp}}</p>
        <p>Duration: {{.ScanDuration}}</p>
    </div>

    <div class="stats">
        <div class="stat-item">
            <h3>Total</h3>
            <p>{{.SummaryStats.TotalFindings}}</p>
        </div>
        <div class="stat-item">
            <h3>Critical</h3>
            <p>{{.SummaryStats.CriticalCount}}</p>
        </div>
        <div class="stat-item">
            <h3>High</h3>
            <p>{{.SummaryStats.HighCount}}</p>
        </div>
        <div class="stat-item">
            <h3>Medium</h3>
            <p>{{.SummaryStats.MediumCount}}</p>
        </div>
        <div class="stat-item">
            <h3>Low</h3>
            <p>{{.SummaryStats.LowCount}}</p>
        </div>
        <div class="stat-item">
            <h3>Info</h3>
            <p>{{.SummaryStats.InfoCount}}</p>
        </div>
    </div>

    <h2>Findings</h2>
    {{range .Findings}}
    <div class="finding {{.Severity | printf "%s" | toLowerCase}}">
        <h3>{{.Title}}</h3>
        <p><strong>Severity:</strong> {{.Severity}}</p>
        <p><strong>Category:</strong> {{.Category}}</p>
        <p><strong>Location:</strong> {{.Location}}</p>
        <p>{{.Description}}</p>
        {{if .CodeSnippet}}
        <code>{{.CodeSnippet}}</code>
        {{end}}
        {{if .Remediation}}
        <p><strong>Remediation:</strong> {{.Remediation}}</p>
        {{end}}
    </div>
    {{end}}
</body>
</html>
`
