package advanced

import (
	"context"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/RevBooyah/tokentropydrift/internal/cache"
	"github.com/RevBooyah/tokentropydrift/internal/config"
	"github.com/RevBooyah/tokentropydrift/internal/metrics"
	"github.com/RevBooyah/tokentropydrift/internal/parallel"
	"github.com/RevBooyah/tokentropydrift/internal/plugins"
	"github.com/RevBooyah/tokentropydrift/internal/streaming"
	"github.com/RevBooyah/tokentropydrift/internal/tokenizers"
)

// AdvancedManager manages all advanced features
type AdvancedManager struct {
	config     *config.Config
	cache      *cache.Cache
	processor  *parallel.Processor
	streamer   *streaming.StreamAnalyzer
	pluginReg  *plugins.Registry
	engine     *metrics.Engine
	tokenizers map[string]tokenizers.Tokenizer
}

// NewAdvancedManager creates a new advanced features manager
func NewAdvancedManager(cfg *config.Config, engine *metrics.Engine) (*AdvancedManager, error) {
	manager := &AdvancedManager{
		config:     cfg,
		engine:     engine,
		tokenizers: make(map[string]tokenizers.Tokenizer),
	}

	// Initialize cache if enabled
	if cfg.Cache.Enabled {
		cacheConfig := cache.CacheConfig{
			MaxSize:         cfg.Cache.MaxSize,
			TTL:             parseDuration(cfg.Cache.TTL),
			CleanupInterval: parseDuration(cfg.Cache.CleanupInterval),
			EnableStats:     cfg.Cache.EnableStats,
		}
		manager.cache = cache.NewCache(cacheConfig)
	}

	// Initialize parallel processor if enabled
	if cfg.Parallel.Enabled {
		processorConfig := parallel.ProcessorConfig{
			MaxWorkers:    cfg.Parallel.MaxWorkers,
			BatchSize:     cfg.Parallel.BatchSize,
			Timeout:       parseDuration(cfg.Parallel.Timeout),
			EnableMetrics: cfg.Parallel.EnableMetrics,
		}
		manager.processor = parallel.NewProcessor(processorConfig)
	}

	// Initialize streaming analyzer if enabled
	if cfg.Streaming.Enabled {
		streamConfig := streaming.StreamConfig{
			ChunkSize:        cfg.Streaming.ChunkSize,
			BufferSize:       cfg.Streaming.BufferSize,
			MaxMemoryMB:      cfg.Streaming.MaxMemoryMB,
			EnableProgress:   cfg.Streaming.EnableProgress,
			ProgressInterval: cfg.Streaming.ProgressInterval,
			Timeout:          parseDuration(cfg.Streaming.Timeout),
		}
		manager.streamer = streaming.NewStreamAnalyzer(streamConfig, engine)
	}

	// Initialize plugin registry if enabled
	if cfg.Plugins.Enabled {
		manager.pluginReg = plugins.NewRegistry()
		if err := manager.loadPlugins(); err != nil {
			return nil, fmt.Errorf("failed to load plugins: %w", err)
		}
	}

	return manager, nil
}

// RegisterTokenizer registers a tokenizer with caching if enabled
func (m *AdvancedManager) RegisterTokenizer(name string, tokenizer tokenizers.Tokenizer) error {
	if m.config.Cache.Enabled && m.cache != nil {
		cacheConfig := cache.CacheConfig{
			MaxSize:         m.config.Cache.MaxSize,
			TTL:             parseDuration(m.config.Cache.TTL),
			CleanupInterval: parseDuration(m.config.Cache.CleanupInterval),
			EnableStats:     m.config.Cache.EnableStats,
		}
		cachedTokenizer := tokenizers.NewCachedTokenizer(tokenizer, cache.NewCache(cacheConfig))
		m.tokenizers[name] = cachedTokenizer
	} else {
		m.tokenizers[name] = tokenizer
	}
	return nil
}

// GetTokenizer retrieves a registered tokenizer
func (m *AdvancedManager) GetTokenizer(name string) (tokenizers.Tokenizer, error) {
	tokenizer, exists := m.tokenizers[name]
	if !exists {
		return nil, fmt.Errorf("tokenizer %s not found", name)
	}
	return tokenizer, nil
}

// AnalyzeWithAdvanced performs analysis using all advanced features
func (m *AdvancedManager) AnalyzeWithAdvanced(
	ctx context.Context,
	texts []string,
	tokenizerName string,
	progressCallback func(int, int, int, time.Duration),
) (*AdvancedAnalysisResult, error) {

	tokenizer, err := m.GetTokenizer(tokenizerName)
	if err != nil {
		return nil, err
	}

	result := &AdvancedAnalysisResult{
		StartTime: time.Now(),
		Config:    m.config,
	}

	// Use parallel processing for large datasets
	if m.config.Parallel.Enabled && len(texts) > m.config.Parallel.BatchSize {
		result.ParallelStats = m.processParallel(ctx, texts, tokenizer)
	} else {
		// Use streaming for very large datasets
		if m.config.Streaming.Enabled && len(texts) > m.config.Streaming.ChunkSize*10 {
			result.StreamingStats = m.processStreaming(ctx, texts, tokenizer, progressCallback)
		} else {
			// Use standard processing
			result.StandardResults = m.processStandard(ctx, texts, tokenizer)
		}
	}

	// Execute plugins if enabled
	if m.config.Plugins.Enabled && m.pluginReg != nil {
		result.PluginResults = m.executePlugins(ctx, texts, tokenizer)
	}

	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	return result, nil
}

// processParallel processes texts using parallel processing
func (m *AdvancedManager) processParallel(
	ctx context.Context,
	texts []string,
	tokenizer tokenizers.Tokenizer,
) *parallel.ProcessingStats {

	results, _, stats := m.processor.ProcessTokenizations(ctx, texts, tokenizer)

	// Convert results to analysis results
	analysisResults := make([]*metrics.AnalysisResult, len(results))
	for i, result := range results {
		analysisResults[i] = &metrics.AnalysisResult{
			Document:      result.Document,
			TokenizerName: result.Tokenizer,
			TokenCount:    len(result.Tokens),
			Tokenization:  result,
		}
	}

	return &stats
}

// processStreaming processes texts using streaming analysis
func (m *AdvancedManager) processStreaming(
	ctx context.Context,
	texts []string,
	tokenizer tokenizers.Tokenizer,
	progressCallback func(int, int, int, time.Duration),
) *streaming.StreamResult {

	// Convert texts to a reader for streaming
	reader := createTextReader(texts)

	streamResult, err := m.streamer.AnalyzeStream(ctx, reader, tokenizer, progressCallback)
	if err != nil {
		// Return empty result on error
		return &streaming.StreamResult{
			Errors: []string{err.Error()},
		}
	}

	return streamResult
}

// processStandard processes texts using standard analysis
func (m *AdvancedManager) processStandard(
	ctx context.Context,
	texts []string,
	tokenizer tokenizers.Tokenizer,
) []*metrics.AnalysisResult {

	results := make([]*metrics.AnalysisResult, len(texts))
	for i, text := range texts {
		result, err := m.engine.AnalyzeDocument(ctx, text, tokenizer)
		if err != nil {
			// Create empty result on error
			results[i] = &metrics.AnalysisResult{
				Document:      text,
				TokenizerName: tokenizer.Name(),
				TokenCount:    0,
				Metrics:       make(map[string]metrics.MetricResult),
			}
		} else {
			results[i] = result
		}
	}

	return results
}

// executePlugins executes all registered plugins
func (m *AdvancedManager) executePlugins(
	ctx context.Context,
	texts []string,
	tokenizer tokenizers.Tokenizer,
) map[string][]plugins.MetricResult {

	// Create analysis context for plugins
	analysisContext := &plugins.AnalysisContext{
		Document:      strings.Join(texts, "\n"),
		TokenizerName: tokenizer.Name(),
		Config:        make(map[string]interface{}),
		Context:       ctx,
	}

	// Execute plugins
	pluginResults, err := m.pluginReg.ExecuteMetrics(analysisContext)
	if err != nil {
		return map[string][]plugins.MetricResult{
			"error": {
				{
					Name:  "plugin_execution_error",
					Value: 0,
					Unit:  "error",
				},
			},
		}
	}

	return pluginResults
}

// loadPlugins loads and registers plugins
func (m *AdvancedManager) loadPlugins() error {
	// TODO: Implement plugin loading from files
	// For now, return nil to indicate no plugins loaded
	return nil
}

// GetCacheStats returns cache statistics if caching is enabled
func (m *AdvancedManager) GetCacheStats() *cache.CacheStats {
	if m.cache != nil {
		return m.cache.GetStats()
	}
	return nil
}

// GetPluginInfo returns information about loaded plugins
func (m *AdvancedManager) GetPluginInfo() []plugins.PluginInfo {
	if m.pluginReg != nil {
		return m.pluginReg.ListInfo()
	}
	return []plugins.PluginInfo{}
}

// Close cleans up all resources
func (m *AdvancedManager) Close() error {
	if m.cache != nil {
		m.cache.Close()
	}
	if m.pluginReg != nil {
		m.pluginReg.Close()
	}
	return nil
}

// AdvancedAnalysisResult represents the result of advanced analysis
type AdvancedAnalysisResult struct {
	StartTime       time.Time                         `json:"start_time"`
	EndTime         time.Time                         `json:"end_time"`
	Duration        time.Duration                     `json:"duration"`
	Config          *config.Config                    `json:"config"`
	StandardResults []*metrics.AnalysisResult         `json:"standard_results,omitempty"`
	ParallelStats   *parallel.ProcessingStats         `json:"parallel_stats,omitempty"`
	StreamingStats  *streaming.StreamResult           `json:"streaming_stats,omitempty"`
	PluginResults   map[string][]plugins.MetricResult `json:"plugin_results,omitempty"`
	CacheStats      *cache.CacheStats                 `json:"cache_stats,omitempty"`
}

// Helper functions

func parseDuration(durationStr string) time.Duration {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		// Return default duration on parse error
		return time.Hour
	}
	return duration
}

func createTextReader(texts []string) io.Reader {
	content := strings.Join(texts, "\n")
	return strings.NewReader(content)
}
