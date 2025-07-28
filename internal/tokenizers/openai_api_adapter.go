package tokenizers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// OpenAITokenizer implements the Tokenizer interface for OpenAI API
type OpenAITokenizer struct {
	*BaseTokenizer
	apiKey     string
	apiBase    string
	modelName  string
	httpClient *http.Client
}

// NewOpenAITokenizer creates a new OpenAI API tokenizer
func NewOpenAITokenizer(name string) *OpenAITokenizer {
	return &OpenAITokenizer{
		BaseTokenizer: NewBaseTokenizer(name),
		apiBase:       "https://api.openai.com/v1",
		modelName:     "gpt-3.5-turbo",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// Initialize sets up the OpenAI API tokenizer
func (o *OpenAITokenizer) Initialize(config TokenizerConfig) error {
	if err := o.BaseTokenizer.Initialize(config); err != nil {
		return err
	}

	// Set API key from config
	if apiKey, ok := config.Parameters["api_key"]; ok {
		o.apiKey = apiKey
	}

	// Set API base URL from config
	if apiBase, ok := config.Parameters["api_base"]; ok {
		o.apiBase = apiBase
	}

	// Set model name from config
	if model, ok := config.Parameters["model"]; ok {
		o.modelName = model
	}

	// Validate required fields
	if o.apiKey == "" {
		return fmt.Errorf("OpenAI API key is required")
	}

	return nil
}

// Tokenize tokenizes a single document using OpenAI API
func (o *OpenAITokenizer) Tokenize(ctx context.Context, text string) (*TokenizationResult, error) {
	// Create request payload
	payload := map[string]interface{}{
		"model": o.modelName,
		"input": text,
	}

	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request payload: %w", err)
	}

	// Create HTTP request
	url := fmt.Sprintf("%s/tokenize", o.apiBase)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+o.apiKey)

	// Make request
	resp, err := o.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to make API request: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check for errors
	if resp.StatusCode != http.StatusOK {
		var errorResp struct {
			Error struct {
				Message string `json:"message"`
			} `json:"error"`
		}
		if err := json.Unmarshal(body, &errorResp); err != nil {
			return nil, fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(body))
		}
		return nil, fmt.Errorf("API error: %s", errorResp.Error.Message)
	}

	// Parse response
	var apiResp struct {
		Tokens []struct {
			ID   int    `json:"id"`
			Text string `json:"text"`
		} `json:"tokens"`
	}

	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}

	// Convert to our token format
	tokens := make([]Token, len(apiResp.Tokens))
	currentPos := 0

	for i, t := range apiResp.Tokens {
		startPos := currentPos
		endPos := startPos + len(t.Text)
		currentPos = endPos

		tokens[i] = Token{
			Text:     t.Text,
			ID:       t.ID,
			StartPos: startPos,
			EndPos:   endPos,
			Metadata: map[string]string{
				"tokenizer": "openai_api",
				"model":     o.modelName,
				"api_base":  o.apiBase,
			},
		}
	}

	return &TokenizationResult{
		Document:  text,
		Tokens:    tokens,
		Tokenizer: o.Name(),
		Metadata: map[string]interface{}{
			"model":     o.modelName,
			"api_base":  o.apiBase,
			"tokenizer": "openai_api",
		},
	}, nil
}

// TokenizeBatch tokenizes multiple documents
func (o *OpenAITokenizer) TokenizeBatch(ctx context.Context, texts []string) ([]*TokenizationResult, error) {
	results := make([]*TokenizationResult, len(texts))
	
	for i, text := range texts {
		result, err := o.Tokenize(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("error tokenizing document %d: %w", i, err)
		}
		results[i] = result
	}
	
	return results, nil
}

// GetVocabSize returns the vocabulary size (approximate for OpenAI models)
func (o *OpenAITokenizer) GetVocabSize() (int, error) {
	// OpenAI doesn't provide vocab size via API, so we return approximate values
	vocabSizes := map[string]int{
		"gpt-3.5-turbo": 100277,
		"gpt-4":         100277,
		"gpt-4-turbo":   100277,
		"gpt-4o":        100277,
	}

	if vocabSize, ok := vocabSizes[o.modelName]; ok {
		return vocabSize, nil
	}

	// Default vocab size for unknown models
	return 100277, nil
}

// Close cleans up resources
func (o *OpenAITokenizer) Close() error {
	// Nothing to clean up for OpenAI API tokenizer
	return nil
}

// RegisterOpenAITokenizer registers the OpenAI API tokenizer
func RegisterOpenAITokenizer() error {
	openAITokenizer := NewOpenAITokenizer("openai-api")
	return RegisterGlobal("openai-api", openAITokenizer)
} 