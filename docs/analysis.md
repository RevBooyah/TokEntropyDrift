# Analysis: Tokenization Metrics & Statistical Processing

This document describes in detail the analysis methods, metrics, and output formats used in **TokEntropyDrift ("ted")** to measure and compare tokenizer behavior.

---

## 🧾 Input Format Support

TokEntropyDrift supports the following input types:

* **Plain Text (.txt):** One or more lines of arbitrary text.
* **CSV (.csv):** One column must be designated as the text source. Others (e.g. source, ID) can be tracked.
* **JSON Lines (.jsonl):** One JSON object per line. Text field can be specified.

Each input line is treated as a **distinct sample** with tracking of:

* File name (source)
* Line number or row index
* Original text string

---

## 📏 Token Count Metrics

* **Total Token Count**: Number of tokens per sample and per tokenizer.
* **Average Token Length**: Total characters / total tokens.
* **Tokens per Byte / per Character**: Compression ratio metric.
* **Token Frequency Histogram**: Global token distribution.
* **Token Coverage**: % of vocab used over a corpus (per tokenizer).

---

## 🧮 Entropy Calculation

### 1. **Global Entropy**

* **Shannon Entropy** over full corpus:

  * Based on token frequencies.
  * Entropy per tokenizer (bits/token and normalized bits/char).

### 2. **Rolling / Local Entropy**

* Sliding window entropy (e.g., 50-token window).
* Token stream broken into windows; entropy computed for each.
* Can reveal local regularity or variability patterns.

### 3. **Bigram/Pairwise Entropy**

* Conditional entropy of token pairs.
* Highlights language model compression opportunities.

---

## 🔁 Token Reuse & Compression Metrics

* **Token Reuse Rate**:

  * Unique tokens / total tokens.
  * Indicates how “bursty” or repetitive tokenization is.

* **Compression Ratio**:

  * Raw text bytes vs. number of tokens.
  * Character/token ratio (higher = more efficient tokenization).

* **Redundancy Factor**: Entropy vs theoretical max.

---

## 📉 Cross-Tokenizer Comparison

For a given sample (line, file):

* Token count delta
* Entropy delta (absolute & %)
* Token alignment visualization
* Entropy-normalized compression metric

For full corpus:

* Side-by-side histograms and bar plots
* Ranking of tokenizers by:

  * Entropy
  * Compression
  * Token count

---

## 🧪 Output Formats

* **Raw Stats** (per sample, per tokenizer): JSON, CSV
* **Aggregate Stats** (corpus-level summaries): CSV, Markdown, LaTeX
* **Visual Assets** (for publication):

  * Entropy graphs (PNG/SVG)
  * Heatmaps (tokens/sentence, entropy/sample)
  * Token boundary maps

All exports include line-level metadata: file, line number, tokenizer ID.

---

## ⚙️ Batch Processing & Scalability

* Default max input size: \~5MB (configurable)
* Streaming line-by-line parser for large corpora
* Tokenizer outputs cached per tokenizer per line (optional)
* Outputs written to disk as chunked CSV/JSON
* Future: Indexed database output for high-volume data (e.g., SQLite)

---

## 🛠 Configuration & CLI Flags (planned)

* `--window-size`: sliding entropy window (default: 50)
* `--min-length`: skip inputs shorter than N tokens
* `--format`: `csv`, `json`, `md`, `latex`
* `--compare`: run pairwise tokenizer comparison mode
* `--normalize`: normalize entropy by char count or byte size
* `--bigram`: include bigram entropy calculation

---

## 🧠 Future Extensions

* Cost estimation mode (OpenAI / Anthropic / custom pricing)
* Per-language entropy baselining
* Zipfian distribution check for vocab usage
* Visualization toolkit hooks
* LLM-friendly Markdown summary report per corpus

---

This file defines the analytical backbone of TokEntropyDrift. See `tokenizers.md` for tokenizer compatibility and `visualizations.md` for rendering the output statistics.
