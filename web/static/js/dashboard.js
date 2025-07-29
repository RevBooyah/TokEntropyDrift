// TokEntropyDrift Dashboard JavaScript

class Dashboard {
    constructor() {
        this.sessionId = null;
        this.documents = [];
        this.tokenizers = [];
        this.analyses = [];
        this.currentAnalysis = null;
        
        this.init();
    }

    async init() {
        await this.loadSession();
        await this.loadDocuments();
        await this.loadTokenizers();
        await this.loadAnalysisHistory();
        this.setupEventListeners();
        this.setupWebSocket();
        
        // Initialize UI state
        this.initializeUIState();
    }

    async loadSession() {
        try {
            const response = await fetch('/api/v1/session');
            const session = await response.json();
            this.sessionId = session.session_id;
            document.getElementById('sessionInfo').textContent = `Session: ${session.session_id.substring(0, 8)}...`;
        } catch (error) {
            console.error('Failed to load session:', error);
        }
    }

    async loadDocuments() {
        try {
            const response = await fetch('/api/v1/documents');
            this.documents = await response.json();
            this.renderDocumentList();
            this.updateDocumentSelect();
        } catch (error) {
            console.error('Failed to load documents:', error);
        }
    }

    async loadTokenizers() {
        try {
            const response = await fetch('/api/v1/tokenizers');
            this.tokenizers = await response.json();
            this.renderTokenizerList();
        } catch (error) {
            console.error('Failed to load tokenizers:', error);
        }
    }

    setupEventListeners() {
        // File upload
        document.getElementById('uploadForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleFileUpload();
        });

        // Analysis form
        document.getElementById('analysisForm').addEventListener('submit', (e) => {
            e.preventDefault();
            this.handleAnalysis();
        });

        // Visualization buttons
        document.getElementById('generateHeatmap').addEventListener('click', () => {
            this.generateVisualization('heatmap');
        });

        document.getElementById('generateDrift').addEventListener('click', () => {
            this.generateVisualization('drift');
        });

        document.getElementById('generateEntropy').addEventListener('click', () => {
            this.generateVisualization('entropy');
        });

        // Document selection
        document.getElementById('documentSelect').addEventListener('change', (e) => {
            this.onDocumentSelect(e.target.value);
        });

        // Export results
        document.getElementById('exportResults').addEventListener('click', () => {
            this.exportResults();
        });
    }

    async handleFileUpload() {
        const form = document.getElementById('uploadForm');
        const formData = new FormData(form);
        const progressBar = document.querySelector('#uploadProgress .progress-bar');
        const progressDiv = document.getElementById('uploadProgress');

        try {
            progressDiv.style.display = 'block';
            progressBar.style.width = '0%';

            const xhr = new XMLHttpRequest();
            
            xhr.upload.addEventListener('progress', (e) => {
                if (e.lengthComputable) {
                    const percentComplete = (e.loaded / e.total) * 100;
                    progressBar.style.width = percentComplete + '%';
                }
            });

            xhr.addEventListener('load', async () => {
                if (xhr.status === 200) {
                    const result = JSON.parse(xhr.responseText);
                    this.showAlert('File uploaded successfully!', 'success');
                    await this.loadDocuments();
                    form.reset();
                } else {
                    this.showAlert('Upload failed: ' + xhr.statusText, 'danger');
                }
                progressDiv.style.display = 'none';
            });

            xhr.addEventListener('error', () => {
                this.showAlert('Upload failed', 'danger');
                progressDiv.style.display = 'none';
            });

            xhr.open('POST', '/api/v1/upload');
            xhr.send(formData);

        } catch (error) {
            console.error('Upload error:', error);
            this.showAlert('Upload failed', 'danger');
            progressDiv.style.display = 'none';
        }
    }

    async handleAnalysis() {
        const documentId = document.getElementById('documentSelect').value;
        if (!documentId) {
            this.showAlert('Please select a document', 'warning');
            return;
        }

        const selectedTokenizers = this.getSelectedTokenizers();
        if (selectedTokenizers.length === 0) {
            // Try to auto-select mock tokenizer
            this.ensureTokenizerSelected();
            const retrySelectedTokenizers = this.getSelectedTokenizers();
            if (retrySelectedTokenizers.length === 0) {
                this.showAlert('Please select at least one tokenizer', 'warning');
                return;
            } else {
                this.showAlert('Auto-selected mock tokenizer (works without Python dependencies)', 'info');
            }
        }

        const selectedMetrics = this.getSelectedMetrics();
        if (selectedMetrics.length === 0) {
            this.showAlert('Please select at least one metric', 'warning');
            return;
        }

        this.showLoading(true);

        try {
            const response = await fetch('/api/v1/analyze', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    document_id: documentId,
                    tokenizer_ids: selectedTokenizers,
                    metrics: selectedMetrics
                })
            });

            if (response.ok) {
                const analysis = await response.json();
                this.currentAnalysis = analysis;
                this.renderAnalysisResults(analysis);
                this.addToAnalysisHistory(analysis);
                this.showAlert('Analysis completed successfully!', 'success');
            } else {
                const error = await response.text();
                this.showAlert('Analysis failed: ' + error, 'danger');
            }
        } catch (error) {
            console.error('Analysis error:', error);
            this.showAlert('Analysis failed', 'danger');
        } finally {
            this.showLoading(false);
        }
    }

    addToAnalysisHistory(analysis) {
        const historyContainer = document.getElementById('analysisHistory');
        
        // Create history item
        const historyItem = document.createElement('div');
        historyItem.className = 'list-group-item list-group-item-action';
        
        const documentName = this.documents.find(doc => doc.id === analysis.document_id)?.filename || analysis.document_id;
        const tokenizerCount = analysis.results.length;
        
        historyItem.innerHTML = `
            <div class="d-flex w-100 justify-content-between">
                <h6 class="mb-1">${documentName}</h6>
                <small class="text-muted">${new Date(analysis.timestamp).toLocaleString()}</small>
            </div>
            <p class="mb-1">${tokenizerCount} tokenizer(s) analyzed</p>
            <small class="text-muted">Analysis ID: ${analysis.id}</small>
        `;
        
        // Add click handler to view this analysis
        historyItem.addEventListener('click', () => {
            this.currentAnalysis = analysis;
            this.renderAnalysisResults(analysis);
        });
        
        // Add to the beginning of the history
        historyContainer.insertBefore(historyItem, historyContainer.firstChild);
        
        // Limit history to 10 items
        const items = historyContainer.querySelectorAll('.list-group-item');
        if (items.length > 10) {
            items[items.length - 1].remove();
        }
    }

    async loadAnalysisHistory() {
        try {
            const response = await fetch('/api/v1/analyses');
            if (response.ok) {
                const analyses = await response.json();
                this.renderAnalysisHistory(analyses);
            }
        } catch (error) {
            console.error('Failed to load analysis history:', error);
        }
    }

    renderAnalysisHistory(analyses) {
        const historyContainer = document.getElementById('analysisHistory');
        historyContainer.innerHTML = '';
        
        if (analyses.length === 0) {
            historyContainer.innerHTML = '<p class="text-muted text-center">No analysis history available</p>';
            return;
        }
        
        analyses.forEach(analysis => {
            this.addToAnalysisHistory(analysis);
        });
    }

    async generateVisualization(type) {
        if (!this.currentAnalysis) {
            this.showAlert('Please run an analysis first', 'warning');
            return;
        }

        this.showLoading(true);

        try {
            const response = await fetch(`/api/v1/visualizations/${type}`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    document_id: this.currentAnalysis.document_id,
                    tokenizers: this.currentAnalysis.results.map(r => r.tokenizer_id),
                    type: type
                })
            });

            if (response.ok) {
                const visualization = await response.json();
                this.renderVisualization(visualization, type);
            } else {
                const error = await response.text();
                this.showAlert(`Failed to generate ${type}: ` + error, 'danger');
            }
        } catch (error) {
            console.error('Visualization error:', error);
            this.showAlert(`Failed to generate ${type}`, 'danger');
        } finally {
            this.showLoading(false);
        }
    }

    renderDocumentList() {
        const container = document.getElementById('documentList');
        container.innerHTML = '';

        if (this.documents.length === 0) {
            container.innerHTML = '<p class="text-muted text-center">No documents uploaded</p>';
            return;
        }

        // Get currently selected document
        const selectedDocId = document.getElementById('documentSelect').value;

        this.documents.forEach(doc => {
            const item = document.createElement('div');
            item.className = 'document-item';
            
            // Add selected class if this document is currently selected
            if (doc.id === selectedDocId) {
                item.classList.add('selected');
            }
            
            item.innerHTML = `
                <div class="document-info">
                    <div>
                        <div class="document-name">${doc.filename}</div>
                        <div class="document-meta">
                            ${this.formatFileSize(doc.size)} • ${doc.lines || 'N/A'} lines • ${doc.chars || 'N/A'} chars
                        </div>
                    </div>
                    <div class="document-actions">
                        <button class="btn btn-sm btn-outline-primary me-1" onclick="dashboard.selectDocument('${doc.id}')" title="Select document">
                            <i class="fas fa-check"></i>
                        </button>
                        <button class="btn btn-sm btn-outline-danger" onclick="dashboard.deleteDocument('${doc.id}')" title="Delete document">
                            <i class="fas fa-trash"></i>
                        </button>
                    </div>
                </div>
            `;
            
            // Make the entire item clickable to select the document
            item.addEventListener('click', (e) => {
                // Don't trigger if clicking on buttons
                if (!e.target.closest('.document-actions')) {
                    this.selectDocument(doc.id);
                }
            });
            
            container.appendChild(item);
        });
    }

    renderTokenizerList() {
        const container = document.getElementById('tokenizerList');
        container.innerHTML = '';

        this.tokenizers.forEach(tokenizer => {
            const item = document.createElement('div');
            item.className = 'tokenizer-item';
            
            // Auto-select mock tokenizer since it works without Python dependencies
            const isChecked = tokenizer.id === 'mock' || tokenizer.enabled;
            
            item.innerHTML = `
                <input type="checkbox" id="tokenizer_${tokenizer.id}" value="${tokenizer.id}" ${isChecked ? 'checked' : ''}>
                <div class="tokenizer-info">
                    <div class="tokenizer-name">${tokenizer.name}</div>
                    <div class="tokenizer-type">${tokenizer.type}</div>
                </div>
            `;
            container.appendChild(item);
        });
        
        // Ensure at least one tokenizer is selected
        this.ensureTokenizerSelected();
    }

    ensureTokenizerSelected() {
        const checkboxes = document.querySelectorAll('#tokenizerList input[type="checkbox"]:checked');
        if (checkboxes.length === 0) {
            // If no tokenizers are selected, select the mock tokenizer
            const mockCheckbox = document.getElementById('tokenizer_mock');
            if (mockCheckbox) {
                mockCheckbox.checked = true;
            } else {
                // If mock doesn't exist, select the first available tokenizer
                const firstCheckbox = document.querySelector('#tokenizerList input[type="checkbox"]');
                if (firstCheckbox) {
                    firstCheckbox.checked = true;
                }
            }
        }
    }

    renderAnalysisResults(analysis) {
        const container = document.getElementById('analysisResults');
        
        if (!analysis || !analysis.results || analysis.results.length === 0) {
            container.innerHTML = `
                <div class="text-center text-muted">
                    <i class="fas fa-chart-line fa-3x mb-3"></i>
                    <p>No analysis results available</p>
                </div>
            `;
            return;
        }
        
        let html = `
            <div class="analysis-result">
                <div class="result-header">
                    <div class="result-title">Analysis Results</div>
                    <small class="text-muted">${new Date(analysis.timestamp).toLocaleString()}</small>
                </div>
                <div class="result-metrics">
        `;

        analysis.results.forEach(result => {
            // Get key metrics from the result
            const tokenCount = result.token_count || 0;
            const entropy = result.metrics && result.metrics["entropy_global_entropy"] ? 
                result.metrics["entropy_global_entropy"].value : 0;
            const compressionRatio = result.metrics && result.metrics["compression_compression_ratio"] ? 
                result.metrics["compression_compression_ratio"].value : 0;
            
            html += `
                <div class="metric-item">
                    <div class="metric-value">${result.tokenizer_name || 'Unknown'}</div>
                    <div class="metric-label">Tokenizer</div>
                </div>
                <div class="metric-item">
                    <div class="metric-value">${tokenCount}</div>
                    <div class="metric-label">Tokens</div>
                </div>
                <div class="metric-item">
                    <div class="metric-value">${entropy.toFixed(3)}</div>
                    <div class="metric-label">Entropy</div>
                </div>
                <div class="metric-item">
                    <div class="metric-value">${compressionRatio.toFixed(3)}</div>
                    <div class="metric-label">Compression</div>
                </div>
            `;
        });

        html += `
                </div>
                <div class="mt-3">
                    <button class="btn btn-sm btn-outline-primary" onclick="dashboard.showDetailedResults()">
                        <i class="fas fa-chart-bar me-1"></i>View Details
                    </button>
                </div>
            </div>
        `;

        container.innerHTML = html;
    }

    renderVisualization(visualization, type) {
        const container = document.getElementById('visualizationContainer');
        
        if (visualization.filepath) {
            const iframe = document.createElement('iframe');
            iframe.src = visualization.filepath;
            iframe.style.width = '100%';
            iframe.style.height = '400px';
            iframe.style.border = 'none';
            iframe.style.borderRadius = '0.375rem';
            
            container.innerHTML = '';
            container.appendChild(iframe);
        } else if (visualization.data) {
            // Handle Plotly.js data
            Plotly.newPlot(container, visualization.data, visualization.layout || {});
        }
    }

    updateDocumentSelect() {
        const select = document.getElementById('documentSelect');
        select.innerHTML = '<option value="">Select a document...</option>';
        
        this.documents.forEach(doc => {
            const option = document.createElement('option');
            option.value = doc.id;
            option.textContent = doc.filename;
            select.appendChild(option);
        });
    }

    getSelectedTokenizers() {
        const checkboxes = document.querySelectorAll('#tokenizerList input[type="checkbox"]:checked');
        return Array.from(checkboxes).map(cb => cb.value);
    }

    getSelectedMetrics() {
        const metrics = [];
        if (document.getElementById('metricEntropy').checked) metrics.push('entropy');
        if (document.getElementById('metricCompression').checked) metrics.push('compression');
        if (document.getElementById('metricReuse').checked) metrics.push('reuse');
        if (document.getElementById('metricDrift').checked) metrics.push('drift');
        return metrics;
    }

    async deleteDocument(docId) {
        if (!confirm('Are you sure you want to delete this document?')) {
            return;
        }

        try {
            const response = await fetch(`/api/v1/documents/${docId}`, {
                method: 'DELETE'
            });

            if (response.ok) {
                this.showAlert('Document deleted successfully', 'success');
                await this.loadDocuments();
            } else {
                this.showAlert('Failed to delete document', 'danger');
            }
        } catch (error) {
            console.error('Delete error:', error);
            this.showAlert('Failed to delete document', 'danger');
        }
    }

    showDetailedResults() {
        if (!this.currentAnalysis) return;

        const modal = new bootstrap.Modal(document.getElementById('resultsModal'));
        const content = document.getElementById('modalContent');
        
        let html = `
            <div class="row">
                <div class="col-md-6">
                    <h6>Analysis Summary</h6>
                    <table class="table table-sm">
                        <tr><td>Document ID:</td><td>${this.currentAnalysis.document_id}</td></tr>
                        <tr><td>Timestamp:</td><td>${new Date(this.currentAnalysis.timestamp).toLocaleString()}</td></tr>
                        <tr><td>Tokenizers:</td><td>${this.currentAnalysis.results.length}</td></tr>
                    </table>
                    
                    <h6 class="mt-4">Detailed Results</h6>
                    <div class="table-responsive">
                        <table class="table table-sm">
                            <thead>
                                <tr>
                                    <th>Tokenizer</th>
                                    <th>Tokens</th>
                                    <th>Entropy</th>
                                    <th>Compression</th>
                                </tr>
                            </thead>
                            <tbody>
        `;

        this.currentAnalysis.results.forEach(result => {
            const tokenCount = result.token_count || 0;
            const entropy = result.metrics && result.metrics["entropy_global_entropy"] ? 
                result.metrics["entropy_global_entropy"].value : 0;
            const compressionRatio = result.metrics && result.metrics["compression_compression_ratio"] ? 
                result.metrics["compression_compression_ratio"].value : 0;
            
            html += `
                <tr>
                    <td>${result.tokenizer_name || 'Unknown'}</td>
                    <td>${tokenCount}</td>
                    <td>${entropy.toFixed(3)}</td>
                    <td>${compressionRatio.toFixed(3)}</td>
                </tr>
            `;
        });

        html += `
                            </tbody>
                        </table>
                    </div>
                </div>
                <div class="col-md-6">
                    <h6>Results Overview</h6>
                    <div id="resultsChart" style="height: 300px;"></div>
                </div>
            </div>
        `;

        content.innerHTML = html;
        modal.show();

        // Generate chart
        this.generateResultsChart();
    }

    generateResultsChart() {
        if (!this.currentAnalysis) return;

        const data = this.currentAnalysis.results.map(result => {
            const entropy = result.metrics && result.metrics["entropy_global_entropy"] ? 
                result.metrics["entropy_global_entropy"].value : 0;
            
            return {
                x: [result.tokenizer_name || 'Unknown'],
                y: [entropy],
                type: 'bar',
                name: 'Entropy'
            };
        });

        const layout = {
            title: 'Entropy Comparison',
            xaxis: { title: 'Tokenizer' },
            yaxis: { title: 'Entropy' }
        };

        Plotly.newPlot('resultsChart', data, layout);
    }

    async exportResults() {
        if (!this.currentAnalysis) {
            this.showAlert('No results to export', 'warning');
            return;
        }

        const dataStr = JSON.stringify(this.currentAnalysis, null, 2);
        const dataBlob = new Blob([dataStr], { type: 'application/json' });
        
        const link = document.createElement('a');
        link.href = URL.createObjectURL(dataBlob);
        link.download = `analysis_${this.currentAnalysis.id}.json`;
        link.click();
    }

    showAlert(message, type) {
        const alertDiv = document.createElement('div');
        alertDiv.className = `alert alert-${type} alert-dismissible fade show`;
        alertDiv.innerHTML = `
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert"></button>
        `;

        const container = document.querySelector('.container-fluid');
        container.insertBefore(alertDiv, container.firstChild);

        // Auto-dismiss after 5 seconds
        setTimeout(() => {
            if (alertDiv.parentNode) {
                alertDiv.remove();
            }
        }, 5000);
    }

    showLoading(show) {
        const overlay = document.getElementById('loadingOverlay');
        overlay.style.display = show ? 'flex' : 'none';
    }

    formatFileSize(bytes) {
        if (bytes === 0) return '0 Bytes';
        const k = 1024;
        const sizes = ['Bytes', 'KB', 'MB', 'GB'];
        const i = Math.floor(Math.log(bytes) / Math.log(k));
        return parseFloat((bytes / Math.pow(k, i)).toFixed(2)) + ' ' + sizes[i];
    }

    setupWebSocket() {
        // TODO: Implement WebSocket for real-time updates
        console.log('WebSocket setup - to be implemented');
    }

    selectDocument(documentId) {
        // Set the document select dropdown to the selected document
        const select = document.getElementById('documentSelect');
        select.value = documentId;
        
        // Trigger the change event to update the UI
        this.onDocumentSelect(documentId);
        
        // Re-render the document list to show the selection
        this.renderDocumentList();
    }

    initializeUIState() {
        // Disable the analysis form submit button initially
        const analysisForm = document.getElementById('analysisForm');
        const submitButton = analysisForm.querySelector('button[type="submit"]');
        submitButton.disabled = true;
        
        // Clear any validation classes from the document select
        const documentSelect = document.getElementById('documentSelect');
        documentSelect.classList.remove('is-valid', 'is-invalid');
    }

    onDocumentSelect(documentId) {
        // Handle document selection change
        console.log('Document selected:', documentId);
        
        // Get the select element
        const select = document.getElementById('documentSelect');
        
        if (documentId) {
            // Find the selected document
            const selectedDoc = this.documents.find(doc => doc.id === documentId);
            
            if (selectedDoc) {
                // Show success feedback
                this.showAlert(`Document "${selectedDoc.filename}" selected successfully!`, 'success');
                
                // Enable the analysis form if it was disabled
                const analysisForm = document.getElementById('analysisForm');
                const submitButton = analysisForm.querySelector('button[type="submit"]');
                submitButton.disabled = false;
                
                // Update the select element styling to show it's selected
                select.classList.add('is-valid');
                select.classList.remove('is-invalid');
                
                // Optionally, you could load document details here
                // this.loadDocumentDetails(documentId);
            }
        } else {
            // No document selected
            select.classList.remove('is-valid', 'is-invalid');
            
            // Disable the analysis form
            const analysisForm = document.getElementById('analysisForm');
            const submitButton = analysisForm.querySelector('button[type="submit"]');
            submitButton.disabled = true;
        }
    }
}

// Initialize dashboard when page loads
let dashboard;
document.addEventListener('DOMContentLoaded', () => {
    dashboard = new Dashboard();
}); 