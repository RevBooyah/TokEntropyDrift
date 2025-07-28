# TokEntropyDrift User Guide

## Table of Contents

1. [Installation](#installation)
2. [Quick Start](#quick-start)
3. [Basic Usage](#basic-usage)
4. [Advanced Features](#advanced-features)
5. [Configuration](#configuration)
6. [Examples and Tutorials](#examples-and-tutorials)
7. [Troubleshooting](#troubleshooting)
8. [API Reference](#api-reference)

## Installation

### Prerequisites

- Go 1.22 or later
- Python 3.8+ (for some tokenizer adapters)
- Git

### Building from Source

```bash
# Clone the repository
git clone https://github.com/RevBooyah/TokEntropyDrift.git
cd tokentropydrift

# Build the binary
go build -o ted cmd/ted/main.go

# Verify installation
./ted --help
```

### Configuration Setup

Copy the example configuration and customize it:

```bash
cp ted.config.yaml.example ted.config.yaml
# Edit ted.config.yaml with your preferences
```

## Quick Start

### 1. Basic Analysis

Analyze a text file with multiple tokenizers:

```bash
./ted analyze examples/english_quotes.txt --tokenizers=gpt2,t5
```

This will:
- Tokenize the text with GPT-2 and T5 tokenizers
- Calculate entropy, compression, and reuse metrics
- Generate a comprehensive analysis report

### 2. Web Dashboard

Launch the interactive web dashboard:

```bash
./ted serve --port=8080
```

Then open `http://localhost:8080` in your browser.

### 3. Visualization

Generate heatmaps and visualizations:

```bash
./ted heatmap examples/tech_stack_descriptions.txt --type=entropy --output=entropy_heatmap.svg
```

## Basic Usage

### Command Structure

All commands follow this pattern:
```bash
./ted <command> [subcommand] [arguments] [flags]
```

### Available Commands

#### `analyze` - Main Analysis Command

```bash
./ted analyze <input-file> [flags]
```

**Flags:**
- `--tokenizers`: Comma-separated list of tokenizers (default: gpt2)
- `--output`: Output directory for results
- `--visualize`: Generate visualizations
- `--format`: Output format (csv, json, markdown)

**Examples:**
```bash
# Basic analysis
./ted analyze my_text.txt

# Multi-tokenizer analysis
./ted analyze my_text.txt --tokenizers=gpt2,bert,t5

# With visualization
./ted analyze my_text.txt --visualize --output=results/
```

#### `serve` - Web Dashboard

```bash
./ted serve [flags]
```

**Flags:**
- `--port`: Port to serve on (default: 8080)
- `--host`: Host to bind to (default: localhost)

**Examples:**
```bash
# Default settings
./ted serve

# Custom port
./ted serve --port=9000

# External access
./ted serve --host=0.0.0.0 --port=8080
```

#### `heatmap` - Generate Visualizations

```bash
./ted heatmap <input-file> [flags]
```

**Flags:**
- `--type`: Heatmap type (entropy, tokens, drift)
- `--output`: Output file path
- `--tokenizers`: Tokenizers to compare

**Examples:**
```bash
# Entropy heatmap
./ted heatmap text.txt --type=entropy --output=entropy.svg

# Token count comparison
./ted heatmap text.txt --type=tokens --tokenizers=gpt2,bert
```

#### `test` - End-to-End Testing

```bash
./ted test <input-file> [flags]
```

**Flags:**
- `--compare-to`: Golden reference file
- `--tokenizers`: Tokenizers to test

**Examples:**
```bash
./ted test examples/english_quotes.txt --compare-to=testdata/golden.csv
```

#### `bench` - Performance Benchmarking

```bash
./ted bench tokenize <input-file> [flags]
```

**Flags:**
- `--tokenizers`: Tokenizers to benchmark
- `--iterations`: Number of iterations

**Examples:**
```bash
./ted bench tokenize large_file.txt --tokenizers=gpt2,t5,bert --iterations=100
```

## Advanced Features

### Caching System

The caching system improves performance by storing tokenization results:

```bash
# Test caching functionality
./ted advanced cache examples/english_quotes.txt
```

**Configuration:**
```yaml
cache:
  enabled: true
  max_size: 10000
  ttl: "1h"
  cleanup_interval: "10m"
  enable_stats: true
```

**Benefits:**
- 90%+ performance improvement for repeated requests
- Automatic cleanup of expired entries
- Cache statistics tracking

### Parallel Processing

Process large datasets efficiently using multiple CPU cores:

```bash
# Test parallel processing
./ted advanced parallel examples/large_dataset.txt
```

**Configuration:**
```yaml
parallel:
  enabled: true
  max_workers: 0  # Auto-detect (75% of CPU cores)
  batch_size: 100
  timeout: "30m"
  enable_metrics: true
```

**Benefits:**
- Linear scaling with CPU cores
- 3-5x improvement for large datasets
- Better resource utilization

### Streaming Analysis

Process large files with minimal memory usage:

```bash
# Test streaming analysis
./ted advanced streaming examples/very_large_file.txt
```

**Configuration:**
```yaml
streaming:
  enabled: true
  chunk_size: 1000
  buffer_size: 65536  # 64KB
  max_memory_mb: 512
  enable_progress: true
  progress_interval: 10
  timeout: "1h"
```

**Benefits:**
- Constant memory usage regardless of file size
- Real-time progress tracking
- Scalable to files of any size

### Plugin System

Extend functionality with custom metrics and analysis:

```bash
# Test plugin system
./ted advanced plugins examples/english_quotes.txt
```

**Configuration:**
```yaml
plugins:
  enabled: true
  auto_load: true
  plugin_directory: "plugins"
  configs:
    token_length_analyzer:
      min_length_threshold: 1
      max_length_threshold: 100
```

**Creating Custom Plugins:**

1. Implement the `Plugin` interface:
```go
type MyPlugin struct {
    *plugins.BasePlugin
}

func (p *MyPlugin) CalculateMetrics(ctx *plugins.AnalysisContext) ([]plugins.MetricResult, error) {
    // Your custom metric calculation logic
    return []plugins.MetricResult{
        {
            Name:  "my_metric",
            Value: 42.0,
            Unit:  "units",
        },
    }, nil
}
```

2. Register your plugin:
```go
registry := plugins.NewRegistry()
registry.Register(&MyPlugin{})
```

## Configuration

### Configuration File Structure

The main configuration file is `ted.config.yaml`:

```yaml
# Input configuration
input:
  source_paths: []
  file_type: "txt"

# Tokenizer configuration
tokenizers:
  enabled: ["mock", "gpt2"]
  configs:
    gpt2:
      type: "bpe"
      library_path: "/usr/local/lib/python3.9/site-packages/tiktoken"
      parameters:
        model: "gpt2"

# Analysis configuration
analysis:
  entropy_window_size: 100
  normalize_entropy: true
  compression_ratio: true
  drift_detection: true

# Advanced Features
cache:
  enabled: true
  max_size: 10000
  ttl: "1h"
  cleanup_interval: "10m"
  enable_stats: true

parallel:
  enabled: true
  max_workers: 0
  batch_size: 100
  timeout: "30m"
  enable_metrics: true

streaming:
  enabled: true
  chunk_size: 1000
  buffer_size: 65536
  max_memory_mb: 512
  enable_progress: true
  progress_interval: 10
  timeout: "1h"

plugins:
  enabled: true
  auto_load: true
  plugin_directory: "plugins"
  configs:
    token_length_analyzer:
      min_length_threshold: 1
      max_length_threshold: 100

# Output configuration
output:
  directory: "output"
  format: "csv"
  include_logs: true
  timestamp_dir: true

# Visualization configuration
visualization:
  theme: "light"
  image_size: "medium"
  file_type: "svg"
  interactive: true

# Server configuration
server:
  port: 8080
  host: "localhost"

# Logging configuration
logging:
  level: "info"
  format: "json"
  file: ""
```

### Environment Variables

You can override configuration with environment variables:

```bash
export TED_SERVER_PORT=9000
export TED_CACHE_ENABLED=true
export TED_PARALLEL_MAX_WORKERS=4
```

### Configuration Validation

The system validates configuration on startup:

```bash
./ted analyze --validate-config
```

## Examples and Tutorials

### Tutorial 1: Basic Tokenization Analysis

**Goal**: Analyze how different tokenizers segment the same text.

**Steps:**
1. Create a test file:
```bash
echo "The quick brown fox jumps over the lazy dog." > test.txt
```

2. Run analysis:
```bash
./ted analyze test.txt --tokenizers=gpt2,bert,t5 --visualize
```

3. Examine results:
```bash
ls output/
# Look for CSV files and visualizations
```

### Tutorial 2: Performance Optimization

**Goal**: Optimize analysis performance for large datasets.

**Steps:**
1. Enable caching:
```yaml
cache:
  enabled: true
  max_size: 50000
```

2. Enable parallel processing:
```yaml
parallel:
  enabled: true
  max_workers: 8
```

3. Run analysis with progress tracking:
```bash
./ted analyze large_file.txt --tokenizers=gpt2,bert
```

### Tutorial 3: Custom Plugin Development

**Goal**: Create a custom metric for token analysis.

**Steps:**
1. Create plugin file `plugins/my_analyzer.go`:
```go
package main

import "github.com/RevBooyah/TokEntropyDrift/internal/plugins"

type MyAnalyzer struct {
    *plugins.BasePlugin
}

func NewMyAnalyzer() *MyAnalyzer {
    return &MyAnalyzer{
        BasePlugin: plugins.NewBasePlugin(plugins.PluginInfo{
            Name:        "my_analyzer",
            Version:     "1.0.0",
            Description: "My custom token analyzer",
        }),
    }
}

func (a *MyAnalyzer) CalculateMetrics(ctx *plugins.AnalysisContext) ([]plugins.MetricResult, error) {
    // Your custom logic here
    return []plugins.MetricResult{
        {
            Name:  "custom_metric",
            Value: 123.45,
            Unit:  "units",
        },
    }, nil
}
```

2. Enable plugin in configuration:
```yaml
plugins:
  enabled: true
  configs:
    my_analyzer:
      parameter1: "value1"
```

3. Test the plugin:
```bash
./ted advanced plugins test.txt
```

### Tutorial 4: Large File Processing

**Goal**: Process a very large file efficiently.

**Steps:**
1. Configure streaming:
```yaml
streaming:
  enabled: true
  chunk_size: 5000
  max_memory_mb: 1024
```

2. Run streaming analysis:
```bash
./ted advanced streaming very_large_file.txt
```

3. Monitor progress and memory usage.

## Troubleshooting

### Common Issues

#### 1. Tokenizer Not Found

**Error**: `tokenizer "gpt2" not found`

**Solution**: Check tokenizer configuration in `ted.config.yaml`:
```yaml
tokenizers:
  enabled: ["gpt2"]
  configs:
    gpt2:
      type: "bpe"
      library_path: "/path/to/tiktoken"
```

#### 2. Memory Issues

**Error**: `out of memory`

**Solution**: Enable streaming analysis:
```yaml
streaming:
  enabled: true
  chunk_size: 1000
  max_memory_mb: 512
```

#### 3. Slow Performance

**Solution**: Enable caching and parallel processing:
```yaml
cache:
  enabled: true
  max_size: 10000

parallel:
  enabled: true
  max_workers: 0  # Auto-detect
```

#### 4. Plugin Loading Errors

**Error**: `failed to load plugin`

**Solution**: Check plugin configuration and file permissions:
```yaml
plugins:
  enabled: true
  plugin_directory: "plugins"
```

### Debug Mode

Enable debug logging:

```bash
export TED_LOGGING_LEVEL=debug
./ted analyze test.txt
```

### Performance Profiling

Profile performance with Go's built-in profiler:

```bash
go build -o ted cmd/ted/main.go
./ted analyze large_file.txt &
go tool pprof http://localhost:6060/debug/pprof/profile
```

## API Reference

### Core Types

#### AnalysisResult
```go
type AnalysisResult struct {
    Document       string
    TokenizerName  string
    TokenCount     int
    Metrics        map[string]MetricResult
    Tokenization   *TokenizationResult
    Metadata       map[string]interface{}
}
```

#### MetricResult
```go
type MetricResult struct {
    MetricName    string
    TokenizerName string
    Value         float64
    Metadata      map[string]interface{}
}
```

### Key Functions

#### AnalyzeDocument
```go
func (e *Engine) AnalyzeDocument(ctx context.Context, document string, tokenizer Tokenizer) (*AnalysisResult, error)
```

#### CompareTokenizers
```go
func (e *Engine) CompareTokenizers(ctx context.Context, document string, tokenizers []Tokenizer) (map[string]interface{}, error)
```

### Plugin Interface

```go
type Plugin interface {
    Info() PluginInfo
    Initialize(config map[string]interface{}) error
    CalculateMetrics(ctx *AnalysisContext) ([]MetricResult, error)
    ValidateConfig(config map[string]interface{}) error
    Cleanup() error
}
```

## Getting Help

- **Documentation**: See the `/docs/` directory for detailed guides
- **Examples**: Check `/examples/` for sample code and data
- **Issues**: Report bugs and feature requests on GitHub
- **Discussions**: Join community discussions for questions and ideas

## Contributing

See `CONTRIBUTING.md` for guidelines on:
- Code style and structure
- Testing requirements
- Documentation standards
- Pull request process 