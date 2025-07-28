package metrics

import (
	"math"

	"github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
)

// EntropyCalculator handles various entropy calculations
type EntropyCalculator struct {
	windowSize int
	normalize  bool
}

// NewEntropyCalculator creates a new entropy calculator
func NewEntropyCalculator(windowSize int, normalize bool) *EntropyCalculator {
	return &EntropyCalculator{
		windowSize: windowSize,
		normalize:  normalize,
	}
}

// CalculateGlobalEntropy calculates Shannon entropy over the entire token sequence
func (e *EntropyCalculator) CalculateGlobalEntropy(tokens []tokenizers.Token) (float64, error) {
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

	// Normalize if requested
	if e.normalize {
		maxEntropy := math.Log2(float64(len(tokenFreq)))
		if maxEntropy > 0 {
			entropy = entropy / maxEntropy
		}
	}

	return entropy, nil
}

// CalculateRollingEntropy calculates entropy over sliding windows
func (e *EntropyCalculator) CalculateRollingEntropy(tokens []tokenizers.Token) ([]float64, error) {
	if len(tokens) == 0 {
		return []float64{}, nil
	}

	windowSize := e.windowSize
	if windowSize <= 0 {
		windowSize = 100 // Default window size
	}

	if windowSize > len(tokens) {
		windowSize = len(tokens)
	}

	var rollingEntropy []float64

	for i := 0; i <= len(tokens)-windowSize; i++ {
		windowTokens := tokens[i : i+windowSize]
		entropy, err := e.CalculateGlobalEntropy(windowTokens)
		if err != nil {
			return nil, err
		}
		rollingEntropy = append(rollingEntropy, entropy)
	}

	return rollingEntropy, nil
}

// CalculateBigramEntropy calculates conditional entropy of token pairs
func (e *EntropyCalculator) CalculateBigramEntropy(tokens []tokenizers.Token) (float64, error) {
	if len(tokens) < 2 {
		return 0.0, nil
	}

	// Count bigram frequencies
	bigramFreq := make(map[string]int)
	unigramFreq := make(map[string]int)

	for i := 0; i < len(tokens)-1; i++ {
		bigram := tokens[i].Text + " " + tokens[i+1].Text
		bigramFreq[bigram]++
		unigramFreq[tokens[i].Text]++
	}
	unigramFreq[tokens[len(tokens)-1].Text]++ // Count last token

	// Calculate conditional entropy
	entropy := 0.0
	totalBigrams := float64(len(tokens) - 1)

	for bigram, freq := range bigramFreq {
		// Split bigram to get first token
		firstToken := tokens[0].Text // Default, will be updated
		for i := 0; i < len(tokens)-1; i++ {
			if tokens[i].Text+" "+tokens[i+1].Text == bigram {
				firstToken = tokens[i].Text
				break
			}
		}

		bigramProb := float64(freq) / totalBigrams
		unigramProb := float64(unigramFreq[firstToken]) / float64(len(tokens))

		if bigramProb > 0 && unigramProb > 0 {
			conditionalProb := bigramProb / unigramProb
			entropy -= bigramProb * math.Log2(conditionalProb)
		}
	}

	return entropy, nil
}

// CalculateNormalizedEntropy calculates entropy normalized by various factors
func (e *EntropyCalculator) CalculateNormalizedEntropy(tokens []tokenizers.Token, normalizationType string) (float64, error) {
	entropy, err := e.CalculateGlobalEntropy(tokens)
	if err != nil {
		return 0.0, err
	}

	switch normalizationType {
	case "vocab_size":
		// Normalize by vocabulary size
		uniqueTokens := make(map[string]bool)
		for _, token := range tokens {
			uniqueTokens[token.Text] = true
		}
		maxEntropy := math.Log2(float64(len(uniqueTokens)))
		if maxEntropy > 0 {
			return entropy / maxEntropy, nil
		}
		return entropy, nil

	case "token_count":
		// Normalize by token count
		maxEntropy := math.Log2(float64(len(tokens)))
		if maxEntropy > 0 {
			return entropy / maxEntropy, nil
		}
		return entropy, nil

	case "character_count":
		// Normalize by character count
		charCount := 0
		for _, token := range tokens {
			charCount += len(token.Text)
		}
		maxEntropy := math.Log2(float64(charCount))
		if maxEntropy > 0 {
			return entropy / maxEntropy, nil
		}
		return entropy, nil

	default:
		return entropy, nil
	}
}

// CalculateEntropyStats calculates comprehensive entropy statistics
func (e *EntropyCalculator) CalculateEntropyStats(tokens []tokenizers.Token) (map[string]float64, error) {
	stats := make(map[string]float64)

	// Global entropy
	if globalEntropy, err := e.CalculateGlobalEntropy(tokens); err == nil {
		stats["global_entropy"] = globalEntropy
	}

	// Bigram entropy
	if bigramEntropy, err := e.CalculateBigramEntropy(tokens); err == nil {
		stats["bigram_entropy"] = bigramEntropy
	}

	// Normalized entropies
	if vocabNormEntropy, err := e.CalculateNormalizedEntropy(tokens, "vocab_size"); err == nil {
		stats["vocab_normalized_entropy"] = vocabNormEntropy
	}

	if tokenNormEntropy, err := e.CalculateNormalizedEntropy(tokens, "token_count"); err == nil {
		stats["token_normalized_entropy"] = tokenNormEntropy
	}

	if charNormEntropy, err := e.CalculateNormalizedEntropy(tokens, "character_count"); err == nil {
		stats["char_normalized_entropy"] = charNormEntropy
	}

	// Rolling entropy statistics
	if rollingEntropy, err := e.CalculateRollingEntropy(tokens); err == nil && len(rollingEntropy) > 0 {
		stats["rolling_entropy_mean"] = calculateMean(rollingEntropy)
		stats["rolling_entropy_std"] = calculateStd(rollingEntropy)
		stats["rolling_entropy_min"] = calculateMin(rollingEntropy)
		stats["rolling_entropy_max"] = calculateMax(rollingEntropy)
	}

	return stats, nil
}

// Helper functions for statistics
func calculateMean(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	sum := 0.0
	for _, v := range values {
		sum += v
	}
	return sum / float64(len(values))
}

func calculateStd(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	mean := calculateMean(values)
	sum := 0.0
	for _, v := range values {
		sum += (v - mean) * (v - mean)
	}
	return math.Sqrt(sum / float64(len(values)))
}

func calculateMin(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	min := values[0]
	for _, v := range values {
		if v < min {
			min = v
		}
	}
	return min
}

func calculateMax(values []float64) float64 {
	if len(values) == 0 {
		return 0.0
	}
	max := values[0]
	for _, v := range values {
		if v > max {
			max = v
		}
	}
	return max
}
