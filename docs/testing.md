# Testing Strategy and Validation: TokEntropyDrift

This document outlines the testing methodology for verifying the correctness, stability, and consistency of TokEntropyDrift’s tokenization analysis and output rendering.

---

## 🧪 Overview

Testing in TokEntropyDrift focuses on:

* Ensuring tokenizer adapters produce consistent and valid outputs
* Verifying metric calculations (entropy, reuse, compression, etc.)
* Catching regressions in visualization and data export
* Validating parsing for `.txt`, `.csv`, and `.jsonl` inputs
* Comparing current outputs to known-good “golden files”

---

## 🧬 Corpus-Level Testing

Test corpora are located in `/examples/` and include:

* `english_quotes.txt`
* `tech_stack_descriptions.txt`
* `structured_news.csv`
* `reddit_threads.jsonl`

Each has associated expected outputs:

```
/testdata/
  english_quotes.gpt2.csv
  tech_stack_descriptions.t5.entropy.csv
  *.tokens.json
```

---

## 🔄 Regression Testing

```bash
$ ted test examples/english_quotes.txt \
    --tokenizers=gpt2 \
    --compare-to=testdata/english_quotes.gpt2.csv
```

* Automatically checks line count, token count, and metric diffs
* Supports `--tolerance` flag for floating-point entropy variations

---

## ✅ Validation Checks

### Tokenizer-Level

* Token count >= 1
* All token start offsets must be increasing
* Token ID presence (if applicable)
* No overlapping token spans

### Metric-Level

* Entropy must be ≥ 0
* Token reuse rate ∈ \[0, 1]
* Compression ratio must be ≥ 1 (in chars/token)

---

## 📤 Visualization Output Testing

* Confirm image files are generated (SVG/PNG)
* Validate graph dimension and presence of expected axes/legends
* Optionally use image checksum validation (planned)

---

## 🧪 Test Utilities

Planned CLI test tools:

```bash
$ ted test run-all
$ ted test visualize-only examples/tech_stack_descriptions.txt
$ ted validate-tokenizer t5
```

---

## 🧩 Future Testing Extensions

* Parallel testing on multiple corpora
* Golden output auto-generation + comparison diff
* Token alignment overlay validation
* Multi-tokenizer drift stress tests
* HTML dashboard snapshot testing

---

TokEntropyDrift prioritizes reproducibility and correctness by validating every major stage of its pipeline—from tokenizer output to entropy export and graphical rendering.
