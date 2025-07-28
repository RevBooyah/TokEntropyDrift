# TokEntropyDrift

*Tokenization Entropy and Drift Explorer ("ted")*

## 📌 Project Overview

**TokEntropyDrift** is a research and development tool for exploring, analyzing, and visualizing the behavior of modern tokenizers across various LLM ecosystems. It helps engineers, researchers, and prompt designers understand how different tokenizers fragment text, how that impacts entropy and compression, and how tokenization patterns vary between models.

This project is designed for:

* NLP researchers comparing tokenization schemes
* LLM developers analyzing tokenizer effects on performance and cost
* Prompt engineers designing for optimal input efficiency
* AI practitioners interested in interpretability of model input structures

## 🎯 Goals

* Build a comparative heatmap visualization tool for multiple tokenizer outputs
* Measure token-level entropy and compression metrics per corpus and tokenizer
* Support common tokenizer types (BPE, SentencePiece, WordPiece, tiktoken, etc.)
* Provide structured, exportable metrics for research use

## 🧩 Core Features

* Multi-tokenizer interface abstraction
* Corpus-level token count and length analytics
* Token boundary visualization and comparison
* Entropy and token reuse statistics
* CSV/JSON/Markdown export
* D3/Plotly-based web dashboard (and CLI option)

## 🧠 Supported Tokenizers

* GPT-2 / OpenAI via `tiktoken`
* HuggingFace BPE (e.g. RoBERTa, GPT-Neo)
* SentencePiece (e.g. T5, mT5)
* WordPiece (e.g. BERT)
* Optional: Custom or user-supplied vocab/tokenizer scripts

## 📊 Visualization Types

* Heatmaps: tokens-per-line, token-length, entropy
* Side-by-side token boundary overlays
* Entropy distribution plots
* Token cost charts (for API budget estimation)

## ⚙️ Technology Stack

* Language: Go (primary CLI and core), Python (tokenizers bridge)
* Web UI: HTML + D3.js or Plotly
* CLI: Cobra / urfave/cli for flags and execution
* Output: JSON, CSV, Markdown

## 🚀 Quick Start

```bash
# Clone the repo
$ git clone https://github.com/RevBooyah/TokEntropyDrift.git
$ cd TokEntropyDrift

# Build CLI (Go)
$ go build -o ted ./cmd/ted

# Analyze sample.txt with GPT-2 and T5
$ ./ted analyze sample.txt --tokenizers=gpt2,t5 --output=out.csv

# Run web dashboard
$ ./ted serve --port=8080
```

## 📁 Example Usage

```bash
# Get entropy heatmap of multiple tokenizers on a dataset
$ ./ted heatmap corpus.txt --tokenizers=gpt2,spiece,custom

# Export results in markdown report format
$ ./ted export --format=md --output=report.md
```

## 🤝 Contribution Guide

* All tokenizer logic is modular (see `internal/tokenizers/`)
* Visualization layer is web-independent (can use D3 or CLI)
* Contributions welcome via issues, PRs, or discussions
* Please follow semantic commits and document new metrics

## 📝 License

MIT License — see `LICENSE` file for details.

---

For detailed documentation, see [`docs/`](./docs/)
