package tokenizers

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// GPT2Tokenizer implements the Tokenizer interface for GPT-2/GPT-3.5/GPT-4 models
type GPT2Tokenizer struct {
	*BaseTokenizer
	modelName string
	pythonPath string
}

// NewGPT2Tokenizer creates a new GPT-2 tokenizer
func NewGPT2Tokenizer(name string) *GPT2Tokenizer {
	return &GPT2Tokenizer{
		BaseTokenizer: NewBaseTokenizer(name),
		modelName:     "gpt2",
		pythonPath:    "python3",
	}
}

// Initialize sets up the GPT-2 tokenizer
func (g *GPT2Tokenizer) Initialize(config TokenizerConfig) error {
	if err := g.BaseTokenizer.Initialize(config); err != nil {
		return err
	}

	// Set model name from config
	if model, ok := config.Parameters["model"]; ok {
		g.modelName = model
	}

	// Set Python path from config
	if pythonPath, ok := config.Parameters["python_path"]; ok {
		g.pythonPath = pythonPath
	}

	// Validate model name
	validModels := map[string]bool{
		"gpt2":        true,
		"gpt2-medium": true,
		"gpt2-large":  true,
		"gpt2-xl":     true,
		"gpt-3.5-turbo": true,
		"gpt-4":       true,
	}

	if !validModels[g.modelName] {
		return fmt.Errorf("invalid GPT model: %s", g.modelName)
	}

	return nil
}

// Tokenize tokenizes a single document using tiktoken
func (g *GPT2Tokenizer) Tokenize(ctx context.Context, text string) (*TokenizationResult, error) {
	// Create Python script for tokenization
	script := fmt.Sprintf(`
import tiktoken
import json
import sys

try:
    # Initialize tokenizer
    encoding = tiktoken.encoding_for_model("%s")
    
    # Tokenize text
    tokens = encoding.encode(text)
    
    # Get token texts
    token_texts = []
    for token_id in tokens:
        token_text = encoding.decode([token_id])
        token_texts.append({
            "id": token_id,
            "text": token_text,
            "start_pos": 0,  # tiktoken doesn't provide position info
            "end_pos": len(token_text)
        })
    
    # Create result
    result = {
        "document": text,
        "tokens": token_texts,
        "tokenizer": "%s",
        "metadata": {
            "model": "%s",
            "vocab_size": encoding.n_vocab
        }
    }
    
    print(json.dumps(result))
    
except Exception as e:
    print(json.dumps({"error": str(e)}), file=sys.stderr)
    sys.exit(1)
`, g.modelName, g.Name(), g.modelName)

	// Execute Python script
	cmd := exec.CommandContext(ctx, g.pythonPath, "-c", script)
	cmd.Stdin = strings.NewReader(text)
	
	output, err := cmd.Output()
	if err != nil {
		// Try to get error output
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("tiktoken error: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to execute tiktoken: %w", err)
	}

	// Parse JSON output
	var result struct {
		Document  string `json:"document"`
		Tokens    []struct {
			ID       int    `json:"id"`
			Text     string `json:"text"`
			StartPos int    `json:"start_pos"`
			EndPos   int    `json:"end_pos"`
		} `json:"tokens"`
		Tokenizer string                 `json:"tokenizer"`
		Metadata  map[string]interface{} `json:"metadata"`
		Error     string                 `json:"error,omitempty"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return nil, fmt.Errorf("failed to parse tiktoken output: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("tiktoken error: %s", result.Error)
	}

	// Convert to our token format
	tokens := make([]Token, len(result.Tokens))
	for i, t := range result.Tokens {
		tokens[i] = Token{
			Text:     t.Text,
			ID:       t.ID,
			StartPos: t.StartPos,
			EndPos:   t.EndPos,
			Metadata: map[string]string{
				"tokenizer": "gpt2",
				"model":     g.modelName,
			},
		}
	}

	return &TokenizationResult{
		Document:  result.Document,
		Tokens:    tokens,
		Tokenizer: result.Tokenizer,
		Metadata:  result.Metadata,
	}, nil
}

// TokenizeBatch tokenizes multiple documents
func (g *GPT2Tokenizer) TokenizeBatch(ctx context.Context, texts []string) ([]*TokenizationResult, error) {
	results := make([]*TokenizationResult, len(texts))
	
	for i, text := range texts {
		result, err := g.Tokenize(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("error tokenizing document %d: %w", i, err)
		}
		results[i] = result
	}
	
	return results, nil
}

// GetVocabSize returns the vocabulary size
func (g *GPT2Tokenizer) GetVocabSize() (int, error) {
	// Create Python script to get vocab size
	script := fmt.Sprintf(`
import tiktoken
import json

try:
    encoding = tiktoken.encoding_for_model("%s")
    print(json.dumps({"vocab_size": encoding.n_vocab}))
except Exception as e:
    print(json.dumps({"error": str(e)}))
`, g.modelName)

	cmd := exec.Command(g.pythonPath, "-c", script)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get vocab size: %w", err)
	}

	var result struct {
		VocabSize int    `json:"vocab_size"`
		Error     string `json:"error,omitempty"`
	}

	if err := json.Unmarshal(output, &result); err != nil {
		return 0, fmt.Errorf("failed to parse vocab size output: %w", err)
	}

	if result.Error != "" {
		return 0, fmt.Errorf("error getting vocab size: %s", result.Error)
	}

	return result.VocabSize, nil
}

// Close cleans up resources
func (g *GPT2Tokenizer) Close() error {
	// Nothing to clean up for GPT-2 tokenizer
	return nil
}

// RegisterGPT2Tokenizer registers the GPT-2 tokenizer with the global registry
func RegisterGPT2Tokenizer() error {
	gpt2Tokenizer := NewGPT2Tokenizer("gpt2")
	return RegisterGlobal("gpt2", gpt2Tokenizer)
}

// RegisterGPT35Tokenizer registers the GPT-3.5 tokenizer
func RegisterGPT35Tokenizer() error {
	gpt35Tokenizer := NewGPT2Tokenizer("gpt-3.5-turbo")
	gpt35Tokenizer.modelName = "gpt-3.5-turbo"
	return RegisterGlobal("gpt-3.5-turbo", gpt35Tokenizer)
}

// RegisterGPT4Tokenizer registers the GPT-4 tokenizer
func RegisterGPT4Tokenizer() error {
	gpt4Tokenizer := NewGPT2Tokenizer("gpt-4")
	gpt4Tokenizer.modelName = "gpt-4"
	return RegisterGlobal("gpt-4", gpt4Tokenizer)
} 