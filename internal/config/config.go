package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Config represents the main application configuration
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

// InputConfig holds input file configuration
type InputConfig struct {
	SourcePaths []string `mapstructure:"source_paths"`
	FileType    string   `mapstructure:"file_type"`
}

// TokenizerConfig holds tokenizer configuration
type TokenizerConfig struct {
	Enabled []string                `mapstructure:"enabled"`
	Configs map[string]TokenizerDef `mapstructure:"configs"`
}

// TokenizerDef represents a tokenizer definition
type TokenizerDef struct {
	Type        string            `mapstructure:"type"`
	LibraryPath string            `mapstructure:"library_path"`
	Parameters  map[string]string `mapstructure:"parameters"`
}

// AnalysisConfig holds analysis parameters
type AnalysisConfig struct {
	EntropyWindowSize int  `mapstructure:"entropy_window_size"`
	NormalizeEntropy  bool `mapstructure:"normalize_entropy"`
	CompressionRatio  bool `mapstructure:"compression_ratio"`
	DriftDetection    bool `mapstructure:"drift_detection"`
}

// CacheConfig holds caching configuration
type CacheConfig struct {
	Enabled         bool   `mapstructure:"enabled"`
	MaxSize         int    `mapstructure:"max_size"`
	TTL             string `mapstructure:"ttl"`
	CleanupInterval string `mapstructure:"cleanup_interval"`
	EnableStats     bool   `mapstructure:"enable_stats"`
}

// ParallelConfig holds parallel processing configuration
type ParallelConfig struct {
	Enabled       bool   `mapstructure:"enabled"`
	MaxWorkers    int    `mapstructure:"max_workers"`
	BatchSize     int    `mapstructure:"batch_size"`
	Timeout       string `mapstructure:"timeout"`
	EnableMetrics bool   `mapstructure:"enable_metrics"`
}

// StreamingConfig holds streaming analysis configuration
type StreamingConfig struct {
	Enabled          bool   `mapstructure:"enabled"`
	ChunkSize        int    `mapstructure:"chunk_size"`
	BufferSize       int    `mapstructure:"buffer_size"`
	MaxMemoryMB      int    `mapstructure:"max_memory_mb"`
	EnableProgress   bool   `mapstructure:"enable_progress"`
	ProgressInterval int    `mapstructure:"progress_interval"`
	Timeout          string `mapstructure:"timeout"`
}

// PluginsConfig holds plugin system configuration
type PluginsConfig struct {
	Enabled         bool                              `mapstructure:"enabled"`
	AutoLoad        bool                              `mapstructure:"auto_load"`
	PluginDirectory string                            `mapstructure:"plugin_directory"`
	Configs         map[string]map[string]interface{} `mapstructure:"configs"`
}

// OutputConfig holds output configuration
type OutputConfig struct {
	Directory    string `mapstructure:"directory"`
	Format       string `mapstructure:"format"`
	IncludeLogs  bool   `mapstructure:"include_logs"`
	TimestampDir bool   `mapstructure:"timestamp_dir"`
}

// VisualizationConfig holds visualization settings
type VisualizationConfig struct {
	Theme       string `mapstructure:"theme"`
	ImageSize   string `mapstructure:"image_size"`
	FileType    string `mapstructure:"file_type"`
	Interactive bool   `mapstructure:"interactive"`
}

// ServerConfig holds web server configuration
type ServerConfig struct {
	Port int    `mapstructure:"port"`
	Host string `mapstructure:"host"`
}

// LoggingConfig holds logging configuration
type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
	File   string `mapstructure:"file"`
}

// LoadConfig loads configuration from file and environment
func LoadConfig(configPath string) (*Config, error) {
	// Set default values
	config := &Config{
		Input: InputConfig{
			FileType: "txt",
		},
		Tokenizers: TokenizerConfig{
			Enabled: []string{"mock", "gpt2"},
		},
		Analysis: AnalysisConfig{
			EntropyWindowSize: 100,
			NormalizeEntropy:  true,
			CompressionRatio:  true,
			DriftDetection:    true,
		},
		Cache: CacheConfig{
			Enabled:         true,
			MaxSize:         10000,
			TTL:             "1h",
			CleanupInterval: "10m",
			EnableStats:     true,
		},
		Parallel: ParallelConfig{
			Enabled:       true,
			MaxWorkers:    0, // Auto-detect
			BatchSize:     100,
			Timeout:       "30m",
			EnableMetrics: true,
		},
		Streaming: StreamingConfig{
			Enabled:          true,
			ChunkSize:        1000,
			BufferSize:       65536, // 64KB
			MaxMemoryMB:      512,
			EnableProgress:   true,
			ProgressInterval: 10,
			Timeout:          "1h",
		},
		Plugins: PluginsConfig{
			Enabled:         true,
			AutoLoad:        true,
			PluginDirectory: "plugins",
			Configs:         make(map[string]map[string]interface{}),
		},
		Output: OutputConfig{
			Directory:    "output",
			Format:       "csv",
			IncludeLogs:  true,
			TimestampDir: true,
		},
		Visualization: VisualizationConfig{
			Theme:       "light",
			ImageSize:   "medium",
			FileType:    "svg",
			Interactive: true,
		},
		Server: ServerConfig{
			Port: 8080,
			Host: "localhost",
		},
		Logging: LoggingConfig{
			Level:  "info",
			Format: "json",
		},
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(config.Output.Directory, 0755); err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create subdirectories
	subdirs := []string{"uploads", "visualizations", "reports", "logs"}
	for _, subdir := range subdirs {
		path := filepath.Join(config.Output.Directory, subdir)
		if err := os.MkdirAll(path, 0755); err != nil {
			return nil, fmt.Errorf("failed to create subdirectory %s: %w", subdir, err)
		}
	}

	return config, nil
}

// ValidateConfig validates the configuration
func (c *Config) ValidateConfig() error {
	// Validate input configuration
	if len(c.Input.SourcePaths) == 0 && c.Input.FileType == "" {
		return fmt.Errorf("input configuration is incomplete")
	}

	// Validate tokenizer configuration
	if len(c.Tokenizers.Enabled) == 0 {
		return fmt.Errorf("no tokenizers enabled")
	}

	// Validate analysis configuration
	if c.Analysis.EntropyWindowSize <= 0 {
		return fmt.Errorf("entropy window size must be positive")
	}

	// Validate output configuration
	if c.Output.Directory == "" {
		return fmt.Errorf("output directory is required")
	}

	// Validate server configuration
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return fmt.Errorf("invalid server port: %d", c.Server.Port)
	}

	return nil
}

// GetOutputPath returns the full path for a given output file
func (c *Config) GetOutputPath(filename string) string {
	return filepath.Join(c.Output.Directory, filename)
}

// GetUploadPath returns the path for uploaded files
func (c *Config) GetUploadPath() string {
	return filepath.Join(c.Output.Directory, "uploads")
}

// GetVisualizationPath returns the path for visualization files
func (c *Config) GetVisualizationPath() string {
	return filepath.Join(c.Output.Directory, "visualizations")
}

// GetReportPath returns the path for report files
func (c *Config) GetReportPath() string {
	return filepath.Join(c.Output.Directory, "reports")
}

// GetLogPath returns the path for log files
func (c *Config) GetLogPath() string {
	return filepath.Join(c.Output.Directory, "logs")
}
