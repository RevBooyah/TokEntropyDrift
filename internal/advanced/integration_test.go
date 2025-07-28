package advanced

import (
	"context"
	"testing"
	"time"

	"github.com/RevBooyah/tokentropydrift/internal/config"
	"github.com/RevBooyah/tokentropydrift/internal/metrics"
	"github.com/RevBooyah/tokentropydrift/internal/tokenizers"
)

func TestAdvancedManagerCreation(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled:         true,
			MaxSize:         1000,
			TTL:             "1h",
			CleanupInterval: "10m",
			EnableStats:     true,
		},
		Parallel: config.ParallelConfig{
			Enabled:       true,
			MaxWorkers:    4,
			BatchSize:     100,
			Timeout:       "30m",
			EnableMetrics: true,
		},
		Streaming: config.StreamingConfig{
			Enabled:          true,
			ChunkSize:        1000,
			BufferSize:       65536,
			MaxMemoryMB:      512,
			EnableProgress:   true,
			ProgressInterval: 10,
			Timeout:          "1h",
		},
		Plugins: config.PluginsConfig{
			Enabled:         true,
			AutoLoad:        true,
			PluginDirectory: "plugins",
			Configs:         make(map[string]interface{}),
		},
	}

	// Create metrics engine
	engine := metrics.NewEngine(metrics.EngineConfig{
		EntropyWindowSize: 100,
		NormalizeEntropy:  true,
		CompressionRatio:  true,
		DriftDetection:    true,
	})

	// Create advanced manager
	manager, err := NewAdvancedManager(cfg, engine)
	if err != nil {
		t.Fatalf("Failed to create AdvancedManager: %v", err)
	}
	defer manager.Close()

	// Test basic functionality
	if manager.config == nil {
		t.Error("Manager config should not be nil")
	}

	if manager.cache == nil {
		t.Error("Cache should be initialized when enabled")
	}

	if manager.processor == nil {
		t.Error("Processor should be initialized when enabled")
	}

	if manager.streamer == nil {
		t.Error("Streamer should be initialized when enabled")
	}

	if manager.pluginReg == nil {
		t.Error("Plugin registry should be initialized when enabled")
	}
}

func TestAdvancedManagerWithMockTokenizer(t *testing.T) {
	// Create test configuration
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled:         true,
			MaxSize:         1000,
			TTL:             "1h",
			CleanupInterval: "10m",
			EnableStats:     true,
		},
		Parallel: config.ParallelConfig{
			Enabled:       true,
			MaxWorkers:    2,
			BatchSize:     10,
			Timeout:       "5m",
			EnableMetrics: true,
		},
		Streaming: config.StreamingConfig{
			Enabled:          true,
			ChunkSize:        100,
			BufferSize:       4096,
			MaxMemoryMB:      100,
			EnableProgress:   true,
			ProgressInterval: 5,
			Timeout:          "10m",
		},
		Plugins: config.PluginsConfig{
			Enabled:         true,
			AutoLoad:        true,
			PluginDirectory: "plugins",
			Configs:         make(map[string]interface{}),
		},
	}

	// Create metrics engine
	engine := metrics.NewEngine(metrics.EngineConfig{
		EntropyWindowSize: 50,
		NormalizeEntropy:  true,
		CompressionRatio:  true,
		DriftDetection:    true,
	})

	// Create advanced manager
	manager, err := NewAdvancedManager(cfg, engine)
	if err != nil {
		t.Fatalf("Failed to create AdvancedManager: %v", err)
	}
	defer manager.Close()

	// Create and register mock tokenizer
	mockTokenizer := tokenizers.NewMockTokenizer()
	mockTokenizer.Initialize(tokenizers.TokenizerConfig{
		Name: "mock",
		Type: "custom",
	})

	err = manager.RegisterTokenizer("mock", mockTokenizer)
	if err != nil {
		t.Fatalf("Failed to register tokenizer: %v", err)
	}

	// Test texts
	texts := []string{
		"The quick brown fox jumps over the lazy dog.",
		"Machine learning models process text efficiently.",
		"Hello, world! This is a test.",
	}

	// Progress callback
	progressCallback := func(chunk, total, lines int, duration time.Duration) {
		t.Logf("Progress: %d/%d chunks, %d lines, %v", chunk, total, lines, duration)
	}

	// Run analysis
	ctx := context.Background()
	result, err := manager.AnalyzeWithAdvanced(ctx, texts, "mock", progressCallback)
	if err != nil {
		t.Fatalf("Analysis failed: %v", err)
	}

	// Verify results
	if result == nil {
		t.Fatal("Result should not be nil")
	}

	if result.StartTime.IsZero() {
		t.Error("Start time should be set")
	}

	if result.EndTime.IsZero() {
		t.Error("End time should be set")
	}

	if result.Duration <= 0 {
		t.Error("Duration should be positive")
	}

	if result.Config == nil {
		t.Error("Config should not be nil")
	}

	// Check cache stats
	cacheStats := manager.GetCacheStats()
	if cacheStats == nil {
		t.Error("Cache stats should not be nil when caching is enabled")
	}

	// Check plugin info
	pluginInfo := manager.GetPluginInfo()
	if pluginInfo == nil {
		t.Error("Plugin info should not be nil")
	}
}

func TestAdvancedManagerWithoutFeatures(t *testing.T) {
	// Create test configuration with all features disabled
	cfg := &config.Config{
		Cache: config.CacheConfig{
			Enabled: false,
		},
		Parallel: config.ParallelConfig{
			Enabled: false,
		},
		Streaming: config.StreamingConfig{
			Enabled: false,
		},
		Plugins: config.PluginsConfig{
			Enabled: false,
		},
	}

	// Create metrics engine
	engine := metrics.NewEngine(metrics.EngineConfig{
		EntropyWindowSize: 100,
		NormalizeEntropy:  true,
		CompressionRatio:  true,
		DriftDetection:    true,
	})

	// Create advanced manager
	manager, err := NewAdvancedManager(cfg, engine)
	if err != nil {
		t.Fatalf("Failed to create AdvancedManager: %v", err)
	}
	defer manager.Close()

	// Verify that optional components are nil when disabled
	if manager.cache != nil {
		t.Error("Cache should be nil when disabled")
	}

	if manager.processor != nil {
		t.Error("Processor should be nil when disabled")
	}

	if manager.streamer != nil {
		t.Error("Streamer should be nil when disabled")
	}

	if manager.pluginReg != nil {
		t.Error("Plugin registry should be nil when disabled")
	}
}
