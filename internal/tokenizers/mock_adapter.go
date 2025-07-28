package tokenizers

import (
	"context"
	"fmt"
	"strings"
)

// MockTokenizer is a simple tokenizer for testing purposes
type MockTokenizer struct {
	*BaseTokenizer
	vocabSize int
}

// NewMockTokenizer creates a new mock tokenizer
func NewMockTokenizer(name string) *MockTokenizer {
	return &MockTokenizer{
		BaseTokenizer: NewBaseTokenizer(name),
		vocabSize:     1000, // Default vocab size
	}
}

// Initialize sets up the mock tokenizer
func (m *MockTokenizer) Initialize(config TokenizerConfig) error {
	if err := m.BaseTokenizer.Initialize(config); err != nil {
		return err
	}
	
	// Set vocab size from config if provided
	if vocabSizeStr, ok := config.Parameters["vocab_size"]; ok {
		if vocabSize, err := fmt.Sscanf(vocabSizeStr, "%d", &m.vocabSize); err != nil {
			return fmt.Errorf("invalid vocab_size parameter: %w", err)
		} else if vocabSize != 1 {
			return fmt.Errorf("failed to parse vocab_size parameter")
		}
	}
	
	return nil
}

// Tokenize tokenizes a single document using simple word splitting
func (m *MockTokenizer) Tokenize(ctx context.Context, text string) (*TokenizationResult, error) {
	// Simple word-based tokenization for testing
	words := strings.Fields(text)
	tokens := make([]Token, len(words))
	
	for i, word := range words {
		// Simple hash-based token ID
		tokenID := 0
		for _, char := range word {
			tokenID = (tokenID*31 + int(char)) % m.vocabSize
		}
		
		// Find position in original text
		startPos := strings.Index(text, word)
		endPos := startPos + len(word)
		
		tokens[i] = Token{
			Text:     word,
			ID:       tokenID,
			StartPos: startPos,
			EndPos:   endPos,
			Metadata: map[string]string{
				"tokenizer": "mock",
				"method":    "word_split",
			},
		}
	}
	
	return &TokenizationResult{
		Document:  text,
		Tokens:    tokens,
		Tokenizer: m.Name(),
		Metadata: map[string]interface{}{
			"tokenizer_type": "mock",
			"vocab_size":     m.vocabSize,
		},
	}, nil
}

// TokenizeBatch tokenizes multiple documents
func (m *MockTokenizer) TokenizeBatch(ctx context.Context, texts []string) ([]*TokenizationResult, error) {
	results := make([]*TokenizationResult, len(texts))
	
	for i, text := range texts {
		result, err := m.Tokenize(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("error tokenizing document %d: %w", i, err)
		}
		results[i] = result
	}
	
	return results, nil
}

// GetVocabSize returns the vocabulary size
func (m *MockTokenizer) GetVocabSize() (int, error) {
	return m.vocabSize, nil
}

// Close cleans up resources
func (m *MockTokenizer) Close() error {
	// Nothing to clean up for mock tokenizer
	return nil
}

// RegisterMockTokenizer registers the mock tokenizer with the global registry
func RegisterMockTokenizer() error {
	mockTokenizer := NewMockTokenizer("mock")
	return RegisterGlobal("mock", mockTokenizer)
} 