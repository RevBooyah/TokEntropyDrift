package examples

import (
	"fmt"
	"math"
	"sort"

	"github.com/RevBooyah/TokEntropyDrift/internal/plugins"
)

// TokenLengthAnalyzer is an example plugin that analyzes token length distributions
type TokenLengthAnalyzer struct {
	*plugins.BasePlugin
}

// NewTokenLengthAnalyzer creates a new token length analyzer plugin
func NewTokenLengthAnalyzer() *TokenLengthAnalyzer {
	info := plugins.PluginInfo{
		Name:        "token_length_analyzer",
		Version:     "1.0.0",
		Description: "Analyzes token length distributions and statistics",
		Author:      "TokEntropyDrift Team",
		Tags:        []string{"analysis", "tokens", "length"},
		Metadata: map[string]string{
			"category": "token_analysis",
		},
	}

	return &TokenLengthAnalyzer{
		BasePlugin: plugins.NewBasePlugin(info),
	}
}

// CalculateMetrics calculates token length metrics
func (t *TokenLengthAnalyzer) CalculateMetrics(ctx *plugins.AnalysisContext) ([]plugins.MetricResult, error) {
	if ctx.Tokenization == nil || len(ctx.Tokenization.Tokens) == 0 {
		return []plugins.MetricResult{}, nil
	}

	tokens := ctx.Tokenization.Tokens
	lengths := make([]int, len(tokens))

	// Calculate token lengths
	for i, token := range tokens {
		lengths[i] = len(token.Text)
	}

	// Sort lengths for percentile calculations
	sort.Ints(lengths)

	// Calculate basic statistics
	totalTokens := len(lengths)
	sum := 0
	for _, length := range lengths {
		sum += length
	}

	mean := float64(sum) / float64(totalTokens)

	// Calculate variance and standard deviation
	variance := 0.0
	for _, length := range lengths {
		diff := float64(length) - mean
		variance += diff * diff
	}
	variance /= float64(totalTokens)
	stdDev := math.Sqrt(variance)

	// Calculate percentiles
	p25 := calculatePercentile(lengths, 25)
	p50 := calculatePercentile(lengths, 50)
	p75 := calculatePercentile(lengths, 75)
	p90 := calculatePercentile(lengths, 90)
	p95 := calculatePercentile(lengths, 95)
	p99 := calculatePercentile(lengths, 99)

	// Calculate min and max
	min := lengths[0]
	max := lengths[len(lengths)-1]

	// Calculate length distribution
	lengthCounts := make(map[int]int)
	for _, length := range lengths {
		lengthCounts[length]++
	}

	// Find most common length
	mostCommonLength := 0
	mostCommonCount := 0
	for length, count := range lengthCounts {
		if count > mostCommonCount {
			mostCommonLength = length
			mostCommonCount = count
		}
	}

	// Calculate entropy of length distribution
	lengthEntropy := calculateEntropy(lengthCounts, totalTokens)

	// Get configuration parameters
	minLengthThreshold := t.GetConfigInt("min_length_threshold", 1)
	maxLengthThreshold := t.GetConfigInt("max_length_threshold", 100)

	// Calculate tokens within thresholds
	tokensWithinThreshold := 0
	for _, length := range lengths {
		if length >= minLengthThreshold && length <= maxLengthThreshold {
			tokensWithinThreshold++
		}
	}
	thresholdPercentage := float64(tokensWithinThreshold) / float64(totalTokens) * 100

	results := []plugins.MetricResult{
		{
			Name:  "token_count",
			Value: float64(totalTokens),
			Unit:  "tokens",
		},
		{
			Name:  "mean_length",
			Value: mean,
			Unit:  "characters",
		},
		{
			Name:  "std_dev_length",
			Value: stdDev,
			Unit:  "characters",
		},
		{
			Name:  "min_length",
			Value: float64(min),
			Unit:  "characters",
		},
		{
			Name:  "max_length",
			Value: float64(max),
			Unit:  "characters",
		},
		{
			Name:  "median_length",
			Value: float64(p50),
			Unit:  "characters",
		},
		{
			Name:  "p25_length",
			Value: float64(p25),
			Unit:  "characters",
		},
		{
			Name:  "p75_length",
			Value: float64(p75),
			Unit:  "characters",
		},
		{
			Name:  "p90_length",
			Value: float64(p90),
			Unit:  "characters",
		},
		{
			Name:  "p95_length",
			Value: float64(p95),
			Unit:  "characters",
		},
		{
			Name:  "p99_length",
			Value: float64(p99),
			Unit:  "characters",
		},
		{
			Name:  "most_common_length",
			Value: float64(mostCommonLength),
			Unit:  "characters",
		},
		{
			Name:  "most_common_count",
			Value: float64(mostCommonCount),
			Unit:  "tokens",
		},
		{
			Name:  "length_entropy",
			Value: lengthEntropy,
			Unit:  "bits",
		},
		{
			Name:  "tokens_within_threshold",
			Value: float64(tokensWithinThreshold),
			Unit:  "tokens",
		},
		{
			Name:  "threshold_percentage",
			Value: thresholdPercentage,
			Unit:  "percent",
		},
	}

	return results, nil
}

// ValidateConfig validates the plugin configuration
func (t *TokenLengthAnalyzer) ValidateConfig(config map[string]interface{}) error {
	// Check for valid threshold values
	if minThreshold, exists := config["min_length_threshold"]; exists {
		if min, ok := minThreshold.(int); !ok || min < 0 {
			return fmt.Errorf("min_length_threshold must be a non-negative integer")
		}
	}

	if maxThreshold, exists := config["max_length_threshold"]; exists {
		if max, ok := maxThreshold.(int); !ok || max <= 0 {
			return fmt.Errorf("max_length_threshold must be a positive integer")
		}
	}

	// Check that min <= max if both are provided
	if minThreshold, minExists := config["min_length_threshold"]; minExists {
		if maxThreshold, maxExists := config["max_length_threshold"]; maxExists {
			if min, ok1 := minThreshold.(int); ok1 {
				if max, ok2 := maxThreshold.(int); ok2 {
					if min > max {
						return fmt.Errorf("min_length_threshold cannot be greater than max_length_threshold")
					}
				}
			}
		}
	}

	return nil
}

// calculatePercentile calculates the nth percentile of a sorted slice
func calculatePercentile(sorted []int, percentile int) int {
	if len(sorted) == 0 {
		return 0
	}

	index := float64(percentile) / 100.0 * float64(len(sorted)-1)
	if index == float64(int(index)) {
		return sorted[int(index)]
	}

	lower := int(index)
	upper := lower + 1
	if upper >= len(sorted) {
		return sorted[lower]
	}

	weight := index - float64(lower)
	return int(float64(sorted[lower])*(1-weight) + float64(sorted[upper])*weight)
}

// calculateEntropy calculates the entropy of a distribution
func calculateEntropy(counts map[int]int, total int) float64 {
	entropy := 0.0
	for _, count := range counts {
		if count > 0 {
			p := float64(count) / float64(total)
			entropy -= p * math.Log2(p)
		}
	}
	return entropy
}
