# Quick Reference

## Essential Commands

### Setup (First Time Only)
```bash
# Create Python environment
python3 -m venv venv
source venv/bin/activate
pip install tiktoken transformers sentencepiece

# Start server
source venv/bin/activate
go run cmd/ted/main.go serve
```

### Daily Usage
```bash
# Start server (always activate venv first)
source venv/bin/activate
go run cmd/ted/main.go serve

# Kill server if port is busy
lsof -ti:8080 | xargs kill -9
```

## Working Tokenizers

| Tokenizer | Type | Status | Notes |
|-----------|------|--------|-------|
| Mock | Word-based | ✅ Working | No dependencies |
| GPT-2 | BPE (tiktoken) | ✅ Working | Fast |
| GPT-3.5-turbo | BPE (tiktoken) | ✅ Working | Fast |
| GPT-4 | BPE (tiktoken) | ✅ Working | Fast |
| BERT | WordPiece (HF) | ✅ Working | Downloads model |
| RoBERTa | BPE (HF) | ✅ Working | Downloads model |
| DistilBERT | WordPiece (HF) | ✅ Working | Downloads model |
| GPT-Neo | BPE (HF) | ✅ Working | Downloads model |

## Common Issues

### "ModuleNotFoundError"
```bash
# Solution: Activate virtual environment
source venv/bin/activate
```

### "Port already in use"
```bash
# Solution: Kill existing process
lsof -ti:8080 | xargs kill -9
```

### "Analysis completed with 0 results"
- Check that virtual environment is activated
- Verify Python packages are installed: `pip list | grep tiktoken`
- Try with mock tokenizer first

## Web Interface

1. **Upload Document**: Use the file upload in the web interface
2. **Select Tokenizers**: Choose from the available tokenizers
3. **Run Analysis**: Click "Run Analysis" button
4. **View Results**: Check Analysis History and Results sections

## API Testing

Test individual tokenizers:
```bash
curl -X POST -H "Content-Type: application/json" \
  -d '{"document_id":"test","tokenizer_ids":["bert-base"],"metrics":["entropy"]}' \
  http://localhost:8080/api/v1/analyze
```

## Configuration

Main config file: `ted.config.yaml`
- Tokenizer settings in `tokenizers` section
- Python paths set to `./venv/bin/python`
- Server port: 8080

## File Locations

- **Virtual Environment**: `./venv/`
- **Configuration**: `ted.config.yaml`
- **Uploads**: `output/` directory
- **Logs**: Server console output 