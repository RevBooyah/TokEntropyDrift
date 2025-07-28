package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/RevBooyah/TokEntropyDrift/internal/config"
	"github.com/RevBooyah/TokEntropyDrift/internal/loader"
	"github.com/RevBooyah/TokEntropyDrift/internal/metrics"
	"github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
	"github.com/RevBooyah/TokEntropyDrift/internal/visualization"
	"github.com/spf13/cobra"
)

// Example CLI command that integrates visualization
var analyzeWithVizCmd = &cobra.Command{
	Use:   "analyze-viz [input-file]",
	Short: "Analyze with visualization output",
	Long: `Analyze tokenization and generate comprehensive visualizations.
This example shows how to integrate the visualization engine with CLI commands.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		runAnalysisWithVisualization(args[0])
	},
}

func main() {
	if err := analyzeWithVizCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func runAnalysisWithVisualization(inputFile string) {
	fmt.Printf("Analyzing %s with visualization output...\n", inputFile)

	// Load configuration
	cfg, err := config.LoadConfig("")
	if err != nil {
		log.Printf("Warning: Using default configuration: %v", err)
		cfg = &config.Config{}
	}

	// Register tokenizers
	fmt.Println("1. Registering tokenizers...")
	if err := tokenizers.RegisterMockTokenizer(); err != nil {
		log.Fatalf("Failed to register mock tokenizer: %v", err)
	}
	if err := tokenizers.RegisterGPT2Tokenizer(); err != nil {
		log.Printf("Warning: GPT-2 tokenizer not available: %v", err)
	}

	// Load documents
	fmt.Println("2. Loading documents...")
	fileType := config.GetFileType(inputFile)
	docLoader := loader.NewLoader(fileType)
	documents, err := docLoader.LoadDocuments(inputFile)
	if err != nil {
		log.Fatalf("Failed to load documents: %v", err)
	}
	fmt.Printf("   Loaded %d documents\n", len(documents))

	// Initialize metrics engine
	fmt.Println("3. Initializing metrics engine...")
	metricsEngine := metrics.NewEngine(metrics.EngineConfig{
		EntropyWindowSize: cfg.Analysis.EntropyWindowSize,
		NormalizeEntropy:  cfg.Analysis.NormalizeEntropy,
	})

	// Get enabled tokenizers
	enabledTokenizers := cfg.Tokenizers.Enabled
	if len(enabledTokenizers) == 0 {
		enabledTokenizers = []string{"mock"}
	}

	// Analyze documents
	fmt.Println("4. Analyzing documents...")
	var analysisResults []*metrics.AnalysisResult

	for i, doc := range documents {
		fmt.Printf("   Analyzing document %d/%d\n", i+1, len(documents))

		for _, tokenizerName := range enabledTokenizers {
			tokenizer, err := tokenizers.GetGlobal(tokenizerName)
			if err != nil {
				log.Printf("Warning: Tokenizer %s not available: %v", tokenizerName, err)
				continue
			}

			result, err := metricsEngine.AnalyzeDocument(context.Background(), doc.Content, tokenizer)
			if err != nil {
				log.Printf("Warning: Failed to analyze with %s: %v", tokenizerName, err)
				continue
			}

			analysisResults = append(analysisResults, result)
		}
	}

	fmt.Printf("   Generated %d analysis results\n", len(analysisResults))

	// Create output directory
	outputDir := cfg.Output.Directory
	if outputDir == "" {
		outputDir = "output"
	}

	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Initialize visualization engine
	fmt.Println("5. Initializing visualization engine...")
	vizEngine := visualization.NewVisualizationEngine(visualization.VisualizationConfig{
		Theme:       cfg.Visualization.Theme,
		ImageSize:   cfg.Visualization.ImageSize,
		FileType:    cfg.Visualization.FileType,
		Interactive: cfg.Visualization.Interactive,
		OutputDir:   outputDir,
	})

	// Generate visualizations
	fmt.Println("6. Generating visualizations...")

	// Comprehensive report
	fmt.Println("   Generating comprehensive report...")
	report, err := vizEngine.GenerateComprehensiveReport(analysisResults)
	if err != nil {
		log.Printf("Warning: Failed to generate comprehensive report: %v", err)
	} else {
		fmt.Printf("   ✅ Report: %s\n", report.Filepath)
	}

	// Individual heatmaps
	heatmapTypes := []string{"token_count", "entropy", "compression", "reuse"}
	for _, heatmapType := range heatmapTypes {
		fmt.Printf("   Generating %s heatmap...\n", heatmapType)
		heatmapData := vizEngine.prepareHeatmapData(analysisResults, heatmapType)
		if heatmapData != nil {
			result, err := vizEngine.GenerateHeatmap(*heatmapData, heatmapType)
			if err != nil {
				log.Printf("   ❌ Failed: %v", err)
			} else {
				fmt.Printf("   ✅ %s heatmap: %s\n", heatmapType, result.Filepath)
			}
		}
	}

	// Token boundary visualization (if multiple tokenizers)
	if len(enabledTokenizers) > 1 && len(documents) > 0 {
		fmt.Println("   Generating token boundary visualization...")

		var tokenizations []*tokenizers.TokenizationResult
		var tokenizerNames []string

		// Get tokenizations for first document
		for _, result := range analysisResults {
			if result.Document == documents[0].Content {
				tokenizations = append(tokenizations, result.Tokenization)
				tokenizerNames = append(tokenizerNames, result.TokenizerName)
			}
		}

		if len(tokenizations) > 0 {
			boundaryData := visualization.TokenBoundaryData{
				DocumentID:     "cli_example",
				Document:       documents[0].Content,
				TokenizerNames: tokenizerNames,
				Tokenizations:  tokenizations,
			}

			result, err := vizEngine.GenerateTokenBoundaryMap(boundaryData)
			if err != nil {
				log.Printf("   ❌ Failed: %v", err)
			} else {
				fmt.Printf("   ✅ Token boundary: %s\n", result.Filepath)
			}
		}
	}

	fmt.Println("\n🎉 Analysis with visualization completed!")
	fmt.Printf("📁 Output directory: %s\n", outputDir)
	fmt.Println("🌐 Open the HTML files in your browser to view the visualizations")
}

// Helper function to get file type from filename
func getFileType(filename string) string {
	// This would be implemented based on file extension
	// For now, default to "txt"
	return "txt"
}
