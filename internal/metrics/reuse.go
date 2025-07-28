package metrics

import (
	"math"
	"sort"

	"github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
)

// ReuseCalculator handles token reuse and frequency analysis
type ReuseCalculator struct {
	includePatterns bool
}

// NewReuseCalculator creates a new reuse calculator
func NewReuseCalculator(includePatterns bool) *ReuseCalculator {
	return &ReuseCalculator{
		includePatterns: includePatterns,
	}
}

// CalculateTokenReuse calculates basic token reuse metrics
func (r *ReuseCalculator) CalculateTokenReuse(tokens []tokenizers.Token) (float64, error) {
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

// CalculateTokenFrequency calculates detailed token frequency statistics
func (r *ReuseCalculator) CalculateTokenFrequency(tokens []tokenizers.Token) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	if len(tokens) == 0 {
		return stats, nil
	}

	// Count token frequencies
	tokenFreq := make(map[string]int)
	for _, token := range tokens {
		tokenFreq[token.Text]++
	}

	// Calculate frequency statistics
	frequencies := make([]int, 0, len(tokenFreq))
	for _, freq := range tokenFreq {
		frequencies = append(frequencies, freq)
	}
	sort.Ints(frequencies)

	// Basic statistics
	stats["unique_tokens"] = len(tokenFreq)
	stats["total_tokens"] = len(tokens)
	stats["reuse_ratio"] = 1.0 - (float64(len(tokenFreq)) / float64(len(tokens)))

	// Frequency distribution
	stats["min_frequency"] = frequencies[0]
	stats["max_frequency"] = frequencies[len(frequencies)-1]
	stats["median_frequency"] = calculateMedian(frequencies)
	stats["mean_frequency"] = calculateMeanInt(frequencies)
	stats["frequency_std"] = calculateStdInt(frequencies)

	// Most frequent tokens
	mostFrequent := r.getMostFrequentTokens(tokenFreq, 10)
	stats["most_frequent_tokens"] = mostFrequent

	// Frequency percentiles
	stats["freq_percentile_25"] = calculatePercentile(frequencies, 25)
	stats["freq_percentile_50"] = calculatePercentile(frequencies, 50)
	stats["freq_percentile_75"] = calculatePercentile(frequencies, 75)
	stats["freq_percentile_90"] = calculatePercentile(frequencies, 90)
	stats["freq_percentile_95"] = calculatePercentile(frequencies, 95)

	return stats, nil
}

// CalculateReusePatterns analyzes patterns in token reuse
func (r *ReuseCalculator) CalculateReusePatterns(tokens []tokenizers.Token) (map[string]interface{}, error) {
	patterns := make(map[string]interface{})

	if len(tokens) == 0 {
		return patterns, nil
	}

	// Analyze consecutive reuse patterns
	consecutivePatterns := r.analyzeConsecutivePatterns(tokens)
	patterns["consecutive_patterns"] = consecutivePatterns

	// Analyze distance patterns
	distancePatterns := r.analyzeDistancePatterns(tokens)
	patterns["distance_patterns"] = distancePatterns

	// Analyze burst patterns
	burstPatterns := r.analyzeBurstPatterns(tokens)
	patterns["burst_patterns"] = burstPatterns

	return patterns, nil
}

// analyzeConsecutivePatterns analyzes consecutive token reuse
func (r *ReuseCalculator) analyzeConsecutivePatterns(tokens []tokenizers.Token) map[string]interface{} {
	patterns := make(map[string]interface{})

	if len(tokens) < 2 {
		return patterns
	}

	consecutiveCount := 0
	totalConsecutive := 0

	for i := 1; i < len(tokens); i++ {
		if tokens[i].Text == tokens[i-1].Text {
			consecutiveCount++
		}
		totalConsecutive++
	}

	patterns["consecutive_reuse_ratio"] = float64(consecutiveCount) / float64(totalConsecutive)
	patterns["consecutive_reuse_count"] = consecutiveCount

	return patterns
}

// analyzeDistancePatterns analyzes distance between token reuse
func (r *ReuseCalculator) analyzeDistancePatterns(tokens []tokenizers.Token) map[string]interface{} {
	patterns := make(map[string]interface{})

	if len(tokens) == 0 {
		return patterns
	}

	// Track last occurrence of each token
	lastOccurrence := make(map[string]int)
	distances := make([]int, 0)

	for i, token := range tokens {
		if lastPos, exists := lastOccurrence[token.Text]; exists {
			distance := i - lastPos
			distances = append(distances, distance)
		}
		lastOccurrence[token.Text] = i
	}

	if len(distances) == 0 {
		patterns["avg_reuse_distance"] = 0.0
		patterns["min_reuse_distance"] = 0
		patterns["max_reuse_distance"] = 0
		return patterns
	}

	// Calculate distance statistics
	sort.Ints(distances)
	patterns["avg_reuse_distance"] = calculateMeanInt(distances)
	patterns["min_reuse_distance"] = distances[0]
	patterns["max_reuse_distance"] = distances[len(distances)-1]
	patterns["median_reuse_distance"] = calculateMedianInt(distances)

	return patterns
}

// analyzeBurstPatterns analyzes burst patterns in token usage
func (r *ReuseCalculator) analyzeBurstPatterns(tokens []tokenizers.Token) map[string]interface{} {
	patterns := make(map[string]interface{})

	if len(tokens) == 0 {
		return patterns
	}

	// Count token frequencies
	tokenFreq := make(map[string]int)
	for _, token := range tokens {
		tokenFreq[token.Text]++
	}

	// Find burst tokens (tokens that appear multiple times in sequence)
	burstTokens := make(map[string]int)
	currentToken := tokens[0].Text
	currentCount := 1

	for i := 1; i < len(tokens); i++ {
		if tokens[i].Text == currentToken {
			currentCount++
		} else {
			if currentCount > 1 {
				burstTokens[currentToken] = currentCount
			}
			currentToken = tokens[i].Text
			currentCount = 1
		}
	}

	// Handle last burst
	if currentCount > 1 {
		burstTokens[currentToken] = currentCount
	}

	// Calculate burst statistics
	if len(burstTokens) == 0 {
		patterns["burst_count"] = 0
		patterns["max_burst_size"] = 0
		patterns["avg_burst_size"] = 0.0
		return patterns
	}

	burstSizes := make([]int, 0, len(burstTokens))
	for _, size := range burstTokens {
		burstSizes = append(burstSizes, size)
	}
	sort.Ints(burstSizes)

	patterns["burst_count"] = len(burstTokens)
	patterns["max_burst_size"] = burstSizes[len(burstSizes)-1]
	patterns["avg_burst_size"] = calculateMeanInt(burstSizes)
	patterns["median_burst_size"] = calculateMedianInt(burstSizes)

	return patterns
}

// CalculateReuseEfficiency calculates efficiency metrics related to token reuse
func (r *ReuseCalculator) CalculateReuseEfficiency(tokens []tokenizers.Token) (map[string]float64, error) {
	efficiency := make(map[string]float64)

	if len(tokens) == 0 {
		return efficiency, nil
	}

	// Count unique tokens
	uniqueTokens := make(map[string]bool)
	tokenFreq := make(map[string]int)

	for _, token := range tokens {
		uniqueTokens[token.Text] = true
		tokenFreq[token.Text]++
	}

	uniqueCount := float64(len(uniqueTokens))
	totalCount := float64(len(tokens))

	// Basic efficiency metrics
	efficiency["vocabulary_efficiency"] = uniqueCount / totalCount
	efficiency["reuse_efficiency"] = 1.0 - efficiency["vocabulary_efficiency"]

	// Entropy-based efficiency
	entropy := 0.0
	for _, freq := range tokenFreq {
		probability := float64(freq) / totalCount
		if probability > 0 {
			entropy -= probability * math.Log2(probability)
		}
	}

	maxEntropy := math.Log2(uniqueCount)
	if maxEntropy > 0 {
		efficiency["entropy_efficiency"] = entropy / maxEntropy
	} else {
		efficiency["entropy_efficiency"] = 0.0
	}

	// Compression efficiency (how well reuse enables compression)
	efficiency["compression_efficiency"] = efficiency["reuse_efficiency"] * efficiency["entropy_efficiency"]

	return efficiency, nil
}

// CalculateReuseStats calculates comprehensive reuse statistics
func (r *ReuseCalculator) CalculateReuseStats(tokens []tokenizers.Token) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// Basic reuse ratio
	if reuseRatio, err := r.CalculateTokenReuse(tokens); err == nil {
		stats["reuse_ratio"] = reuseRatio
	}

	// Token frequency analysis
	if freqStats, err := r.CalculateTokenFrequency(tokens); err == nil {
		for k, v := range freqStats {
			stats["freq_"+k] = v
		}
	}

	// Reuse patterns
	if r.includePatterns {
		if patternStats, err := r.CalculateReusePatterns(tokens); err == nil {
			for k, v := range patternStats {
				stats["pattern_"+k] = v
			}
		}
	}

	// Reuse efficiency
	if efficiencyStats, err := r.CalculateReuseEfficiency(tokens); err == nil {
		for k, v := range efficiencyStats {
			stats["efficiency_"+k] = v
		}
	}

	return stats, nil
}

// Helper functions
func (r *ReuseCalculator) getMostFrequentTokens(tokenFreq map[string]int, count int) []map[string]interface{} {
	type tokenFreqPair struct {
		token string
		freq  int
	}

	pairs := make([]tokenFreqPair, 0, len(tokenFreq))
	for token, freq := range tokenFreq {
		pairs = append(pairs, tokenFreqPair{token, freq})
	}

	// Sort by frequency (descending)
	sort.Slice(pairs, func(i, j int) bool {
		return pairs[i].freq > pairs[j].freq
	})

	// Take top N
	if count > len(pairs) {
		count = len(pairs)
	}

	result := make([]map[string]interface{}, count)
	for i := 0; i < count; i++ {
		result[i] = map[string]interface{}{
			"token":     pairs[i].token,
			"frequency": pairs[i].freq,
		}
	}

	return result
}

func calculateMedian(values []int) float64 {
	if len(values) == 0 {
		return 0.0
	}

	sort.Ints(values)
	if len(values)%2 == 0 {
		return float64(values[len(values)/2-1]+values[len(values)/2]) / 2.0
	}
	return float64(values[len(values)/2])
}

func calculateMedianInt(values []int) int {
	if len(values) == 0 {
		return 0
	}

	sort.Ints(values)
	if len(values)%2 == 0 {
		return (values[len(values)/2-1] + values[len(values)/2]) / 2
	}
	return values[len(values)/2]
}

func calculateMeanInt(values []int) float64 {
	if len(values) == 0 {
		return 0.0
	}
	sum := 0
	for _, v := range values {
		sum += v
	}
	return float64(sum) / float64(len(values))
}

func calculatePercentile(values []int, percentile int) int {
	if len(values) == 0 {
		return 0
	}

	sort.Ints(values)
	index := int(float64(percentile) / 100.0 * float64(len(values)-1))
	return values[index]
}

func calculateStdInt(values []int) float64 {
	if len(values) == 0 {
		return 0.0
	}
	mean := calculateMeanInt(values)
	sum := 0.0
	for _, v := range values {
		sum += (float64(v) - mean) * (float64(v) - mean)
	}
	return math.Sqrt(sum / float64(len(values)))
}
