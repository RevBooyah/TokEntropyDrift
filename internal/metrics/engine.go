package metrics

import (
	"context"
	"fmt"
	"math"

	"github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
)

// MetricResult represents the result of a metric calculation
type MetricResult struct {
	MetricName    string                 `json:"metric_name"`
	TokenizerName string                 `json:"tokenizer_name"`
	Value         float64                `json:"value"`
	Metadata      map[string]interface{} `json:"metadata,omitempty"`
}

// AnalysisResult represents the complete analysis results for a document
type AnalysisResult struct {
	Document      string                         `json:"document"`
	TokenizerName string                         `json:"tokenizer_name"`
	TokenCount    int                            `json:"token_count"`
	Metrics       map[string]MetricResult        `json:"metrics"`
	Tokenization  *tokenizers.TokenizationResult `json:"tokenization"`
	Metadata      map[string]interface{}         `json:"metadata,omitempty"`
}

// Engine handles metric calculations for tokenization analysis
type Engine struct {
	config EngineConfig
}

// EngineConfig holds configuration for the metric engine
type EngineConfig struct {
	EntropyWindowSize int  `json:"entropy_window_size"`
	NormalizeEntropy  bool `json:"normalize_entropy"`
	CompressionRatio  bool `json:"compression_ratio"`
	DriftDetection    bool `json:"drift_detection"`
}

// NewEngine creates a new metric engine with the given configuration
func NewEngine(config EngineConfig) *Engine {
	return &Engine{
		config: config,
	}
}

// AnalyzeDocument performs complete analysis on a single document
func (e *Engine) AnalyzeDocument(ctx context.Context, document string, tokenizer tokenizers.Tokenizer) (*AnalysisResult, error) {
	// Tokenize the document
	tokenization, err := tokenizer.Tokenize(ctx, document)
	if err != nil {
		return nil, fmt.Errorf("error tokenizing document: %w", err)
	}

	// Calculate metrics
	metrics := make(map[string]MetricResult)

	// Token count
	tokenCount := len(tokenization.Tokens)
	metrics["token_count"] = MetricResult{
		MetricName:    "token_count",
		TokenizerName: tokenizer.Name(),
		Value:         float64(tokenCount),
	}

	// Enhanced entropy calculations
	entropyCalc := NewEntropyCalculator(e.config.EntropyWindowSize, e.config.NormalizeEntropy)
	if entropyStats, err := entropyCalc.CalculateEntropyStats(tokenization.Tokens); err == nil {
		for metricName, value := range entropyStats {
			metrics["entropy_"+metricName] = MetricResult{
				MetricName:    "entropy_" + metricName,
				TokenizerName: tokenizer.Name(),
				Value:         value,
			}
		}
	}

	// Enhanced compression calculations
	compressionCalc := NewCompressionCalculator(true)
	if compressionStats, err := compressionCalc.CalculateCompressionStats(document, tokenization.Tokens, 0.0); err == nil {
		for metricName, value := range compressionStats {
			metrics["compression_"+metricName] = MetricResult{
				MetricName:    "compression_" + metricName,
				TokenizerName: tokenizer.Name(),
				Value:         value,
			}
		}
	}

	// Enhanced reuse calculations
	reuseCalc := NewReuseCalculator(true)
	if reuseStats, err := reuseCalc.CalculateReuseStats(tokenization.Tokens); err == nil {
		// Convert interface{} values to float64 for metrics
		for metricName, value := range reuseStats {
			if floatValue, ok := value.(float64); ok {
				metrics["reuse_"+metricName] = MetricResult{
					MetricName:    "reuse_" + metricName,
					TokenizerName: tokenizer.Name(),
					Value:         floatValue,
				}
			}
		}
	}

	return &AnalysisResult{
		Document:      document,
		TokenizerName: tokenizer.Name(),
		TokenCount:    tokenCount,
		Metrics:       metrics,
		Tokenization:  tokenization,
	}, nil
}

// AnalyzeBatch performs analysis on multiple documents
func (e *Engine) AnalyzeBatch(ctx context.Context, documents []string, tokenizer tokenizers.Tokenizer) ([]*AnalysisResult, error) {
	var results []*AnalysisResult

	for _, document := range documents {
		result, err := e.AnalyzeDocument(ctx, document, tokenizer)
		if err != nil {
			return nil, fmt.Errorf("error analyzing document: %w", err)
		}
		results = append(results, result)
	}

	return results, nil
}

// CalculateEntropy calculates Shannon entropy for the given tokens
func (e *Engine) CalculateEntropy(tokens []tokenizers.Token) (float64, error) {
	if len(tokens) == 0 {
		return 0.0, nil
	}

	// Count token frequencies
	tokenFreq := make(map[string]int)
	for _, token := range tokens {
		tokenFreq[token.Text]++
	}

	// Calculate entropy
	entropy := 0.0
	totalTokens := float64(len(tokens))

	for _, freq := range tokenFreq {
		probability := float64(freq) / totalTokens
		if probability > 0 {
			entropy -= probability * math.Log2(probability)
		}
	}

	// Normalize if configured
	if e.config.NormalizeEntropy {
		maxEntropy := math.Log2(float64(len(tokenFreq)))
		if maxEntropy > 0 {
			entropy = entropy / maxEntropy
		}
	}

	return entropy, nil
}

// CalculateRollingEntropy calculates entropy over sliding windows
func (e *Engine) CalculateRollingEntropy(tokens []tokenizers.Token) ([]float64, error) {
	if len(tokens) == 0 {
		return []float64{}, nil
	}

	windowSize := e.config.EntropyWindowSize
	if windowSize <= 0 {
		windowSize = 100 // Default window size
	}

	var rollingEntropy []float64

	for i := 0; i <= len(tokens)-windowSize; i++ {
		windowTokens := tokens[i : i+windowSize]
		entropy, err := e.CalculateEntropy(windowTokens)
		if err != nil {
			return nil, err
		}
		rollingEntropy = append(rollingEntropy, entropy)
	}

	return rollingEntropy, nil
}

// CalculateCompressionRatio calculates the compression ratio
func (e *Engine) CalculateCompressionRatio(originalText string, tokens []tokenizers.Token) (float64, error) {
	if len(originalText) == 0 {
		return 0.0, nil
	}

	// Calculate token representation size (assuming each token ID is 4 bytes)
	tokenSize := len(tokens) * 4

	// Calculate original text size in bytes
	originalSize := len([]byte(originalText))

	// Calculate compression ratio
	compressionRatio := float64(tokenSize) / float64(originalSize)

	return compressionRatio, nil
}

// CalculateTokenReuse calculates token reuse metrics
func (e *Engine) CalculateTokenReuse(tokens []tokenizers.Token) (float64, error) {
	if len(tokens) == 0 {
		return 0.0, nil
	}

	// Count unique tokens
	uniqueTokens := make(map[string]bool)
	for _, token := range tokens {
		uniqueTokens[token.Text] = true
	}

	// Calculate reuse ratio
	uniqueCount := float64(len(uniqueTokens))
	totalCount := float64(len(tokens))

	reuseRatio := 1.0 - (uniqueCount / totalCount)

	return reuseRatio, nil
}

// CalculateDrift calculates drift between two tokenization results
func (e *Engine) CalculateDrift(result1, result2 *tokenizers.TokenizationResult) (float64, error) {
	if result1 == nil || result2 == nil {
		return 0.0, fmt.Errorf("both tokenization results must be provided")
	}

	// Extract token texts
	tokens1 := make([]string, len(result1.Tokens))
	tokens2 := make([]string, len(result2.Tokens))

	for i, token := range result1.Tokens {
		tokens1[i] = token.Text
	}
	for i, token := range result2.Tokens {
		tokens2[i] = token.Text
	}

	// Calculate Jaccard distance as a measure of drift
	drift := e.calculateJaccardDistance(tokens1, tokens2)

	return drift, nil
}

// calculateJaccardDistance calculates the Jaccard distance between two token sets
func (e *Engine) calculateJaccardDistance(tokens1, tokens2 []string) float64 {
	// Create sets
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, token := range tokens1 {
		set1[token] = true
	}
	for _, token := range tokens2 {
		set2[token] = true
	}

	// Calculate intersection and union
	intersection := 0
	union := 0

	for token := range set1 {
		if set2[token] {
			intersection++
		}
		union++
	}

	for token := range set2 {
		if !set1[token] {
			union++
		}
	}

	// Calculate Jaccard distance
	if union == 0 {
		return 0.0
	}

	jaccardSimilarity := float64(intersection) / float64(union)
	return 1.0 - jaccardSimilarity
}

// GetMetricNames returns the list of available metrics
func (e *Engine) GetMetricNames() []string {
	return []string{
		"token_count",
		"entropy_global_entropy",
		"entropy_bigram_entropy",
		"entropy_vocab_normalized_entropy",
		"entropy_token_normalized_entropy",
		"entropy_char_normalized_entropy",
		"entropy_rolling_entropy_mean",
		"entropy_rolling_entropy_std",
		"compression_compression_ratio",
		"compression_compression_efficiency",
		"compression_space_savings_percent",
		"compression_avg_token_size",
		"compression_token_density",
		"compression_char_density",
		"reuse_reuse_ratio",
		"reuse_vocabulary_efficiency",
		"reuse_reuse_efficiency",
		"reuse_entropy_efficiency",
		"reuse_compression_efficiency",
		"drift_jaccard_distance",
		"drift_alignment_score",
		"drift_position_drift",
		"drift_length_drift",
		"drift_content_similarity",
	}
}

// CompareTokenizers performs cross-tokenizer comparison analysis
func (e *Engine) CompareTokenizers(ctx context.Context, document string, tokenizers []tokenizers.Tokenizer) (map[string]interface{}, error) {
	if len(tokenizers) < 2 {
		return nil, fmt.Errorf("at least 2 tokenizers required for comparison")
	}

	// Analyze document with each tokenizer
	results := make([]*AnalysisResult, len(tokenizers))
	for i, tokenizer := range tokenizers {
		result, err := e.AnalyzeDocument(ctx, document, tokenizer)
		if err != nil {
			return nil, fmt.Errorf("error analyzing with tokenizer %s: %w", tokenizer.Name(), err)
		}
		results[i] = result
	}

	// Calculate drift between tokenizers
	driftCalc := NewDriftCalculator(0.5)
	comparison := make(map[string]interface{})

	// Compare each pair of tokenizers
	for i := 0; i < len(results); i++ {
		for j := i + 1; j < len(results); j++ {
			pairName := fmt.Sprintf("%s_vs_%s", results[i].TokenizerName, results[j].TokenizerName)

			if driftStats, err := driftCalc.CalculateDriftStats(results[i].Tokenization, results[j].Tokenization); err == nil {
				comparison[pairName] = driftStats
			}
		}
	}

	// Add individual results
	comparison["individual_results"] = results

	return comparison, nil
}

// ValidateConfig validates the engine configuration
func (e *Engine) ValidateConfig() error {
	if e.config.EntropyWindowSize < 0 {
		return fmt.Errorf("entropy window size must be non-negative")
	}

	return nil
}
