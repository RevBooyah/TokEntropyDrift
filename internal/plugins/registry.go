package plugins

import (
	"fmt"
	"sync"
	"time"
)

// Registry manages plugin registration and execution
type Registry struct {
	plugins map[string]Plugin
	configs map[string]map[string]interface{}
	mu      sync.RWMutex
}

// NewRegistry creates a new plugin registry
func NewRegistry() *Registry {
	return &Registry{
		plugins: make(map[string]Plugin),
		configs: make(map[string]map[string]interface{}),
	}
}

// Register adds a plugin to the registry
func (r *Registry) Register(plugin Plugin) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	info := plugin.Info()
	if info.Name == "" {
		return fmt.Errorf("plugin name cannot be empty")
	}

	if _, exists := r.plugins[info.Name]; exists {
		return fmt.Errorf("plugin %s is already registered", info.Name)
	}

	r.plugins[info.Name] = plugin
	return nil
}

// Unregister removes a plugin from the registry
func (r *Registry) Unregister(name string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s is not registered", name)
	}

	// Cleanup the plugin
	if err := plugin.Cleanup(); err != nil {
		return fmt.Errorf("error cleaning up plugin %s: %w", name, err)
	}

	delete(r.plugins, name)
	delete(r.configs, name)
	return nil
}

// Get retrieves a plugin by name
func (r *Registry) Get(name string) (Plugin, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s is not registered", name)
	}

	return plugin, nil
}

// List returns all registered plugin names
func (r *Registry) List() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.plugins))
	for name := range r.plugins {
		names = append(names, name)
	}

	return names
}

// ListInfo returns information about all registered plugins
func (r *Registry) ListInfo() []PluginInfo {
	r.mu.RLock()
	defer r.mu.RUnlock()

	infos := make([]PluginInfo, 0, len(r.plugins))
	for _, plugin := range r.plugins {
		infos = append(infos, plugin.Info())
	}

	return infos
}

// Configure sets configuration for a plugin
func (r *Registry) Configure(name string, config map[string]interface{}) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s is not registered", name)
	}

	// Validate configuration
	if err := plugin.ValidateConfig(config); err != nil {
		return fmt.Errorf("invalid configuration for plugin %s: %w", name, err)
	}

	// Store configuration
	r.configs[name] = config

	// Initialize plugin with new configuration
	if err := plugin.Initialize(config); err != nil {
		return fmt.Errorf("error initializing plugin %s: %w", name, err)
	}

	return nil
}

// ExecuteMetrics runs metric calculations for all plugins
func (r *Registry) ExecuteMetrics(ctx *AnalysisContext) (map[string][]MetricResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	results := make(map[string][]MetricResult)

	for name, plugin := range r.plugins {
		metrics, err := plugin.CalculateMetrics(ctx)
		if err != nil {
			return nil, fmt.Errorf("error executing plugin %s: %w", name, err)
		}

		// Add timestamp to metrics if not present
		for i := range metrics {
			if metrics[i].Timestamp.IsZero() {
				metrics[i].Timestamp = time.Now()
			}
		}

		results[name] = metrics
	}

	return results, nil
}

// ExecuteMetricsForPlugin runs metric calculations for a specific plugin
func (r *Registry) ExecuteMetricsForPlugin(name string, ctx *AnalysisContext) ([]MetricResult, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	plugin, exists := r.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s is not registered", name)
	}

	metrics, err := plugin.CalculateMetrics(ctx)
	if err != nil {
		return nil, fmt.Errorf("error executing plugin %s: %w", name, err)
	}

	// Add timestamp to metrics if not present
	for i := range metrics {
		if metrics[i].Timestamp.IsZero() {
			metrics[i].Timestamp = time.Now()
		}
	}

	return metrics, nil
}

// Close cleans up all plugins
func (r *Registry) Close() error {
	r.mu.Lock()
	defer r.mu.Unlock()

	var errors []error
	for name, plugin := range r.plugins {
		if err := plugin.Cleanup(); err != nil {
			errors = append(errors, fmt.Errorf("error cleaning up plugin %s: %w", name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("errors during cleanup: %v", errors)
	}

	return nil
}

// GetPluginCount returns the number of registered plugins
func (r *Registry) GetPluginCount() int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	return len(r.plugins)
}

// IsRegistered checks if a plugin is registered
func (r *Registry) IsRegistered(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	_, exists := r.plugins[name]
	return exists
}
