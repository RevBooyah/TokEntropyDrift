# TokEntropyDrift Setup Guide

This guide covers setting up TokEntropyDrift with working tokenizers on a fresh installation.

## Prerequisites

- Go 1.19+ installed
- Python 3.7+ installed
- Git (for cloning the repository)

## Quick Start

1. **Clone and navigate to the repository:**
   ```bash
   git clone <repository-url>
   cd ted
   ```

2. **Set up Python virtual environment:**
   ```bash
   python3 -m venv venv
   source venv/bin/activate
   pip install tiktoken transformers sentencepiece
   ```

3. **Start the server:**
   ```bash
   source venv/bin/activate
   go run cmd/ted/main.go serve
   ```

4. **Access the web interface:**
   Open http://localhost:8080 in your browser

## Detailed Setup

### Python Environment Setup

The tokenizers require Python packages that aren't available in the system Python. We use a virtual environment to avoid conflicts.

#### Step 1: Create Virtual Environment
```bash
python3 -m venv venv
```

#### Step 2: Activate Virtual Environment
```bash
source venv/bin/activate
```

#### Step 3: Install Required Packages
```bash
pip install tiktoken transformers sentencepiece
```

**Package Details:**
- `tiktoken`: Required for GPT-2, GPT-3.5-turbo, GPT-4 tokenizers
- `transformers`: Required for HuggingFace tokenizers (BERT, RoBERTa, DistilBERT, GPT-Neo)
- `sentencepiece`: Required for SentencePiece tokenizers (T5, mT5, ALBERT)

### Server Configuration

The server is configured to use the virtual environment's Python interpreter. The configuration is already set in `ted.config.yaml`:

```yaml
tokenizers:
  enabled: ["mock", "gpt2", "gpt-3.5-turbo", "gpt-4", "roberta-base", "bert-base", "distilbert-base"]
  configs:
    gpt2:
      type: "bpe"
      parameters:
        model: "gpt2"
        python_path: "./venv/bin/python"
    bert-base:
      type: "wordpiece"
      parameters:
        model: "bert-base-uncased"
        python_path: "./venv/bin/python"
    # ... other tokenizers
```

### Starting the Server

Always activate the virtual environment before starting the server:

```bash
source venv/bin/activate
go run cmd/ted/main.go serve
```

**Important:** The virtual environment must be activated for the tokenizers to work properly.

## Working Tokenizers

After setup, the following tokenizers will be available:

### ✅ Fully Working (No Additional Setup)
- **Mock**: Word-based tokenizer for testing
- **GPT-2**: BPE tokenizer using tiktoken
- **GPT-3.5-turbo**: BPE tokenizer using tiktoken  
- **GPT-4**: BPE tokenizer using tiktoken
- **BERT**: WordPiece tokenizer using HuggingFace
- **RoBERTa**: BPE tokenizer using HuggingFace
- **DistilBERT**: WordPiece tokenizer using HuggingFace
- **GPT-Neo**: BPE tokenizer using HuggingFace

### ⚠️ Requires Additional Setup
- **T5, mT5, ALBERT**: Require model downloads (SentencePiece)
- **OpenAI API**: Requires API key configuration

## Troubleshooting

### Port Already in Use
If you get "address already in use" error:
```bash
lsof -ti:8080 | xargs kill -9
```

### Python Packages Not Found
If tokenizers fail with "ModuleNotFoundError":
1. Ensure virtual environment is activated: `source venv/bin/activate`
2. Verify packages are installed: `pip list | grep -E "(tiktoken|transformers|sentencepiece)"`
3. Reinstall if needed: `pip install --force-reinstall tiktoken transformers sentencepiece`

### HuggingFace Model Downloads
Some tokenizers may need to download models on first use. This happens automatically but may take time depending on your internet connection.

## Development Notes

### Virtual Environment Management
- The virtual environment is stored in `./venv/`
- Always activate before running the server
- Add `venv/` to `.gitignore` (already included)

### Configuration Files
- `ted.config.yaml`: Main configuration file
- Tokenizer settings are in the `tokenizers` section
- Python paths are set to `./venv/bin/python`

### Tokenizer Architecture
- Go server calls Python scripts via subprocess
- Python scripts read input from stdin
- Results are returned as JSON
- Each tokenizer type has its own adapter in `internal/tokenizers/`

## Next Steps

1. **Upload documents** through the web interface
2. **Select tokenizers** for analysis
3. **Run analysis** to compare different tokenization methods
4. **View results** in the Analysis History section

For advanced usage, see the main README.md file. 