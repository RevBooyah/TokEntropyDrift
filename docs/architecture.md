# Architecture Overview: TokEntropyDrift ("ted")

This document outlines the system design, execution flow, modular components, and extension points for the TokEntropyDrift platform.

---

## 🚀 Execution Modes

TokEntropyDrift supports two primary execution modes:

1. **Command-Line Interface (CLI)**

   * All actions initiated from the CLI (tokenization, analysis, visualization, export)
   * Configurable via flags or config file
2. **Web Dashboard**

   * Visualizations are served as static or interactive content via a lightweight web server
   * `ted serve` starts the server to host latest results

---

## 🧩 System Modules (Pipeline)

```text
[ CLI ]
   |
   v
[ Input Loader ] → [ Tokenizer Adapter(s) ] → [ Metric Engine ] → [ Export Manager ]
                                                              ↘
                                                           [ Visualization Server ]
```

Each box is a discrete, swappable module.

### Core Modules

* **Input Loader**: Handles plain text, JSONL, or CSV. Tracks file/line metadata.
* **Tokenizer Adapter**: Wraps any tokenizer into a unified interface.
* **Metric Engine**: Computes all statistics (entropy, reuse, token counts, compression).
* **Export Manager**: Handles writing CSV/JSON/Markdown/LaTeX + visual asset placement.
* **Visualization Server**: Serves static UI with Plotly.js/D3.js assets.

---

## 🔌 Tokenizer Plugin System

* Tokenizers are loaded via YAML or runtime-registered via CLI
* Adapter pattern allows calling Python scripts, subprocess binaries, or Go-wrapped native libs
* All adapters must conform to a shared Go interface
* Tokenizer config includes:

  * name
  * type (bpe, spiece, etc.)
  * command or library path
  * vocab/model files (if needed)

```bash
$ ted tokenizer add ./tokenizers/my_custom_model --name=custom1
```

---

## ⚙️ Config System

TokEntropyDrift uses a single `ted.config.yaml` or `.toml` file per project run:

* Input source path(s)
* Tokenizers to use
* Entropy parameters (window size, normalization)
* Output format preferences (image size, filetype, export toggles)
* Visualization defaults (theme, layout)

Overrides can be passed via CLI flags.

---

## 🗃 Output File Layout

All outputs from a run are written into a timestamped folder:

```
/output/
  run_2025-07-27_13-42-15/
    config.yaml
    metrics.csv
    entropy.json
    visualizations/
    reports/
    logs/
```

---

## 🧠 Visualization Engine

* Powered by Plotly.js (interactive) + D3 for advanced maps
* `ted serve` starts a static file server from an output directory
* Visuals are generated from metric exports (no runtime recomputation)

---

## 🧰 Logging & Debugging

* Logs are emitted in structured **JSON** format for each run
* Per-batch debug output (line-level optional)
* Logged stages:

  * Input ingestion
  * Tokenization success/fail per tokenizer
  * Metric generation stats
  * Export/visualization confirmation

---

## 🧪 Testing & Sample Corpora

* Unit testing is minimal — focus is on **end-to-end testing** with full corpora
* Multiple corpora will be provided:

  * `samples/english_quotes.txt`
  * `samples/code_snippets.jsonl`
  * `samples/multilingual.csv`
* Golden outputs stored in `/testdata/`
* CLI test runner planned:

```bash
$ ted test corpus.jsonl --tokenizers=gpt2,custom --compare-to=./testdata/gpt2_ref.csv
```

---

## 🔧 Extensibility Roadmap

* All analysis steps modularized under `/modules/`
* Future support for:

  * Cost modeling modules
  * Language-specific entropy modules
  * Prompt sensitivity visual diff tools
  * CLI-only plugin loader (experimental flags)

---

This architecture is designed for research-grade analysis with dev-grade extensibility. All components are modular and testable at the batch level.
