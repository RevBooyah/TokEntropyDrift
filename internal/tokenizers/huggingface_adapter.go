package tokenizers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// HuggingFaceTokenizer implements the Tokenizer interface for HuggingFace tokenizers
type HuggingFaceTokenizer struct {
	*BaseTokenizer
	modelName     string
	pythonPath    string
	modelPath     string
	tokenizerType string
}

// NewHuggingFaceTokenizer creates a new HuggingFace tokenizer
func NewHuggingFaceTokenizer(name string) *HuggingFaceTokenizer {
	return &HuggingFaceTokenizer{
		BaseTokenizer: NewBaseTokenizer(name),
		pythonPath:    "python3",
		tokenizerType: "bpe",
	}
}

// Initialize sets up the HuggingFace tokenizer
func (h *HuggingFaceTokenizer) Initialize(config TokenizerConfig) error {
	if err := h.BaseTokenizer.Initialize(config); err != nil {
		return err
	}

	// Set model name from config
	if model, ok := config.Parameters["model"]; ok {
		h.modelName = model
	}

	// Set model path from config
	if modelPath, ok := config.Parameters["model_path"]; ok {
		h.modelPath = modelPath
	}

	// Set Python path from config
	if pythonPath, ok := config.Parameters["python_path"]; ok {
		h.pythonPath = pythonPath
	}

	// Set tokenizer type from config
	if tokenizerType, ok := config.Parameters["tokenizer_type"]; ok {
		h.tokenizerType = tokenizerType
	}

	// Validate tokenizer type
	validTypes := map[string]bool{
		"bpe":       true,
		"wordpiece": true,
		"spiece":    true,
	}

	if !validTypes[h.tokenizerType] {
		return fmt.Errorf("invalid tokenizer type: %s", h.tokenizerType)
	}

	return nil
}

// Tokenize tokenizes a single document using HuggingFace tokenizers
func (h *HuggingFaceTokenizer) Tokenize(ctx context.Context, text string) (*TokenizationResult, error) {
	// Create Python script for tokenization
	script := fmt.Sprintf(`
from transformers import AutoTokenizer
import json
import sys

try:
    # Read text from stdin
    text = sys.stdin.read()
    
    # Initialize tokenizer
    if "%s":
        tokenizer = AutoTokenizer.from_pretrained("%s")
    else:
        tokenizer = AutoTokenizer.from_pretrained("%s")
    
    # Tokenize text
    encoding = tokenizer(text, return_offsets_mapping=True, add_special_tokens=False)
    
    # Extract tokens and positions
    tokens = encoding.tokens()
    offset_mapping = encoding.offset_mapping
    input_ids = encoding.input_ids
    
    # Create token objects
    token_objects = []
    for i, (token, (start, end)) in enumerate(zip(tokens, offset_mapping)):
        token_objects.append({
            "id": input_ids[i] if i < len(input_ids) else 0,
            "text": token,
            "start_pos": start,
            "end_pos": end
        })
    
    # Create result
    result = {
        "document": text,
        "tokens": token_objects,
        "tokenizer": "%s",
        "metadata": {
            "model": "%s",
            "tokenizer_type": "%s",
            "vocab_size": tokenizer.vocab_size
        }
    }
    
    print(json.dumps(result))
    
except Exception as e:
    print(json.dumps({"error": str(e)}), file=sys.stderr)
    sys.exit(1)
`, h.modelPath, h.modelPath, h.modelName, h.Name(), h.modelName, h.tokenizerType)

	// Execute Python script with virtual environment
	cmd := exec.CommandContext(ctx, h.pythonPath, "-c", script)
	cmd.Stdin = strings.NewReader(text)

	// Set virtual environment variables
	cmd.Env = append(os.Environ(),
		"VIRTUAL_ENV="+filepath.Join(".", "venv"),
		"PATH="+filepath.Join(".", "venv", "bin")+":"+os.Getenv("PATH"),
	)

	output, err := cmd.Output()
	if err != nil {
		// Try to get error output
		if exitErr, ok := err.(*exec.ExitError); ok {
			return nil, fmt.Errorf("huggingface tokenizer error: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to execute huggingface tokenizer: %w", err)
	}

	// Parse JSON output
	var result struct {
		Document string `json:"document"`
		Tokens   []struct {
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
		return nil, fmt.Errorf("failed to parse huggingface tokenizer output: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("huggingface tokenizer error: %s", result.Error)
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
				"tokenizer":      "huggingface",
				"model":          h.modelName,
				"tokenizer_type": h.tokenizerType,
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
func (h *HuggingFaceTokenizer) TokenizeBatch(ctx context.Context, texts []string) ([]*TokenizationResult, error) {
	results := make([]*TokenizationResult, len(texts))

	for i, text := range texts {
		result, err := h.Tokenize(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("error tokenizing document %d: %w", i, err)
		}
		results[i] = result
	}

	return results, nil
}

// GetVocabSize returns the vocabulary size
func (h *HuggingFaceTokenizer) GetVocabSize() (int, error) {
	// Create Python script to get vocab size
	script := fmt.Sprintf(`
from transformers import AutoTokenizer
import json

try:
    if "%s":
        tokenizer = AutoTokenizer.from_pretrained("%s")
    else:
        tokenizer = AutoTokenizer.from_pretrained("%s")
    
    print(json.dumps({"vocab_size": tokenizer.vocab_size}))
except Exception as e:
    print(json.dumps({"error": str(e)}))
`, h.modelPath, h.modelPath, h.modelName)

	cmd := exec.Command(h.pythonPath, "-c", script)

	// Set virtual environment variables
	cmd.Env = append(os.Environ(),
		"VIRTUAL_ENV="+filepath.Join(".", "venv"),
		"PATH="+filepath.Join(".", "venv", "bin")+":"+os.Getenv("PATH"),
	)

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
func (h *HuggingFaceTokenizer) Close() error {
	// Nothing to clean up for HuggingFace tokenizer
	return nil
}

// RegisterRoBERTaTokenizer registers the RoBERTa tokenizer
func RegisterRoBERTaTokenizer() error {
	robertaTokenizer := NewHuggingFaceTokenizer("roberta-base")
	robertaTokenizer.modelName = "roberta-base"
	robertaTokenizer.tokenizerType = "bpe"
	return RegisterGlobal("roberta-base", robertaTokenizer)
}

// RegisterGPTNeoTokenizer registers the GPT-Neo tokenizer
func RegisterGPTNeoTokenizer() error {
	gptNeoTokenizer := NewHuggingFaceTokenizer("EleutherAI/gpt-neo-125M")
	gptNeoTokenizer.modelName = "EleutherAI/gpt-neo-125M"
	gptNeoTokenizer.tokenizerType = "bpe"
	return RegisterGlobal("gpt-neo", gptNeoTokenizer)
}

// RegisterBERTTokenizer registers the BERT tokenizer
func RegisterBERTTokenizer() error {
	bertTokenizer := NewHuggingFaceTokenizer("bert-base-uncased")
	bertTokenizer.modelName = "bert-base-uncased"
	bertTokenizer.tokenizerType = "wordpiece"
	return RegisterGlobal("bert-base", bertTokenizer)
}

// RegisterDistilBERTTokenizer registers the DistilBERT tokenizer
func RegisterDistilBERTTokenizer() error {
	distilBertTokenizer := NewHuggingFaceTokenizer("distilbert-base-uncased")
	distilBertTokenizer.modelName = "distilbert-base-uncased"
	distilBertTokenizer.tokenizerType = "wordpiece"
	return RegisterGlobal("distilbert-base", distilBertTokenizer)
}
