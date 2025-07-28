package tokenizers

import (
	"fmt"
	"strings"
)

// RegisterAllTokenizers registers all available tokenizers with the global registry
func RegisterAllTokenizers() error {
	tokenizers := []struct {
		name   string
		register func() error
	}{
		{"mock", RegisterMockTokenizer},
		{"gpt2", RegisterGPT2Tokenizer},
		{"gpt-3.5-turbo", RegisterGPT35Tokenizer},
		{"gpt-4", RegisterGPT4Tokenizer},
		{"roberta-base", RegisterRoBERTaTokenizer},
		{"gpt-neo", RegisterGPTNeoTokenizer},
		{"bert-base", RegisterBERTTokenizer},
		{"distilbert-base", RegisterDistilBERTTokenizer},
		{"t5-base", RegisterT5Tokenizer},
		{"mt5-base", RegisterMT5Tokenizer},
		{"albert-base", RegisterALBERTTokenizer},
		{"openai-api", RegisterOpenAITokenizer},
	}

	var errors []string
	for _, t := range tokenizers {
		if err := t.register(); err != nil {
			errors = append(errors, fmt.Sprintf("%s: %v", t.name, err))
		}
	}

	if len(errors) > 0 {
		return fmt.Errorf("failed to register some tokenizers: %s", strings.Join(errors, "; "))
	}

	return nil
}

// GetAvailableTokenizers returns a list of all available tokenizer names
func GetAvailableTokenizers() []string {
	return []string{
		"mock",
		"gpt2",
		"gpt-3.5-turbo", 
		"gpt-4",
		"roberta-base",
		"gpt-neo",
		"bert-base",
		"distilbert-base",
		"t5-base",
		"mt5-base",
		"albert-base",
		"openai-api",
	}
}

// ValidateTokenizerName checks if a tokenizer name is valid
func ValidateTokenizerName(name string) bool {
	available := GetAvailableTokenizers()
	for _, availableName := range available {
		if availableName == name {
			return true
		}
	}
	return false
}

// GetTokenizerDescription returns a description of the tokenizer
func GetTokenizerDescription(name string) string {
	descriptions := map[string]string{
		"mock":           "Mock tokenizer for testing (word-based)",
		"gpt2":           "GPT-2 tokenizer using tiktoken (BPE)",
		"gpt-3.5-turbo":  "GPT-3.5 Turbo tokenizer using tiktoken (BPE)",
		"gpt-4":          "GPT-4 tokenizer using tiktoken (BPE)",
		"roberta-base":   "RoBERTa tokenizer using HuggingFace (BPE)",
		"gpt-neo":        "GPT-Neo tokenizer using HuggingFace (BPE)",
		"bert-base":      "BERT tokenizer using HuggingFace (WordPiece)",
		"distilbert-base": "DistilBERT tokenizer using HuggingFace (WordPiece)",
		"t5-base":        "T5 tokenizer using SentencePiece (Unigram)",
		"mt5-base":       "mT5 tokenizer using SentencePiece (Unigram)",
		"albert-base":    "ALBERT tokenizer using SentencePiece (WordPiece)",
		"openai-api":     "OpenAI API tokenizer (requires API key)",
	}

	if desc, ok := descriptions[name]; ok {
		return desc
	}
	return "Unknown tokenizer"
}

// GetTokenizerRequirements returns the requirements for a tokenizer
func GetTokenizerRequirements(name string) map[string]string {
	requirements := map[string]map[string]string{
		"mock": {},
		"gpt2": {
			"python": "Python 3.7+ with tiktoken package",
		},
		"gpt-3.5-turbo": {
			"python": "Python 3.7+ with tiktoken package",
		},
		"gpt-4": {
			"python": "Python 3.7+ with tiktoken package",
		},
		"roberta-base": {
			"python": "Python 3.7+ with transformers package",
		},
		"gpt-neo": {
			"python": "Python 3.7+ with transformers package",
		},
		"bert-base": {
			"python": "Python 3.7+ with transformers package",
		},
		"distilbert-base": {
			"python": "Python 3.7+ with transformers package",
		},
		"t5-base": {
			"python": "Python 3.7+ with sentencepiece package",
		},
		"mt5-base": {
			"python": "Python 3.7+ with sentencepiece package",
		},
		"albert-base": {
			"python": "Python 3.7+ with sentencepiece package",
		},
		"openai-api": {
			"api_key": "OpenAI API key required",
		},
	}

	if req, ok := requirements[name]; ok {
		return req
	}
	return map[string]string{}
}

// GetTokenizerType returns the type of a tokenizer
func GetTokenizerType(name string) string {
	types := map[string]string{
		"mock":           "custom",
		"gpt2":           "bpe",
		"gpt-3.5-turbo":  "bpe",
		"gpt-4":          "bpe",
		"roberta-base":   "bpe",
		"gpt-neo":        "bpe",
		"bert-base":      "wordpiece",
		"distilbert-base": "wordpiece",
		"t5-base":        "unigram",
		"mt5-base":       "unigram",
		"albert-base":    "wordpiece",
		"openai-api":     "bpe",
	}

	if tokenizerType, ok := types[name]; ok {
		return tokenizerType
	}
	return "unknown"
}

// GetTokenizerBackend returns the backend used by a tokenizer
func GetTokenizerBackend(name string) string {
	backends := map[string]string{
		"mock":           "go",
		"gpt2":           "tiktoken",
		"gpt-3.5-turbo":  "tiktoken",
		"gpt-4":          "tiktoken",
		"roberta-base":   "transformers",
		"gpt-neo":        "transformers",
		"bert-base":      "transformers",
		"distilbert-base": "transformers",
		"t5-base":        "sentencepiece",
		"mt5-base":       "sentencepiece",
		"albert-base":    "sentencepiece",
		"openai-api":     "api",
	}

	if backend, ok := backends[name]; ok {
		return backend
	}
	return "unknown"
} 