# Tokenizer Integration and Specification

This document outlines the design, implementation, and extensibility of tokenizer support within **TokEntropyDrift ("ted")**.

---

## 🧠 Tokenizer Architecture

All tokenizer modules conform to a **unified interface** and are pluggable into the TED analysis pipeline. Each tokenizer implementation must:

* Accept a raw text string as input
* Return a list of tokens (with optional metadata)
* Support batch and line-by-line tokenization
* Output normalized tokenization metadata for analysis and alignment

Supported via:

* Native bindings (Python-based or binary CLI interface)
* Go wrapper with optional Python bridge (via cgo or subprocess)
* Easily extensible with config/registry file for adding new tokenizers

---

## 🔌 Built-in Tokenizers (Initial Support)

| Tokenizer Name          | Type            | Backing Library              | Notes                    |
| ----------------------- | --------------- | ---------------------------- | ------------------------ |
| GPT-2 / GPT-3.5 / GPT-4 | BPE             | `tiktoken`                   | Fast, OpenAI-compatible  |
| HuggingFace BPE         | BPE             | `transformers`, `tokenizers` | RoBERTa, GPT-Neo, etc.   |
| SentencePiece           | Unigram/BPE     | `sentencepiece`              | T5, mT5, ALBERT          |
| WordPiece               | WordPiece       | `transformers`               | BERT, DistilBERT         |
| OpenAI API              | BPE             | REST API                     | Optional via config flag |
| Claude / PaLM           | Approximate BPE | Custom mappings              | TBD                      |
| Custom                  | Any             | Configured by user           | Via vocab/model files    |

---

## 🧾 Required Token Metadata

Each tokenizer must return the following per-token:

* **Token ID** (if available)
* **Text Segment** (string of token)
* **Start Offset** (character index in original input)
* **Byte Length** (of tokenized segment)
* **Tokenizer-specific Metadata** (optional flags, vocab info)

These values are stored in structured output and used to compute entropy, alignment, and reuse stats.

---

## ⚙️ Tokenizer Configuration and Extensibility

Each tokenizer is registered via a config entry:

```json
{
  "name": "gpt2",
  "type": "bpe",
  "backend": "python",
  "command": "python3 scripts/gpt2_tokenizer.py",
  "vocab_path": "vocab/gpt2-vocab.json",
  "model_path": "vocab/gpt2-merges.txt"
}
```

### Custom Tokenizer Support

* Users may drop `.model`, `.vocab`, `.json`, or other files into `tokenizers/`
* CLI command `ted tokenizer build` can train BPE or unigram models using a sample corpus
* Support for SentencePiece training and HuggingFace tokenizer training planned

---

## 🔃 Normalization & Preprocessing

Tokenizer adapters will:

* Apply native tokenizer preprocessing (e.g. space normalization, lowercasing)
* Optionally override defaults via config flags
* Support Unicode normalization (`NFC`, `NFKC`, etc)

Standardized preprocessing will be tracked per tokenizer and included in metadata output.

---

## 📊 Tokenizer-Level Metrics

Each tokenizer is analyzed across the following metrics:

* Vocabulary size
* Avg. token length
* Token reuse rate
* Entropy (global, local, normalized)
* UNK or OOV token frequency
* Token density (tokens per char, per byte)
* Compression ratio

---

## 🧩 Token Alignment and Mapping

When comparing tokenizers on the same input, TED computes:

* Token boundary overlaps
* Token count difference per line
* Sequence distance (e.g., Levenshtein, alignment score)
* Shared prefix/suffix token analysis

This allows:

* Side-by-side visual comparisons
* Entropy drift graphs
* Boundary heatmaps

---

## 🔄 Batch and Streaming Support

* Tokenizers must support:

  * Single string input
  * Batch mode (array of strings)
  * Streamed line-by-line file tokenization
* Caching optional for repeated runs

---

## 🔧 Planned CLI Options

* `ted tokenizer list` — list available tokenizers
* `ted tokenizer build` — train a new tokenizer from corpus
* `ted tokenizer test` — validate tokenizer output vs expected
* `ted tokenizer info <name>` — show vocab size, type, metadata

---

See `analysis.md` for downstream processing of tokenizer output and `visualizations.md` for how outputs are rendered.
