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
        this.setupEventListeners();
        this.setupWebSocket();
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
            this.showAlert('Please select at least one tokenizer', 'warning');
            return;
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

        this.documents.forEach(doc => {
            const item = document.createElement('div');
            item.className = 'document-item';
            item.innerHTML = `
                <div class="document-info">
                    <div>
                        <div class="document-name">${doc.filename}</div>
                        <div class="document-meta">
                            ${this.formatFileSize(doc.size)} • ${doc.lines} lines • ${doc.chars} chars
                        </div>
                    </div>
                    <button class="btn btn-sm btn-outline-danger" onclick="dashboard.deleteDocument('${doc.id}')">
                        <i class="fas fa-trash"></i>
                    </button>
                </div>
            `;
            container.appendChild(item);
        });
    }

    renderTokenizerList() {
        const container = document.getElementById('tokenizerList');
        container.innerHTML = '';

        this.tokenizers.forEach(tokenizer => {
            const item = document.createElement('div');
            item.className = 'tokenizer-item';
            item.innerHTML = `
                <input type="checkbox" id="tokenizer_${tokenizer.id}" value="${tokenizer.id}" ${tokenizer.enabled ? 'checked' : ''}>
                <div class="tokenizer-info">
                    <div class="tokenizer-name">${tokenizer.name}</div>
                    <div class="tokenizer-type">${tokenizer.type}</div>
                </div>
            `;
            container.appendChild(item);
        });
    }

    renderAnalysisResults(analysis) {
        const container = document.getElementById('analysisResults');
        
        let html = `
            <div class="analysis-result">
                <div class="result-header">
                    <div class="result-title">Analysis Results</div>
                    <small class="text-muted">${new Date(analysis.timestamp).toLocaleString()}</small>
                </div>
                <div class="result-metrics">
        `;

        analysis.results.forEach(result => {
            html += `
                <div class="metric-item">
                    <div class="metric-value">${result.tokenizer_id}</div>
                    <div class="metric-label">Tokenizer</div>
                </div>
                <div class="metric-item">
                    <div class="metric-value">${result.token_count}</div>
                    <div class="metric-label">Tokens</div>
                </div>
                <div class="metric-item">
                    <div class="metric-value">${result.entropy.toFixed(3)}</div>
                    <div class="metric-label">Entropy</div>
                </div>
                <div class="metric-item">
                    <div class="metric-value">${result.compression_ratio.toFixed(3)}</div>
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

        const data = this.currentAnalysis.results.map(result => ({
            x: [result.tokenizer_id],
            y: [result.entropy],
            type: 'bar',
            name: 'Entropy'
        }));

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

    onDocumentSelect(documentId) {
        // Handle document selection change
        console.log('Document selected:', documentId);
    }
}

// Initialize dashboard when page loads
let dashboard;
document.addEventListener('DOMContentLoaded', () => {
    dashboard = new Dashboard();
}); 