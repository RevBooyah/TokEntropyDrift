# Contributing to TokEntropyDrift

Thank you for your interest in contributing to TokEntropyDrift! This document provides guidelines and information for contributors.

## Table of Contents

1. [Getting Started](#getting-started)
2. [Development Setup](#development-setup)
3. [Code Style and Standards](#code-style-and-standards)
4. [Testing Guidelines](#testing-guidelines)
5. [Documentation Standards](#documentation-standards)
6. [Pull Request Process](#pull-request-process)
7. [Issue Reporting](#issue-reporting)
8. [Community Guidelines](#community-guidelines)

## Getting Started

### Before You Start

1. **Check existing issues**: Search for existing issues or discussions related to your contribution
2. **Read the documentation**: Familiarize yourself with the project structure and goals
3. **Join the community**: Participate in discussions and ask questions

### Types of Contributions

We welcome various types of contributions:

- **Bug fixes**: Fix issues and improve reliability
- **Feature additions**: Add new functionality and capabilities
- **Documentation**: Improve guides, tutorials, and API documentation
- **Examples**: Create new examples and use cases
- **Testing**: Add tests and improve test coverage
- **Performance**: Optimize code and improve efficiency
- **Plugins**: Create custom plugins and extensions

## Development Setup

### Prerequisites

- Go 1.22 or later
- Python 3.8+ (for some tokenizer adapters)
- Git
- Make (optional, for build scripts)

### Local Development Environment

1. **Fork the repository**:
   ```bash
   # Fork on GitHub, then clone your fork
   git clone https://github.com/RevBooyah/tokentropydrift.git
   cd tokentropydrift
   ```

2. **Set up the upstream remote**:
   ```bash
   git remote add upstream https://github.com/original-owner/tokentropydrift.git
   ```

3. **Create a development branch**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

4. **Install dependencies**:
   ```bash
   # Go dependencies
   go mod download
   
   # Python dependencies (if needed)
   pip install tiktoken transformers sentencepiece
   ```

5. **Build the application**:
   ```bash
   go build -o ted cmd/ted/main.go
   ```

6. **Run tests**:
   ```bash
   go test ./...
   ```

### Development Workflow

1. **Keep your fork updated**:
   ```bash
   git fetch upstream
   git checkout main
   git merge upstream/main
   ```

2. **Create feature branches**:
   ```bash
   git checkout -b feature/your-feature-name
   ```

3. **Make your changes**:
   - Write code following the style guidelines
   - Add tests for new functionality
   - Update documentation as needed

4. **Test your changes**:
   ```bash
   # Run all tests
   go test ./...
   
   # Run specific tests
   go test ./internal/metrics/...
   
   # Run with race detection
   go test -race ./...
   
   # Run benchmarks
   go test -bench=. ./...
   ```

5. **Commit your changes**:
   ```bash
   git add .
   git commit -m "feat: add new tokenizer support"
   ```

## Code Style and Standards

### Go Code Style

Follow the official Go style guide and use `gofmt`:

```bash
# Format code
gofmt -w .

# Run linter
golangci-lint run
```

### Code Organization

1. **Package structure**:
   ```
   internal/
   ├── tokenizers/     # Tokenizer implementations
   ├── metrics/        # Analysis metrics
   ├── cache/          # Caching layer
   ├── parallel/       # Parallel processing
   ├── streaming/      # Streaming analysis
   ├── plugins/        # Plugin system
   └── advanced/       # Advanced features integration
   ```

2. **File naming**:
   - Use snake_case for file names
   - Group related functionality in the same package
   - Keep files focused and not too large

3. **Function organization**:
   ```go
   // Package comment
   package tokenizers
   
   // Imports
   import (
       "context"
       "fmt"
   )
   
   // Types and interfaces
   type Tokenizer interface {
       // Interface methods
   }
   
   // Implementation
   type MockTokenizer struct {
       // Fields
   }
   
   // Methods
   func (t *MockTokenizer) Tokenize(ctx context.Context, text string) (*TokenizationResult, error) {
       // Implementation
   }
   ```

### Error Handling

1. **Use proper error handling**:
   ```go
   // Good
   result, err := tokenizer.Tokenize(ctx, text)
   if err != nil {
       return fmt.Errorf("tokenization failed: %w", err)
   }
   
   // Avoid
   result, _ := tokenizer.Tokenize(ctx, text)
   ```

2. **Create meaningful error messages**:
   ```go
   // Good
   return fmt.Errorf("failed to initialize tokenizer %s: %w", name, err)
   
   // Avoid
   return fmt.Errorf("error: %v", err)
   ```

3. **Use custom error types when appropriate**:
   ```go
   type TokenizerError struct {
       Tokenizer string
       Message   string
       Cause     error
   }
   
   func (e *TokenizerError) Error() string {
       return fmt.Sprintf("tokenizer %s error: %s", e.Tokenizer, e.Message)
   }
   
   func (e *TokenizerError) Unwrap() error {
       return e.Cause
   }
   ```

### Documentation

1. **Package documentation**:
   ```go
   // Package tokenizers provides interfaces and implementations for text tokenization.
   // It supports various tokenizer types including BPE, WordPiece, and SentencePiece.
   package tokenizers
   ```

2. **Function documentation**:
   ```go
   // Tokenize splits the input text into tokens using the specified tokenizer.
   // It returns a TokenizationResult containing the tokens and metadata.
   // The context can be used for cancellation and timeouts.
   func (t *Tokenizer) Tokenize(ctx context.Context, text string) (*TokenizationResult, error) {
       // Implementation
   }
   ```

3. **Example usage**:
   ```go
   func ExampleTokenizer_Tokenize() {
       tokenizer := NewMockTokenizer()
       result, err := tokenizer.Tokenize(context.Background(), "Hello, world!")
       if err != nil {
           log.Fatal(err)
       }
       fmt.Printf("Tokenized %d tokens\n", len(result.Tokens))
       // Output: Tokenized 3 tokens
   }
   ```

## Testing Guidelines

### Test Structure

1. **Unit tests**:
   ```go
   func TestMockTokenizer_Tokenize(t *testing.T) {
       tests := []struct {
           name     string
           input    string
           expected int
           wantErr  bool
       }{
           {
               name:     "simple text",
               input:    "Hello, world!",
               expected: 3,
               wantErr:  false,
           },
           {
               name:     "empty text",
               input:    "",
               expected: 0,
               wantErr:  false,
           },
       }
       
       for _, tt := range tests {
           t.Run(tt.name, func(t *testing.T) {
               tokenizer := NewMockTokenizer()
               result, err := tokenizer.Tokenize(context.Background(), tt.input)
               
               if (err != nil) != tt.wantErr {
                   t.Errorf("Tokenize() error = %v, wantErr %v", err, tt.wantErr)
                   return
               }
               
               if len(result.Tokens) != tt.expected {
                   t.Errorf("Tokenize() got %d tokens, want %d", len(result.Tokens), tt.expected)
               }
           })
       }
   }
   ```

2. **Integration tests**:
   ```go
   func TestIntegration_TokenizerPipeline(t *testing.T) {
       // Test complete pipeline
       engine := metrics.NewEngine(metrics.EngineConfig{})
       tokenizer := tokenizers.NewMockTokenizer()
       
       result, err := engine.AnalyzeDocument(context.Background(), "test text", tokenizer)
       if err != nil {
           t.Fatalf("Analysis failed: %v", err)
       }
       
       if result.TokenCount == 0 {
           t.Error("Expected non-zero token count")
       }
   }
   ```

3. **Benchmark tests**:
   ```go
   func BenchmarkMockTokenizer_Tokenize(b *testing.B) {
       tokenizer := NewMockTokenizer()
       text := "This is a benchmark test with multiple words and punctuation marks!"
       
       b.ResetTimer()
       for i := 0; i < b.N; i++ {
           _, err := tokenizer.Tokenize(context.Background(), text)
           if err != nil {
               b.Fatal(err)
           }
       }
   }
   ```

### Test Coverage

1. **Run coverage analysis**:
   ```bash
   go test -cover ./...
   go test -coverprofile=coverage.out ./...
   go tool cover -html=coverage.out -o coverage.html
   ```

2. **Maintain high coverage**:
   - Aim for at least 80% code coverage
   - Focus on critical paths and edge cases
   - Test error conditions and boundary values

### Test Data

1. **Use test fixtures**:
   ```go
   // testdata/sample_texts.txt
   The quick brown fox jumps over the lazy dog.
   Machine learning models process text efficiently.
   Hello, world! This is a test.
   ```

2. **Create test helpers**:
   ```go
   func createTestTokenizer(t *testing.T) Tokenizer {
       tokenizer := NewMockTokenizer()
       err := tokenizer.Initialize(TokenizerConfig{
           Name: "test",
           Type: "mock",
       })
       if err != nil {
           t.Fatalf("Failed to initialize test tokenizer: %v", err)
       }
       return tokenizer
   }
   ```

## Documentation Standards

### Code Documentation

1. **Follow Go documentation conventions**:
   - Start with the name of the thing being documented
   - Use complete sentences
   - Be concise but informative

2. **Document public APIs**:
   ```go
   // Tokenizer defines the interface for text tokenization.
   // Implementations should be thread-safe and support cancellation via context.
   type Tokenizer interface {
       // Tokenize splits the input text into tokens.
       // Returns a TokenizationResult with tokens and metadata.
       Tokenize(ctx context.Context, text string) (*TokenizationResult, error)
   }
   ```

### User Documentation

1. **Update relevant documentation**:
   - README.md for major changes
   - User guide for new features
   - API reference for interface changes
   - Tutorials for new use cases

2. **Write clear examples**:
   ```markdown
   ## Using the New Feature
   
   To use the new tokenizer:
   
   ```bash
   ./ted analyze input.txt --tokenizers=new-tokenizer
   ```
   
   Configuration:
   ```yaml
   tokenizers:
     enabled: ["new-tokenizer"]
     configs:
       new-tokenizer:
         type: "custom"
         parameters:
           model: "path/to/model"
   ```
   ```

### API Documentation

1. **Update API reference**:
   - Document new types and interfaces
   - Provide usage examples
   - Include error conditions

2. **Generate documentation**:
   ```bash
   # Generate Go documentation
   go doc ./...
   
   # Generate API reference
   godoc -http=:6060
   ```

## Pull Request Process

### Before Submitting

1. **Ensure code quality**:
   ```bash
   # Format code
   gofmt -w .
   
   # Run linter
   golangci-lint run
   
   # Run tests
   go test ./...
   
   # Run benchmarks
   go test -bench=. ./...
   ```

2. **Update documentation**:
   - Update relevant documentation files
   - Add examples for new features
   - Update API reference if needed

3. **Test your changes**:
   - Test with different input types
   - Test error conditions
   - Test performance impact

### Pull Request Guidelines

1. **Title and description**:
   ```
   feat: add support for new tokenizer type
   
   - Implements CustomTokenizer for specialized text processing
   - Adds configuration options for custom parameters
   - Includes comprehensive tests and documentation
   - Fixes #123 (if applicable)
   ```

2. **Keep PRs focused**:
   - One feature or fix per PR
   - Keep changes reasonably sized
   - Break large changes into multiple PRs

3. **Include tests**:
   - Add unit tests for new functionality
   - Add integration tests for complex features
   - Update existing tests if needed

4. **Update examples**:
   - Add examples for new features
   - Update existing examples if needed
   - Test examples work correctly

### Review Process

1. **Code review checklist**:
   - [ ] Code follows style guidelines
   - [ ] Tests are included and pass
   - [ ] Documentation is updated
   - [ ] Examples work correctly
   - [ ] No performance regressions
   - [ ] Error handling is appropriate

2. **Address feedback**:
   - Respond to review comments
   - Make requested changes
   - Explain decisions when needed
   - Be open to suggestions

## Issue Reporting

### Bug Reports

When reporting bugs, include:

1. **Environment information**:
   - Operating system and version
   - Go version
   - Python version (if applicable)
   - TokEntropyDrift version

2. **Steps to reproduce**:
   ```bash
   # Commands that reproduce the issue
   ./ted analyze input.txt --tokenizers=gpt2
   ```

3. **Expected vs actual behavior**:
   - What you expected to happen
   - What actually happened
   - Any error messages

4. **Additional context**:
   - Input files (if relevant)
   - Configuration files
   - Logs or error output

### Feature Requests

When requesting features, include:

1. **Problem description**:
   - What problem you're trying to solve
   - Current limitations or pain points

2. **Proposed solution**:
   - How you envision the feature working
   - Any specific requirements or constraints

3. **Use cases**:
   - Real-world scenarios where this would be useful
   - Examples of how you would use it

4. **Alternatives considered**:
   - Other approaches you've considered
   - Why this solution is preferred

## Community Guidelines

### Communication

1. **Be respectful and inclusive**:
   - Treat all contributors with respect
   - Welcome newcomers and help them get started
   - Use inclusive language

2. **Provide constructive feedback**:
   - Focus on the code, not the person
   - Suggest improvements, not just point out problems
   - Be specific and actionable

3. **Ask questions**:
   - Don't hesitate to ask for clarification
   - Help others understand your perspective
   - Be open to learning from others

### Contribution Recognition

1. **Credit contributors**:
   - All contributors are listed in CONTRIBUTORS.md
   - Significant contributions are acknowledged in release notes
   - Contributors are mentioned in relevant documentation

2. **Celebrate contributions**:
   - Thank contributors for their work
   - Highlight interesting contributions
   - Share success stories

### Getting Help

1. **Documentation**:
   - Check the documentation first
   - Look for existing issues or discussions
   - Search the codebase for examples

2. **Community channels**:
   - GitHub Discussions for questions
   - GitHub Issues for bugs and features
   - Pull requests for code reviews

3. **Mentorship**:
   - Experienced contributors are available to help
   - Don't hesitate to ask for guidance
   - Offer to help others when you can

## Development Tools

### Recommended Tools

1. **IDE/Editor**:
   - GoLand (JetBrains)
   - Visual Studio Code with Go extension
   - Vim/Neovim with Go plugins

2. **Linting and Formatting**:
   ```bash
   # Install golangci-lint
   go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
   
   # Run linter
   golangci-lint run
   ```

3. **Testing Tools**:
   ```bash
   # Install testify for testing
   go get github.com/stretchr/testify
   
   # Run tests with coverage
   go test -cover ./...
   ```

### Development Scripts

1. **Makefile targets**:
   ```makefile
   .PHONY: build test lint clean
   
   build:
       go build -o ted cmd/ted/main.go
   
   test:
       go test ./...
   
   test-race:
       go test -race ./...
   
   test-coverage:
       go test -coverprofile=coverage.out ./...
       go tool cover -html=coverage.out
   
   lint:
       golangci-lint run
   
   clean:
       rm -f ted coverage.out
   ```

2. **Pre-commit hooks**:
   ```bash
   # .git/hooks/pre-commit
   #!/bin/bash
   set -e
   
   echo "Running tests..."
   go test ./...
   
   echo "Running linter..."
   golangci-lint run
   
   echo "Formatting code..."
   gofmt -w .
   ```

## Conclusion

Thank you for contributing to TokEntropyDrift! Your contributions help make this project better for everyone.

Remember:
- Start small and build up
- Ask questions when you need help
- Be patient with the review process
- Celebrate your contributions

We look forward to working with you!
