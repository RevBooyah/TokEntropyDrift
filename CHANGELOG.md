# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial release of TokEntropyDrift
- CLI-driven tokenizer analysis
- Multi-tokenizer comparison
- Entropy & compression metrics
- Token reuse & drift detection
- Token boundary overlays
- Rolling entropy visualizations
- Heatmaps (tokens, entropy)
- Static & headless image export
- Tokenizer plugin architecture
- Caching layer for performance
- Parallel processing
- Streaming analysis
- Plugin system for custom metrics
- Web dashboard with interactive visualizations

### Supported Tokenizers
- GPT-2 / GPT-3.5 / GPT-4 (via `tiktoken`)
- HuggingFace BPE (e.g. RoBERTa, GPT-Neo)
- SentencePiece (e.g. T5, mT5)
- WordPiece (e.g. BERT, DistilBERT)
- OpenAI API tokenizer
- Custom tokenizers (via config + vocab/model files)

## [1.0.0] - 2025-01-XX

### Added
- Initial release
- Core tokenization analysis functionality
- Basic visualization capabilities
- CLI interface
- Configuration system

### Changed
- N/A

### Deprecated
- N/A

### Removed
- N/A

### Fixed
- N/A

### Security
- N/A 