package metrics

import (
	"math"

	"github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
)

// CompressionCalculator handles various compression metrics
type CompressionCalculator struct {
	includeMetadata bool
}

// NewCompressionCalculator creates a new compression calculator
func NewCompressionCalculator(includeMetadata bool) *CompressionCalculator {
	return &CompressionCalculator{
		includeMetadata: includeMetadata,
	}
}

// CalculateCompressionRatio calculates the basic compression ratio
func (c *CompressionCalculator) CalculateCompressionRatio(originalText string, tokens []tokenizers.Token) (float64, error) {
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

// CalculateByteLevelCompression calculates detailed byte-level compression metrics
func (c *CompressionCalculator) CalculateByteLevelCompression(originalText string, tokens []tokenizers.Token) (map[string]float64, error) {
	metrics := make(map[string]float64)

	originalBytes := []byte(originalText)
	originalSize := len(originalBytes)

	// Calculate token sizes
	tokenSizes := make([]int, len(tokens))
	totalTokenSize := 0

	for i, token := range tokens {
		// Estimate token size (ID + metadata)
		tokenSize := 4 // 4 bytes for token ID
		if c.includeMetadata {
			tokenSize += len(token.Text) // Include token text for reconstruction
		}
		tokenSizes[i] = tokenSize
		totalTokenSize += tokenSize
	}

	// Basic compression ratio
	metrics["compression_ratio"] = float64(totalTokenSize) / float64(originalSize)

	// Compression efficiency (lower is better)
	metrics["compression_efficiency"] = float64(originalSize) / float64(totalTokenSize)

	// Space savings percentage
	metrics["space_savings_percent"] = (1.0 - float64(totalTokenSize)/float64(originalSize)) * 100

	// Average token size
	metrics["avg_token_size"] = float64(totalTokenSize) / float64(len(tokens))

	// Token density (tokens per byte)
	metrics["token_density"] = float64(len(tokens)) / float64(originalSize)

	// Character density (characters per token)
	totalChars := 0
	for _, token := range tokens {
		totalChars += len(token.Text)
	}
	metrics["char_density"] = float64(totalChars) / float64(len(tokens))

	return metrics, nil
}

// CalculateTokenLevelCompression calculates token-level compression metrics
func (c *CompressionCalculator) CalculateTokenLevelCompression(tokens []tokenizers.Token) (map[string]float64, error) {
	metrics := make(map[string]float64)

	if len(tokens) == 0 {
		return metrics, nil
	}

	// Token length statistics
	tokenLengths := make([]int, len(tokens))
	totalLength := 0

	for i, token := range tokens {
		length := len(token.Text)
		tokenLengths[i] = length
		totalLength += length
	}

	// Average token length
	metrics["avg_token_length"] = float64(totalLength) / float64(len(tokens))

	// Token length variance
	meanLength := metrics["avg_token_length"]
	variance := 0.0
	for _, length := range tokenLengths {
		variance += math.Pow(float64(length)-meanLength, 2)
	}
	metrics["token_length_variance"] = variance / float64(len(tokens))
	metrics["token_length_std"] = math.Sqrt(metrics["token_length_variance"])

	// Token length distribution
	metrics["min_token_length"] = float64(calculateMinInt(tokenLengths))
	metrics["max_token_length"] = float64(calculateMaxInt(tokenLengths))

	// Token efficiency (characters per token)
	metrics["token_efficiency"] = float64(totalLength) / float64(len(tokens))

	return metrics, nil
}

// CalculateRedundancyFactor calculates redundancy and efficiency metrics
func (c *CompressionCalculator) CalculateRedundancyFactor(tokens []tokenizers.Token, entropy float64) (map[string]float64, error) {
	metrics := make(map[string]float64)

	if len(tokens) == 0 {
		return metrics, nil
	}

	// Count unique tokens
	uniqueTokens := make(map[string]bool)
	for _, token := range tokens {
		uniqueTokens[token.Text] = true
	}

	// Calculate redundancy metrics
	totalTokens := float64(len(tokens))
	uniqueTokenCount := float64(len(uniqueTokens))

	// Redundancy factor (how much repetition exists)
	metrics["redundancy_factor"] = 1.0 - (uniqueTokenCount / totalTokens)

	// Vocabulary utilization
	metrics["vocab_utilization"] = uniqueTokenCount / totalTokens

	// Theoretical maximum entropy
	maxEntropy := math.Log2(uniqueTokenCount)
	if maxEntropy > 0 {
		// Entropy efficiency (how close to maximum entropy)
		metrics["entropy_efficiency"] = entropy / maxEntropy

		// Redundancy based on entropy
		metrics["entropy_redundancy"] = 1.0 - (entropy / maxEntropy)
	}

	// Token diversity (unique tokens per total tokens)
	metrics["token_diversity"] = uniqueTokenCount / totalTokens

	// Compression potential (how much more compression is possible)
	metrics["compression_potential"] = 1.0 - metrics["entropy_efficiency"]

	return metrics, nil
}

// CalculateCompressionStats calculates comprehensive compression statistics
func (c *CompressionCalculator) CalculateCompressionStats(originalText string, tokens []tokenizers.Token, entropy float64) (map[string]float64, error) {
	stats := make(map[string]float64)

	// Basic compression ratio
	if compressionRatio, err := c.CalculateCompressionRatio(originalText, tokens); err == nil {
		stats["compression_ratio"] = compressionRatio
	}

	// Byte-level compression
	if byteStats, err := c.CalculateByteLevelCompression(originalText, tokens); err == nil {
		for k, v := range byteStats {
			stats["byte_"+k] = v
		}
	}

	// Token-level compression
	if tokenStats, err := c.CalculateTokenLevelCompression(tokens); err == nil {
		for k, v := range tokenStats {
			stats["token_"+k] = v
		}
	}

	// Redundancy factor
	if redundancyStats, err := c.CalculateRedundancyFactor(tokens, entropy); err == nil {
		for k, v := range redundancyStats {
			stats["redundancy_"+k] = v
		}
	}

	return stats, nil
}

// Helper functions for integer statistics
func calculateMinInt(values []int) int {
	if len(values) == 0 {
		return 0
	}
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

func calculateMaxInt(values []int) int {
	if len(values) == 0 {
		return 0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}
