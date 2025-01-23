package models

import (
	"time"
)

// Severity represents the severity level of a security finding
type Severity string

const (
	SeverityCritical Severity = "CRITICAL"
	SeverityHigh     Severity = "HIGH"
	SeverityMedium   Severity = "MEDIUM"
	SeverityLow      Severity = "LOW"
	SeverityInfo     Severity = "INFO"
)

// Finding represents a security finding or vulnerability
type Finding struct {
	ID          string    `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Severity    Severity  `json:"severity"`
	Category    string    `json:"category"`
	Location    string    `json:"location"`
	CodeSnippet string    `json:"codeSnippet,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
	Remediation string    `json:"remediation,omitempty"`
	Confidence  float64   `json:"confidence"`
}
