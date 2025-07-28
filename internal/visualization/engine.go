package visualization

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/RevBooyah/TokEntropyDrift/internal/metrics"
	"github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
)

// VisualizationEngine handles generation of various visualizations
type VisualizationEngine struct {
	config VisualizationConfig
}

// VisualizationConfig holds configuration for visualization generation
type VisualizationConfig struct {
	Theme       string `json:"theme"`      // light, dark
	ImageSize   string `json:"image_size"` // small, medium, large
	FileType    string `json:"file_type"`  // svg, png, html
	Interactive bool   `json:"interactive"`
	OutputDir   string `json:"output_dir"`
}

// NewVisualizationEngine creates a new visualization engine
func NewVisualizationEngine(config VisualizationConfig) *VisualizationEngine {
	return &VisualizationEngine{
		config: config,
	}
}

// GenerateHeatmap generates a heatmap visualization
func (v *VisualizationEngine) GenerateHeatmap(data HeatmapData, vizType string) (*VisualizationResult, error) {
	switch vizType {
	case "token_count":
		return v.generateTokenCountHeatmap(data)
	case "entropy":
		return v.generateEntropyHeatmap(data)
	case "compression":
		return v.generateCompressionHeatmap(data)
	case "reuse":
		return v.generateReuseHeatmap(data)
	default:
		return nil, fmt.Errorf("unsupported heatmap type: %s", vizType)
	}
}

// GenerateTokenBoundaryMap generates a token boundary visualization
func (v *VisualizationEngine) GenerateTokenBoundaryMap(data TokenBoundaryData) (*VisualizationResult, error) {
	// Create Plotly.js visualization
	plotData := v.createTokenBoundaryPlotData(data)

	layout := map[string]interface{}{
		"title": map[string]interface{}{
			"text": "Token Boundary Analysis",
			"x":    0.5,
		},
		"xaxis": map[string]interface{}{
			"title":    "Character Position",
			"showgrid": true,
		},
		"yaxis": map[string]interface{}{
			"title":    "Tokenizer",
			"showgrid": true,
		},
		"height":   v.getHeight(),
		"width":    v.getWidth(),
		"template": v.getTemplate(),
	}

	// Generate HTML
	html := v.generatePlotlyHTML(plotData, layout, "token_boundary")

	// Save to file
	filename := fmt.Sprintf("token_boundary_%s.%s", data.DocumentID, v.config.FileType)
	filepath := filepath.Join(v.config.OutputDir, filename)

	if err := os.WriteFile(filepath, []byte(html), 0644); err != nil {
		return nil, fmt.Errorf("error writing visualization file: %w", err)
	}

	return &VisualizationResult{
		Type:     "token_boundary",
		Filepath: filepath,
		Data:     plotData,
		Metadata: map[string]interface{}{
			"document_id": data.DocumentID,
			"tokenizers":  data.TokenizerNames,
		},
	}, nil
}

// GenerateDriftVisualization generates drift comparison visualizations
func (v *VisualizationEngine) GenerateDriftVisualization(data DriftData) (*VisualizationResult, error) {
	// Create multiple plots for drift analysis
	plots := make([]map[string]interface{}, 0)

	// Token count drift
	if tokenCountPlot := v.createTokenCountDriftPlot(data); tokenCountPlot != nil {
		plots = append(plots, tokenCountPlot)
	}

	// Entropy drift
	if entropyPlot := v.createEntropyDriftPlot(data); entropyPlot != nil {
		plots = append(plots, entropyPlot)
	}

	// Alignment score
	if alignmentPlot := v.createAlignmentPlot(data); alignmentPlot != nil {
		plots = append(plots, alignmentPlot)
	}

	// Generate HTML with subplots
	html := v.generateMultiPlotHTML(plots, "drift_analysis")

	// Save to file
	filename := fmt.Sprintf("drift_analysis_%s.%s", data.ComparisonID, v.config.FileType)
	filepath := filepath.Join(v.config.OutputDir, filename)

	if err := os.WriteFile(filepath, []byte(html), 0644); err != nil {
		return nil, fmt.Errorf("error writing visualization file: %w", err)
	}

	return &VisualizationResult{
		Type:     "drift_analysis",
		Filepath: filepath,
		Data:     plots,
		Metadata: map[string]interface{}{
			"comparison_id": data.ComparisonID,
			"tokenizer1":    data.Tokenizer1,
			"tokenizer2":    data.Tokenizer2,
		},
	}, nil
}

// GenerateRollingEntropyPlot generates rolling entropy visualization
func (v *VisualizationEngine) GenerateRollingEntropyPlot(data RollingEntropyData) (*VisualizationResult, error) {
	// Create Plotly.js line plot
	plotData := v.createRollingEntropyPlotData(data)

	layout := map[string]interface{}{
		"title": map[string]interface{}{
			"text": "Rolling Entropy Analysis",
			"x":    0.5,
		},
		"xaxis": map[string]interface{}{
			"title":    "Window Position",
			"showgrid": true,
		},
		"yaxis": map[string]interface{}{
			"title":    "Entropy",
			"showgrid": true,
		},
		"height":   v.getHeight(),
		"width":    v.getWidth(),
		"template": v.getTemplate(),
	}

	// Generate HTML
	html := v.generatePlotlyHTML(plotData, layout, "rolling_entropy")

	// Save to file
	filename := fmt.Sprintf("rolling_entropy_%s.%s", data.DocumentID, v.config.FileType)
	filepath := filepath.Join(v.config.OutputDir, filename)

	if err := os.WriteFile(filepath, []byte(html), 0644); err != nil {
		return nil, fmt.Errorf("error writing visualization file: %w", err)
	}

	return &VisualizationResult{
		Type:     "rolling_entropy",
		Filepath: filepath,
		Data:     plotData,
		Metadata: map[string]interface{}{
			"document_id": data.DocumentID,
			"window_size": data.WindowSize,
		},
	}, nil
}

// GenerateComprehensiveReport generates a comprehensive visualization report
func (v *VisualizationEngine) GenerateComprehensiveReport(analysisResults []*metrics.AnalysisResult) (*VisualizationResult, error) {
	// Generate multiple visualizations
	visualizations := make([]*VisualizationResult, 0)

	// Token count heatmap
	if heatmapData := v.prepareHeatmapData(analysisResults, "token_count"); heatmapData != nil {
		if heatmap, err := v.GenerateHeatmap(*heatmapData, "token_count"); err == nil {
			visualizations = append(visualizations, heatmap)
		}
	}

	// Entropy heatmap
	if entropyData := v.prepareHeatmapData(analysisResults, "entropy"); entropyData != nil {
		if entropyHeatmap, err := v.GenerateHeatmap(*entropyData, "entropy"); err == nil {
			visualizations = append(visualizations, entropyHeatmap)
		}
	}

	// Compression heatmap
	if compressionData := v.prepareHeatmapData(analysisResults, "compression"); compressionData != nil {
		if compressionHeatmap, err := v.GenerateHeatmap(*compressionData, "compression"); err == nil {
			visualizations = append(visualizations, compressionHeatmap)
		}
	}

	// Generate report HTML
	html := v.generateReportHTML(visualizations)

	// Save to file
	filename := fmt.Sprintf("comprehensive_report.%s", v.config.FileType)
	filepath := filepath.Join(v.config.OutputDir, filename)

	if err := os.WriteFile(filepath, []byte(html), 0644); err != nil {
		return nil, fmt.Errorf("error writing report file: %w", err)
	}

	return &VisualizationResult{
		Type:     "comprehensive_report",
		Filepath: filepath,
		Data:     visualizations,
		Metadata: map[string]interface{}{
			"visualization_count": len(visualizations),
			"analysis_results":    len(analysisResults),
		},
	}, nil
}

// Helper methods for configuration
func (v *VisualizationEngine) getHeight() int {
	switch v.config.ImageSize {
	case "small":
		return 400
	case "large":
		return 800
	default:
		return 600
	}
}

func (v *VisualizationEngine) getWidth() int {
	switch v.config.ImageSize {
	case "small":
		return 600
	case "large":
		return 1200
	default:
		return 900
	}
}

func (v *VisualizationEngine) getTemplate() string {
	if v.config.Theme == "dark" {
		return "plotly_dark"
	}
	return "plotly_white"
}

// VisualizationResult represents the result of a visualization generation
type VisualizationResult struct {
	Type     string                 `json:"type"`
	Filepath string                 `json:"filepath"`
	Data     interface{}            `json:"data"`
	Metadata map[string]interface{} `json:"metadata"`
}

// Data structures for different visualization types
type HeatmapData struct {
	XLabels    []string    `json:"x_labels"`
	YLabels    []string    `json:"y_labels"`
	Values     [][]float64 `json:"values"`
	ColorScale string      `json:"color_scale"`
	Title      string      `json:"title"`
}

type TokenBoundaryData struct {
	DocumentID     string                           `json:"document_id"`
	Document       string                           `json:"document"`
	TokenizerNames []string                         `json:"tokenizer_names"`
	Tokenizations  []*tokenizers.TokenizationResult `json:"tokenizations"`
}

type DriftData struct {
	ComparisonID string               `json:"comparison_id"`
	Tokenizer1   string               `json:"tokenizer1"`
	Tokenizer2   string               `json:"tokenizer2"`
	Documents    []string             `json:"documents"`
	DriftMetrics map[string][]float64 `json:"drift_metrics"`
}

type RollingEntropyData struct {
	DocumentID    string    `json:"document_id"`
	TokenizerName string    `json:"tokenizer_name"`
	WindowSize    int       `json:"window_size"`
	EntropyValues []float64 `json:"entropy_values"`
}
