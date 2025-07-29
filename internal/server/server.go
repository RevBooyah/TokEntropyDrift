package server

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/RevBooyah/TokEntropyDrift/internal/config"
	"github.com/RevBooyah/TokEntropyDrift/internal/loader"
	"github.com/RevBooyah/TokEntropyDrift/internal/metrics"
	"github.com/RevBooyah/TokEntropyDrift/internal/tokenizers"
	"github.com/RevBooyah/TokEntropyDrift/internal/visualization"
	"github.com/gorilla/mux"
)

// Server represents the web dashboard server
type Server struct {
	config            *config.Config
	router            *mux.Router
	tokenizerRegistry *tokenizers.TokenizerRegistry
	metricsEngine     *metrics.Engine
	vizEngine         *visualization.VisualizationEngine
	uploadDir         string
	sessions          map[string]*Session
}

// Session represents a user session
type Session struct {
	ID       string
	Created  time.Time
	LastSeen time.Time
	Uploads  []string
	Analyses []string
}

// AnalysisRequest represents a request for analysis
type AnalysisRequest struct {
	DocumentID   string   `json:"document_id"`
	TokenizerIDs []string `json:"tokenizer_ids"`
	Metrics      []string `json:"metrics"`
}

// AnalysisResponse represents the response from analysis
type AnalysisResponse struct {
	ID             string                               `json:"id"`
	DocumentID     string                               `json:"document_id"`
	Results        []*metrics.AnalysisResult            `json:"results"`
	Visualizations []*visualization.VisualizationResult `json:"visualizations"`
	Timestamp      time.Time                            `json:"timestamp"`
}

// NewServer creates a new web server instance
func NewServer(cfg *config.Config) *Server {
	// Create upload directory
	uploadDir := filepath.Join(cfg.Output.Directory, "uploads")
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		log.Fatalf("Failed to create upload directory: %v", err)
	}

	// Register all available tokenizers with the global registry
	if err := tokenizers.RegisterAllTokenizers(); err != nil {
		log.Printf("Warning: Failed to register some tokenizers: %v", err)
	}

	metricsEngine := metrics.NewEngine(metrics.EngineConfig{
		EntropyWindowSize: cfg.Analysis.EntropyWindowSize,
		NormalizeEntropy:  cfg.Analysis.NormalizeEntropy,
		CompressionRatio:  cfg.Analysis.CompressionRatio,
		DriftDetection:    cfg.Analysis.DriftDetection,
	})
	vizEngine := visualization.NewVisualizationEngine(visualization.VisualizationConfig{
		Theme:       cfg.Visualization.Theme,
		ImageSize:   cfg.Visualization.ImageSize,
		FileType:    cfg.Visualization.FileType,
		Interactive: cfg.Visualization.Interactive,
		OutputDir:   filepath.Join(cfg.Output.Directory, "visualizations"),
	})

	server := &Server{
		config:            cfg,
		router:            mux.NewRouter(),
		tokenizerRegistry: tokenizers.GlobalRegistry,
		metricsEngine:     metricsEngine,
		vizEngine:         vizEngine,
		uploadDir:         uploadDir,
		sessions:          make(map[string]*Session),
	}

	server.setupRoutes()
	return server
}

// setupRoutes configures all the HTTP routes
func (s *Server) setupRoutes() {
	// Static file serving
	s.router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	s.router.PathPrefix("/uploads/").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir(s.uploadDir))))
	s.router.PathPrefix("/visualizations/").Handler(http.StripPrefix("/visualizations/", http.FileServer(http.Dir(filepath.Join(s.config.Output.Directory, "visualizations")))))

	// API routes
	api := s.router.PathPrefix("/api/v1").Subrouter()

	// File upload and management
	api.HandleFunc("/upload", s.handleFileUpload).Methods("POST")
	api.HandleFunc("/documents", s.handleListDocuments).Methods("GET")
	api.HandleFunc("/documents/{id}", s.handleGetDocument).Methods("GET")
	api.HandleFunc("/documents/{id}", s.handleDeleteDocument).Methods("DELETE")

	// Tokenizer management
	api.HandleFunc("/tokenizers", s.handleListTokenizers).Methods("GET")
	api.HandleFunc("/tokenizers/{id}", s.handleGetTokenizer).Methods("GET")

	// Analysis endpoints
	api.HandleFunc("/analyze", s.handleAnalyze).Methods("POST")
	api.HandleFunc("/analyses", s.handleListAnalyses).Methods("GET")
	api.HandleFunc("/analyses/{id}", s.handleGetAnalysis).Methods("GET")

	// Visualization endpoints
	api.HandleFunc("/visualizations/heatmap", s.handleGenerateHeatmap).Methods("POST")
	api.HandleFunc("/visualizations/drift", s.handleGenerateDriftViz).Methods("POST")
	api.HandleFunc("/visualizations/entropy", s.handleGenerateEntropyViz).Methods("POST")

	// Session management
	api.HandleFunc("/session", s.handleGetSession).Methods("GET")
	api.HandleFunc("/session", s.handleCreateSession).Methods("POST")

	// WebSocket for real-time updates
	api.HandleFunc("/ws", s.handleWebSocket)

	// Main dashboard route
	s.router.HandleFunc("/", s.handleDashboard).Methods("GET")
	s.router.HandleFunc("/dashboard", s.handleDashboard).Methods("GET")
	s.router.HandleFunc("/compare", s.handleCompareView).Methods("GET")
	s.router.HandleFunc("/visualize", s.handleVisualizeView).Methods("GET")
}

// Start starts the web server
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.config.Server.Host, s.config.Server.Port)
	log.Printf("Starting TokEntropyDrift dashboard server on %s", addr)
	return http.ListenAndServe(addr, s.router)
}

// handleDashboard serves the main dashboard page
func (s *Server) handleDashboard(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/dashboard.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":  "TokEntropyDrift Dashboard",
		"Config": s.config,
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

// handleCompareView serves the comparison view page
func (s *Server) handleCompareView(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/compare.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":  "Tokenizer Comparison",
		"Config": s.config,
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

// handleVisualizeView serves the visualization view page
func (s *Server) handleVisualizeView(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/visualize.html")
	if err != nil {
		http.Error(w, "Template not found", http.StatusInternalServerError)
		return
	}

	data := map[string]interface{}{
		"Title":  "Visualization Studio",
		"Config": s.config,
	}

	w.Header().Set("Content-Type", "text/html")
	tmpl.Execute(w, data)
}

// handleFileUpload handles file uploads
func (s *Server) handleFileUpload(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form
	if err := r.ParseMultipartForm(32 << 20); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "No file uploaded", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), header.Filename)
	filepath := filepath.Join(s.uploadDir, filename)

	// Create file
	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Failed to create file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy uploaded file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	// Load and validate document
	docLoader := loader.NewLoader(s.config.Input.FileType)
	documents, err := docLoader.LoadDocuments(filepath)
	if err != nil {
		os.Remove(filepath) // Clean up invalid file
		http.Error(w, fmt.Sprintf("Invalid file: %v", err), http.StatusBadRequest)
		return
	}

	if len(documents) == 0 {
		os.Remove(filepath) // Clean up empty file
		http.Error(w, "File contains no valid documents", http.StatusBadRequest)
		return
	}

	response := map[string]interface{}{
		"id":       filename,
		"filename": filename,
		"size":     header.Size,
		"type":     header.Header.Get("Content-Type"),
		"lines":    len(documents),
		"chars":    len(documents[0].Content),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleListDocuments lists uploaded documents
func (s *Server) handleListDocuments(w http.ResponseWriter, r *http.Request) {
	files, err := os.ReadDir(s.uploadDir)
	if err != nil {
		http.Error(w, "Failed to read upload directory", http.StatusInternalServerError)
		return
	}

	var documents []map[string]interface{}
	for _, file := range files {
		if !file.IsDir() {
			info, err := file.Info()
			if err != nil {
				continue
			}

			documents = append(documents, map[string]interface{}{
				"id":       strings.TrimSuffix(file.Name(), filepath.Ext(file.Name())),
				"filename": file.Name(),
				"size":     info.Size(),
				"modified": info.ModTime(),
			})
		}
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents)
}

// handleGetDocument retrieves a specific document
func (s *Server) handleGetDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	// Find file by ID
	files, err := os.ReadDir(s.uploadDir)
	if err != nil {
		http.Error(w, "Failed to read upload directory", http.StatusInternalServerError)
		return
	}

	var documents []loader.Document
	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), docID) {
			filepath := filepath.Join(s.uploadDir, file.Name())
			docLoader := loader.NewLoader(s.config.Input.FileType)
			documents, err = docLoader.LoadDocuments(filepath)
			if err != nil {
				http.Error(w, "Failed to load document", http.StatusInternalServerError)
				return
			}
			break
		}
	}

	if len(documents) == 0 {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(documents[0])
}

// handleDeleteDocument deletes a document
func (s *Server) handleDeleteDocument(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	docID := vars["id"]

	// Find and delete file
	files, err := os.ReadDir(s.uploadDir)
	if err != nil {
		http.Error(w, "Failed to read upload directory", http.StatusInternalServerError)
		return
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), docID) {
			filepath := filepath.Join(s.uploadDir, file.Name())
			if err := os.Remove(filepath); err != nil {
				http.Error(w, "Failed to delete file", http.StatusInternalServerError)
				return
			}
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}

	http.Error(w, "Document not found", http.StatusNotFound)
}

// handleListTokenizers lists available tokenizers
func (s *Server) handleListTokenizers(w http.ResponseWriter, r *http.Request) {
	availableTokenizers := tokenizers.GetAvailableTokenizers()

	var response []map[string]interface{}
	for _, tokenizerID := range availableTokenizers {
		response = append(response, map[string]interface{}{
			"id":          tokenizerID,
			"name":        tokenizerID,
			"type":        tokenizers.GetTokenizerType(tokenizerID),
			"description": tokenizers.GetTokenizerDescription(tokenizerID),
			"enabled":     true,
		})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleGetTokenizer retrieves tokenizer details
func (s *Server) handleGetTokenizer(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	tokenizerID := vars["id"]

	if !tokenizers.ValidateTokenizerName(tokenizerID) {
		http.Error(w, "Tokenizer not found", http.StatusNotFound)
		return
	}

	tokenizer := map[string]interface{}{
		"id":           tokenizerID,
		"name":         tokenizerID,
		"type":         tokenizers.GetTokenizerType(tokenizerID),
		"description":  tokenizers.GetTokenizerDescription(tokenizerID),
		"backend":      tokenizers.GetTokenizerBackend(tokenizerID),
		"requirements": tokenizers.GetTokenizerRequirements(tokenizerID),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tokenizer)
}

// handleAnalyze performs analysis on uploaded documents
func (s *Server) handleAnalyze(w http.ResponseWriter, r *http.Request) {
	var req AnalysisRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	log.Printf("Analysis request: DocumentID=%s, TokenizerIDs=%v, Metrics=%v", req.DocumentID, req.TokenizerIDs, req.Metrics)

	// Load document
	documents, err := s.loadDocumentByID(req.DocumentID)
	if err != nil {
		log.Printf("Failed to load document %s: %v", req.DocumentID, err)
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	document := documents[0].Content
	log.Printf("Loaded document with %d characters", len(document))

	// Perform analysis
	results := make([]*metrics.AnalysisResult, 0)
	ctx := context.Background()

	for _, tokenizerID := range req.TokenizerIDs {
		log.Printf("Processing tokenizer: %s", tokenizerID)

		if !tokenizers.ValidateTokenizerName(tokenizerID) {
			log.Printf("Invalid tokenizer name: %s", tokenizerID)
			continue
		}

		// Get tokenizer from registry
		tokenizer, err := s.tokenizerRegistry.Get(tokenizerID)
		if err != nil {
			log.Printf("Tokenizer %s not found in registry, creating new one: %v", tokenizerID, err)
			// Try to create and register the tokenizer
			tokenizer, err = s.createTokenizer(tokenizerID)
			if err != nil {
				log.Printf("Failed to create tokenizer %s: %v", tokenizerID, err)
				continue
			}
		}

		log.Printf("Using tokenizer: %s", tokenizer.Name())

		// Analyze document
		result, err := s.metricsEngine.AnalyzeDocument(ctx, document, tokenizer)
		if err != nil {
			log.Printf("Failed to analyze document with tokenizer %s: %v", tokenizerID, err)
			continue
		}

		log.Printf("Analysis successful for tokenizer %s: %d tokens", tokenizerID, result.TokenCount)
		results = append(results, result)
	}

	log.Printf("Analysis completed with %d results", len(results))

	// Generate visualizations
	visualizations := make([]*visualization.VisualizationResult, 0)
	for _, result := range results {
		// Generate heatmap
		heatmapData := visualization.HeatmapData{
			XLabels:    []string{"Tokens", "Entropy", "Compression"},
			YLabels:    []string{result.TokenizerName},
			Values:     [][]float64{{float64(result.TokenCount), result.Metrics["entropy_shannon"].Value, result.Metrics["compression_ratio"].Value}},
			ColorScale: "Viridis",
			Title:      "Analysis Results",
		}

		viz, err := s.vizEngine.GenerateHeatmap(heatmapData, "entropy")
		if err == nil {
			visualizations = append(visualizations, viz)
		}
	}

	response := AnalysisResponse{
		ID:             fmt.Sprintf("analysis_%d", time.Now().Unix()),
		DocumentID:     req.DocumentID,
		Results:        results,
		Visualizations: visualizations,
		Timestamp:      time.Now(),
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// loadDocumentByID loads a document by its ID
func (s *Server) loadDocumentByID(docID string) ([]loader.Document, error) {
	files, err := os.ReadDir(s.uploadDir)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		if !file.IsDir() && strings.HasPrefix(file.Name(), docID) {
			filepath := filepath.Join(s.uploadDir, file.Name())
			docLoader := loader.NewLoader(s.config.Input.FileType)
			return docLoader.LoadDocuments(filepath)
		}
	}

	return nil, fmt.Errorf("document not found")
}

// createTokenizer creates and registers a new tokenizer
func (s *Server) createTokenizer(tokenizerID string) (tokenizers.Tokenizer, error) {
	var tokenizer tokenizers.Tokenizer

	switch tokenizerID {
	case "mock":
		tokenizer = tokenizers.NewMockTokenizer(tokenizerID)
	case "gpt2":
		// For now, fall back to mock tokenizer for GPT-2
		// In a real implementation, you would create the actual GPT-2 tokenizer
		tokenizer = tokenizers.NewMockTokenizer(tokenizerID)
	default:
		// For any other tokenizer, try to create a mock tokenizer as fallback
		tokenizer = tokenizers.NewMockTokenizer(tokenizerID)
	}

	// Initialize the tokenizer with default config
	config := tokenizers.TokenizerConfig{
		Name: tokenizerID,
		Type: tokenizers.GetTokenizerType(tokenizerID),
		Parameters: map[string]string{
			"vocab_size": "1000",
		},
	}

	if err := tokenizer.Initialize(config); err != nil {
		return nil, fmt.Errorf("failed to initialize tokenizer %s: %w", tokenizerID, err)
	}

	// Register the tokenizer with the registry
	if err := s.tokenizerRegistry.Register(tokenizerID, tokenizer); err != nil {
		return nil, fmt.Errorf("failed to register tokenizer %s: %w", tokenizerID, err)
	}

	return tokenizer, nil
}

// handleListAnalyses lists previous analyses
func (s *Server) handleListAnalyses(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement analysis storage and retrieval
	analyses := []map[string]interface{}{}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(analyses)
}

// handleGetAnalysis retrieves a specific analysis
func (s *Server) handleGetAnalysis(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	_ = vars["id"] // analysisID

	// TODO: Implement analysis storage and retrieval
	http.Error(w, "Analysis not found", http.StatusNotFound)
}

// handleGenerateHeatmap generates heatmap visualizations
func (s *Server) handleGenerateHeatmap(w http.ResponseWriter, r *http.Request) {
	var req struct {
		DocumentID string   `json:"document_id"`
		Tokenizers []string `json:"tokenizers"`
		Type       string   `json:"type"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Load document
	documents, err := s.loadDocumentByID(req.DocumentID)
	if err != nil {
		http.Error(w, "Document not found", http.StatusNotFound)
		return
	}

	document := documents[0].Content

	// Generate heatmap data
	var xLabels []string
	var yLabels []string
	var values [][]float64

	ctx := context.Background()
	for _, tokenizerID := range req.Tokenizers {
		if !tokenizers.ValidateTokenizerName(tokenizerID) {
			continue
		}

		tokenizer, err := s.tokenizerRegistry.Get(tokenizerID)
		if err != nil {
			continue
		}

		result, err := s.metricsEngine.AnalyzeDocument(ctx, document, tokenizer)
		if err != nil {
			continue
		}

		yLabels = append(yLabels, tokenizerID)
		if len(xLabels) == 0 {
			xLabels = []string{"Token Count", "Entropy", "Compression Ratio"}
		}
		values = append(values, []float64{
			float64(result.TokenCount),
			result.Metrics["entropy_shannon"].Value,
			result.Metrics["compression_ratio"].Value,
		})
	}

	heatmapData := visualization.HeatmapData{
		XLabels:    xLabels,
		YLabels:    yLabels,
		Values:     values,
		ColorScale: "Viridis",
		Title:      fmt.Sprintf("Analysis Heatmap - %s", req.Type),
	}

	viz, err := s.vizEngine.GenerateHeatmap(heatmapData, req.Type)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to generate heatmap: %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(viz)
}

// handleGenerateDriftViz generates drift visualizations
func (s *Server) handleGenerateDriftViz(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement drift visualization generation
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// handleGenerateEntropyViz generates entropy visualizations
func (s *Server) handleGenerateEntropyViz(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement entropy visualization generation
	http.Error(w, "Not implemented", http.StatusNotImplemented)
}

// handleGetSession retrieves or creates a user session
func (s *Server) handleGetSession(w http.ResponseWriter, r *http.Request) {
	sessionID := r.URL.Query().Get("session_id")

	var session *Session
	if sessionID != "" {
		session = s.sessions[sessionID]
	}

	if session == nil {
		// Create new session
		sessionID = fmt.Sprintf("session_%d", time.Now().Unix())
		session = &Session{
			ID:       sessionID,
			Created:  time.Now(),
			LastSeen: time.Now(),
			Uploads:  []string{},
			Analyses: []string{},
		}
		s.sessions[sessionID] = session
	} else {
		session.LastSeen = time.Now()
	}

	response := map[string]interface{}{
		"session_id": session.ID,
		"created":    session.Created,
		"last_seen":  session.LastSeen,
		"uploads":    session.Uploads,
		"analyses":   session.Analyses,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleCreateSession creates a new session
func (s *Server) handleCreateSession(w http.ResponseWriter, r *http.Request) {
	sessionID := fmt.Sprintf("session_%d", time.Now().Unix())
	session := &Session{
		ID:       sessionID,
		Created:  time.Now(),
		LastSeen: time.Now(),
		Uploads:  []string{},
		Analyses: []string{},
	}
	s.sessions[sessionID] = session

	response := map[string]interface{}{
		"session_id": session.ID,
		"created":    session.Created,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// handleWebSocket handles WebSocket connections for real-time updates
func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	// TODO: Implement WebSocket support for real-time updates
	http.Error(w, "WebSocket not implemented", http.StatusNotImplemented)
}
