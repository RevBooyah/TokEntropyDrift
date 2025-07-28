package streaming

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/RevBooyah/TokEntropyDrift/internal/metrics"
	"github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
)

// StreamConfig holds configuration for streaming analysis
type StreamConfig struct {
	ChunkSize        int           `json:"chunk_size"`        // Number of lines per chunk
	BufferSize       int           `json:"buffer_size"`       // Buffer size for reading
	MaxMemoryMB      int           `json:"max_memory_mb"`     // Maximum memory usage in MB
	EnableProgress   bool          `json:"enable_progress"`   // Whether to show progress updates
	ProgressInterval int           `json:"progress_interval"` // Progress update interval in chunks
	Timeout          time.Duration `json:"timeout"`           // Timeout for processing
}

// StreamResult represents the result of streaming analysis
type StreamResult struct {
	TotalChunks       int                       `json:"total_chunks"`
	ProcessedChunks   int                       `json:"processed_chunks"`
	FailedChunks      int                       `json:"failed_chunks"`
	TotalLines        int                       `json:"total_lines"`
	ProcessedLines    int                       `json:"processed_lines"`
	StartTime         time.Time                 `json:"start_time"`
	EndTime           time.Time                 `json:"end_time"`
	Duration          time.Duration             `json:"duration"`
	ChunkResults      []*metrics.AnalysisResult `json:"chunk_results"`
	AggregatedMetrics map[string]float64        `json:"aggregated_metrics"`
	Errors            []string                  `json:"errors"`
}

// ProgressCallback is called to report progress during streaming analysis
type ProgressCallback func(chunk int, total int, lines int, duration time.Duration)

// StreamAnalyzer provides streaming analysis capabilities
type StreamAnalyzer struct {
	config StreamConfig
	engine *metrics.Engine
}

// NewStreamAnalyzer creates a new streaming analyzer
func NewStreamAnalyzer(config StreamConfig, engine *metrics.Engine) *StreamAnalyzer {
	// Set reasonable defaults
	if config.ChunkSize <= 0 {
		config.ChunkSize = 1000
	}
	if config.BufferSize <= 0 {
		config.BufferSize = 64 * 1024 // 64KB
	}
	if config.MaxMemoryMB <= 0 {
		config.MaxMemoryMB = 512 // 512MB default
	}
	if config.ProgressInterval <= 0 {
		config.ProgressInterval = 10
	}

	return &StreamAnalyzer{
		config: config,
		engine: engine,
	}
}

// AnalyzeStream analyzes a stream of text data
func (s *StreamAnalyzer) AnalyzeStream(
	ctx context.Context,
	reader io.Reader,
	tokenizer tokenizers.Tokenizer,
	progressCallback ProgressCallback,
) (*StreamResult, error) {

	result := &StreamResult{
		StartTime:         time.Now(),
		ChunkResults:      make([]*metrics.AnalysisResult, 0),
		AggregatedMetrics: make(map[string]float64),
		Errors:            make([]string, 0),
	}

	// Create buffered reader
	bufReader := bufio.NewReaderSize(reader, s.config.BufferSize)

	// Process chunks
	chunkNum := 0
	lineCount := 0

	for {
		select {
		case <-ctx.Done():
			return result, ctx.Err()
		default:
			// Read chunk
			chunk, err := s.readChunk(bufReader)
			if err == io.EOF {
				break
			}
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Error reading chunk %d: %v", chunkNum, err))
				result.FailedChunks++
				chunkNum++
				continue
			}

			if len(chunk) == 0 {
				break
			}

			// Process chunk
			chunkResult, err := s.processChunk(ctx, chunk, tokenizer, chunkNum)
			if err != nil {
				result.Errors = append(result.Errors, fmt.Sprintf("Error processing chunk %d: %v", chunkNum, err))
				result.FailedChunks++
			} else {
				result.ChunkResults = append(result.ChunkResults, chunkResult)
				result.ProcessedChunks++
			}

			lineCount += len(chunk)
			chunkNum++

			// Report progress
			if s.config.EnableProgress && progressCallback != nil && chunkNum%s.config.ProgressInterval == 0 {
				progressCallback(chunkNum, -1, lineCount, time.Since(result.StartTime))
			}
		}
	}

	result.TotalChunks = chunkNum
	result.TotalLines = lineCount
	result.EndTime = time.Now()
	result.Duration = result.EndTime.Sub(result.StartTime)

	// Aggregate metrics
	s.aggregateMetrics(result)

	return result, nil
}

// AnalyzeFile analyzes a file using streaming
func (s *StreamAnalyzer) AnalyzeFile(
	ctx context.Context,
	filePath string,
	tokenizer tokenizers.Tokenizer,
	progressCallback ProgressCallback,
) (*StreamResult, error) {

	file, err := openFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	return s.AnalyzeStream(ctx, file, tokenizer, progressCallback)
}

// readChunk reads a chunk of lines from the reader
func (s *StreamAnalyzer) readChunk(reader *bufio.Reader) ([]string, error) {
	var chunk []string

	for len(chunk) < s.config.ChunkSize {
		line, err := reader.ReadString('\n')
		if err == io.EOF {
			if line != "" {
				chunk = append(chunk, line)
			}
			break
		}
		if err != nil {
			return chunk, err
		}

		// Remove trailing newline
		if len(line) > 0 && line[len(line)-1] == '\n' {
			line = line[:len(line)-1]
		}

		chunk = append(chunk, line)
	}

	return chunk, nil
}

// processChunk processes a single chunk of text
func (s *StreamAnalyzer) processChunk(
	ctx context.Context,
	chunk []string,
	tokenizer tokenizers.Tokenizer,
	chunkNum int,
) (*metrics.AnalysisResult, error) {

	// Combine chunk lines into a single document
	document := ""
	for i, line := range chunk {
		if i > 0 {
			document += "\n"
		}
		document += line
	}

	// Analyze the chunk
	result, err := s.engine.AnalyzeDocument(ctx, document, tokenizer)
	if err != nil {
		return nil, err
	}

	// Add chunk metadata
	if result.Metadata == nil {
		result.Metadata = make(map[string]interface{})
	}
	result.Metadata["chunk_number"] = chunkNum
	result.Metadata["chunk_size"] = len(chunk)
	result.Metadata["chunk_lines"] = chunk

	return result, nil
}

// aggregateMetrics aggregates metrics across all chunks
func (s *StreamAnalyzer) aggregateMetrics(result *StreamResult) {
	if len(result.ChunkResults) == 0 {
		return
	}

	// Initialize aggregation maps
	metricSums := make(map[string]float64)
	metricCounts := make(map[string]int)

	// Aggregate metrics from all chunks
	for _, chunkResult := range result.ChunkResults {
		for metricName, metric := range chunkResult.Metrics {
			metricSums[metricName] += metric.Value
			metricCounts[metricName]++
		}
	}

	// Calculate averages
	for metricName, sum := range metricSums {
		if count := metricCounts[metricName]; count > 0 {
			result.AggregatedMetrics[metricName] = sum / float64(count)
		}
	}

	// Add summary metrics
	result.AggregatedMetrics["total_chunks"] = float64(result.TotalChunks)
	result.AggregatedMetrics["processed_chunks"] = float64(result.ProcessedChunks)
	result.AggregatedMetrics["failed_chunks"] = float64(result.FailedChunks)
	result.AggregatedMetrics["success_rate"] = float64(result.ProcessedChunks) / float64(result.TotalChunks) * 100
}

// GetConfig returns the current configuration
func (s *StreamAnalyzer) GetConfig() StreamConfig {
	return s.config
}

// SetConfig updates the configuration
func (s *StreamAnalyzer) SetConfig(config StreamConfig) {
	s.config = config
}

// openFile opens a file for reading (placeholder for actual implementation)
func openFile(filePath string) (io.ReadCloser, error) {
	// This would be implemented to actually open files
	// For now, return an error to indicate it needs implementation
	return nil, fmt.Errorf("file opening not implemented yet")
}
