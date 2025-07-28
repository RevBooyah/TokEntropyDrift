package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/RevBooyah/TokEntropyDrift/internal/loader"
	"github.com/RevBooyah/TokEntropyDrift/internal/metrics"
	"github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
	"github.com/RevBooyah/TokEntropyDrift/internal/visualization"
)

func main() {
	fmt.Println("TokEntropyDrift Visualization Example")
	fmt.Println("=====================================")

	// Create output directory
	if err := os.MkdirAll("output", 0755); err != nil {
		log.Fatalf("Failed to create output directory: %v", err)
	}

	// Step 1: Initialize tokenizers
	fmt.Println("\n1. Initializing tokenizers...")
	if err := tokenizers.RegisterMockTokenizer(); err != nil {
		log.Fatalf("Failed to register mock tokenizer: %v", err)
	}
	if err := tokenizers.RegisterGPT2Tokenizer(); err != nil {
		log.Printf("Warning: Failed to register GPT-2 tokenizer: %v", err)
	}

	// Step 2: Load sample documents
	fmt.Println("2. Loading sample documents...")
	docLoader := loader.NewLoader("txt")
	documents, err := docLoader.LoadDocuments("examples/english_quotes.txt")
	if err != nil {
		log.Fatalf("Failed to load documents: %v", err)
	}
	fmt.Printf("   Loaded %d documents\n", len(documents))

	// Step 3: Initialize metrics engine
	fmt.Println("3. Initializing metrics engine...")
	metricsEngine := metrics.NewEngine(metrics.EngineConfig{
		EntropyWindowSize: 50,
		NormalizeEntropy:  true,
	})

	// Step 4: Analyze documents with multiple tokenizers
	fmt.Println("4. Analyzing documents...")
	var analysisResults []*metrics.AnalysisResult
	tokenizerNames := []string{"mock"}

	// Try to add GPT-2 if available
	if _, err := tokenizers.GetGlobal("gpt2"); err == nil {
		tokenizerNames = append(tokenizerNames, "gpt2")
	}

	for i, doc := range documents[:5] { // Use first 5 documents
		fmt.Printf("   Analyzing document %d: %s\n", i+1, truncateString(doc.Content, 50))

		for _, tokenizerName := range tokenizerNames {
			tokenizer, err := tokenizers.GetGlobal(tokenizerName)
			if err != nil {
				log.Printf("Warning: Failed to get tokenizer %s: %v", tokenizerName, err)
				continue
			}

			result, err := metricsEngine.AnalyzeDocument(context.Background(), doc.Content, tokenizer)
			if err != nil {
				log.Printf("Warning: Failed to analyze document with %s: %v", tokenizerName, err)
				continue
			}

			analysisResults = append(analysisResults, result)
		}
	}

	fmt.Printf("   Generated %d analysis results\n", len(analysisResults))

	// Step 5: Initialize visualization engine
	fmt.Println("5. Initializing visualization engine...")
	vizEngine := visualization.NewVisualizationEngine(visualization.VisualizationConfig{
		Theme:       "light",
		ImageSize:   "medium",
		FileType:    "html",
		Interactive: true,
		OutputDir:   "output",
	})

	// Step 6: Generate comprehensive report
	fmt.Println("6. Generating comprehensive report...")
	report, err := vizEngine.GenerateComprehensiveReport(analysisResults)
	if err != nil {
		log.Fatalf("Failed to generate comprehensive report: %v", err)
	}
	fmt.Printf("   Report generated: %s\n", report.Filepath)
	fmt.Printf("   Visualizations included: %d\n", report.Metadata["visualization_count"])

	// Step 7: Generate individual heatmaps
	fmt.Println("7. Generating individual heatmaps...")
	heatmapTypes := []string{"token_count", "entropy", "compression", "reuse"}

	for _, heatmapType := range heatmapTypes {
		heatmapData := vizEngine.prepareHeatmapData(analysisResults, heatmapType)
		if heatmapData != nil {
			result, err := vizEngine.GenerateHeatmap(*heatmapData, heatmapType)
			if err != nil {
				log.Printf("Warning: Failed to generate %s heatmap: %v", heatmapType, err)
			} else {
				fmt.Printf("   Generated %s heatmap: %s\n", heatmapType, result.Filepath)
			}
		}
	}

	// Step 8: Generate token boundary visualization
	fmt.Println("8. Generating token boundary visualization...")
	if len(analysisResults) > 0 {
		// Get tokenizations for the first document
		var tokenizations []*tokenizers.TokenizationResult
		var tokenizerNames []string

		for _, result := range analysisResults {
			if result.Document == documents[0].Content {
				tokenizations = append(tokenizations, result.Tokenization)
				tokenizerNames = append(tokenizerNames, result.TokenizerName)
			}
		}

		if len(tokenizations) > 0 {
			boundaryData := visualization.TokenBoundaryData{
				DocumentID:     "sample_quote",
				Document:       documents[0].Content,
				TokenizerNames: tokenizerNames,
				Tokenizations:  tokenizations,
			}

			boundaryViz, err := vizEngine.GenerateTokenBoundaryMap(boundaryData)
			if err != nil {
				log.Printf("Warning: Failed to generate token boundary visualization: %v", err)
			} else {
				fmt.Printf("   Generated token boundary visualization: %s\n", boundaryViz.Filepath)
			}
		}
	}

	// Step 9: Generate drift analysis (if multiple tokenizers)
	fmt.Println("9. Generating drift analysis...")
	if len(tokenizerNames) > 1 {
		// Prepare drift data
		driftData := prepareDriftData(analysisResults, tokenizerNames)
		if driftData != nil {
			driftViz, err := vizEngine.GenerateDriftVisualization(*driftData)
			if err != nil {
				log.Printf("Warning: Failed to generate drift visualization: %v", err)
			} else {
				fmt.Printf("   Generated drift visualization: %s\n", driftViz.Filepath)
			}
		}
	}

	// Step 10: Generate rolling entropy plot
	fmt.Println("10. Generating rolling entropy plot...")
	if len(analysisResults) > 0 {
		// Use the first result for rolling entropy
		result := analysisResults[0]
		if entropyMetric, exists := result.Metrics["entropy_rolling_entropy_mean"]; exists {
			rollingData := visualization.RollingEntropyData{
				DocumentID:    "sample_doc",
				TokenizerName: result.TokenizerName,
				WindowSize:    50,
				EntropyValues: []float64{entropyMetric.Value}, // Simplified for example
			}

			rollingViz, err := vizEngine.GenerateRollingEntropyPlot(rollingData)
			if err != nil {
				log.Printf("Warning: Failed to generate rolling entropy plot: %v", err)
			} else {
				fmt.Printf("   Generated rolling entropy plot: %s\n", rollingViz.Filepath)
			}
		}
	}

	fmt.Println("\n✅ Visualization example completed successfully!")
	fmt.Println("\n📁 Generated files in 'output/' directory:")
	fmt.Println("   - comprehensive_report.html (main report)")
	fmt.Println("   - Various heatmaps and visualizations")
	fmt.Println("\n🌐 Open the HTML files in your browser to view the interactive visualizations!")
}

// Helper function to truncate strings
func truncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// Helper function to prepare drift data
func prepareDriftData(analysisResults []*metrics.AnalysisResult, tokenizerNames []string) *visualization.DriftData {
	if len(tokenizerNames) < 2 {
		return nil
	}

	// Group results by document
	docResults := make(map[string]map[string]*metrics.AnalysisResult)
	for _, result := range analysisResults {
		if _, exists := docResults[result.Document]; !exists {
			docResults[result.Document] = make(map[string]*metrics.AnalysisResult)
		}
		docResults[result.Document][result.TokenizerName] = result
	}

	// Prepare drift metrics
	var documents []string
	var tokenCountDeltas []float64
	var entropyDeltas []float64
	var alignmentScores []float64

	for doc, results := range docResults {
		if len(results) >= 2 {
			documents = append(documents, truncateString(doc, 30))

			// Calculate deltas between first two tokenizers
			tokenizer1 := tokenizerNames[0]
			tokenizer2 := tokenizerNames[1]

			if result1, exists1 := results[tokenizer1]; exists1 {
				if result2, exists2 := results[tokenizer2]; exists2 {
					// Token count delta
					tokenDelta := float64(result2.TokenCount - result1.TokenCount)
					tokenCountDeltas = append(tokenCountDeltas, tokenDelta)

					// Entropy delta
					entropy1 := result1.Metrics["entropy_global_entropy"].Value
					entropy2 := result2.Metrics["entropy_global_entropy"].Value
					entropyDeltas = append(entropyDeltas, entropy2-entropy1)

					// Simple alignment score (placeholder)
					alignmentScores = append(alignmentScores, 0.85)
				}
			}
		}
	}

	if len(documents) == 0 {
		return nil
	}

	return &visualization.DriftData{
		ComparisonID: fmt.Sprintf("%s_vs_%s", tokenizerNames[0], tokenizerNames[1]),
		Tokenizer1:   tokenizerNames[0],
		Tokenizer2:   tokenizerNames[1],
		Documents:    documents,
		DriftMetrics: map[string][]float64{
			"token_count_delta": tokenCountDeltas,
			"entropy_delta":     entropyDeltas,
			"alignment_score":   alignmentScores,
		},
	}
}
