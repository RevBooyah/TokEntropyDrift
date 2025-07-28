package visualization

import (
	"fmt"

	"github.com/RevBooyah/TokEntropyDrift/internal/metrics"
)

// generateTokenCountHeatmap generates a heatmap showing token counts
func (v *VisualizationEngine) generateTokenCountHeatmap(data HeatmapData) (*VisualizationResult, error) {
	// Create Plotly.js heatmap
	plotData := map[string]interface{}{
		"type":       "heatmap",
		"x":          data.XLabels,
		"y":          data.YLabels,
		"z":          data.Values,
		"colorscale": "Viridis",
		"colorbar": map[string]interface{}{
			"title": "Token Count",
		},
	}

	layout := map[string]interface{}{
		"title": map[string]interface{}{
			"text": "Token Count Heatmap",
			"x":    0.5,
		},
		"xaxis": map[string]interface{}{
			"title":     "Document",
			"tickangle": -45,
		},
		"yaxis": map[string]interface{}{
			"title": "Tokenizer",
		},
		"height":   v.getHeight(),
		"width":    v.getWidth(),
		"template": v.getTemplate(),
	}

	// Generate HTML
	html := v.generatePlotlyHTML([]map[string]interface{}{plotData}, layout, "token_count_heatmap")

	// Save to file
	filename := fmt.Sprintf("token_count_heatmap.%s", v.config.FileType)
	filepath := fmt.Sprintf("%s/%s", v.config.OutputDir, filename)

	if err := v.saveHTML(filepath, html); err != nil {
		return nil, err
	}

	return &VisualizationResult{
		Type:     "token_count_heatmap",
		Filepath: filepath,
		Data:     plotData,
		Metadata: map[string]interface{}{
			"x_labels_count": len(data.XLabels),
			"y_labels_count": len(data.YLabels),
			"min_value":      v.getMinValue(data.Values),
			"max_value":      v.getMaxValue(data.Values),
		},
	}, nil
}

// generateEntropyHeatmap generates a heatmap showing entropy values
func (v *VisualizationEngine) generateEntropyHeatmap(data HeatmapData) (*VisualizationResult, error) {
	// Create Plotly.js heatmap
	plotData := map[string]interface{}{
		"type":       "heatmap",
		"x":          data.XLabels,
		"y":          data.YLabels,
		"z":          data.Values,
		"colorscale": "Plasma",
		"colorbar": map[string]interface{}{
			"title": "Entropy",
		},
	}

	layout := map[string]interface{}{
		"title": map[string]interface{}{
			"text": "Entropy Heatmap",
			"x":    0.5,
		},
		"xaxis": map[string]interface{}{
			"title":     "Document",
			"tickangle": -45,
		},
		"yaxis": map[string]interface{}{
			"title": "Tokenizer",
		},
		"height":   v.getHeight(),
		"width":    v.getWidth(),
		"template": v.getTemplate(),
	}

	// Generate HTML
	html := v.generatePlotlyHTML([]map[string]interface{}{plotData}, layout, "entropy_heatmap")

	// Save to file
	filename := fmt.Sprintf("entropy_heatmap.%s", v.config.FileType)
	filepath := fmt.Sprintf("%s/%s", v.config.OutputDir, filename)

	if err := v.saveHTML(filepath, html); err != nil {
		return nil, err
	}

	return &VisualizationResult{
		Type:     "entropy_heatmap",
		Filepath: filepath,
		Data:     plotData,
		Metadata: map[string]interface{}{
			"x_labels_count": len(data.XLabels),
			"y_labels_count": len(data.YLabels),
			"min_value":      v.getMinValue(data.Values),
			"max_value":      v.getMaxValue(data.Values),
		},
	}, nil
}

// generateCompressionHeatmap generates a heatmap showing compression ratios
func (v *VisualizationEngine) generateCompressionHeatmap(data HeatmapData) (*VisualizationResult, error) {
	// Create Plotly.js heatmap
	plotData := map[string]interface{}{
		"type":       "heatmap",
		"x":          data.XLabels,
		"y":          data.YLabels,
		"z":          data.Values,
		"colorscale": "RdYlBu_r", // Red for high compression, blue for low
		"colorbar": map[string]interface{}{
			"title": "Compression Ratio",
		},
	}

	layout := map[string]interface{}{
		"title": map[string]interface{}{
			"text": "Compression Ratio Heatmap",
			"x":    0.5,
		},
		"xaxis": map[string]interface{}{
			"title":     "Document",
			"tickangle": -45,
		},
		"yaxis": map[string]interface{}{
			"title": "Tokenizer",
		},
		"height":   v.getHeight(),
		"width":    v.getWidth(),
		"template": v.getTemplate(),
	}

	// Generate HTML
	html := v.generatePlotlyHTML([]map[string]interface{}{plotData}, layout, "compression_heatmap")

	// Save to file
	filename := fmt.Sprintf("compression_heatmap.%s", v.config.FileType)
	filepath := fmt.Sprintf("%s/%s", v.config.OutputDir, filename)

	if err := v.saveHTML(filepath, html); err != nil {
		return nil, err
	}

	return &VisualizationResult{
		Type:     "compression_heatmap",
		Filepath: filepath,
		Data:     plotData,
		Metadata: map[string]interface{}{
			"x_labels_count": len(data.XLabels),
			"y_labels_count": len(data.YLabels),
			"min_value":      v.getMinValue(data.Values),
			"max_value":      v.getMaxValue(data.Values),
		},
	}, nil
}

// generateReuseHeatmap generates a heatmap showing token reuse rates
func (v *VisualizationEngine) generateReuseHeatmap(data HeatmapData) (*VisualizationResult, error) {
	// Create Plotly.js heatmap
	plotData := map[string]interface{}{
		"type":       "heatmap",
		"x":          data.XLabels,
		"y":          data.YLabels,
		"z":          data.Values,
		"colorscale": "Greens", // Green for high reuse
		"colorbar": map[string]interface{}{
			"title": "Reuse Rate",
		},
	}

	layout := map[string]interface{}{
		"title": map[string]interface{}{
			"text": "Token Reuse Heatmap",
			"x":    0.5,
		},
		"xaxis": map[string]interface{}{
			"title":     "Document",
			"tickangle": -45,
		},
		"yaxis": map[string]interface{}{
			"title": "Tokenizer",
		},
		"height":   v.getHeight(),
		"width":    v.getWidth(),
		"template": v.getTemplate(),
	}

	// Generate HTML
	html := v.generatePlotlyHTML([]map[string]interface{}{plotData}, layout, "reuse_heatmap")

	// Save to file
	filename := fmt.Sprintf("reuse_heatmap.%s", v.config.FileType)
	filepath := fmt.Sprintf("%s/%s", v.config.OutputDir, filename)

	if err := v.saveHTML(filepath, html); err != nil {
		return nil, err
	}

	return &VisualizationResult{
		Type:     "reuse_heatmap",
		Filepath: filepath,
		Data:     plotData,
		Metadata: map[string]interface{}{
			"x_labels_count": len(data.XLabels),
			"y_labels_count": len(data.YLabels),
			"min_value":      v.getMinValue(data.Values),
			"max_value":      v.getMaxValue(data.Values),
		},
	}, nil
}

// prepareHeatmapData prepares data for heatmap generation from analysis results
func (v *VisualizationEngine) prepareHeatmapData(analysisResults []*metrics.AnalysisResult, metricType string) *HeatmapData {
	if len(analysisResults) == 0 {
		return nil
	}

	// Group results by tokenizer and document
	tokenizerMap := make(map[string]map[string]float64)
	documents := make(map[string]bool)

	for _, result := range analysisResults {
		if _, exists := tokenizerMap[result.TokenizerName]; !exists {
			tokenizerMap[result.TokenizerName] = make(map[string]float64)
		}

		// Extract metric value
		var value float64
		switch metricType {
		case "token_count":
			value = float64(result.TokenCount)
		case "entropy":
			if metric, exists := result.Metrics["entropy_global_entropy"]; exists {
				value = metric.Value
			}
		case "compression":
			if metric, exists := result.Metrics["compression_compression_ratio"]; exists {
				value = metric.Value
			}
		case "reuse":
			if metric, exists := result.Metrics["reuse_reuse_ratio"]; exists {
				value = metric.Value
			}
		}

		tokenizerMap[result.TokenizerName][result.Document] = value
		documents[result.Document] = true
	}

	// Create ordered lists
	tokenizers := make([]string, 0, len(tokenizerMap))
	for tokenizer := range tokenizerMap {
		tokenizers = append(tokenizers, tokenizer)
	}

	docList := make([]string, 0, len(documents))
	for doc := range documents {
		docList = append(docList, doc)
	}

	// Create values matrix
	values := make([][]float64, len(tokenizers))
	for i, tokenizer := range tokenizers {
		values[i] = make([]float64, len(docList))
		for j, doc := range docList {
			if val, exists := tokenizerMap[tokenizer][doc]; exists {
				values[i][j] = val
			} else {
				values[i][j] = 0.0
			}
		}
	}

	return &HeatmapData{
		XLabels: docList,
		YLabels: tokenizers,
		Values:  values,
		Title:   fmt.Sprintf("%s Heatmap", metricType),
	}
}

// Helper functions
func (v *VisualizationEngine) getMinValue(values [][]float64) float64 {
	if len(values) == 0 || len(values[0]) == 0 {
		return 0.0
	}

	min := values[0][0]
	for _, row := range values {
		for _, val := range row {
			if val < min {
				min = val
			}
		}
	}
	return min
}

func (v *VisualizationEngine) getMaxValue(values [][]float64) float64 {
	if len(values) == 0 || len(values[0]) == 0 {
		return 0.0
	}

	max := values[0][0]
	for _, row := range values {
		for _, val := range row {
			if val > max {
				max = val
			}
		}
	}
	return max
}
