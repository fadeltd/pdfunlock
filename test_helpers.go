//go:build test
// +build test

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

// TestHelper provides utilities for testing PDFUnlock
type TestHelper struct {
	tempDir string
	t       *testing.T
}

// NewTestHelper creates a new test helper with a temporary directory
func NewTestHelper(t *testing.T) *TestHelper {
	tempDir, err := os.MkdirTemp("", "pdfunlock_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	return &TestHelper{
		tempDir: tempDir,
		t:       t,
	}
}

// Cleanup removes the temporary directory
func (th *TestHelper) Cleanup() {
	os.RemoveAll(th.tempDir)
}

// TempDir returns the temporary directory path
func (th *TestHelper) TempDir() string {
	return th.tempDir
}

// CopyProtectedPDF copies the actual protected PDF file for testing
func (th *TestHelper) CopyProtectedPDF(filename string) string {
	path := filepath.Join(th.tempDir, filename)
	sourcePath := "testdata/report_protected.pdf"

	// Read the source file
	sourceData, err := os.ReadFile(sourcePath)
	if err != nil {
		th.t.Fatalf("Failed to read source PDF %s: %v", sourcePath, err)
	}

	// Write to destination
	if err := os.WriteFile(path, sourceData, 0644); err != nil {
		th.t.Fatalf("Failed to copy protected PDF to %s: %v", filename, err)
	}

	return path
}

// CreateTextFile creates a text file for testing
func (th *TestHelper) CreateTextFile(filename, content string) string {
	path := filepath.Join(th.tempDir, filename)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		th.t.Fatalf("Failed to create text file %s: %v", filename, err)
	}
	return path
}

// CreateSubdir creates a subdirectory in the temp directory
func (th *TestHelper) CreateSubdir(dirname string) string {
	path := filepath.Join(th.tempDir, dirname)
	if err := os.MkdirAll(path, 0755); err != nil {
		th.t.Fatalf("Failed to create subdirectory %s: %v", dirname, err)
	}
	return path
}

// FileExists checks if a file exists
func (th *TestHelper) FileExists(filename string) bool {
	path := filepath.Join(th.tempDir, filename)
	_, err := os.Stat(path)
	return err == nil
}

// GetFilePath returns the full path to a file in the temp directory
func (th *TestHelper) GetFilePath(filename string) string {
	return filepath.Join(th.tempDir, filename)
}

// CreateTestScenario creates a complete test scenario with multiple files
func (th *TestHelper) CreateTestScenario() {
	// Create main directory structure
	th.CreateSubdir("pdfs")
	th.CreateSubdir("output")
	th.CreateSubdir("mixed")

	// Create PDF files in pdfs directory using actual protected PDF
	pdfDir := filepath.Join(th.tempDir, "pdfs")
	for i := 1; i <= 3; i++ {
		filename := fmt.Sprintf("document%d.pdf", i)
		th.CopyProtectedPDF(filepath.Join("pdfs", filename))
	}

	// Create mixed files in mixed directory
	mixedDir := filepath.Join(th.tempDir, "mixed")
	// Copy one actual protected PDF to mixed directory
	th.CopyProtectedPDF(filepath.Join("mixed", "document.pdf"))

	// Create non-PDF files
	nonPDFFiles := map[string]string{
		"readme.txt": "This is a text file",
		"image.jpg":  "fake jpg content",
		"data.json":  `{"test": "data"}`,
	}

	for filename, content := range nonPDFFiles {
		path := filepath.Join(mixedDir, filename)
		if err := os.WriteFile(path, []byte(content), 0644); err != nil {
			th.t.Fatalf("Failed to create file %s: %v", filename, err)
		}
	}
}

// MockPasswordInput simulates password input for testing
// This is a placeholder - in real testing you might use dependency injection
// or interfaces to mock the password input functionality
func (th *TestHelper) MockPasswordInput(password string) {
	// In a real implementation, you would mock the getPassword() function
	// This could be done by making getPassword() an interface or using
	// dependency injection patterns
	th.t.Logf("Mock password input: %s", password)
}

// AssertFileCount checks that a directory contains the expected number of files
func (th *TestHelper) AssertFileCount(dir string, expectedCount int, fileType string) {
	files, err := os.ReadDir(dir)
	if err != nil {
		th.t.Fatalf("Failed to read directory %s: %v", dir, err)
	}

	count := 0
	for _, file := range files {
		if file.Type().IsRegular() {
			if fileType == "" || filepath.Ext(file.Name()) == fileType {
				count++
			}
		}
	}

	if count != expectedCount {
		th.t.Errorf("Expected %d %s files in %s, got %d", expectedCount, fileType, dir, count)
	}
}
