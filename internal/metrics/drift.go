package metrics

import (
	"fmt"
	"math"

	"github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
)

// DriftCalculator handles drift detection and cross-tokenizer comparison
type DriftCalculator struct {
	alignmentThreshold float64
}

// NewDriftCalculator creates a new drift calculator
func NewDriftCalculator(alignmentThreshold float64) *DriftCalculator {
	return &DriftCalculator{
		alignmentThreshold: alignmentThreshold,
	}
}

// CalculateJaccardDistance calculates the Jaccard distance between two token sets
func (d *DriftCalculator) CalculateJaccardDistance(tokens1, tokens2 []tokenizers.Token) (float64, error) {
	if len(tokens1) == 0 && len(tokens2) == 0 {
		return 0.0, nil
	}

	// Create sets
	set1 := make(map[string]bool)
	set2 := make(map[string]bool)

	for _, token := range tokens1 {
		set1[token.Text] = true
	}
	for _, token := range tokens2 {
		set2[token.Text] = true
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
		return 0.0, nil
	}

	jaccardSimilarity := float64(intersection) / float64(union)
	return 1.0 - jaccardSimilarity, nil
}

// CalculateTokenAlignment calculates alignment between token sequences
func (d *DriftCalculator) CalculateTokenAlignment(tokens1, tokens2 []tokenizers.Token) (map[string]float64, error) {
	metrics := make(map[string]float64)

	// Extract token texts
	texts1 := make([]string, len(tokens1))
	texts2 := make([]string, len(tokens2))

	for i, token := range tokens1 {
		texts1[i] = token.Text
	}
	for i, token := range tokens2 {
		texts2[i] = token.Text
	}

	// Calculate alignment metrics
	alignmentScore := d.calculateAlignmentScore(texts1, texts2)
	metrics["alignment_score"] = alignmentScore

	// Position-based drift
	positionDrift := d.calculatePositionDrift(texts1, texts2)
	metrics["position_drift"] = positionDrift

	// Length-based drift
	lengthDrift := math.Abs(float64(len(texts1)-len(texts2))) / math.Max(float64(len(texts1)), float64(len(texts2)))
	metrics["length_drift"] = lengthDrift

	// Content similarity
	contentSimilarity := d.calculateContentSimilarity(texts1, texts2)
	metrics["content_similarity"] = contentSimilarity

	return metrics, nil
}

// calculateAlignmentScore calculates how well tokens align between sequences
func (d *DriftCalculator) calculateAlignmentScore(texts1, texts2 []string) float64 {
	if len(texts1) == 0 || len(texts2) == 0 {
		return 0.0
	}

	// Use dynamic programming to find longest common subsequence
	lcs := d.longestCommonSubsequence(texts1, texts2)

	// Calculate alignment score based on LCS length
	maxLength := math.Max(float64(len(texts1)), float64(len(texts2)))
	return float64(lcs) / maxLength
}

// longestCommonSubsequence finds the longest common subsequence
func (d *DriftCalculator) longestCommonSubsequence(texts1, texts2 []string) int {
	if len(texts1) == 0 || len(texts2) == 0 {
		return 0
	}

	// Create DP table
	dp := make([][]int, len(texts1)+1)
	for i := range dp {
		dp[i] = make([]int, len(texts2)+1)
	}

	// Fill DP table
	for i := 1; i <= len(texts1); i++ {
		for j := 1; j <= len(texts2); j++ {
			if texts1[i-1] == texts2[j-1] {
				dp[i][j] = dp[i-1][j-1] + 1
			} else {
				dp[i][j] = int(math.Max(float64(dp[i-1][j]), float64(dp[i][j-1])))
			}
		}
	}

	return dp[len(texts1)][len(texts2)]
}

// calculatePositionDrift calculates drift based on token positions
func (d *DriftCalculator) calculatePositionDrift(texts1, texts2 []string) float64 {
	if len(texts1) == 0 || len(texts2) == 0 {
		return 0.0
	}

	// Find common tokens and their positions
	commonTokens := make(map[string][]int)

	// Record positions in sequence 1
	for i, text := range texts1 {
		if _, exists := commonTokens[text]; !exists {
			commonTokens[text] = make([]int, 0)
		}
		commonTokens[text] = append(commonTokens[text], i)
	}

	// Calculate position differences
	totalDrift := 0.0
	count := 0

	for i, text := range texts2 {
		if positions, exists := commonTokens[text]; exists && len(positions) > 0 {
			// Find closest position
			minDiff := math.Abs(float64(i - positions[0]))
			for _, pos := range positions {
				diff := math.Abs(float64(i - pos))
				if diff < minDiff {
					minDiff = diff
				}
			}
			totalDrift += minDiff
			count++
		}
	}

	if count == 0 {
		return 0.0
	}

	return totalDrift / float64(count)
}

// calculateContentSimilarity calculates content similarity between token sequences
func (d *DriftCalculator) calculateContentSimilarity(texts1, texts2 []string) float64 {
	if len(texts1) == 0 || len(texts2) == 0 {
		return 0.0
	}

	// Create frequency maps
	freq1 := make(map[string]int)
	freq2 := make(map[string]int)

	for _, text := range texts1 {
		freq1[text]++
	}
	for _, text := range texts2 {
		freq2[text]++
	}

	// Calculate cosine similarity
	dotProduct := 0.0
	magnitude1 := 0.0
	magnitude2 := 0.0

	// Get all unique tokens
	allTokens := make(map[string]bool)
	for token := range freq1 {
		allTokens[token] = true
	}
	for token := range freq2 {
		allTokens[token] = true
	}

	for token := range allTokens {
		count1 := float64(freq1[token])
		count2 := float64(freq2[token])

		dotProduct += count1 * count2
		magnitude1 += count1 * count1
		magnitude2 += count2 * count2
	}

	if magnitude1 == 0 || magnitude2 == 0 {
		return 0.0
	}

	return dotProduct / (math.Sqrt(magnitude1) * math.Sqrt(magnitude2))
}

// CalculateCrossTokenizerDrift calculates drift between two tokenization results
func (d *DriftCalculator) CalculateCrossTokenizerDrift(result1, result2 *tokenizers.TokenizationResult) (map[string]float64, error) {
	if result1 == nil || result2 == nil {
		return nil, fmt.Errorf("both tokenization results must be provided")
	}

	metrics := make(map[string]float64)

	// Jaccard distance
	if jaccardDistance, err := d.CalculateJaccardDistance(result1.Tokens, result2.Tokens); err == nil {
		metrics["jaccard_distance"] = jaccardDistance
	}

	// Token alignment
	if alignmentMetrics, err := d.CalculateTokenAlignment(result1.Tokens, result2.Tokens); err == nil {
		for k, v := range alignmentMetrics {
			metrics["alignment_"+k] = v
		}
	}

	// Token count drift
	tokenCount1 := float64(len(result1.Tokens))
	tokenCount2 := float64(len(result2.Tokens))
	maxCount := math.Max(tokenCount1, tokenCount2)
	if maxCount > 0 {
		metrics["token_count_drift"] = math.Abs(tokenCount1-tokenCount2) / maxCount
	}

	// Average token length drift
	avgLength1 := d.calculateAverageTokenLength(result1.Tokens)
	avgLength2 := d.calculateAverageTokenLength(result2.Tokens)
	maxLength := math.Max(avgLength1, avgLength2)
	if maxLength > 0 {
		metrics["avg_length_drift"] = math.Abs(avgLength1-avgLength2) / maxLength
	}

	// Vocabulary overlap
	vocabOverlap := d.calculateVocabularyOverlap(result1.Tokens, result2.Tokens)
	metrics["vocab_overlap"] = vocabOverlap

	return metrics, nil
}

// calculateAverageTokenLength calculates the average length of tokens
func (d *DriftCalculator) calculateAverageTokenLength(tokens []tokenizers.Token) float64 {
	if len(tokens) == 0 {
		return 0.0
	}

	totalLength := 0
	for _, token := range tokens {
		totalLength += len(token.Text)
	}

	return float64(totalLength) / float64(len(tokens))
}

// calculateVocabularyOverlap calculates the overlap between vocabularies
func (d *DriftCalculator) calculateVocabularyOverlap(tokens1, tokens2 []tokenizers.Token) float64 {
	vocab1 := make(map[string]bool)
	vocab2 := make(map[string]bool)

	for _, token := range tokens1 {
		vocab1[token.Text] = true
	}
	for _, token := range tokens2 {
		vocab2[token.Text] = true
	}

	intersection := 0
	union := 0

	for token := range vocab1 {
		if vocab2[token] {
			intersection++
		}
		union++
	}

	for token := range vocab2 {
		if !vocab1[token] {
			union++
		}
	}

	if union == 0 {
		return 0.0
	}

	return float64(intersection) / float64(union)
}

// CalculateDriftStats calculates comprehensive drift statistics
func (d *DriftCalculator) CalculateDriftStats(result1, result2 *tokenizers.TokenizationResult) (map[string]float64, error) {
	stats := make(map[string]float64)

	// Cross-tokenizer drift
	if driftMetrics, err := d.CalculateCrossTokenizerDrift(result1, result2); err == nil {
		for k, v := range driftMetrics {
			stats["drift_"+k] = v
		}
	}

	// Individual tokenizer statistics
	if result1 != nil {
		stats["tokenizer1_token_count"] = float64(len(result1.Tokens))
		stats["tokenizer1_avg_length"] = d.calculateAverageTokenLength(result1.Tokens)
	}

	if result2 != nil {
		stats["tokenizer2_token_count"] = float64(len(result2.Tokens))
		stats["tokenizer2_avg_length"] = d.calculateAverageTokenLength(result2.Tokens)
	}

	return stats, nil
}
