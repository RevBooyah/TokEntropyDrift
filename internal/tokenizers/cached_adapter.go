package tokenizers

import (
	"context"
	"fmt"

	"github.com/RevBooyah/TokEntropyDrift/internal/cache"
)

// CachedTokenizer wraps a tokenizer with caching functionality
type CachedTokenizer struct {
	tokenizer Tokenizer
	cache     *cache.Cache
	name      string
}

// NewCachedTokenizer creates a new cached tokenizer wrapper
func NewCachedTokenizer(tokenizer Tokenizer, cacheConfig cache.CacheConfig) *CachedTokenizer {
	return &CachedTokenizer{
		tokenizer: tokenizer,
		cache:     cache.NewCache(cacheConfig),
		name:      fmt.Sprintf("cached_%s", tokenizer.Name()),
	}
}

// Name returns the cached tokenizer name
func (c *CachedTokenizer) Name() string {
	return c.name
}

// Type returns the underlying tokenizer type
func (c *CachedTokenizer) Type() string {
	return c.tokenizer.Type()
}

// Initialize initializes the underlying tokenizer
func (c *CachedTokenizer) Initialize(config TokenizerConfig) error {
	return c.tokenizer.Initialize(config)
}

// Tokenize tokenizes text with caching
func (c *CachedTokenizer) Tokenize(ctx context.Context, text string) (*TokenizationResult, error) {
	// Generate cache key
	cacheKey := cache.GenerateKey(c.tokenizer.Name(), text)

	// Try to get from cache first
	if cached, found := c.cache.Get(cacheKey); found {
		if result, ok := cached.(*TokenizationResult); ok {
			return result, nil
		}
	}

	// Not in cache, tokenize and cache the result
	result, err := c.tokenizer.Tokenize(ctx, text)
	if err != nil {
		return nil, err
	}

	// Cache the result
	c.cache.Set(cacheKey, result)

	return result, nil
}

// TokenizeBatch tokenizes multiple texts with caching
func (c *CachedTokenizer) TokenizeBatch(ctx context.Context, texts []string) ([]*TokenizationResult, error) {
	results := make([]*TokenizationResult, len(texts))
	uncachedIndices := make([]int, 0)

	// Check cache for each text
	for i, text := range texts {
		cacheKey := cache.GenerateKey(c.tokenizer.Name(), text)
		if cached, found := c.cache.Get(cacheKey); found {
			if result, ok := cached.(*TokenizationResult); ok {
				results[i] = result
				continue
			}
		}
		uncachedIndices = append(uncachedIndices, i)
	}

	// If all texts were cached, return results
	if len(uncachedIndices) == 0 {
		return results, nil
	}

	// Tokenize uncached texts
	uncachedTexts := make([]string, len(uncachedIndices))
	for i, idx := range uncachedIndices {
		uncachedTexts[i] = texts[idx]
	}

	uncachedResults, err := c.tokenizer.TokenizeBatch(ctx, uncachedTexts)
	if err != nil {
		return nil, err
	}

	// Cache results and fill in the results slice
	for i, idx := range uncachedIndices {
		result := uncachedResults[i]
		results[idx] = result

		// Cache the result
		cacheKey := cache.GenerateKey(c.tokenizer.Name(), texts[idx])
		c.cache.Set(cacheKey, result)
	}

	return results, nil
}

// GetVocabSize returns the vocabulary size of the underlying tokenizer
func (c *CachedTokenizer) GetVocabSize() (int, error) {
	return c.tokenizer.GetVocabSize()
}

// Close closes both the cache and the underlying tokenizer
func (c *CachedTokenizer) Close() error {
	c.cache.Close()
	return c.tokenizer.Close()
}

// GetCacheStats returns cache statistics
func (c *CachedTokenizer) GetCacheStats() cache.CacheStats {
	return c.cache.GetStats()
}

// ClearCache clears the tokenizer cache
func (c *CachedTokenizer) ClearCache() {
	c.cache.Clear()
}
