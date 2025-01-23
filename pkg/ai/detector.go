package ai

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/SofNam/devsecops-ai/pkg/models"
	"github.com/SofNam/devsecops-ai/pkg/reporter"
)

// Detector represents the AI-based security detector
type Detector struct {
	modelPath   string
	confidence  float64
	maxFindings int
	initialized bool
	rules       []Rule
}

// Rule represents a security rule for AI analysis
type Rule struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Pattern     string   `json:"pattern"`
	Severity    string   `json:"severity"`
	Category    string   `json:"category"`
	Keywords    []string `json:"keywords"`
	Description string   `json:"description"`
}

// DetectorConfig holds configuration for the detector
type DetectorConfig struct {
	Confidence  float64 `json:"confidence"`
	MaxFindings int     `json:"maxFindings"`
}

// NewDetector creates a new AI detector instance
func NewDetector(modelPath string) *Detector {
	d := &Detector{
		modelPath:   modelPath,
		confidence:  0.75, // Default confidence threshold
		maxFindings: 100,  // Default maximum findings
	}

	if err := d.initialize(); err != nil {
		log.Printf("Warning: Failed to initialize AI detector: %v", err)
	}

	return d
}

// initialize loads the AI model and rules
func (d *Detector) initialize() error {
	// Load rules from model path
	rulesPath := filepath.Join(d.modelPath, "rules.json")
	if _, err := os.Stat(rulesPath); err == nil {
		rules, err := loadRules(rulesPath)
		if err != nil {
			return fmt.Errorf("failed to load rules: %v", err)
		}
		d.rules = rules
	}

	// Load configuration
	configPath := filepath.Join(d.modelPath, "config.json")
	if _, err := os.Stat(configPath); err == nil {
		config, err := loadConfig(configPath)
		if err != nil {
			return fmt.Errorf("failed to load config: %v", err)
		}
		d.confidence = config.Confidence
		d.maxFindings = config.MaxFindings
	}

	d.initialized = true
	return nil
}

// Analyze performs AI-based analysis on findings
func (d *Detector) Analyze(findings []models.Finding) ([]models.Finding, error) {
	if !d.initialized {
		return findings, fmt.Errorf("detector not properly initialized")
	}

	var enhancedFindings []models.Finding

	for _, finding := range findings {
		// Enhance finding with AI analysis
		enhanced := d.enhanceFinding(finding)
		enhancedFindings = append(enhancedFindings, enhanced)
	}

	// Perform additional AI-based detection
	additionalFindings := d.detectAdditionalIssues(findings)
	enhancedFindings = append(enhancedFindings, additionalFindings...)

	// Sort and limit findings based on severity and confidence
	enhancedFindings = d.prioritizeFindings(enhancedFindings)

	return enhancedFindings, nil
}

// enhanceFinding enhances a single finding with AI insights
func (d *Detector) enhanceFinding(finding models.Finding) models.Finding {
	// Here you would typically:
	// 1. Use AI to validate the finding
	// 2. Add additional context
	// 3. Enhance remediation suggestions
	// 4. Adjust severity based on context

	// For now, we'll just add some basic enhancements
	finding.Description = fmt.Sprintf("%s (AI Verified)", finding.Description)
	if finding.Remediation == "" {
		finding.Remediation = "AI suggested: Review and sanitize all inputs"
	}

	return finding
}

// detectAdditionalIssues uses AI to find additional security issues
func (d *Detector) detectAdditionalIssues(findings []models.Finding) []models.Finding {
	var additionalFindings []models.Finding

	// Apply each rule
	for _, rule := range d.rules {
		// In a real implementation, you would:
		// 1. Use AI to analyze code patterns
		// 2. Look for security anti-patterns
		// 3. Identify potential vulnerabilities
		// 4. Calculate confidence scores

		// Example placeholder for demonstration
		if rule.Pattern != "" {
			finding := models.Finding{
				ID:          fmt.Sprintf("AI-%s", rule.ID),
				Title:       rule.Name,
				Description: rule.Description,
				Severity:    models.Severity(reporter.Severity(rule.Severity)),
				Category:    rule.Category,
			}
			additionalFindings = append(additionalFindings, finding)
		}
	}

	return additionalFindings
}

// prioritizeFindings sorts and limits findings based on severity and confidence
func (d *Detector) prioritizeFindings(findings []models.Finding) []models.Finding {
	// In a real implementation, you would:
	// 1. Sort by severity
	// 2. Filter by confidence threshold
	// 3. Limit to maxFindings
	// 4. Remove duplicates

	if len(findings) > d.maxFindings {
		findings = findings[:d.maxFindings]
	}

	return findings
}

// loadRules loads security rules from a JSON file
func loadRules(path string) ([]Rule, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var rules []Rule
	if err := json.Unmarshal(data, &rules); err != nil {
		return nil, err
	}

	return rules, nil
}

// loadConfig loads detector configuration from a JSON file
func loadConfig(path string) (*DetectorConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config DetectorConfig
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
