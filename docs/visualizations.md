# Visualization and Rendering Specification

This document describes the visual output formats, rendering options, and planned visualizations for **TokEntropyDrift ("ted")**.

---

## 🎯 Goals

* Provide **interactive and exportable** visualizations of tokenizer behavior
* Support both **live dashboards** and **headless output** for CI/reports
* Align visual analysis with entropy, token drift, and token boundary exploration

---

## 🖥️ Output Targets

* **Web Dashboard** (HTML + Plotly.js or D3.js)
* **Static Export** (PNG, SVG, PDF)
* **Markdown & LaTeX snippets** (for research papers)
* **Headless CLI mode** for report automation

---

## 📦 Rendering Technology Stack

* **Plotly.js**: Primary visualization engine (clean, interactive, supports export)
* **Go/CLI Backend**: Generates JSON/CSV for web front-end or Python plots
* **Matplotlib (fallback)**: For headless plot generation in CI/testing contexts
* **Optional D3.js modules**: For advanced boundary overlays or animated sequences

---

## 📊 Visualization Types

### 1. **Token Count Heatmap**

* **Axes**: Line/sample vs Tokenizer
* **Color**: Token count per line
* Export: Heatmap (PNG/SVG) + CSV raw values

### 2. **Average Token Length Heatmap**

* **Axes**: Line vs Tokenizer
* Color indicates average character length per token

### 3. **Entropy Heatmap**

* Rolling entropy window per line/tokenizer
* Supports global and local entropy views

### 4. **Compression Ratio Map**

* Visualize compression (char count / token count)
* Highlights dense vs sparse tokenization

---

## 🧱 Token Boundary Visualizations

### 5. **Token Boundary Map**

* Horizontal bars showing token splits by character offset
* Supports multiple tokenizers aligned vertically
* Hover: token text, ID, offset

### 6. **Alignment Overlay**

* Compare tokenization boundaries for multiple tokenizers on the same input
* Highlight differences or mismatches in boundary segmentation
* Supports interactive diff mode

---

## 🔁 Token Drift Visualizations

### 7. **Token Count Delta Line Chart**

* Shows difference in token count between tokenizers per input line
* Useful for tracking tokenization inflation or compression

### 8. **Entropy Drift Line Chart**

* Shows per-line entropy delta between tokenizers
* Can be normalized by byte, char, or line length

---

## ⌛ Token Timeline Visualizations

### 9. **Token Stream Plot**

* Plot each token across sequence for one or more tokenizers
* Token length vs position, color by token ID or rank
* Exportable for per-sample views

### 10. **Animated Token Timeline** *(optional)*

* Shows token progression over time
* Can visualize how slight text variations alter token boundaries over frames

---

## 🧾 Export Formats

* **Static**: PNG, SVG, PDF
* **Data**: CSV, JSON
* **Rich Reports**:

  * Markdown tables with embedded sparkline images
  * LaTeX + TikZ snippets for academic papers

---

## 🔄 CLI/Headless Support

* All visualizations can be generated in non-interactive mode:

```bash
$ ted heatmap corpus.txt --output=out.png --mode=headless
$ ted compare entropy --visualize --format=svg
```

---

## 📁 File Structure

```
/output
  /heatmaps
    entropy_gpt2_t5.svg
  /boundaries
    sentence_45_comparison.png
  /charts
    token_drift.csv
    token_count_delta.svg
  /reports
    corpus_analysis.md
```

---

Future visualizations may include:

* Zipf plots for vocab usage
* Vocabulary overlap Venn diagrams
* Attention overlays (if integrated with model output)

See `analysis.md` for the statistical backing and `tokenizers.md` for tokenizer output structure.
