package visualization

import (
	"fmt"
)

// createTokenCountDriftPlot creates a plot showing token count differences
func (v *VisualizationEngine) createTokenCountDriftPlot(data DriftData) map[string]interface{} {
	if len(data.Documents) == 0 || len(data.DriftMetrics["token_count_delta"]) == 0 {
		return nil
	}

	// Create line plot data
	plotData := map[string]interface{}{
		"type": "scatter",
		"mode": "lines+markers",
		"x":    data.Documents,
		"y":    data.DriftMetrics["token_count_delta"],
		"name": "Token Count Delta",
		"line": map[string]interface{}{
			"color": "#1f77b4",
			"width": 2,
		},
		"marker": map[string]interface{}{
			"size":  6,
			"color": "#1f77b4",
		},
	}

	// Add subplot layout
	plotData["xaxis"] = "x"
	plotData["yaxis"] = "y"
	plotData["subplot"] = "xy"

	return plotData
}

// createEntropyDriftPlot creates a plot showing entropy differences
func (v *VisualizationEngine) createEntropyDriftPlot(data DriftData) map[string]interface{} {
	if len(data.Documents) == 0 || len(data.DriftMetrics["entropy_delta"]) == 0 {
		return nil
	}

	// Create line plot data
	plotData := map[string]interface{}{
		"type": "scatter",
		"mode": "lines+markers",
		"x":    data.Documents,
		"y":    data.DriftMetrics["entropy_delta"],
		"name": "Entropy Delta",
		"line": map[string]interface{}{
			"color": "#ff7f0e",
			"width": 2,
		},
		"marker": map[string]interface{}{
			"size":  6,
			"color": "#ff7f0e",
		},
	}

	// Add subplot layout
	plotData["xaxis"] = "x2"
	plotData["yaxis"] = "y2"
	plotData["subplot"] = "xy"

	return plotData
}

// createAlignmentPlot creates a plot showing alignment scores
func (v *VisualizationEngine) createAlignmentPlot(data DriftData) map[string]interface{} {
	if len(data.Documents) == 0 || len(data.DriftMetrics["alignment_score"]) == 0 {
		return nil
	}

	// Create bar plot data
	plotData := map[string]interface{}{
		"type": "bar",
		"x":    data.Documents,
		"y":    data.DriftMetrics["alignment_score"],
		"name": "Alignment Score",
		"marker": map[string]interface{}{
			"color": "#2ca02c",
		},
	}

	// Add subplot layout
	plotData["xaxis"] = "x3"
	plotData["yaxis"] = "y3"
	plotData["subplot"] = "xy"

	return plotData
}

// createTokenBoundaryPlotData creates data for token boundary visualization
func (v *VisualizationEngine) createTokenBoundaryPlotData(data TokenBoundaryData) []map[string]interface{} {
	var plotData []map[string]interface{}

	for i, tokenization := range data.Tokenizations {
		// Create horizontal bar for each tokenizer
		xPositions := make([]float64, 0)
		yPositions := make([]float64, 0)
		colors := make([]string, 0)
		text := make([]string, 0)

		currentPos := 0.0
		for _, token := range tokenization.Tokens {
			// Add start position
			xPositions = append(xPositions, currentPos)
			yPositions = append(yPositions, float64(i))
			colors = append(colors, "#1f77b4")
			text = append(text, fmt.Sprintf("Start: %s", token.Text))

			// Add end position
			currentPos += float64(len(token.Text))
			xPositions = append(xPositions, currentPos)
			yPositions = append(yPositions, float64(i))
			colors = append(colors, "#ff7f0e")
			text = append(text, fmt.Sprintf("End: %s", token.Text))
		}

		// Create scatter plot for this tokenizer
		plot := map[string]interface{}{
			"type": "scatter",
			"mode": "markers",
			"x":    xPositions,
			"y":    yPositions,
			"text": text,
			"name": data.TokenizerNames[i],
			"marker": map[string]interface{}{
				"size":   8,
				"color":  colors,
				"symbol": "circle",
			},
			"hovertemplate": "<b>%{text}</b><br>Position: %{x}<br>Tokenizer: %{fullData.name}<extra></extra>",
		}

		plotData = append(plotData, plot)
	}

	return plotData
}

// createRollingEntropyPlotData creates data for rolling entropy visualization
func (v *VisualizationEngine) createRollingEntropyPlotData(data RollingEntropyData) []map[string]interface{} {
	// Create x-axis positions
	xPositions := make([]int, len(data.EntropyValues))
	for i := range xPositions {
		xPositions[i] = i
	}

	// Create line plot
	plotData := []map[string]interface{}{
		{
			"type": "scatter",
			"mode": "lines+markers",
			"x":    xPositions,
			"y":    data.EntropyValues,
			"name": fmt.Sprintf("%s (window=%d)", data.TokenizerName, data.WindowSize),
			"line": map[string]interface{}{
				"color": "#1f77b4",
				"width": 2,
			},
			"marker": map[string]interface{}{
				"size":  4,
				"color": "#1f77b4",
			},
			"hovertemplate": "<b>%{fullData.name}</b><br>Window: %{x}<br>Entropy: %{y:.4f}<extra></extra>",
		},
	}

	return plotData
}
