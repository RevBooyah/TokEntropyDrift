# Visualization Engine

This package provides comprehensive visualization capabilities for the TokEntropyDrift project, enabling interactive analysis of tokenization behavior across different models and datasets.

## 📁 Package Structure

```
internal/visualization/
├── engine.go          # Main visualization engine and core types
├── heatmaps.go        # Heatmap generation functions
├── plotly.go          # Plotly.js HTML generation
├── drift.go           # Drift analysis visualizations
└── README.md          # This file
```

## 🏗️ Architecture

### Core Components

1. **VisualizationEngine**: Main orchestrator for all visualization types
2. **Data Structures**: Type-safe structures for different visualization data
3. **HTML Generators**: Plotly.js integration for interactive charts
4. **Export Handlers**: Support for multiple output formats

### Design Principles

- **Modular**: Each visualization type is self-contained
- **Extensible**: Easy to add new visualization types
- **Interactive**: Rich Plotly.js-based visualizations
- **Exportable**: Support for static and interactive formats
- **Configurable**: Theme, size, and format options

## 🔧 Core Types

### VisualizationEngine

```go
type VisualizationEngine struct {
    config VisualizationConfig
}
```

Main engine that coordinates all visualization generation.

### VisualizationConfig

```go
type VisualizationConfig struct {
    Theme         string // "light" or "dark"
    ImageSize     string // "small", "medium", "large"
    FileType      string // "html", "svg", "png"
    Interactive   bool   // Enable interactive features
    OutputDir     string // Directory for output files
}
```

Configuration for visualization generation.

### Data Structures

- **HeatmapData**: For heatmap visualizations
- **TokenBoundaryData**: For token boundary analysis
- **DriftData**: For cross-tokenizer comparison
- **RollingEntropyData**: For entropy pattern analysis

## 📊 Visualization Types

### 1. Heatmaps

**Purpose**: Compare metrics across tokenizers and documents

**Types**:
- `token_count`: Number of tokens per document/tokenizer
- `entropy`: Entropy values per document/tokenizer
- `compression`: Compression ratios per document/tokenizer
- `reuse`: Token reuse rates per document/tokenizer

**Usage**:
```go
heatmapData := vizEngine.prepareHeatmapData(analysisResults, "token_count")
result, err := vizEngine.GenerateHeatmap(*heatmapData, "token_count")
```

### 2. Token Boundary Visualizations

**Purpose**: Visualize how different tokenizers segment text

**Features**:
- Horizontal bars showing token boundaries
- Color-coded start/end positions
- Hover tooltips with token details
- Multiple tokenizer comparison

**Usage**:
```go
boundaryData := visualization.TokenBoundaryData{
    DocumentID:     "sample_doc",
    Document:       "The quick brown fox...",
    TokenizerNames: []string{"gpt2", "bert", "t5"},
    Tokenizations:  []*tokenizers.TokenizationResult{...},
}
result, err := vizEngine.GenerateTokenBoundaryMap(boundaryData)
```

### 3. Drift Analysis

**Purpose**: Compare tokenization behavior between models

**Generated Plots**:
- Token count delta line chart
- Entropy drift line chart
- Alignment score bar chart

**Usage**:
```go
driftData := visualization.DriftData{
    ComparisonID: "gpt2_vs_bert",
    Tokenizer1:   "gpt2",
    Tokenizer2:   "bert",
    Documents:    []string{"doc1", "doc2", "doc3"},
    DriftMetrics: map[string][]float64{...},
}
result, err := vizEngine.GenerateDriftVisualization(driftData)
```

### 4. Rolling Entropy Plots

**Purpose**: Analyze entropy patterns over sliding windows

**Usage**:
```go
rollingData := visualization.RollingEntropyData{
    DocumentID:    "sample_doc",
    TokenizerName: "gpt2",
    WindowSize:    100,
    EntropyValues: []float64{2.1, 2.3, 1.9, 2.5, 2.0},
}
result, err := vizEngine.GenerateRollingEntropyPlot(rollingData)
```

### 5. Comprehensive Reports

**Purpose**: Multi-page HTML reports with all visualizations

**Features**:
- Navigation menu for different visualizations
- Summary page with analysis overview
- Interactive iframe-based visualization display
- Export capabilities for each chart

**Usage**:
```go
report, err := vizEngine.GenerateComprehensiveReport(analysisResults)
```

## 🎨 Customization

### Themes

- **Light**: Clean white background with dark text
- **Dark**: Dark background with light text

### Image Sizes

- **Small**: 400x600 pixels
- **Medium**: 600x900 pixels (default)
- **Large**: 800x1200 pixels

### File Types

- **HTML**: Interactive Plotly.js visualizations
- **SVG**: Static SVG images
- **PNG**: Static PNG images

### Custom Color Schemes

```go
plotData := map[string]interface{}{
    "type": "heatmap",
    "colorscale": "Viridis", // Options: Viridis, Plasma, RdYlBu, etc.
}
```

## 🔌 Extending the Engine

### Adding New Visualization Types

1. **Define Data Structure**:
```go
type NewVizData struct {
    // Your data fields
}
```

2. **Add Method to Engine**:
```go
func (v *VisualizationEngine) GenerateNewViz(data NewVizData) (*VisualizationResult, error) {
    // Implementation
}
```

3. **Create Plot Data**:
```go
func (v *VisualizationEngine) createNewVizPlotData(data NewVizData) []map[string]interface{} {
    // Return Plotly.js compatible data
}
```

4. **Update Documentation**: Add usage examples and descriptions

### Example: Adding a Scatter Plot

```go
// 1. Define data structure
type ScatterData struct {
    XValues []float64 `json:"x_values"`
    YValues []float64 `json:"y_values"`
    Labels  []string  `json:"labels"`
}

// 2. Add engine method
func (v *VisualizationEngine) GenerateScatterPlot(data ScatterData) (*VisualizationResult, error) {
    plotData := v.createScatterPlotData(data)
    
    layout := map[string]interface{}{
        "title": "Scatter Plot",
        "xaxis": map[string]interface{}{"title": "X Axis"},
        "yaxis": map[string]interface{}{"title": "Y Axis"},
    }
    
    html := v.generatePlotlyHTML(plotData, layout, "scatter_plot")
    
    filename := fmt.Sprintf("scatter_plot.%s", v.config.FileType)
    filepath := filepath.Join(v.config.OutputDir, filename)
    
    if err := v.saveHTML(filepath, html); err != nil {
        return nil, err
    }
    
    return &VisualizationResult{
        Type:     "scatter_plot",
        Filepath: filepath,
        Data:     plotData,
    }, nil
}

// 3. Create plot data
func (v *VisualizationEngine) createScatterPlotData(data ScatterData) []map[string]interface{} {
    return []map[string]interface{}{
        {
            "type": "scatter",
            "mode": "markers",
            "x":    data.XValues,
            "y":    data.YValues,
            "text": data.Labels,
            "marker": map[string]interface{}{
                "size":  8,
                "color": "#1f77b4",
            },
        },
    }
}
```

## 🧪 Testing

### Unit Tests

Create tests for each visualization type:

```go
func TestGenerateHeatmap(t *testing.T) {
    vizEngine := NewVisualizationEngine(VisualizationConfig{
        Theme:       "light",
        ImageSize:   "medium",
        FileType:    "html",
        Interactive: true,
        OutputDir:   "test_output",
    })
    
    heatmapData := HeatmapData{
        XLabels: []string{"doc1", "doc2"},
        YLabels: []string{"gpt2", "bert"},
        Values:  [][]float64{{10, 15}, {12, 18}},
    }
    
    result, err := vizEngine.GenerateHeatmap(heatmapData, "token_count")
    assert.NoError(t, err)
    assert.NotNil(t, result)
    assert.Equal(t, "token_count_heatmap", result.Type)
}
```

### Integration Tests

Test the full pipeline:

```go
func TestVisualizationPipeline(t *testing.T) {
    // Setup analysis results
    // Generate visualizations
    // Verify output files exist
    // Check HTML content
}
```

## 🐛 Debugging

### Common Issues

1. **Empty visualizations**: Check that analysis results contain expected metrics
2. **Missing data**: Verify tokenizers are properly registered
3. **File permission errors**: Ensure output directory is writable
4. **Large file sizes**: Use smaller datasets or reduce image quality

### Debug Mode

```go
vizEngine := NewVisualizationEngine(VisualizationConfig{
    // ... config
    Debug: true, // Add debug field to config
})
```

## 📈 Performance Considerations

1. **Batch Processing**: Generate multiple visualizations in one call
2. **Caching**: Reuse prepared data for multiple visualization types
3. **Memory Management**: Close visualization engine when done
4. **Optimization**: Use appropriate image sizes for your use case

## 🔗 Dependencies

- **Plotly.js**: For interactive visualizations
- **Go standard library**: For file operations and JSON handling
- **Internal packages**: metrics, tokenizers for data integration

## 📚 Related Documentation

- [Visualization Usage Guide](../docs/visualization_usage.md)
- [Main Project README](../../README.md)
- [Architecture Documentation](../../docs/architecture.md)

## 🤝 Contributing

When contributing to the visualization engine:

1. Follow the existing code style and patterns
2. Add comprehensive tests for new features
3. Update this README with new functionality
4. Ensure all visualizations are accessible and responsive
5. Test with different themes and configurations

For questions or issues, please refer to the main project documentation or create an issue on GitHub. 