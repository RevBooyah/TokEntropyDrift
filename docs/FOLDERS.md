# Project Folder Structure: TokEntropyDrift ("ted")

This document outlines the standard folder layout for the TokEntropyDrift codebase. It is designed to help contributors and AI development tools (e.g., Cursor, Claude) quickly understand where to place files, add new components, or locate outputs.

---

## 📂 Core Project Structure

```
/cmd/ted/               # CLI entrypoint and orchestration
/internal/
  tokenizers/           # Tokenizer adapters (GPT-2, T5, custom)
  metrics/              # Core metric functions (entropy, compression, reuse)
  parser/               # Input loaders for .txt, .csv, .jsonl
  exporter/             # Handles JSON, CSV, Markdown, LaTeX output
  visualizer/           # Generates visual data and image files
  logging/              # Structured logging utilities (JSON logs)

/modules/               # Optional/experimental extensions and plugins
/tokenizers/            # User-supplied tokenizer config, vocab, model files
/examples/              # Input corpora for testing and validation
/testdata/              # Expected output for golden tests
/output/                # All runtime-generated output (one folder per run)
  run_YYYY-MM-DD_HHMMSS/
    config.yaml         # Snapshot of CLI/config used
    metrics.csv         # Full corpus analysis output
    entropy.json        # Per-token or per-line entropy metadata
    visualizations/     # PNG/SVG renderings (heatmaps, timelines, overlays)
    logs/               # JSON logs for debugging and traceability
/docs/                  # Markdown documentation and module references
```

---

## 🧠 Detailed Folder Guide

### `/cmd/ted/`

* Entry point for the CLI tool
* Only orchestration logic — delegates all real work to modules

### `/internal/tokenizers/`

* One file per tokenizer adapter (`gpt2.go`, `t5.go`, etc.)
* Must implement the common Tokenizer interface

### `/internal/metrics/`

* Functions for entropy, compression, reuse, and token alignment
* Group by conceptual domain to avoid unnecessary fragmentation

### `/internal/visualizer/`

* Responsible for rendering images (Plotly.js configs, SVG/PNG generation)
* Reads data from disk (JSON/CSV) and exports visuals

### `/modules/`

* Experimental or advanced extensions
* Examples: cost overlays, prompt perturbation tools, language-specific add-ons

### `/tokenizers/`

* Contains vocab/model/config files for runtime-registered tokenizers
* YAML/TOML files should be autodetected by the CLI

### `/examples/`

* Input corpora (plain text, CSV, JSONL)
* One folder per use case (e.g., `english/`, `code/`, `news/`)

### `/testdata/`

* Golden file outputs to validate tokenizer adapters and metrics
* Use for regression detection and output comparison

### `/output/`

* Each run creates a timestamped folder with all outputs
* Keeps metrics, images, config, logs grouped for reproducibility
* Makes headless mode and CI comparisons straightforward

---

## 🔍 Tips for AI/Cursor Use

* Never put tokenizer logic or metrics in `/cmd/` — delegate to internal modules
* Use `/tokenizers/` for any non-Go resources related to tokenization (vocab/model)
* Place all runtime outputs under `/output/` to keep structure clean
* If in doubt, keep components modular, well-documented, and testable

This structure ensures TokEntropyDrift remains scalable, research-friendly, and LLM-co-developer ready.
