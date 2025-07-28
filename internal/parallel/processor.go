package parallel

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
)

// ProcessorConfig holds configuration for parallel processing
type ProcessorConfig struct {
	MaxWorkers    int           `json:"max_workers"`    // Maximum number of worker goroutines
	BatchSize     int           `json:"batch_size"`     // Number of items per batch
	Timeout       time.Duration `json:"timeout"`        // Timeout for processing
	EnableMetrics bool          `json:"enable_metrics"` // Whether to collect processing metrics
}

// ProcessingStats holds statistics about parallel processing
type ProcessingStats struct {
	TotalItems     int           `json:"total_items"`
	ProcessedItems int           `json:"processed_items"`
	FailedItems    int           `json:"failed_items"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	Duration       time.Duration `json:"duration"`
	WorkersUsed    int           `json:"workers_used"`
}

// Processor provides parallel processing capabilities
type Processor struct {
	config ProcessorConfig
	stats  ProcessingStats
}

// NewProcessor creates a new parallel processor
func NewProcessor(config ProcessorConfig) *Processor {
	// Set reasonable defaults
	if config.MaxWorkers <= 0 {
		config.MaxWorkers = runtime.NumCPU()
	}
	if config.BatchSize <= 0 {
		config.BatchSize = 100
	}

	return &Processor{
		config: config,
	}
}

// processItems processes items in parallel using the provided function
func (p *Processor) processItems(
	ctx context.Context,
	items []string,
	processFunc func(context.Context, string) (*tokenizers.TokenizationResult, error),
) ([]*tokenizers.TokenizationResult, []error, ProcessingStats) {

	p.stats = ProcessingStats{
		TotalItems:  len(items),
		StartTime:   time.Now(),
		WorkersUsed: p.config.MaxWorkers,
	}

	if len(items) == 0 {
		p.stats.EndTime = time.Now()
		p.stats.Duration = p.stats.EndTime.Sub(p.stats.StartTime)
		return []*tokenizers.TokenizationResult{}, []error{}, p.stats
	}

	// Create channels for results and errors
	resultChan := make(chan *tokenizers.TokenizationResult, len(items))
	errorChan := make(chan error, len(items))

	// Create context with timeout if specified
	if p.config.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, p.config.Timeout)
		defer cancel()
	}

	// Start worker goroutines
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, p.config.MaxWorkers)

	for i := 0; i < len(items); i += p.config.BatchSize {
		end := i + p.config.BatchSize
		if end > len(items) {
			end = len(items)
		}

		batch := items[i:end]

		wg.Add(1)
		go func(batchItems []string) {
			defer wg.Done()
			semaphore <- struct{}{}        // Acquire semaphore
			defer func() { <-semaphore }() // Release semaphore

			for _, item := range batchItems {
				select {
				case <-ctx.Done():
					errorChan <- ctx.Err()
					return
				default:
					result, err := processFunc(ctx, item)
					if err != nil {
						errorChan <- err
						p.stats.FailedItems++
					} else {
						resultChan <- result
						p.stats.ProcessedItems++
					}
				}
			}
		}(batch)
	}

	// Wait for all workers to complete
	wg.Wait()
	close(resultChan)
	close(errorChan)

	// Collect results and errors
	var results []*tokenizers.TokenizationResult
	var errors []error

	for result := range resultChan {
		results = append(results, result)
	}

	for err := range errorChan {
		errors = append(errors, err)
	}

	p.stats.EndTime = time.Now()
	p.stats.Duration = p.stats.EndTime.Sub(p.stats.StartTime)

	return results, errors, p.stats
}

// ProcessTokenizations processes tokenization in parallel
func (p *Processor) ProcessTokenizations(
	ctx context.Context,
	texts []string,
	tokenizer tokenizers.Tokenizer,
) ([]*tokenizers.TokenizationResult, []error, ProcessingStats) {

	processFunc := func(ctx context.Context, text string) (*tokenizers.TokenizationResult, error) {
		return tokenizer.Tokenize(ctx, text)
	}

	results, errors, stats := p.processItems(ctx, texts, processFunc)
	return results, errors, stats
}

// ProcessTokenizationsBatch processes tokenization in parallel using batch processing
func (p *Processor) ProcessTokenizationsBatch(
	ctx context.Context,
	texts []string,
	tokenizer tokenizers.Tokenizer,
) ([]*tokenizers.TokenizationResult, []error, ProcessingStats) {

	// For now, just use the regular processing method
	// TODO: Implement proper batch processing
	return p.ProcessTokenizations(ctx, texts, tokenizer)
}

// createBatches splits a slice into batches of the specified size
func (p *Processor) createBatches(items []string, batchSize int) [][]string {
	var batches [][]string
	for i := 0; i < len(items); i += batchSize {
		end := i + batchSize
		if end > len(items) {
			end = len(items)
		}
		batches = append(batches, items[i:end])
	}
	return batches
}

// GetStats returns the current processing statistics
func (p *Processor) GetStats() ProcessingStats {
	return p.stats
}

// GetOptimalWorkerCount returns the optimal number of workers based on system resources
func GetOptimalWorkerCount() int {
	cpuCount := runtime.NumCPU()
	// Use 75% of available CPUs to avoid overwhelming the system
	return int(float64(cpuCount) * 0.75)
}

// GetOptimalBatchSize returns the optimal batch size based on item count and worker count
func GetOptimalBatchSize(itemCount, workerCount int) int {
	if itemCount <= workerCount {
		return 1
	}

	// Aim for 2-4 batches per worker for good load balancing
	batchSize := itemCount / (workerCount * 3)
	if batchSize < 1 {
		batchSize = 1
	}
	if batchSize > 1000 {
		batchSize = 1000
	}

	return batchSize
}
