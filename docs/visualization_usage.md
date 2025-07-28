# Visualization Engine Usage Guide

This document provides comprehensive guidance on using the TokEntropyDrift visualization engine to create interactive charts, heatmaps, and reports for tokenization analysis.

---

## 🎯 Overview

The visualization engine provides multiple ways to visualize tokenization analysis results:

- **Heatmaps**: Compare metrics across tokenizers and documents
- **Token Boundary Visualizations**: See how different tokenizers segment text
- **Drift Analysis**: Compare tokenization behavior between models
- **Rolling Entropy Plots**: Analyze entropy patterns over text windows
- **Comprehensive Reports**: Multi-page HTML reports with all visualizations

---

## 🚀 Quick Start

### Basic Usage

```go
package main

import (
    "github.com/RevBooyah/TokEntropyDrift/internal/visualization"
"github.com/RevBooyah/TokEntropyDrift/internal/metrics"
)

func main() {
    // Initialize visualization engine
    vizEngine := visualization.NewVisualizationEngine(visualization.VisualizationConfig{
        Theme:       "light",
        ImageSize:   "medium",
        FileType:    "html",
        Interactive: true,
        OutputDir:   "output",
    })

    // Generate heatmap from analysis results
    heatmapData := vizEngine.prepareHeatmapData(analysisResults, "token_count")
    result, err := vizEngine.GenerateHeatmap(*heatmapData, "token_count")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Generated heatmap: %s\n", result.Filepath)
}
```

---

## 📊 Visualization Types

### 1. Heatmaps

Heatmaps are perfect for comparing metrics across multiple tokenizers and documents.

#### Available Heatmap Types

- **`token_count`**: Number of tokens per document/tokenizer
- **`entropy`**: Entropy values per document/tokenizer  
- **`compression`**: Compression ratios per document/tokenizer
- **`reuse`**: Token reuse rates per document/tokenizer

#### Example: Token Count Heatmap

```go
// Prepare data from analysis results
heatmapData := vizEngine.prepareHeatmapData(analysisResults, "token_count")

// Generate heatmap
result, err := vizEngine.GenerateHeatmap(*heatmapData, "token_count")
if err != nil {
    log.Fatal(err)
}

// Result contains filepath and metadata
fmt.Printf("Heatmap saved to: %s\n", result.Filepath)
```

#### Example: Entropy Heatmap

```go
// Generate entropy heatmap
entropyData := vizEngine.prepareHeatmapData(analysisResults, "entropy")
result, err := vizEngine.GenerateHeatmap(*entropyData, "entropy")
```

### 2. Token Boundary Visualizations

Visualize how different tokenizers segment the same text.

```go
boundaryData := visualization.TokenBoundaryData{
    DocumentID:     "sample_doc",
    Document:       "The quick brown fox jumps over the lazy dog.",
    TokenizerNames: []string{"gpt2", "bert", "t5"},
    Tokenizations:  []*tokenizers.TokenizationResult{gpt2Result, bertResult, t5Result},
}

result, err := vizEngine.GenerateTokenBoundaryMap(boundaryData)
if err != nil {
    log.Fatal(err)
}
```

**Features:**
- Horizontal bars showing token boundaries
- Color-coded start/end positions
- Hover tooltips with token details
- Multiple tokenizer comparison

### 3. Drift Analysis Visualizations

Compare tokenization behavior between different models.

```go
driftData := visualization.DriftData{
    ComparisonID: "gpt2_vs_bert",
    Tokenizer1:   "gpt2",
    Tokenizer2:   "bert",
    Documents:    []string{"doc1", "doc2", "doc3"},
    DriftMetrics: map[string][]float64{
        "token_count_delta": {5.2, -3.1, 8.7},
        "entropy_delta":     {0.15, -0.08, 0.22},
        "alignment_score":   {0.85, 0.92, 0.78},
    },
}

result, err := vizEngine.GenerateDriftVisualization(driftData)
if err != nil {
    log.Fatal(err)
}
```

**Generated Plots:**
- Token count delta line chart
- Entropy drift line chart  
- Alignment score bar chart

### 4. Rolling Entropy Plots

Analyze entropy patterns over sliding windows.

```go
rollingData := visualization.RollingEntropyData{
    DocumentID:    "sample_doc",
    TokenizerName: "gpt2",
    WindowSize:    100,
    EntropyValues: []float64{2.1, 2.3, 1.9, 2.5, 2.0}, // ... more values
}

result, err := vizEngine.GenerateRollingEntropyPlot(rollingData)
if err != nil {
    log.Fatal(err)
}
```

### 5. Comprehensive Reports

Generate multi-page HTML reports with all visualizations.

```go
// Generate comprehensive report from analysis results
report, err := vizEngine.GenerateComprehensiveReport(analysisResults)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Report generated: %s\n", report.Filepath)
fmt.Printf("Visualizations included: %d\n", report.Metadata["visualization_count"])
```

**Report Features:**
- Navigation menu for different visualizations
- Summary page with analysis overview
- Interactive iframe-based visualization display
- Export capabilities for each chart

---

## ⚙️ Configuration Options

### VisualizationConfig

```go
type VisualizationConfig struct {
    Theme         string // "light" or "dark"
    ImageSize     string // "small", "medium", or "large"
    FileType      string // "html", "svg", or "png"
    Interactive   bool   // Enable interactive features
    OutputDir     string // Directory for output files
}
```

### Theme Options

- **`light`**: Clean white background with dark text
- **`dark`**: Dark background with light text

### Image Size Options

- **`small`**: 400x600 pixels
- **`medium`**: 600x900 pixels (default)
- **`large`**: 800x1200 pixels

### File Type Options

- **`html`**: Interactive Plotly.js visualizations
- **`svg`**: Static SVG images
- **`png`**: Static PNG images

---

## 🎨 Customization

### Custom Color Schemes

```go
// Custom heatmap colorscale
plotData := map[string]interface{}{
    "type": "heatmap",
    "x":    data.XLabels,
    "y":    data.YLabels,
    "z":    data.Values,
    "colorscale": "Viridis", // Options: Viridis, Plasma, RdYlBu, etc.
}
```

### Custom Layouts

```go
layout := map[string]interface{}{
    "title": map[string]interface{}{
        "text": "Custom Title",
        "x":    0.5,
        "font": map[string]interface{}{
            "size": 20,
            "color": "#333",
        },
    },
    "xaxis": map[string]interface{}{
        "title": "Custom X Axis",
        "tickangle": -45,
    },
    "yaxis": map[string]interface{}{
        "title": "Custom Y Axis",
    },
    "height": 800,
    "width":  1200,
    "template": "plotly_white",
}
```

---

## 📁 Output File Structure

```
output/
├── comprehensive_report.html          # Multi-page report
├── token_count_heatmap.html          # Token count heatmap
├── entropy_heatmap.html              # Entropy heatmap
├── compression_heatmap.html          # Compression heatmap
├── reuse_heatmap.html                # Reuse rate heatmap
├── token_boundary_test_doc.html      # Token boundary visualization
├── drift_analysis_gpt2_vs_bert.html  # Drift analysis
└── rolling_entropy_sample_doc.html   # Rolling entropy plot
```

---

## 🔧 Integration with CLI

### Using with `ted analyze` command

```bash
# Analyze with visualization output
ted analyze examples/english_quotes.txt --output output --visualize

# Generate specific heatmap
ted heatmap examples/english_quotes.txt --type entropy --output entropy_heatmap.html
```

### Configuration File

```yaml
# ted.config.yaml
visualization:
  theme: "light"
  image_size: "medium"
  file_type: "html"
  interactive: true
  output_dir: "output"
```

---

## 🐛 Troubleshooting

### Common Issues

1. **Empty visualizations**: Ensure analysis results contain the expected metrics
2. **Missing data**: Check that tokenizers are properly registered and initialized
3. **File permission errors**: Ensure output directory is writable
4. **Large file sizes**: Use smaller datasets or reduce image quality

### Debug Mode

```go
// Enable debug logging
vizEngine := visualization.NewVisualizationEngine(visualization.VisualizationConfig{
    // ... config
    Debug: true, // Add debug field to config
})
```

---

## 📈 Performance Tips

1. **Batch Processing**: Generate multiple visualizations in one call
2. **Caching**: Reuse prepared data for multiple visualization types
3. **Optimization**: Use appropriate image sizes for your use case
4. **Memory Management**: Close visualization engine when done

```go
// Efficient batch processing
results := []*visualization.VisualizationResult{}

// Generate all heatmaps at once
for _, vizType := range []string{"token_count", "entropy", "compression", "reuse"} {
    data := vizEngine.prepareHeatmapData(analysisResults, vizType)
    if data != nil {
        result, err := vizEngine.GenerateHeatmap(*data, vizType)
        if err == nil {
            results = append(results, result)
        }
    }
}
```

---

## 🔗 Advanced Usage

### Custom Plotly.js Integration

```go
// Create custom Plotly.js visualization
customData := []map[string]interface{}{
    {
        "type": "scatter",
        "mode": "lines+markers",
        "x":    []float64{1, 2, 3, 4, 5},
        "y":    []float64{2, 4, 1, 3, 5},
        "name": "Custom Data",
    },
}

customLayout := map[string]interface{}{
    "title": "Custom Visualization",
    "xaxis": map[string]interface{}{"title": "X Axis"},
    "yaxis": map[string]interface{}{"title": "Y Axis"},
}

html := vizEngine.generatePlotlyHTML(customData, customLayout, "custom_viz")
```

### Integration with Web Server

```go
// Serve visualizations via HTTP
http.HandleFunc("/visualization", func(w http.ResponseWriter, r *http.Request) {
    // Generate visualization on-demand
    result, err := vizEngine.GenerateHeatmap(heatmapData, "token_count")
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    
    // Serve the HTML file
    http.ServeFile(w, r, result.Filepath)
})
```

---

## 📚 Examples

See the `examples/` directory for complete working examples:

- `examples/visualization_basic.go` - Basic heatmap generation
- `examples/visualization_advanced.go` - Advanced multi-plot reports
- `examples/visualization_cli.go` - CLI integration examples

---

## 🤝 Contributing

When adding new visualization types:

1. Add new methods to `VisualizationEngine`
2. Create corresponding data structures
3. Update this documentation
4. Add tests for new functionality
5. Update the CLI integration

For questions or issues, please refer to the main project documentation or create an issue on GitHub. 