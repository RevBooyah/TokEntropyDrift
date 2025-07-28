package plugins

import (
	"context"
	"time"

	"github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
)

// PluginInfo contains metadata about a plugin
type PluginInfo struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description"`
	Author      string            `json:"author"`
	Tags        []string          `json:"tags"`
	Metadata    map[string]string `json:"metadata,omitempty"`
}

// MetricResult represents the result of a custom metric calculation
type MetricResult struct {
	Name      string                 `json:"name"`
	Value     float64                `json:"value"`
	Unit      string                 `json:"unit,omitempty"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
	Timestamp time.Time              `json:"timestamp"`
}

// AnalysisContext provides context for metric calculations
type AnalysisContext struct {
	Document      string                         `json:"document"`
	Tokenization  *tokenizers.TokenizationResult `json:"tokenization"`
	TokenizerName string                         `json:"tokenizer_name"`
	Config        map[string]interface{}         `json:"config"`
	Context       context.Context                `json:"-"`
}

// Plugin defines the interface that all plugins must implement
type Plugin interface {
	// Info returns information about the plugin
	Info() PluginInfo

	// Initialize is called when the plugin is loaded
	Initialize(config map[string]interface{}) error

	// CalculateMetrics calculates custom metrics for the given context
	CalculateMetrics(ctx *AnalysisContext) ([]MetricResult, error)

	// ValidateConfig validates the plugin configuration
	ValidateConfig(config map[string]interface{}) error

	// Cleanup is called when the plugin is unloaded
	Cleanup() error
}

// BasePlugin provides common functionality for plugin implementations
type BasePlugin struct {
	info   PluginInfo
	config map[string]interface{}
}

// NewBasePlugin creates a new base plugin
func NewBasePlugin(info PluginInfo) *BasePlugin {
	return &BasePlugin{
		info:   info,
		config: make(map[string]interface{}),
	}
}

// Info returns the plugin information
func (b *BasePlugin) Info() PluginInfo {
	return b.info
}

// Initialize provides a default implementation
func (b *BasePlugin) Initialize(config map[string]interface{}) error {
	b.config = config
	return nil
}

// ValidateConfig provides a default implementation
func (b *BasePlugin) ValidateConfig(config map[string]interface{}) error {
	// Default implementation accepts any config
	return nil
}

// Cleanup provides a default implementation
func (b *BasePlugin) Cleanup() error {
	// Default implementation does nothing
	return nil
}

// GetConfig returns the plugin configuration
func (b *BasePlugin) GetConfig() map[string]interface{} {
	return b.config
}

// GetConfigValue retrieves a configuration value with type assertion
func (b *BasePlugin) GetConfigValue(key string, defaultValue interface{}) interface{} {
	if value, exists := b.config[key]; exists {
		return value
	}
	return defaultValue
}

// GetConfigString retrieves a string configuration value
func (b *BasePlugin) GetConfigString(key string, defaultValue string) string {
	value := b.GetConfigValue(key, defaultValue)
	if str, ok := value.(string); ok {
		return str
	}
	return defaultValue
}

// GetConfigInt retrieves an integer configuration value
func (b *BasePlugin) GetConfigInt(key string, defaultValue int) int {
	value := b.GetConfigValue(key, defaultValue)
	if num, ok := value.(int); ok {
		return num
	}
	if num, ok := value.(float64); ok {
		return int(num)
	}
	return defaultValue
}

// GetConfigFloat retrieves a float configuration value
func (b *BasePlugin) GetConfigFloat(key string, defaultValue float64) float64 {
	value := b.GetConfigValue(key, defaultValue)
	if num, ok := value.(float64); ok {
		return num
	}
	if num, ok := value.(int); ok {
		return float64(num)
	}
	return defaultValue
}

// GetConfigBool retrieves a boolean configuration value
func (b *BasePlugin) GetConfigBool(key string, defaultValue bool) bool {
	value := b.GetConfigValue(key, defaultValue)
	if b, ok := value.(bool); ok {
		return b
	}
	return defaultValue
}
