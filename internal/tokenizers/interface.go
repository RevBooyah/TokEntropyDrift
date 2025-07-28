package tokenizers

import (
	"context"
	"fmt"
)

// Token represents a single token with metadata
type Token struct {
	Text      string            `json:"text"`
	ID        int               `json:"id"`
	StartPos  int               `json:"start_pos"`
	EndPos    int               `json:"end_pos"`
	Metadata  map[string]string `json:"metadata,omitempty"`
}

// TokenizationResult represents the result of tokenizing a document
type TokenizationResult struct {
	Document  string   `json:"document"`
	Tokens    []Token  `json:"tokens"`
	Tokenizer string   `json:"tokenizer"`
	Metadata  map[string]interface{} `json:"metadata,omitempty"`
}

// TokenizerConfig represents configuration for a tokenizer
type TokenizerConfig struct {
	Name       string            `json:"name"`
	Type       string            `json:"type"` // bpe, spiece, wordpiece, custom
	Command    string            `json:"command,omitempty"`
	LibraryPath string           `json:"library_path,omitempty"`
	VocabFile  string            `json:"vocab_file,omitempty"`
	ModelFile  string            `json:"model_file,omitempty"`
	Parameters map[string]string `json:"parameters,omitempty"`
}

// Tokenizer defines the interface that all tokenizer adapters must implement
type Tokenizer interface {
	// Name returns the name of the tokenizer
	Name() string
	
	// Type returns the type of the tokenizer (bpe, spiece, wordpiece, custom)
	Type() string
	
	// Initialize prepares the tokenizer for use
	Initialize(config TokenizerConfig) error
	
	// Tokenize tokenizes a single document
	Tokenize(ctx context.Context, text string) (*TokenizationResult, error)
	
	// TokenizeBatch tokenizes multiple documents
	TokenizeBatch(ctx context.Context, texts []string) ([]*TokenizationResult, error)
	
	// GetVocabSize returns the vocabulary size of the tokenizer
	GetVocabSize() (int, error)
	
	// Close cleans up any resources used by the tokenizer
	Close() error
}

// BaseTokenizer provides common functionality for tokenizer implementations
type BaseTokenizer struct {
	name   string
	config TokenizerConfig
}

// NewBaseTokenizer creates a new base tokenizer
func NewBaseTokenizer(name string) *BaseTokenizer {
	return &BaseTokenizer{
		name: name,
	}
}

// Name returns the tokenizer name
func (b *BaseTokenizer) Name() string {
	return b.name
}

// Type returns the tokenizer type
func (b *BaseTokenizer) Type() string {
	return b.config.Type
}

// Initialize sets up the base tokenizer configuration
func (b *BaseTokenizer) Initialize(config TokenizerConfig) error {
	b.config = config
	return nil
}

// Close provides a default implementation for cleanup
func (b *BaseTokenizer) Close() error {
	// Default implementation does nothing
	return nil
}

// ValidateConfig validates the tokenizer configuration
func ValidateConfig(config TokenizerConfig) error {
	if config.Name == "" {
		return fmt.Errorf("tokenizer name is required")
	}
	
	if config.Type == "" {
		return fmt.Errorf("tokenizer type is required")
	}
	
	// Validate type
	validTypes := map[string]bool{
		"bpe":       true,
		"spiece":    true,
		"wordpiece": true,
		"custom":    true,
	}
	
	if !validTypes[config.Type] {
		return fmt.Errorf("invalid tokenizer type: %s", config.Type)
	}
	
	return nil
}

// TokenizerRegistry manages available tokenizers
type TokenizerRegistry struct {
	tokenizers map[string]Tokenizer
}

// NewTokenizerRegistry creates a new tokenizer registry
func NewTokenizerRegistry() *TokenizerRegistry {
	return &TokenizerRegistry{
		tokenizers: make(map[string]Tokenizer),
	}
}

// Register registers a tokenizer with the registry
func (r *TokenizerRegistry) Register(name string, tokenizer Tokenizer) error {
	if name == "" {
		return fmt.Errorf("tokenizer name cannot be empty")
	}
	
	if tokenizer == nil {
		return fmt.Errorf("tokenizer cannot be nil")
	}
	
	if _, exists := r.tokenizers[name]; exists {
		return fmt.Errorf("tokenizer %s already registered", name)
	}
	
	r.tokenizers[name] = tokenizer
	return nil
}

// Get retrieves a tokenizer by name
func (r *TokenizerRegistry) Get(name string) (Tokenizer, error) {
	tokenizer, exists := r.tokenizers[name]
	if !exists {
		return nil, fmt.Errorf("tokenizer %s not found", name)
	}
	
	return tokenizer, nil
}

// List returns all registered tokenizer names
func (r *TokenizerRegistry) List() []string {
	names := make([]string, 0, len(r.tokenizers))
	for name := range r.tokenizers {
		names = append(names, name)
	}
	return names
}

// Unregister removes a tokenizer from the registry
func (r *TokenizerRegistry) Unregister(name string) error {
	if _, exists := r.tokenizers[name]; !exists {
		return fmt.Errorf("tokenizer %s not found", name)
	}
	
	delete(r.tokenizers, name)
	return nil
}

// Global registry instance
var GlobalRegistry = NewTokenizerRegistry()

// RegisterGlobal registers a tokenizer with the global registry
func RegisterGlobal(name string, tokenizer Tokenizer) error {
	return GlobalRegistry.Register(name, tokenizer)
}

// GetGlobal retrieves a tokenizer from the global registry
func GetGlobal(name string) (Tokenizer, error) {
	return GlobalRegistry.Get(name)
}

// ListGlobal returns all registered tokenizer names from the global registry
func ListGlobal() []string {
	return GlobalRegistry.List()
} 