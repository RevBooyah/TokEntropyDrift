# TokEntropyDrift Tutorials

This document contains step-by-step tutorials for various use cases of TokEntropyDrift.

## Table of Contents

1. [Tutorial 1: Getting Started](#tutorial-1-getting-started)
2. [Tutorial 2: Multi-Tokenizer Comparison](#tutorial-2-multi-tokenizer-comparison)
3. [Tutorial 3: Performance Optimization](#tutorial-3-performance-optimization)
4. [Tutorial 4: Custom Plugin Development](#tutorial-4-custom-plugin-development)
5. [Tutorial 5: Large Dataset Processing](#tutorial-5-large-dataset-processing)
6. [Tutorial 6: Web Dashboard Usage](#tutorial-6-web-dashboard-usage)
7. [Tutorial 7: Visualization and Reporting](#tutorial-7-visualization-and-reporting)
8. [Tutorial 8: Integration with CI/CD](#tutorial-8-integration-with-cicd)

---

## Tutorial 1: Getting Started

**Duration**: 15 minutes  
**Difficulty**: Beginner  
**Prerequisites**: Basic command line knowledge

### Objective
Learn the basics of TokEntropyDrift by analyzing a simple text file.

### Step 1: Setup

1. **Build the application**:
   ```bash
   git clone https://github.com/RevBooyah/TokEntropyDrift.git
   cd tokentropydrift
   go build -o ted cmd/ted/main.go
   ```

2. **Create a test file**:
   ```bash
   echo "The quick brown fox jumps over the lazy dog. This is a simple test sentence." > tutorial1.txt
   ```

### Step 2: Basic Analysis

1. **Run your first analysis**:
   ```bash
   ./ted analyze tutorial1.txt --tokenizers=gpt2
   ```

2. **Check the output**:
   ```bash
   ls output/
   # You should see CSV files with analysis results
   ```

3. **Examine the results**:
   ```bash
   cat output/*.csv
   ```

### Step 3: Understanding the Output

The analysis produces several metrics:
- **Token Count**: Number of tokens generated
- **Entropy**: Information content of the tokenization
- **Compression Ratio**: How much the text is compressed
- **Token Reuse**: How efficiently tokens are reused

### Expected Results
- Token count: ~15-20 tokens
- Entropy: ~3-4 bits per token
- Compression ratio: ~0.6-0.8

### Next Steps
- Try different tokenizers: `--tokenizers=bert,t5`
- Add visualization: `--visualize`
- Explore the web dashboard: `./ted serve`

---

## Tutorial 2: Multi-Tokenizer Comparison

**Duration**: 30 minutes  
**Difficulty**: Beginner  
**Prerequisites**: Tutorial 1

### Objective
Compare how different tokenizers handle the same text and understand their differences.

### Step 1: Prepare Test Data

1. **Create a diverse test file**:
   ```bash
   cat > tutorial2.txt << 'EOF'
   The quick brown fox jumps over the lazy dog.
   Python is a programming language: def hello(): print("Hello, World!")
   Machine learning models like GPT-3 and BERT are transformers.
   Special characters: @#$%^&*() and emojis: 😀🎉🚀
   EOF
   ```

### Step 2: Run Multi-Tokenizer Analysis

1. **Analyze with multiple tokenizers**:
   ```bash
   ./ted analyze tutorial2.txt --tokenizers=gpt2,bert,t5 --visualize
   ```

2. **Generate comparison heatmap**:
   ```bash
   ./ted heatmap tutorial2.txt --type=entropy --tokenizers=gpt2,bert,t5 --output=comparison.svg
   ```

### Step 3: Analyze the Differences

1. **Token Count Comparison**:
   - GPT-2: Usually more tokens (subword tokenization)
   - BERT: Moderate token count (WordPiece)
   - T5: Fewer tokens (SentencePiece)

2. **Entropy Patterns**:
   - Higher entropy = more unpredictable tokenization
   - Lower entropy = more consistent patterns

3. **Compression Efficiency**:
   - Better compression = fewer tokens for same text
   - Trade-off between compression and vocabulary size

### Step 4: Visual Analysis

1. **Open the web dashboard**:
   ```bash
   ./ted serve --port=8080
   ```

2. **Upload your file and explore**:
   - Compare token boundaries
   - View entropy heatmaps
   - Analyze drift patterns

### Expected Insights
- Different tokenizers handle special characters differently
- Code snippets may have varying compression ratios
- Emojis and special symbols affect tokenization patterns

---

## Tutorial 3: Performance Optimization

**Duration**: 45 minutes  
**Difficulty**: Intermediate  
**Prerequisites**: Tutorial 1, basic understanding of performance concepts

### Objective
Learn how to optimize TokEntropyDrift for large datasets and repeated analysis.

### Step 1: Create Large Test Dataset

1. **Generate a large test file**:
   ```bash
   # Create a 10MB test file
   for i in {1..10000}; do
     echo "This is test sentence number $i with some technical terms like machine learning, artificial intelligence, and natural language processing." >> large_dataset.txt
   done
   ```

### Step 2: Baseline Performance Test

1. **Run analysis without optimizations**:
   ```bash
   time ./ted analyze large_dataset.txt --tokenizers=gpt2,bert,t5
   ```

2. **Note the execution time and memory usage**

### Step 3: Enable Caching

1. **Configure caching in `ted.config.yaml`**:
   ```yaml
   cache:
     enabled: true
     max_size: 50000
     ttl: "2h"
     cleanup_interval: "15m"
     enable_stats: true
   ```

2. **Test caching performance**:
   ```bash
   # First run (cache miss)
   time ./ted analyze large_dataset.txt --tokenizers=gpt2
   
   # Second run (cache hit)
   time ./ted analyze large_dataset.txt --tokenizers=gpt2
   ```

3. **Check cache statistics**:
   ```bash
   ./ted advanced cache large_dataset.txt
   ```

### Step 4: Enable Parallel Processing

1. **Configure parallel processing**:
   ```yaml
   parallel:
     enabled: true
     max_workers: 0  # Auto-detect
     batch_size: 1000
     timeout: "30m"
     enable_metrics: true
   ```

2. **Test parallel performance**:
   ```bash
   time ./ted analyze large_dataset.txt --tokenizers=gpt2,bert,t5
   ```

### Step 5: Enable Streaming for Very Large Files

1. **Create a very large file** (if you have the space):
   ```bash
   # Create a 100MB+ file
   for i in {1..100000}; do
     echo "Large dataset line $i with technical content about machine learning and artificial intelligence." >> very_large_dataset.txt
   done
   ```

2. **Configure streaming**:
   ```yaml
   streaming:
     enabled: true
     chunk_size: 5000
     buffer_size: 131072  # 128KB
     max_memory_mb: 1024
     enable_progress: true
     progress_interval: 5
     timeout: "2h"
   ```

3. **Test streaming analysis**:
   ```bash
   ./ted advanced streaming very_large_dataset.txt
   ```

### Step 6: Performance Comparison

Create a performance report:

```bash
echo "Performance Comparison Report" > performance_report.md
echo "=============================" >> performance_report.md
echo "" >> performance_report.md
echo "Baseline (no optimizations):" >> performance_report.md
echo "- Time: [your baseline time]" >> performance_report.md
echo "- Memory: [your baseline memory]" >> performance_report.md
echo "" >> performance_report.md
echo "With caching:" >> performance_report.md
echo "- Time: [cached time]" >> performance_report.md
echo "- Memory: [cached memory]" >> performance_report.md
echo "" >> performance_report.md
echo "With parallel processing:" >> performance_report.md
echo "- Time: [parallel time]" >> performance_report.md
echo "- Memory: [parallel memory]" >> performance_report.md
```

### Expected Improvements
- **Caching**: 80-90% speedup for repeated analysis
- **Parallel Processing**: 3-5x speedup for large datasets
- **Streaming**: Constant memory usage regardless of file size

---

## Tutorial 4: Custom Plugin Development

**Duration**: 60 minutes  
**Difficulty**: Advanced  
**Prerequisites**: Tutorial 1, basic Go programming knowledge

### Objective
Learn how to create custom plugins to extend TokEntropyDrift with your own metrics and analysis.

### Step 1: Plugin Structure

1. **Create plugin directory**:
   ```bash
   mkdir -p plugins/tutorial4
   cd plugins/tutorial4
   ```

2. **Create your first plugin**:
   ```go
   package main

   import (
       "math"
       "github.com/RevBooyah/TokEntropyDrift/internal/plugins"
   )

   // WordLengthAnalyzer analyzes word length patterns
   type WordLengthAnalyzer struct {
       *plugins.BasePlugin
   }

   func NewWordLengthAnalyzer() *WordLengthAnalyzer {
       info := plugins.PluginInfo{
           Name:        "word_length_analyzer",
           Version:     "1.0.0",
           Description: "Analyzes word length patterns in tokenized text",
           Author:      "Tutorial User",
           Tags:        []string{"analysis", "words", "length"},
       }
       
       return &WordLengthAnalyzer{
           BasePlugin: plugins.NewBasePlugin(info),
       }
   }

   func (w *WordLengthAnalyzer) CalculateMetrics(ctx *plugins.AnalysisContext) ([]plugins.MetricResult, error) {
       if ctx.Tokenization == nil || len(ctx.Tokenization.Tokens) == 0 {
           return []plugins.MetricResult{}, nil
       }

       tokens := ctx.Tokenization.Tokens
       
       // Calculate word lengths
       var lengths []int
       for _, token := range tokens {
           lengths = append(lengths, len(token.Text))
       }

       // Calculate statistics
       total := len(lengths)
       if total == 0 {
           return []plugins.MetricResult{}, nil
       }

       sum := 0
       for _, length := range lengths {
           sum += length
       }
       mean := float64(sum) / float64(total)

       // Calculate variance
       variance := 0.0
       for _, length := range lengths {
           diff := float64(length) - mean
           variance += diff * diff
       }
       variance /= float64(total)
       stdDev := math.Sqrt(variance)

       // Find min and max
       min := lengths[0]
       max := lengths[0]
       for _, length := range lengths {
           if length < min {
               min = length
           }
           if length > max {
               max = length
           }
       }

       return []plugins.MetricResult{
           {
               Name:  "mean_word_length",
               Value: mean,
               Unit:  "characters",
           },
           {
               Name:  "std_dev_word_length",
               Value: stdDev,
               Unit:  "characters",
           },
           {
               Name:  "min_word_length",
               Value: float64(min),
               Unit:  "characters",
           },
           {
               Name:  "max_word_length",
               Value: float64(max),
               Unit:  "characters",
           },
           {
               Name:  "total_words",
               Value: float64(total),
               Unit:  "words",
           },
       }, nil
   }

   func (w *WordLengthAnalyzer) ValidateConfig(config map[string]interface{}) error {
       // Add validation logic here
       return nil
   }
   ```

### Step 2: Build and Test Plugin

1. **Build the plugin**:
   ```bash
   go build -o word_length_analyzer.so -buildmode=plugin word_length_analyzer.go
   ```

2. **Configure the plugin**:
   ```yaml
   # In ted.config.yaml
   plugins:
     enabled: true
     auto_load: true
     plugin_directory: "plugins/tutorial4"
     configs:
       word_length_analyzer:
         enabled: true
   ```

3. **Test the plugin**:
   ```bash
   ./ted advanced plugins tutorial1.txt
   ```

### Step 3: Advanced Plugin Features

1. **Add configuration support**:
   ```go
   func (w *WordLengthAnalyzer) CalculateMetrics(ctx *plugins.AnalysisContext) ([]plugins.MetricResult, error) {
       // Get configuration
       minThreshold := w.GetConfigInt("min_length_threshold", 1)
       maxThreshold := w.GetConfigInt("max_length_threshold", 100)
       
       // ... existing code ...
       
       // Add threshold-based metrics
       withinThreshold := 0
       for _, length := range lengths {
           if length >= minThreshold && length <= maxThreshold {
               withinThreshold++
           }
       }
       
       thresholdPercentage := float64(withinThreshold) / float64(total) * 100
       
       results = append(results, plugins.MetricResult{
           Name:  "words_within_threshold",
           Value: float64(withinThreshold),
           Unit:  "words",
       })
       
       results = append(results, plugins.MetricResult{
           Name:  "threshold_percentage",
           Value: thresholdPercentage,
           Unit:  "percent",
       })
       
       return results, nil
   }
   ```

2. **Add error handling**:
   ```go
   func (w *WordLengthAnalyzer) ValidateConfig(config map[string]interface{}) error {
       if minThreshold, exists := config["min_length_threshold"]; exists {
           if min, ok := minThreshold.(int); !ok || min < 0 {
               return fmt.Errorf("min_length_threshold must be a non-negative integer")
           }
       }
       
       if maxThreshold, exists := config["max_length_threshold"]; exists {
           if max, ok := maxThreshold.(int); !ok || max <= 0 {
               return fmt.Errorf("max_length_threshold must be a positive integer")
           }
       }
       
       return nil
   }
   ```

### Step 4: Plugin Integration

1. **Test with different configurations**:
   ```yaml
   plugins:
     configs:
       word_length_analyzer:
         min_length_threshold: 3
         max_length_threshold: 15
   ```

2. **Run analysis with your plugin**:
   ```bash
   ./ted analyze tutorial2.txt --tokenizers=gpt2 --visualize
   ```

3. **Check plugin results in output**:
   ```bash
   grep "word_length" output/*.csv
   ```

### Step 5: Plugin Distribution

1. **Create plugin documentation**:
   ```markdown
   # Word Length Analyzer Plugin

   ## Description
   Analyzes word length patterns in tokenized text.

   ## Configuration
   - `min_length_threshold`: Minimum word length to consider (default: 1)
   - `max_length_threshold`: Maximum word length to consider (default: 100)

   ## Metrics
   - `mean_word_length`: Average word length in characters
   - `std_dev_word_length`: Standard deviation of word lengths
   - `min_word_length`: Shortest word length
   - `max_word_length`: Longest word length
   - `total_words`: Total number of words analyzed
   - `words_within_threshold`: Words within configured length range
   - `threshold_percentage`: Percentage of words within threshold
   ```

### Expected Results
- Custom metrics appear in analysis output
- Configuration changes affect plugin behavior
- Plugin integrates seamlessly with existing analysis

---

## Tutorial 5: Large Dataset Processing

**Duration**: 45 minutes  
**Difficulty**: Intermediate  
**Prerequisites**: Tutorial 3

### Objective
Learn how to efficiently process very large datasets using streaming and parallel processing.

### Step 1: Prepare Large Dataset

1. **Create a large dataset** (adjust size based on your system):
   ```bash
   # Create a 50MB dataset
   for i in {1..50000}; do
     echo "Dataset entry $i: This is a comprehensive analysis of machine learning models including transformers, convolutional neural networks, and recurrent neural networks. The text contains technical terminology and domain-specific vocabulary." >> large_dataset.txt
   done
   ```

2. **Check file size**:
   ```bash
   ls -lh large_dataset.txt
   ```

### Step 2: Configure for Large Dataset Processing

1. **Update configuration for large datasets**:
   ```yaml
   # Optimize for large files
   streaming:
     enabled: true
     chunk_size: 10000
     buffer_size: 262144  # 256KB
     max_memory_mb: 2048
     enable_progress: true
     progress_interval: 5
     timeout: "4h"

   parallel:
     enabled: true
     max_workers: 0  # Auto-detect
     batch_size: 5000
     timeout: "2h"
     enable_metrics: true

   cache:
     enabled: true
     max_size: 100000
     ttl: "4h"
     cleanup_interval: "30m"
     enable_stats: true
   ```

### Step 3: Streaming Analysis

1. **Run streaming analysis**:
   ```bash
   ./ted advanced streaming large_dataset.txt
   ```

2. **Monitor progress and memory usage**:
   ```bash
   # In another terminal, monitor memory
   watch -n 1 'ps aux | grep ted'
   ```

### Step 4: Parallel Processing

1. **Run parallel analysis**:
   ```bash
   ./ted advanced parallel large_dataset.txt
   ```

2. **Compare performance**:
   ```bash
   # Time the analysis
   time ./ted analyze large_dataset.txt --tokenizers=gpt2,bert
   ```

### Step 5: Memory-Efficient Processing

1. **Create a memory monitoring script**:
   ```bash
   cat > monitor_memory.sh << 'EOF'
   #!/bin/bash
   while true; do
     echo "$(date): $(ps aux | grep ted | grep -v grep | awk '{print $6/1024 " MB"}')"
     sleep 5
   done
   EOF
   chmod +x monitor_memory.sh
   ```

2. **Run with memory monitoring**:
   ```bash
   ./monitor_memory.sh &
   ./ted analyze large_dataset.txt --tokenizers=gpt2,bert
   kill %1
   ```

### Step 6: Results Analysis

1. **Generate performance report**:
   ```bash
   echo "Large Dataset Processing Report" > large_dataset_report.md
   echo "===============================" >> large_dataset_report.md
   echo "" >> large_dataset_report.md
   echo "Dataset Size: $(ls -lh large_dataset.txt | awk '{print $5}')" >> large_dataset_report.md
   echo "Processing Time: [your time]" >> large_dataset_report.md
   echo "Peak Memory Usage: [your memory]" >> large_dataset_report.md
   echo "Tokens Generated: [count from output]" >> large_dataset_report.md
   ```

2. **Analyze results**:
   ```bash
   # Check output files
   ls -la output/
   
   # Examine metrics
   head -20 output/*.csv
   ```

### Expected Results
- Constant memory usage with streaming
- Faster processing with parallel execution
- Comprehensive analysis results
- Progress tracking for long-running operations

---

## Tutorial 6: Web Dashboard Usage

**Duration**: 30 minutes  
**Difficulty**: Beginner  
**Prerequisites**: Tutorial 1

### Objective
Learn how to use the interactive web dashboard for visual analysis and exploration.

### Step 1: Launch Dashboard

1. **Start the web server**:
   ```bash
   ./ted serve --port=8080 --host=localhost
   ```

2. **Open in browser**:
   ```
   http://localhost:8080
   ```

### Step 2: Upload and Analyze

1. **Upload a file**:
   - Click "Upload File" or drag and drop
   - Select `tutorial2.txt` from previous tutorials
   - Choose tokenizers: GPT-2, BERT, T5

2. **Run analysis**:
   - Click "Analyze" to start processing
   - Watch real-time progress updates

### Step 3: Explore Visualizations

1. **Token Boundary Visualization**:
   - Compare how different tokenizers segment text
   - Hover over tokens to see details
   - Zoom and pan to explore

2. **Entropy Heatmaps**:
   - View entropy patterns across tokenizers
   - Identify high/low entropy regions
   - Export visualizations

3. **Drift Analysis**:
   - Compare tokenization behavior
   - Identify divergence points
   - Analyze consistency

### Step 4: Interactive Features

1. **Real-time Filtering**:
   - Filter by tokenizer
   - Filter by metric ranges
   - Search for specific tokens

2. **Export Options**:
   - Download CSV data
   - Export visualizations as PNG/SVG
   - Generate comprehensive reports

### Step 5: Advanced Dashboard Features

1. **Batch Processing**:
   - Upload multiple files
   - Queue analysis jobs
   - Monitor progress

2. **Custom Analysis**:
   - Configure analysis parameters
   - Select specific metrics
   - Set visualization options

### Expected Experience
- Intuitive web interface
- Real-time analysis updates
- Interactive visualizations
- Easy data export

---

## Tutorial 7: Visualization and Reporting

**Duration**: 40 minutes  
**Difficulty**: Intermediate  
**Prerequisites**: Tutorial 1, basic understanding of data visualization

### Objective
Learn how to create comprehensive visualizations and reports for analysis results.

### Step 1: Generate Base Analysis

1. **Create a diverse test dataset**:
   ```bash
   cat > visualization_test.txt << 'EOF'
   The quick brown fox jumps over the lazy dog.
   Machine learning models like GPT-3 and BERT are transformers.
   Python code: def hello(): print("Hello, World!")
   Special characters: @#$%^&*() and emojis: 😀🎉🚀
   Technical terms: convolutional neural networks, recurrent neural networks
   Natural language processing involves tokenization and embedding.
   EOF
   ```

2. **Run comprehensive analysis**:
   ```bash
   ./ted analyze visualization_test.txt --tokenizers=gpt2,bert,t5 --visualize
   ```

### Step 2: Create Heatmaps

1. **Entropy heatmap**:
   ```bash
   ./ted heatmap visualization_test.txt --type=entropy --tokenizers=gpt2,bert,t5 --output=entropy_heatmap.svg
   ```

2. **Token count heatmap**:
   ```bash
   ./ted heatmap visualization_test.txt --type=tokens --tokenizers=gpt2,bert,t5 --output=token_heatmap.svg
   ```

3. **Compression heatmap**:
   ```bash
   ./ted heatmap visualization_test.txt --type=compression --tokenizers=gpt2,bert,t5 --output=compression_heatmap.svg
   ```

### Step 3: Generate Reports

1. **Create a report script**:
   ```bash
   cat > generate_report.sh << 'EOF'
   #!/bin/bash
   
   echo "Generating TokEntropyDrift Analysis Report"
   echo "=========================================="
   echo ""
   echo "Analysis Date: $(date)"
   echo "Input File: $1"
   echo "Tokenizers: $2"
   echo ""
   
   # Run analysis
   ./ted analyze "$1" --tokenizers="$2" --visualize
   
   # Generate heatmaps
   ./ted heatmap "$1" --type=entropy --tokenizers="$2" --output=report_entropy.svg
   ./ted heatmap "$1" --type=tokens --tokenizers="$2" --output=report_tokens.svg
   ./ted heatmap "$1" --type=compression --tokenizers="$2" --output=report_compression.svg
   
   echo "Report generated successfully!"
   echo "Check output/ directory for results."
   EOF
   
   chmod +x generate_report.sh
   ```

2. **Generate comprehensive report**:
   ```bash
   ./generate_report.sh visualization_test.txt "gpt2,bert,t5"
   ```

### Step 4: Custom Visualizations

1. **Create a custom visualization script**:
   ```python
   #!/usr/bin/env python3
   import pandas as pd
   import matplotlib.pyplot as plt
   import seaborn as sns
   
   # Load analysis results
   df = pd.read_csv('output/analysis_results.csv')
   
   # Create comparison plot
   plt.figure(figsize=(12, 8))
   
   # Token count comparison
   plt.subplot(2, 2, 1)
   df.groupby('tokenizer')['token_count'].mean().plot(kind='bar')
   plt.title('Average Token Count by Tokenizer')
   plt.ylabel('Token Count')
   
   # Entropy comparison
   plt.subplot(2, 2, 2)
   df.groupby('tokenizer')['entropy'].mean().plot(kind='bar')
   plt.title('Average Entropy by Tokenizer')
   plt.ylabel('Entropy (bits)')
   
   # Compression ratio comparison
   plt.subplot(2, 2, 3)
   df.groupby('tokenizer')['compression_ratio'].mean().plot(kind='bar')
   plt.title('Average Compression Ratio by Tokenizer')
   plt.ylabel('Compression Ratio')
   
   # Token reuse comparison
   plt.subplot(2, 2, 4)
   df.groupby('tokenizer')['token_reuse'].mean().plot(kind='bar')
   plt.title('Average Token Reuse by Tokenizer')
   plt.ylabel('Token Reuse')
   
   plt.tight_layout()
   plt.savefig('custom_analysis.png', dpi=300, bbox_inches='tight')
   plt.show()
   ```

### Step 5: Interactive Reports

1. **Create an HTML report**:
   ```html
   <!DOCTYPE html>
   <html>
   <head>
       <title>TokEntropyDrift Analysis Report</title>
       <script src="https://cdn.plot.ly/plotly-latest.min.js"></script>
   </head>
   <body>
       <h1>TokEntropyDrift Analysis Report</h1>
       <div id="entropy-chart"></div>
       <div id="token-chart"></div>
       
       <script>
           // Load your data and create interactive charts
           // This would use the CSV output from your analysis
       </script>
   </body>
   </html>
   ```

### Expected Output
- Multiple visualization formats (SVG, PNG, HTML)
- Comprehensive analysis reports
- Interactive charts and graphs
- Professional presentation-ready output

---

## Tutorial 8: Integration with CI/CD

**Duration**: 30 minutes  
**Difficulty**: Advanced  
**Prerequisites**: Tutorial 1, basic CI/CD knowledge

### Objective
Learn how to integrate TokEntropyDrift into continuous integration and deployment pipelines.

### Step 1: Create Analysis Script

1. **Create a CI analysis script**:
   ```bash
   cat > ci_analysis.sh << 'EOF'
   #!/bin/bash
   
   set -e  # Exit on any error
   
   echo "Starting TokEntropyDrift CI Analysis"
   echo "===================================="
   
   # Configuration
   INPUT_FILE="${1:-examples/english_quotes.txt}"
   TOKENIZERS="${2:-gpt2,bert}"
   OUTPUT_DIR="ci_output"
   THRESHOLD_ENTROPY=4.0
   THRESHOLD_COMPRESSION=0.7
   
   # Create output directory
   mkdir -p "$OUTPUT_DIR"
   
   # Run analysis
   echo "Running analysis on $INPUT_FILE with tokenizers: $TOKENIZERS"
   ./ted analyze "$INPUT_FILE" --tokenizers="$TOKENIZERS" --output="$OUTPUT_DIR"
   
   # Check results
   echo "Checking analysis results..."
   
   # Extract metrics from CSV output
   ENTROPY=$(grep "entropy" "$OUTPUT_DIR"/*.csv | tail -1 | cut -d',' -f3)
   COMPRESSION=$(grep "compression" "$OUTPUT_DIR"/*.csv | tail -1 | cut -d',' -f3)
   
   echo "Average Entropy: $ENTROPY"
   echo "Average Compression: $COMPRESSION"
   
   # Validate thresholds
   if (( $(echo "$ENTROPY > $THRESHOLD_ENTROPY" | bc -l) )); then
       echo "WARNING: Entropy ($ENTROPY) exceeds threshold ($THRESHOLD_ENTROPY)"
       exit 1
   fi
   
   if (( $(echo "$COMPRESSION < $THRESHOLD_COMPRESSION" | bc -l) )); then
       echo "WARNING: Compression ($COMPRESSION) below threshold ($THRESHOLD_COMPRESSION)"
       exit 1
   fi
   
   echo "Analysis passed all checks!"
   echo "Results saved to $OUTPUT_DIR/"
   
   # Generate summary report
   cat > "$OUTPUT_DIR/ci_summary.md" << SUMMARY_EOF
   # CI Analysis Summary
   
   - **Input File**: $INPUT_FILE
   - **Tokenizers**: $TOKENIZERS
   - **Average Entropy**: $ENTROPY
   - **Average Compression**: $COMPRESSION
   - **Status**: PASSED
   - **Timestamp**: $(date)
   
   SUMMARY_EOF
   
   echo "CI analysis completed successfully!"
   EOF
   
   chmod +x ci_analysis.sh
   ```

### Step 2: GitHub Actions Integration

1. **Create GitHub Actions workflow**:
   ```yaml
   # .github/workflows/tokenizer-analysis.yml
   name: Tokenizer Analysis
   
   on:
     push:
       paths:
         - 'examples/**'
         - 'testdata/**'
     pull_request:
       paths:
         - 'examples/**'
         - 'testdata/**'
   
   jobs:
     analyze:
       runs-on: ubuntu-latest
       
       steps:
       - uses: actions/checkout@v3
       
       - name: Set up Go
         uses: actions/setup-go@v4
         with:
           go-version: '1.22'
       
       - name: Build TokEntropyDrift
         run: |
           go build -o ted cmd/ted/main.go
       
       - name: Run Analysis
         run: |
           ./ci_analysis.sh examples/english_quotes.txt "gpt2,bert"
       
       - name: Upload Results
         uses: actions/upload-artifact@v3
         with:
           name: analysis-results
           path: ci_output/
       
       - name: Comment Results
         if: github.event_name == 'pull_request'
         uses: actions/github-script@v6
         with:
           script: |
             const fs = require('fs');
             const summary = fs.readFileSync('ci_output/ci_summary.md', 'utf8');
             github.rest.issues.createComment({
               issue_number: context.issue.number,
               owner: context.repo.owner,
               repo: context.repo.repo,
               body: summary
             });
   ```

### Step 3: Jenkins Integration

1. **Create Jenkins pipeline**:
   ```groovy
   // Jenkinsfile
   pipeline {
       agent any
       
       stages {
           stage('Setup') {
               steps {
                   sh 'go build -o ted cmd/ted/main.go'
               }
           }
           
           stage('Analysis') {
               steps {
                   sh './ci_analysis.sh examples/english_quotes.txt "gpt2,bert"'
               }
           }
           
           stage('Archive Results') {
               steps {
                   archiveArtifacts artifacts: 'ci_output/**/*', fingerprint: true
               }
           }
           
           stage('Publish Report') {
               steps {
                   publishHTML([
                       allowMissing: false,
                       alwaysLinkToLastBuild: true,
                       keepAll: true,
                       reportDir: 'ci_output',
                       reportFiles: '*.html',
                       reportName: 'Tokenizer Analysis Report'
                   ])
               }
           }
       }
       
       post {
           always {
               cleanWs()
           }
       }
   }
   ```

### Step 4: Docker Integration

1. **Create Dockerfile**:
   ```dockerfile
   FROM golang:1.22-alpine AS builder
   
   WORKDIR /app
   COPY . .
   RUN go build -o ted cmd/ted/main.go
   
   FROM alpine:latest
   RUN apk --no-cache add ca-certificates
   
   WORKDIR /root/
   COPY --from=builder /app/ted .
   COPY --from=builder /app/ted.config.yaml .
   COPY --from=builder /app/examples ./examples
   
   ENTRYPOINT ["./ted"]
   ```

2. **Create docker-compose for analysis**:
   ```yaml
   # docker-compose.yml
   version: '3.8'
   
   services:
     ted-analysis:
       build: .
       volumes:
         - ./data:/app/data
         - ./output:/app/output
       command: analyze /app/data/input.txt --tokenizers=gpt2,bert --output=/app/output
   
     ted-dashboard:
       build: .
       ports:
         - "8080:8080"
       volumes:
         - ./data:/app/data
         - ./output:/app/output
       command: serve --host=0.0.0.0 --port=8080
   ```

### Step 5: Monitoring and Alerting

1. **Create monitoring script**:
   ```bash
   cat > monitor_analysis.sh << 'EOF'
   #!/bin/bash
   
   # Monitor analysis performance and alert on issues
   
   ANALYSIS_FILE="$1"
   ALERT_EMAIL="$2"
   
   # Run analysis and capture metrics
   start_time=$(date +%s)
   ./ted analyze "$ANALYSIS_FILE" --tokenizers=gpt2,bert --output=monitor_output
   end_time=$(date +%s)
   
   duration=$((end_time - start_time))
   
   # Check for issues
   if [ $duration -gt 300 ]; then  # More than 5 minutes
       echo "ALERT: Analysis took $duration seconds" | mail -s "TokEntropyDrift Performance Alert" "$ALERT_EMAIL"
   fi
   
   # Check for errors
   if [ -f "monitor_output/error.log" ]; then
       echo "ALERT: Analysis errors detected" | mail -s "TokEntropyDrift Error Alert" "$ALERT_EMAIL"
   fi
   
   echo "Analysis completed in $duration seconds"
   EOF
   
   chmod +x monitor_analysis.sh
   ```

### Expected Integration Benefits
- Automated analysis in CI/CD pipelines
- Quality gates based on tokenization metrics
- Automated reporting and alerting
- Scalable analysis infrastructure

---

## Conclusion

These tutorials provide a comprehensive introduction to TokEntropyDrift's capabilities. Each tutorial builds upon the previous ones, gradually introducing more advanced features and use cases.

### Next Steps

1. **Experiment**: Try the tutorials with your own data
2. **Customize**: Adapt the examples to your specific needs
3. **Extend**: Create custom plugins and visualizations
4. **Integrate**: Incorporate TokEntropyDrift into your workflows
5. **Contribute**: Share your improvements with the community

### Resources

- **Documentation**: See `/docs/` for detailed guides
- **Examples**: Check `/examples/` for sample code and data
- **API Reference**: See the user guide for complete API documentation
- **Community**: Join discussions for questions and ideas

Happy analyzing! 🚀 