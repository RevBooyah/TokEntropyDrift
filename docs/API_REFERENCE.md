# TokEntropyDrift API Reference

This document provides comprehensive API documentation for TokEntropyDrift, including all public interfaces, types, and functions.

## Table of Contents

1. [Core Types](#core-types)
2. [Tokenizer Interface](#tokenizer-interface)
3. [Metrics Engine](#metrics-engine)
4. [Advanced Features](#advanced-features)
5. [Configuration](#configuration)
6. [Plugin System](#plugin-system)
7. [Error Handling](#error-handling)
8. [Examples](#examples)

## Core Types

### Token

Represents a single token with metadata.

```go
type Token struct {
    Text      string            `json:"text"`
    ID        int               `json:"id"`
    StartPos  int               `json:"start_pos"`
    EndPos    int               `json:"end_pos"`
    Metadata  map[string]string `json:"metadata,omitempty"`
}
```

**Fields:**
- `Text`: The actual token text
- `ID`: Unique identifier for the token
- `StartPos`: Starting position in the original text
- `EndPos`: Ending position in the original text
- `Metadata`: Additional token metadata

### TokenizationResult

Represents the complete result of tokenizing a document.

```go
type TokenizationResult struct {
    Document  string                 `json:"document"`
    Tokens    []Token                `json:"tokens"`
    Tokenizer string                 `json:"tokenizer"`
    Metadata  map[string]interface{} `json:"metadata,omitempty"`
}
```

**Fields:**
- `Document`: The original input text
- `Tokens`: Array of tokens
- `Tokenizer`: Name of the tokenizer used
- `Metadata`: Additional tokenization metadata

### AnalysisResult

Represents the complete analysis results for a document.

```go
type AnalysisResult struct {
    Document       string                    `json:"document"`
    TokenizerName  string                    `json:"tokenizer_name"`
    TokenCount     int                       `json:"token_count"`
    Metrics        map[string]MetricResult   `json:"metrics"`
    Tokenization   *TokenizationResult       `json:"tokenization"`
    Metadata       map[string]interface{}    `json:"metadata,omitempty"`
}
```

**Fields:**
- `Document`: The original input text
- `TokenizerName`: Name of the tokenizer used
- `TokenCount`: Number of tokens generated
- `Metrics`: Map of calculated metrics
- `Tokenization`: Complete tokenization result
- `Metadata`: Additional analysis metadata

### MetricResult

Represents a single metric calculation result.

```go
type MetricResult struct {
    MetricName    string                 `json:"metric_name"`
    TokenizerName string                 `json:"tokenizer_name"`
    Value         float64                `json:"value"`
    Metadata      map[string]interface{} `json:"metadata,omitempty"`
}
```

**Fields:**
- `MetricName`: Name of the metric
- `TokenizerName`: Name of the tokenizer
- `Value`: Calculated metric value
- `Metadata`: Additional metric metadata

## Tokenizer Interface

### Tokenizer

The main interface that all tokenizers must implement.

```go
type Tokenizer interface {
    // Name returns the name of the tokenizer
    Name() string
    
    // Type returns the type of the tokenizer (bpe, spiece, wordpiece, custom)
    Type() string
    
    // Initialize prepares the tokenizer for use
    Initialize(config TokenizerConfig) error
    
    // Tokenize tokenizes a single document
    Tokenize(ctx context.Context, text string) (*TokenizationResult, error)
    
    // TokenizeBatch tokenizes multiple documents
    TokenizeBatch(ctx context.Context, texts []string) ([]*TokenizationResult, error)
    
    // GetVocabSize returns the vocabulary size of the tokenizer
    GetVocabSize() (int, error)
    
    // Close cleans up any resources used by the tokenizer
    Close() error
}
```

### TokenizerConfig

Configuration for a tokenizer.

```go
type TokenizerConfig struct {
    Name       string            `json:"name"`
    Type       string            `json:"type"`
    Command    string            `json:"command,omitempty"`
    LibraryPath string           `json:"library_path,omitempty"`
    VocabFile  string            `json:"vocab_file,omitempty"`
    ModelFile  string            `json:"model_file,omitempty"`
    Parameters map[string]string `json:"parameters,omitempty"`
}
```

### Usage Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
)

func main() {
    // Create a mock tokenizer
    tokenizer := tokenizers.NewMockTokenizer()
    
    // Initialize with configuration
    config := tokenizers.TokenizerConfig{
        Name: "mock",
        Type: "custom",
    }
    
    if err := tokenizer.Initialize(config); err != nil {
        log.Fatal(err)
    }
    defer tokenizer.Close()
    
    // Tokenize text
    ctx := context.Background()
    result, err := tokenizer.Tokenize(ctx, "Hello, world!")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Tokenized %d tokens\n", len(result.Tokens))
    for _, token := range result.Tokens {
        fmt.Printf("Token: %s (ID: %d)\n", token.Text, token.ID)
    }
}
```

## Metrics Engine

### Engine

The main metrics calculation engine.

```go
type Engine struct {
    config EngineConfig
}

func NewEngine(config EngineConfig) *Engine
```

### EngineConfig

Configuration for the metrics engine.

```go
type EngineConfig struct {
    EntropyWindowSize int  `json:"entropy_window_size"`
    NormalizeEntropy  bool `json:"normalize_entropy"`
    CompressionRatio  bool `json:"compression_ratio"`
    DriftDetection    bool `json:"drift_detection"`
}
```

### Key Methods

#### AnalyzeDocument

Performs complete analysis on a single document.

```go
func (e *Engine) AnalyzeDocument(ctx context.Context, document string, tokenizer Tokenizer) (*AnalysisResult, error)
```

#### AnalyzeBatch

Performs analysis on multiple documents.

```go
func (e *Engine) AnalyzeBatch(ctx context.Context, documents []string, tokenizer Tokenizer) ([]*AnalysisResult, error)
```

#### CompareTokenizers

Compares multiple tokenizers on the same document.

```go
func (e *Engine) CompareTokenizers(ctx context.Context, document string, tokenizers []Tokenizer) (map[string]interface{}, error)
```

### Usage Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/RevBooyah/TokEntropyDrift/internal/metrics"
    "github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
)

func main() {
    // Create metrics engine
    config := metrics.EngineConfig{
        EntropyWindowSize: 100,
        NormalizeEntropy:  true,
        CompressionRatio:  true,
        DriftDetection:    true,
    }
    engine := metrics.NewEngine(config)
    
    // Create tokenizer
    tokenizer := tokenizers.NewMockTokenizer()
    tokenizerConfig := tokenizers.TokenizerConfig{
        Name: "mock",
        Type: "custom",
    }
    tokenizer.Initialize(tokenizerConfig)
    defer tokenizer.Close()
    
    // Analyze document
    ctx := context.Background()
    result, err := engine.AnalyzeDocument(ctx, "Hello, world!", tokenizer)
    if err != nil {
        log.Fatal(err)
    }
    
    // Print results
    fmt.Printf("Document: %s\n", result.Document)
    fmt.Printf("Tokenizer: %s\n", result.TokenizerName)
    fmt.Printf("Token Count: %d\n", result.TokenCount)
    
    for metricName, metric := range result.Metrics {
        fmt.Printf("%s: %.4f\n", metricName, metric.Value)
    }
}
```

## Advanced Features

### Cache System

#### Cache

In-memory cache with TTL and size limits.

```go
type Cache struct {
    config CacheConfig
    data   map[string]CacheEntry
    mu     sync.RWMutex
    stats  CacheStats
    stop   chan struct{}
}

func NewCache(config CacheConfig) *Cache
```

#### CacheConfig

Configuration for the cache.

```go
type CacheConfig struct {
    MaxSize         int           `json:"max_size"`
    TTL             time.Duration `json:"ttl"`
    CleanupInterval time.Duration `json:"cleanup_interval"`
    EnableStats     bool          `json:"enable_stats"`
}
```

#### CacheStats

Cache performance statistics.

```go
type CacheStats struct {
    Hits      int64 `json:"hits"`
    Misses    int64 `json:"misses"`
    Evictions int64 `json:"evictions"`
    Size      int   `json:"size"`
    MaxSize   int   `json:"max_size"`
}
```

#### Key Methods

```go
// Get retrieves a value from the cache
func (c *Cache) Get(key string) (interface{}, bool)

// Set stores a value in the cache
func (c *Cache) Set(key string, value interface{})

// GetStats returns cache statistics
func (c *Cache) GetStats() CacheStats

// Close stops the cache and cleans up resources
func (c *Cache) Close()
```

### Parallel Processing

#### Processor

Generic parallel processing framework.

```go
type Processor struct {
    config ProcessorConfig
    stats  ProcessingStats
}

func NewProcessor(config ProcessorConfig) *Processor
```

#### ProcessorConfig

Configuration for parallel processing.

```go
type ProcessorConfig struct {
    MaxWorkers    int           `json:"max_workers"`
    BatchSize     int           `json:"batch_size"`
    Timeout       time.Duration `json:"timeout"`
    EnableMetrics bool          `json:"enable_metrics"`
}
```

#### ProcessingStats

Parallel processing statistics.

```go
type ProcessingStats struct {
    TotalItems     int           `json:"total_items"`
    ProcessedItems int           `json:"processed_items"`
    FailedItems    int           `json:"failed_items"`
    StartTime      time.Time     `json:"start_time"`
    EndTime        time.Time     `json:"end_time"`
    Duration       time.Duration `json:"duration"`
    WorkersUsed    int           `json:"workers_used"`
}
```

#### Key Methods

```go
// ProcessTokenizations processes tokenization in parallel
func (p *Processor) ProcessTokenizations(ctx context.Context, texts []string, tokenizer Tokenizer) ([]*TokenizationResult, []error, ProcessingStats)

// GetStats returns the current processing statistics
func (p *Processor) GetStats() ProcessingStats
```

### Streaming Analysis

#### StreamAnalyzer

Memory-efficient streaming analysis.

```go
type StreamAnalyzer struct {
    config StreamConfig
    engine *metrics.Engine
}

func NewStreamAnalyzer(config StreamConfig, engine *metrics.Engine) *StreamAnalyzer
```

#### StreamConfig

Configuration for streaming analysis.

```go
type StreamConfig struct {
    ChunkSize        int           `json:"chunk_size"`
    BufferSize       int           `json:"buffer_size"`
    MaxMemoryMB      int           `json:"max_memory_mb"`
    EnableProgress   bool          `json:"enable_progress"`
    ProgressInterval int           `json:"progress_interval"`
    Timeout          time.Duration `json:"timeout"`
}
```

#### StreamResult

Result of streaming analysis.

```go
type StreamResult struct {
    TotalChunks      int                    `json:"total_chunks"`
    ProcessedChunks  int                    `json:"processed_chunks"`
    FailedChunks     int                    `json:"failed_chunks"`
    TotalLines       int                    `json:"total_lines"`
    ProcessedLines   int                    `json:"processed_lines"`
    StartTime        time.Time              `json:"start_time"`
    EndTime          time.Time              `json:"end_time"`
    Duration         time.Duration          `json:"duration"`
    ChunkResults     []*AnalysisResult      `json:"chunk_results"`
    AggregatedMetrics map[string]float64    `json:"aggregated_metrics"`
    Errors           []string               `json:"errors"`
}
```

#### Key Methods

```go
// AnalyzeStream analyzes a stream of text data
func (s *StreamAnalyzer) AnalyzeStream(ctx context.Context, reader io.Reader, tokenizer Tokenizer, progressCallback ProgressCallback) (*StreamResult, error)

// GetConfig returns the current configuration
func (s *StreamAnalyzer) GetConfig() StreamConfig
```



## Configuration

### Config

Main application configuration.

```go
type Config struct {
    Input         InputConfig         `mapstructure:"input"`
    Tokenizers    TokenizerConfig     `mapstructure:"tokenizers"`
    Analysis      AnalysisConfig      `mapstructure:"analysis"`
    Cache         CacheConfig         `mapstructure:"cache"`
    Parallel      ParallelConfig      `mapstructure:"parallel"`
    Streaming     StreamingConfig     `mapstructure:"streaming"`
    Plugins       PluginsConfig       `mapstructure:"plugins"`
    Output        OutputConfig        `mapstructure:"output"`
    Visualization VisualizationConfig `mapstructure:"visualization"`
    Server        ServerConfig        `mapstructure:"server"`
    Logging       LoggingConfig       `mapstructure:"logging"`
}
```

### Configuration Subtypes

#### InputConfig
```go
type InputConfig struct {
    SourcePaths []string `mapstructure:"source_paths"`
    FileType    string   `mapstructure:"file_type"`
}
```

#### AnalysisConfig
```go
type AnalysisConfig struct {
    EntropyWindowSize int  `mapstructure:"entropy_window_size"`
    NormalizeEntropy  bool `mapstructure:"normalize_entropy"`
    CompressionRatio  bool `mapstructure:"compression_ratio"`
    DriftDetection    bool `mapstructure:"drift_detection"`
}
```

#### OutputConfig
```go
type OutputConfig struct {
    Directory    string `mapstructure:"directory"`
    Format       string `mapstructure:"format"`
    IncludeLogs  bool   `mapstructure:"include_logs"`
    TimestampDir bool   `mapstructure:"timestamp_dir"`
}
```

### Loading Configuration

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/RevBooyah/TokEntropyDrift/internal/config"
)

func main() {
    // Load configuration from file
    cfg, err := config.LoadConfig("ted.config.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Access configuration
    fmt.Printf("Analysis window size: %d\n", cfg.Analysis.EntropyWindowSize)
    fmt.Printf("Cache enabled: %t\n", cfg.Cache.Enabled)
    fmt.Printf("Parallel workers: %d\n", cfg.Parallel.MaxWorkers)
}
```

## Plugin System

### Plugin Interface

Interface that all plugins must implement.

```go
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
```

### PluginInfo

Plugin metadata.

```go
type PluginInfo struct {
    Name        string            `json:"name"`
    Version     string            `json:"version"`
    Description string            `json:"description"`
    Author      string            `json:"author"`
    Tags        []string          `json:"tags"`
    Metadata    map[string]string `json:"metadata,omitempty"`
}
```

### AnalysisContext

Context for metric calculations.

```go
type AnalysisContext struct {
    Document       string                        `json:"document"`
    Tokenization   *TokenizationResult           `json:"tokenization"`
    TokenizerName  string                        `json:"tokenizer_name"`
    Config         map[string]interface{}        `json:"config"`
    Context        context.Context               `json:"-"`
}
```

### Registry

Plugin registration and management.

```go
type Registry struct {
    plugins map[string]Plugin
    configs map[string]map[string]interface{}
    mu      sync.RWMutex
}

func NewRegistry() *Registry
```

#### Key Methods

```go
// Register adds a plugin to the registry
func (r *Registry) Register(plugin Plugin) error

// Get retrieves a plugin by name
func (r *Registry) Get(name string) (Plugin, error)

// List returns all registered plugin names
func (r *Registry) List() []string

// ExecuteMetrics runs metric calculations for all plugins
func (r *Registry) ExecuteMetrics(ctx *AnalysisContext) (map[string][]MetricResult, error)

// Close cleans up all plugins
func (r *Registry) Close() error
```

### Plugin Development Example

```go
package main

import (
    "math"
    "github.com/RevBooyah/TokEntropyDrift/internal/plugins"
)

// MyCustomPlugin is a custom plugin example
type MyCustomPlugin struct {
    *plugins.BasePlugin
}

func NewMyCustomPlugin() *MyCustomPlugin {
    info := plugins.PluginInfo{
        Name:        "my_custom_plugin",
        Version:     "1.0.0",
        Description: "A custom plugin example",
        Author:      "Your Name",
        Tags:        []string{"example", "custom"},
    }
    
    return &MyCustomPlugin{
        BasePlugin: plugins.NewBasePlugin(info),
    }
}

func (p *MyCustomPlugin) CalculateMetrics(ctx *plugins.AnalysisContext) ([]plugins.MetricResult, error) {
    if ctx.Tokenization == nil || len(ctx.Tokenization.Tokens) == 0 {
        return []plugins.MetricResult{}, nil
    }
    
    tokens := ctx.Tokenization.Tokens
    
    // Calculate custom metric
    totalLength := 0
    for _, token := range tokens {
        totalLength += len(token.Text)
    }
    
    avgLength := float64(totalLength) / float64(len(tokens))
    
    return []plugins.MetricResult{
        {
            Name:  "average_token_length",
            Value: avgLength,
            Unit:  "characters",
        },
    }, nil
}

func (p *MyCustomPlugin) ValidateConfig(config map[string]interface{}) error {
    // Add validation logic here
    return nil
}

// Usage
func main() {
    registry := plugins.NewRegistry()
    
    plugin := NewMyCustomPlugin()
    registry.Register(plugin)
    
    // Use the plugin...
}
```

## Error Handling

### Common Error Types

```go
// Tokenizer errors
var (
    ErrTokenizerNotFound = errors.New("tokenizer not found")
    ErrTokenizerInit     = errors.New("failed to initialize tokenizer")
    ErrTokenization      = errors.New("tokenization failed")
)

// Configuration errors
var (
    ErrConfigNotFound = errors.New("configuration file not found")
    ErrConfigInvalid  = errors.New("invalid configuration")
)

// Plugin errors
var (
    ErrPluginNotFound = errors.New("plugin not found")
    ErrPluginInit     = errors.New("failed to initialize plugin")
    ErrPluginExec     = errors.New("plugin execution failed")
)
```

### Error Handling Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/RevBooyah/TokEntropyDrift/internal/metrics"
    "github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
)

func main() {
    // Create engine
    engine := metrics.NewEngine(metrics.EngineConfig{})
    
    // Create tokenizer
    tokenizer := tokenizers.NewMockTokenizer()
    
    // Initialize with error handling
    if err := tokenizer.Initialize(tokenizers.TokenizerConfig{
        Name: "mock",
        Type: "custom",
    }); err != nil {
        log.Fatalf("Failed to initialize tokenizer: %v", err)
    }
    defer tokenizer.Close()
    
    // Analyze with error handling
    ctx := context.Background()
    result, err := engine.AnalyzeDocument(ctx, "Hello, world!", tokenizer)
    if err != nil {
        log.Fatalf("Analysis failed: %v", err)
    }
    
    fmt.Printf("Analysis successful: %d tokens\n", result.TokenCount)
}
```

## Examples

### Complete Analysis Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/RevBooyah/TokEntropyDrift/internal/config"
    "github.com/RevBooyah/TokEntropyDrift/internal/metrics"
    "github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
)

func main() {
    // Load configuration
    cfg, err := config.LoadConfig("ted.config.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create metrics engine
    engine := metrics.NewEngine(metrics.EngineConfig{
        EntropyWindowSize: cfg.Analysis.EntropyWindowSize,
        NormalizeEntropy:  cfg.Analysis.NormalizeEntropy,
        CompressionRatio:  cfg.Analysis.CompressionRatio,
        DriftDetection:    cfg.Analysis.DriftDetection,
    })
    
    // Create and initialize tokenizer
    mockTokenizer := tokenizers.NewMockTokenizer()
    mockTokenizer.Initialize(tokenizers.TokenizerConfig{
        Name: "mock",
        Type: "custom",
    })
    defer mockTokenizer.Close()
    
    // Prepare test data
    texts := []string{
        "The quick brown fox jumps over the lazy dog.",
        "Machine learning models process text efficiently.",
        "Hello, world! This is a test.",
    }
    
    // Run analysis on each text
    ctx := context.Background()
    startTime := time.Now()
    
    for i, text := range texts {
        result, err := engine.AnalyzeDocument(ctx, text, mockTokenizer)
        if err != nil {
            log.Fatal(err)
        }
        
        fmt.Printf("Text %d: %d tokens\n", i+1, result.TokenCount)
        for metricName, metric := range result.Metrics {
            fmt.Printf("  %s: %.4f\n", metricName, metric.Value)
        }
    }
    
    duration := time.Since(startTime)
    fmt.Printf("Analysis completed in %v\n", duration)
    fmt.Printf("Processed %d texts\n", len(texts))
}
```

### Plugin Development Example

```go
package main

import (
    "context"
    "fmt"
    "math"
    
    "github.com/RevBooyah/TokEntropyDrift/internal/plugins"
)

// TokenComplexityPlugin analyzes token complexity
type TokenComplexityPlugin struct {
    *plugins.BasePlugin
}

func NewTokenComplexityPlugin() *TokenComplexityPlugin {
    info := plugins.PluginInfo{
        Name:        "token_complexity",
        Version:     "1.0.0",
        Description: "Analyzes token complexity patterns",
        Author:      "Developer",
        Tags:        []string{"analysis", "complexity"},
    }
    
    return &TokenComplexityPlugin{
        BasePlugin: plugins.NewBasePlugin(info),
    }
}

func (p *TokenComplexityPlugin) CalculateMetrics(ctx *plugins.AnalysisContext) ([]plugins.MetricResult, error) {
    if ctx.Tokenization == nil || len(ctx.Tokenization.Tokens) == 0 {
        return []plugins.MetricResult{}, nil
    }
    
    tokens := ctx.Tokenization.Tokens
    
    // Calculate complexity metrics
    var lengths []int
    var hasUpperCase, hasLowerCase, hasDigits, hasSpecial int
    
    for _, token := range tokens {
        length := len(token.Text)
        lengths = append(lengths, length)
        
        for _, char := range token.Text {
            switch {
            case char >= 'A' && char <= 'Z':
                hasUpperCase++
            case char >= 'a' && char <= 'z':
                hasLowerCase++
            case char >= '0' && char <= '9':
                hasDigits++
            default:
                hasSpecial++
            }
        }
    }
    
    // Calculate statistics
    totalTokens := len(tokens)
    avgLength := float64(sum(lengths)) / float64(totalTokens)
    
    complexity := float64(hasUpperCase+hasLowerCase+hasDigits+hasSpecial) / float64(totalTokens)
    
    return []plugins.MetricResult{
        {
            Name:  "average_token_length",
            Value: avgLength,
            Unit:  "characters",
        },
        {
            Name:  "token_complexity_score",
            Value: complexity,
            Unit:  "score",
        },
        {
            Name:  "uppercase_ratio",
            Value: float64(hasUpperCase) / float64(totalTokens),
            Unit:  "ratio",
        },
        {
            Name:  "lowercase_ratio",
            Value: float64(hasLowerCase) / float64(totalTokens),
            Unit:  "ratio",
        },
        {
            Name:  "digit_ratio",
            Value: float64(hasDigits) / float64(totalTokens),
            Unit:  "ratio",
        },
        {
            Name:  "special_char_ratio",
            Value: float64(hasSpecial) / float64(totalTokens),
            Unit:  "ratio",
        },
    }, nil
}

func sum(nums []int) int {
    total := 0
    for _, num := range nums {
        total += num
    }
    return total
}

func (p *TokenComplexityPlugin) ValidateConfig(config map[string]interface{}) error {
    // Add validation logic here
    return nil
}

// Usage example
func main() {
    registry := plugins.NewRegistry()
    
    plugin := NewTokenComplexityPlugin()
    registry.Register(plugin)
    
    // Configure plugin
    registry.Configure("token_complexity", map[string]interface{}{
        "enabled": true,
    })
    
    fmt.Println("Token complexity plugin registered successfully!")
}
```

## Best Practices

### 1. Resource Management

Always close resources properly:

```go
tokenizer := tokenizers.NewMockTokenizer()
defer tokenizer.Close()

manager, err := advanced.NewAdvancedManager(cfg, engine)
if err != nil {
    return err
}
defer manager.Close()
```

### 2. Error Handling

Use proper error handling and logging:

```go
result, err := engine.AnalyzeDocument(ctx, document, tokenizer)
if err != nil {
    log.Printf("Analysis failed: %v", err)
    return fmt.Errorf("analysis failed: %w", err)
}
```

### 3. Configuration

Use configuration files for flexibility:

```go
cfg, err := config.LoadConfig("ted.config.yaml")
if err != nil {
    log.Fatal(err)
}

// Use configuration values
engine := metrics.NewEngine(metrics.EngineConfig{
    EntropyWindowSize: cfg.Analysis.EntropyWindowSize,
    NormalizeEntropy:  cfg.Analysis.NormalizeEntropy,
})
```

### 4. Context Usage

Always use context for cancellation and timeouts:

```go
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()

result, err := engine.AnalyzeDocument(ctx, document, tokenizer)
```

### 5. Plugin Development

Follow the plugin interface correctly:

```go
// Always implement all interface methods
type MyPlugin struct {
    *plugins.BasePlugin
}

func (p *MyPlugin) CalculateMetrics(ctx *plugins.AnalysisContext) ([]plugins.MetricResult, error) {
    // Your implementation
}

func (p *MyPlugin) ValidateConfig(config map[string]interface{}) error {
    // Validation logic
    return nil
}
```

## Conclusion

This API reference provides comprehensive documentation for using TokEntropyDrift programmatically. The library is designed to be flexible, extensible, and easy to integrate into existing applications.

For more examples and use cases, see the tutorials and examples directories. For questions and support, please refer to the project documentation or community resources. 