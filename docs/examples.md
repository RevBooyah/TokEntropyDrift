# Examples and Sample Corpora

This document outlines the included examples and how they demonstrate the functionality, testing capability, and visual output of **TokEntropyDrift ("ted")**.

---

## 📁 Directory Structure

```
/examples
  english_quotes.txt
  tech_stack_descriptions.txt
  source_code_snippets.txt
  structured_news.csv
  reddit_threads.jsonl
  annotated/
    english_quotes.tokens.json
    tech_stack_descriptions.entropy.csv
/output/examples/
  *.png
  *.svg
  *.csv
```

---

## 📚 Sample Corpora

### 1. `english_quotes.txt`

* Short, single-sentence quotes (\~50 lines)
* Ideal for:

  * Quick tokenization comparisons
  * Rolling entropy tests
  * Token reuse metrics

### 2. `tech_stack_descriptions.txt`

* Medium-length paragraphs describing various software systems
* Demonstrates:

  * Token length analysis
  * Compression ratio
  * Cross-tokenizer boundary misalignment

### 3. `source_code_snippets.txt`

* Code blocks from real open source projects (Python, JS, Go)
* Useful for:

  * BPE behavior with symbols
  * WordPiece fragmentation
  * Token inflation analysis

### 4. `structured_news.csv`

* Realistic headline + summary fields
* Used to test CSV ingestion and per-column tokenization

### 5. `reddit_threads.jsonl`

* JSON lines: `title`, `body`, `tags`
* Multiline, informal language; great for drift analysis

---

## 🧪 Annotated Output Examples

Each primary corpus is paired with:

* **Tokenized outputs** (JSON with per-token data)
* **Entropy metrics** (CSV)
* **Cross-tokenizer deltas** (e.g., GPT-2 vs T5)

Example:

```json
{
  "input": "I love transformers!",
  "tokenizers": {
    "gpt2": ["I", " love", " transform", "ers", "!"],
    "t5": ["▁I", "▁love", "▁transformers", "!"],
    "delta": {
      "token_count": 5 vs 4,
      "entropy_diff": 0.271 bits
    }
  }
}
```

---

## 🎨 Sample Visualizations (Output Gallery)

All visual assets for examples are stored under `/output/examples/`

* `english_quotes_entropy.svg`
* `tech_stack_token_drift.png`
* `source_code_boundary_overlay.png`
* `reddit_threads_heatmap.svg`

---

## 🧪 CLI Demonstrations

```bash
# Run tokenization + analysis on short text
$ ted analyze examples/english_quotes.txt --tokenizers=gpt2,t5

# Compare token boundary and entropy between GPT-2 and T5
$ ted compare --input examples/tech_stack_descriptions.txt --tokenizers=gpt2,t5 --visualize

# Run heatmap and save visualization
$ ted heatmap examples/source_code_snippets.txt --output=output/examples/code_entropy.svg

# Run test comparison against golden output
$ ted test examples/english_quotes.txt --compare-to=output/gold/english_quotes.gpt2.csv
```

---

## 🔍 Feature Coverage by Example

| Example File                  | Token Drift | Entropy | Compression | Token Boundaries | CSV/JSON Parsing |
| ----------------------------- | ----------- | ------- | ----------- | ---------------- | ---------------- |
| `english_quotes.txt`          | ✅           | ✅       | ✅           | ✅                | ❌                |
| `tech_stack_descriptions.txt` | ✅           | ✅       | ✅           | ✅                | ❌                |
| `source_code_snippets.txt`    | ✅           | ✅       | ✅           | ✅                | ❌                |
| `structured_news.csv`         | ✅           | ✅       | ✅           | ✅                | ✅                |
| `reddit_threads.jsonl`        | ✅           | ✅       | ✅           | ✅                | ✅                |

---

These corpora serve both as validation for testing pipelines and demonstrations of the platform’s core functionality.
