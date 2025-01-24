package ai

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/SofNam/devsecops-ai/pkg/models"
)

// Classifier represents the AI-based security classifier
type Classifier struct {
	modelPath    string
	threshold    float64
	categories   []string
	initialized  bool
	modelConfig  ModelConfig
	categoryData map[string]CategoryFeatures
}

// ModelConfig holds AI model configuration
type ModelConfig struct {
	Threshold   float64 `json:"threshold"`
	BatchSize   int     `json:"batchSize"`
	EnableCache bool    `json:"enableCache"`
}

// CategoryFeatures holds feature data for each security category
type CategoryFeatures struct {
	Patterns  []string  `json:"patterns"`
	Keywords  []string  `json:"keywords"`
	Weights   []float64 `json:"weights"`
	Threshold float64   `json:"threshold"`
}

// NewClassifier creates a new AI classifier instance
func NewClassifier(modelPath string) *Classifier {
	c := &Classifier{
		modelPath:    modelPath,
		threshold:    0.8,
		categoryData: make(map[string]CategoryFeatures),
	}

	if err := c.initialize(); err != nil {
		return c
	}

	return c
}

// initialize loads model configuration and category data
func (c *Classifier) initialize() error {
	// Load model configuration
	configPath := filepath.Join(c.modelPath, "config.json")
	if err := c.loadConfig(configPath); err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Load category data
	categoryPath := filepath.Join(c.modelPath, "rules.json")
	if err := c.loadCategories(categoryPath); err != nil {
		return fmt.Errorf("failed to load categories: %v", err)
	}

	c.initialized = true
	return nil
}

// loadConfig loads model configuration from JSON
func (c *Classifier) loadConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var config struct {
		ModelSettings ModelConfig `json:"modelSettings"`
		Categories    []string    `json:"categories"`
	}

	if err := json.Unmarshal(data, &config); err != nil {
		return err
	}

	c.modelConfig = config.ModelSettings
	c.categories = config.Categories
	c.threshold = config.ModelSettings.Threshold

	return nil
}

// loadCategories loads category feature data
func (c *Classifier) loadCategories(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var rulesData struct {
		Rules []Rule `json:"rules"`
	}

	if err := json.Unmarshal(data, &rulesData); err != nil {
		return err
	}

	// Process rules into category features
	for _, rule := range rulesData.Rules {
		features := c.categoryData[rule.Category]
		features.Patterns = append(features.Patterns, rule.Pattern)
		features.Keywords = append(features.Keywords, rule.Keywords...)
		features.Weights = append(features.Weights, 1.0) // Default weight
		features.Threshold = c.threshold
		c.categoryData[rule.Category] = features
	}

	return nil
}

// Classify performs classification on a finding
func (c *Classifier) Classify(finding *models.Finding) error {
	if !c.initialized {
		return fmt.Errorf("classifier not initialized")
	}

	// Calculate confidence scores for each category
	scores := make(map[string]float64)
	for category, features := range c.categoryData {
		score := c.calculateScore(finding, features)
		scores[category] = score
	}

	// Get highest scoring category
	bestCategory, bestScore := c.getBestCategory(scores)

	// Update finding if confidence threshold is met
	if bestScore >= c.threshold {
		finding.Category = bestCategory
		finding.Confidence = bestScore
	}

	return nil
}

// calculateScore calculates confidence score for a category
func (c *Classifier) calculateScore(finding *models.Finding, features CategoryFeatures) float64 {
	var score float64

	// Pattern matching
	for i, pattern := range features.Patterns {
		if strings.Contains(finding.CodeSnippet, pattern) {
			score += features.Weights[i]
		}
	}

	// Keyword matching
	for _, keyword := range features.Keywords {
		if strings.Contains(strings.ToLower(finding.Description), strings.ToLower(keyword)) {
			score += 0.5 // Lower weight for keyword matches
		}
	}

	// Normalize score
	maxScore := float64(len(features.Patterns)) + (float64(len(features.Keywords)) * 0.5)
	if maxScore > 0 {
		score /= maxScore
	}

	return score
}

// getBestCategory returns highest scoring category and score
func (c *Classifier) getBestCategory(scores map[string]float64) (string, float64) {
	type categoryScore struct {
		category string
		score    float64
	}

	var sorted []categoryScore
	for category, score := range scores {
		sorted = append(sorted, categoryScore{category, score})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].score > sorted[j].score
	})

	if len(sorted) > 0 {
		return sorted[0].category, sorted[0].score
	}

	return "", 0
}

// GetCategories returns list of supported categories
func (c *Classifier) GetCategories() []string {
	return c.categories
}

// UpdateThreshold updates classification threshold
func (c *Classifier) UpdateThreshold(threshold float64) {
	c.threshold = threshold
}
