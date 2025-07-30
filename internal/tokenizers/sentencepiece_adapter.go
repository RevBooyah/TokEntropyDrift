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

// SentencePieceTokenizer implements the Tokenizer interface for SentencePiece models
type SentencePieceTokenizer struct {
	*BaseTokenizer
	modelPath  string
	pythonPath string
	modelType  string
}

// NewSentencePieceTokenizer creates a new SentencePiece tokenizer
func NewSentencePieceTokenizer(name string) *SentencePieceTokenizer {
	return &SentencePieceTokenizer{
		BaseTokenizer: NewBaseTokenizer(name),
		pythonPath:    "python3",
		modelType:     "unigram",
	}
}

// Initialize sets up the SentencePiece tokenizer
func (s *SentencePieceTokenizer) Initialize(config TokenizerConfig) error {
	if err := s.BaseTokenizer.Initialize(config); err != nil {
		return err
	}

	// Set model path from config
	if modelPath, ok := config.Parameters["model_path"]; ok {
		s.modelPath = modelPath
	}

	// Set Python path from config
	if pythonPath, ok := config.Parameters["python_path"]; ok {
		s.pythonPath = pythonPath
	}

	// Set model type from config
	if modelType, ok := config.Parameters["model_type"]; ok {
		s.modelType = modelType
	}

	// Validate model type
	validTypes := map[string]bool{
		"unigram": true,
		"bpe":     true,
		"char":    true,
		"word":    true,
	}

	if !validTypes[s.modelType] {
		return fmt.Errorf("invalid sentencepiece model type: %s", s.modelType)
	}

	return nil
}

// Tokenize tokenizes a single document using SentencePiece
func (s *SentencePieceTokenizer) Tokenize(ctx context.Context, text string) (*TokenizationResult, error) {
	// Create Python script for tokenization
	script := fmt.Sprintf(`
import sentencepiece as spm
import json
import sys

try:
    # Read text from stdin
    text = sys.stdin.read()
    
    # Initialize tokenizer
    sp = spm.SentencePieceProcessor()
    sp.load("%s")
    
    # Tokenize text
    pieces = sp.encode_as_pieces(text)
    ids = sp.encode_as_ids(text)
    
    # Get token positions (approximate)
    token_objects = []
    current_pos = 0
    
    for i, (piece, token_id) in enumerate(zip(pieces, ids)):
        # Estimate position based on piece length
        start_pos = current_pos
        end_pos = start_pos + len(piece)
        current_pos = end_pos
        
        token_objects.append({
            "id": token_id,
            "text": piece,
            "start_pos": start_pos,
            "end_pos": end_pos
        })
    
    # Create result
    result = {
        "document": text,
        "tokens": token_objects,
        "tokenizer": "%s",
        "metadata": {
            "model_path": "%s",
            "model_type": "%s",
            "vocab_size": sp.get_piece_size()
        }
    }
    
    print(json.dumps(result))
    
except Exception as e:
    print(json.dumps({"error": str(e)}), file=sys.stderr)
    sys.exit(1)
`, s.modelPath, s.Name(), s.modelPath, s.modelType)

	// Execute Python script with virtual environment
	cmd := exec.CommandContext(ctx, s.pythonPath, "-c", script)
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
			return nil, fmt.Errorf("sentencepiece error: %s", string(exitErr.Stderr))
		}
		return nil, fmt.Errorf("failed to execute sentencepiece: %w", err)
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
		return nil, fmt.Errorf("failed to parse sentencepiece output: %w", err)
	}

	if result.Error != "" {
		return nil, fmt.Errorf("sentencepiece error: %s", result.Error)
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
				"tokenizer":  "sentencepiece",
				"model_path": s.modelPath,
				"model_type": s.modelType,
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
func (s *SentencePieceTokenizer) TokenizeBatch(ctx context.Context, texts []string) ([]*TokenizationResult, error) {
	results := make([]*TokenizationResult, len(texts))

	for i, text := range texts {
		result, err := s.Tokenize(ctx, text)
		if err != nil {
			return nil, fmt.Errorf("error tokenizing document %d: %w", i, err)
		}
		results[i] = result
	}

	return results, nil
}

// GetVocabSize returns the vocabulary size
func (s *SentencePieceTokenizer) GetVocabSize() (int, error) {
	// Create Python script to get vocab size
	script := fmt.Sprintf(`
import sentencepiece as spm
import json

try:
    sp = spm.SentencePieceProcessor()
    sp.load("%s")
    print(json.dumps({"vocab_size": sp.get_piece_size()}))
except Exception as e:
    print(json.dumps({"error": str(e)}))
`, s.modelPath)

	cmd := exec.Command(s.pythonPath, "-c", script)

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
func (s *SentencePieceTokenizer) Close() error {
	// Nothing to clean up for SentencePiece tokenizer
	return nil
}

// RegisterT5Tokenizer registers the T5 tokenizer
func RegisterT5Tokenizer() error {
	t5Tokenizer := NewSentencePieceTokenizer("t5-base")
	t5Tokenizer.modelPath = "t5-base"
	t5Tokenizer.modelType = "unigram"
	return RegisterGlobal("t5-base", t5Tokenizer)
}

// RegisterMT5Tokenizer registers the mT5 tokenizer
func RegisterMT5Tokenizer() error {
	mt5Tokenizer := NewSentencePieceTokenizer("mt5-base")
	mt5Tokenizer.modelPath = "mt5-base"
	mt5Tokenizer.modelType = "unigram"
	return RegisterGlobal("mt5-base", mt5Tokenizer)
}

// RegisterALBERTTokenizer registers the ALBERT tokenizer
func RegisterALBERTTokenizer() error {
	albertTokenizer := NewSentencePieceTokenizer("albert-base-v2")
	albertTokenizer.modelPath = "albert-base-v2"
	albertTokenizer.modelType = "wordpiece"
	return RegisterGlobal("albert-base", albertTokenizer)
}
