package visualization

import (
	"encoding/json"
	"fmt"
	"os"
)

// generatePlotlyHTML generates HTML with Plotly.js visualization
func (v *VisualizationEngine) generatePlotlyHTML(data []map[string]interface{}, layout map[string]interface{}, id string) string {
	// Convert data to JSON
	dataJSON, _ := json.Marshal(data)
	layoutJSON, _ := json.Marshal(layout)

	// Generate HTML template
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>TokEntropyDrift Visualization</title>
    <script src="https://cdn.plot.ly/plotly-latest.min.js"></script>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: %s;
        }
        .container {
            max-width: 1200px;
            margin: 0 auto;
        }
        .plot-container {
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            padding: 20px;
            margin: 20px 0;
        }
        .title {
            text-align: center;
            color: #333;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="title">TokEntropyDrift Analysis</h1>
        <div class="plot-container">
            <div id="%s"></div>
        </div>
    </div>
    
    <script>
        var data = %s;
        var layout = %s;
        
        Plotly.newPlot('%s', data, layout, {
            responsive: true,
            displayModeBar: true,
            modeBarButtonsToRemove: ['pan2d', 'lasso2d', 'select2d'],
            toImageButtonOptions: {
                format: 'png',
                filename: '%s',
                height: %d,
                width: %d,
                scale: 2
            }
        });
    </script>
</body>
</html>`, v.getBackgroundColor(), id, string(dataJSON), string(layoutJSON), id, id, v.getHeight(), v.getWidth())

	return html
}

// generateMultiPlotHTML generates HTML with multiple subplots
func (v *VisualizationEngine) generateMultiPlotHTML(plots []map[string]interface{}, id string) string {
	// Convert plots to JSON
	plotsJSON, _ := json.Marshal(plots)

	// Create subplot layout
	rows := 1
	cols := len(plots)
	if len(plots) > 2 {
		rows = 2
		cols = (len(plots) + 1) / 2
	}

	layout := map[string]interface{}{
		"grid": map[string]interface{}{
			"rows":    rows,
			"columns": cols,
			"pattern": "independent",
		},
		"height":   v.getHeight() * rows,
		"width":    v.getWidth(),
		"template": v.getTemplate(),
	}

	layoutJSON, _ := json.Marshal(layout)

	// Generate HTML template
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>TokEntropyDrift Multi-Plot Visualization</title>
    <script src="https://cdn.plot.ly/plotly-latest.min.js"></script>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 20px;
            background-color: %s;
        }
        .container {
            max-width: 1400px;
            margin: 0 auto;
        }
        .plot-container {
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            padding: 20px;
            margin: 20px 0;
        }
        .title {
            text-align: center;
            color: #333;
            margin-bottom: 20px;
        }
    </style>
</head>
<body>
    <div class="container">
        <h1 class="title">TokEntropyDrift Multi-Plot Analysis</h1>
        <div class="plot-container">
            <div id="%s"></div>
        </div>
    </div>
    
    <script>
        var plots = %s;
        var layout = %s;
        
        Plotly.newPlot('%s', plots, layout, {
            responsive: true,
            displayModeBar: true,
            modeBarButtonsToRemove: ['pan2d', 'lasso2d', 'select2d'],
            toImageButtonOptions: {
                format: 'png',
                filename: '%s',
                height: %d,
                width: %d,
                scale: 2
            }
        });
    </script>
</body>
</html>`, v.getBackgroundColor(), id, string(plotsJSON), string(layoutJSON), id, id, v.getHeight()*rows, v.getWidth())

	return html
}

// generateReportHTML generates a comprehensive report HTML
func (v *VisualizationEngine) generateReportHTML(visualizations []*VisualizationResult) string {
	// Create navigation and iframe structure
	navItems := ""
	iframeContent := ""

	for i, viz := range visualizations {
		navItems += fmt.Sprintf(`
            <li><a href="#viz%d" onclick="showVisualization(%d)">%s</a></li>`, i, i, viz.Type)

		iframeContent += fmt.Sprintf(`
            <div id="viz%d" class="viz-frame" style="display: %s;">
                <iframe src="%s" width="100%%" height="600px" frameborder="0"></iframe>
            </div>`, i, func() string {
			if i == 0 {
				return "block"
			} else {
				return "none"
			}
		}(), viz.Filepath)
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>TokEntropyDrift Comprehensive Report</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            margin: 0;
            padding: 0;
            background-color: %s;
        }
        .header {
            background-color: #2c3e50;
            color: white;
            padding: 20px;
            text-align: center;
        }
        .nav {
            background-color: #34495e;
            padding: 10px;
        }
        .nav ul {
            list-style: none;
            margin: 0;
            padding: 0;
            display: flex;
            justify-content: center;
            flex-wrap: wrap;
        }
        .nav li {
            margin: 0 10px;
        }
        .nav a {
            color: white;
            text-decoration: none;
            padding: 8px 16px;
            border-radius: 4px;
            transition: background-color 0.3s;
        }
        .nav a:hover {
            background-color: #5a6c7d;
        }
        .content {
            padding: 20px;
            max-width: 1400px;
            margin: 0 auto;
        }
        .viz-frame {
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            margin: 20px 0;
            overflow: hidden;
        }
        .summary {
            background-color: white;
            border-radius: 8px;
            box-shadow: 0 2px 4px rgba(0,0,0,0.1);
            padding: 20px;
            margin: 20px 0;
        }
    </style>
</head>
<body>
    <div class="header">
        <h1>TokEntropyDrift Comprehensive Analysis Report</h1>
        <p>Generated on %s</p>
    </div>
    
    <div class="nav">
        <ul>
            <li><a href="#summary" onclick="showSummary()">Summary</a></li>%s
        </ul>
    </div>
    
    <div class="content">
        <div id="summary" class="summary">
            <h2>Analysis Summary</h2>
            <p>This report contains %d visualizations analyzing tokenization behavior across different tokenizers and documents.</p>
            <ul>
                <li><strong>Total Visualizations:</strong> %d</li>
                <li><strong>Generated:</strong> %s</li>
                <li><strong>Theme:</strong> %s</li>
            </ul>
        </div>
        
        %s
    </div>
    
    <script>
        function showVisualization(index) {
            // Hide all visualizations
            var frames = document.querySelectorAll('.viz-frame');
            for (var i = 0; i < frames.length; i++) {
                frames[i].style.display = 'none';
            }
            
            // Show selected visualization
            document.getElementById('viz' + index).style.display = 'block';
            
            // Update navigation
            var navLinks = document.querySelectorAll('.nav a');
            for (var i = 0; i < navLinks.length; i++) {
                navLinks[i].style.backgroundColor = '';
            }
            event.target.style.backgroundColor = '#5a6c7d';
        }
        
        function showSummary() {
            // Hide all visualizations
            var frames = document.querySelectorAll('.viz-frame');
            for (var i = 0; i < frames.length; i++) {
                frames[i].style.display = 'none';
            }
            
            // Show summary
            document.getElementById('summary').style.display = 'block';
            
            // Update navigation
            var navLinks = document.querySelectorAll('.nav a');
            for (var i = 0; i < navLinks.length; i++) {
                navLinks[i].style.backgroundColor = '';
            }
            event.target.style.backgroundColor = '#5a6c7d';
        }
    </script>
</body>
</html>`, v.getBackgroundColor(), v.getCurrentTimestamp(), navItems, len(visualizations), len(visualizations), v.getCurrentTimestamp(), v.config.Theme, iframeContent)

	return html
}

// Helper methods for HTML generation
func (v *VisualizationEngine) getBackgroundColor() string {
	if v.config.Theme == "dark" {
		return "#1a1a1a"
	}
	return "#f5f5f5"
}

func (v *VisualizationEngine) getCurrentTimestamp() string {
	// This would be better implemented with actual time formatting
	return "2024-01-01 12:00:00"
}

func (v *VisualizationEngine) saveHTML(filepath string, html string) error {
	return os.WriteFile(filepath, []byte(html), 0644)
}
