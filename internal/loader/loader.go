package loader

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Document represents a single document with metadata
type Document struct {
	Content    string            `json:"content"`
	LineNumber int               `json:"line_number"`
	FilePath   string            `json:"file_path"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

// Loader handles loading documents from various file formats
type Loader struct {
	fileType string
}

// NewLoader creates a new loader for the specified file type
func NewLoader(fileType string) *Loader {
	return &Loader{
		fileType: strings.ToLower(fileType),
	}
}

// LoadDocuments loads all documents from the given file path
func (l *Loader) LoadDocuments(filePath string) ([]Document, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("error opening file %s: %w", filePath, err)
	}
	defer file.Close()

	switch l.fileType {
	case "txt", "text":
		return l.loadTextFile(file, filePath)
	case "jsonl", "json":
		return l.loadJSONLFile(file, filePath)
	case "csv":
		return l.loadCSVFile(file, filePath)
	default:
		return nil, fmt.Errorf("unsupported file type: %s", l.fileType)
	}
}

// loadTextFile loads documents from a plain text file
func (l *Loader) loadTextFile(file *os.File, filePath string) ([]Document, error) {
	var documents []Document
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines
		if line == "" {
			continue
		}

		doc := Document{
			Content:    line,
			LineNumber: lineNumber,
			FilePath:   filePath,
			Metadata: map[string]string{
				"file_type": "text",
				"file_name": filepath.Base(filePath),
			},
		}
		documents = append(documents, doc)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading text file: %w", err)
	}

	return documents, nil
}

// loadJSONLFile loads documents from a JSONL (JSON Lines) file
func (l *Loader) loadJSONLFile(file *os.File, filePath string) ([]Document, error) {
	var documents []Document
	scanner := bufio.NewScanner(file)
	lineNumber := 0

	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(scanner.Text())
		
		// Skip empty lines
		if line == "" {
			continue
		}

		// Parse JSON line
		var jsonData map[string]interface{}
		if err := json.Unmarshal([]byte(line), &jsonData); err != nil {
			return nil, fmt.Errorf("error parsing JSON at line %d: %w", lineNumber, err)
		}

		// Extract content field (default to "text" or "content")
		content, ok := jsonData["text"].(string)
		if !ok {
			content, ok = jsonData["content"].(string)
		}
		if !ok {
			// If no text/content field, use the entire JSON as string
			content = line
		}

		// Extract metadata
		metadata := make(map[string]string)
		for k, v := range jsonData {
			if k != "text" && k != "content" {
				if str, ok := v.(string); ok {
					metadata[k] = str
				} else {
					metadata[k] = fmt.Sprintf("%v", v)
				}
			}
		}
		metadata["file_type"] = "jsonl"
		metadata["file_name"] = filepath.Base(filePath)

		doc := Document{
			Content:    content,
			LineNumber: lineNumber,
			FilePath:   filePath,
			Metadata:   metadata,
		}
		documents = append(documents, doc)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading JSONL file: %w", err)
	}

	return documents, nil
}

// loadCSVFile loads documents from a CSV file
func (l *Loader) loadCSVFile(file *os.File, filePath string) ([]Document, error) {
	var documents []Document
	reader := csv.NewReader(file)
	lineNumber := 0

	// Read header
	header, err := reader.Read()
	if err != nil {
		return nil, fmt.Errorf("error reading CSV header: %w", err)
	}
	lineNumber++

	// Find content column (default to "text" or "content")
	contentColIndex := -1
	for i, col := range header {
		if col == "text" || col == "content" {
			contentColIndex = i
			break
		}
	}
	if contentColIndex == -1 {
		// Use first column as content if no text/content column found
		contentColIndex = 0
	}

	// Read data rows
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("error reading CSV row %d: %w", lineNumber+1, err)
		}
		lineNumber++

		if len(record) == 0 {
			continue
		}

		// Extract content
		content := record[contentColIndex]

		// Extract metadata from other columns
		metadata := make(map[string]string)
		for i, value := range record {
			if i != contentColIndex && i < len(header) {
				metadata[header[i]] = value
			}
		}
		metadata["file_type"] = "csv"
		metadata["file_name"] = filepath.Base(filePath)

		doc := Document{
			Content:    content,
			LineNumber: lineNumber,
			FilePath:   filePath,
			Metadata:   metadata,
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

// GetFileType returns the detected file type based on extension
func GetFileType(filePath string) string {
	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".txt", ".text":
		return "txt"
	case ".jsonl", ".json":
		return "jsonl"
	case ".csv":
		return "csv"
	default:
		return "txt" // Default to text
	}
}

// ValidateFile checks if the file exists and is readable
func ValidateFile(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("file does not exist or is not readable: %w", err)
	}
	defer file.Close()
	return nil
} 