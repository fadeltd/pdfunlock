package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// TestFindPDFFiles tests the findPDFFiles function
func TestFindPDFFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "pdfunlock_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []struct {
		name     string
		expected bool
	}{
		{"test1.pdf", true},
		{"test2.PDF", true}, // Test case insensitive
		{"test3.txt", false},
		{"document.pdf", true},
		{"readme.md", false},
	}

	for _, tf := range testFiles {
		filePath := filepath.Join(tempDir, tf.name)
		if err = os.WriteFile(filePath, []byte("dummy content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", tf.name, err)
		}
	}

	// Test findPDFFiles
	pdfFiles, err := findPDFFiles(tempDir)
	if err != nil {
		t.Fatalf("findPDFFiles failed: %v", err)
	}

	// Count expected PDF files
	expectedCount := 0
	for _, tf := range testFiles {
		if tf.expected {
			expectedCount++
		}
	}

	if len(pdfFiles) != expectedCount {
		t.Errorf("Expected %d PDF files, got %d", expectedCount, len(pdfFiles))
	}

	// Verify all found files are PDFs
	for _, file := range pdfFiles {
		if !strings.HasSuffix(strings.ToLower(file), ".pdf") {
			t.Errorf("Non-PDF file found: %s", file)
		}
	}
}

// TestGenerateOutputPath tests the generateOutputPath function
func TestGenerateOutputPath(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"/path/to/document.pdf", "/path/to/document_unlocked.pdf"},
		{"simple.pdf", "simple_unlocked.pdf"},
		{"/home/user/test.PDF", "/home/user/test_unlocked.PDF"},
		{"./relative/path/file.pdf", "relative/path/file_unlocked.pdf"},
	}

	for _, test := range tests {
		result := generateOutputPath(test.input)
		if result != test.expected {
			t.Errorf("generateOutputPath(%q) = %q, expected %q", test.input, result, test.expected)
		}
	}
}

// TestIsAuthenticationError tests the isAuthenticationError function
func TestIsAuthenticationError(t *testing.T) {
	tests := []struct {
		errorMsg string
		expected bool
	}{
		{"authentication failed", true},
		{"Authentication Failed", true}, // Case insensitive
		{"AUTHENTICATION FAILED", true},
		{"wrong password", true},
		{"Wrong Password", true},
		{"invalid password", true},
		{"file not found", false},
		{"permission denied", false},
		{"corrupted file", false},
		{"", false},
	}

	for _, test := range tests {
		err := fmt.Errorf("%s", test.errorMsg)
		result := isAuthenticationError(err)
		if result != test.expected {
			t.Errorf("isAuthenticationError(%q) = %v, expected %v", test.errorMsg, result, test.expected)
		}
	}
}

// TestIsNotEncrypted tests the isNotEncrypted function
func TestIsNotEncrypted(t *testing.T) {
	tests := []struct {
		errorMsg string
		expected bool
	}{
		{"not encrypted", true},
		{"Not Encrypted", true}, // Case insensitive
		{"NOT ENCRYPTED", true},
		{"file is not encrypted", true},
		{"authentication failed", false},
		{"file not found", false},
		{"corrupted file", false},
		{"", false},
	}

	for _, test := range tests {
		err := fmt.Errorf("%s", test.errorMsg)
		result := isNotEncrypted(err)
		if result != test.expected {
			t.Errorf("isNotEncrypted(%q) = %v, expected %v", test.errorMsg, result, test.expected)
		}
	}
}

// TestClearPasswordCache tests the clearPasswordCache function
func TestClearPasswordCache(t *testing.T) {
	// Setup: Add some passwords to cache
	passwordCache["/path/to/dir1"] = "password1"
	passwordCache["/path/to/dir2"] = "password2"

	// Test clearing specific directory
	clearPasswordCache("/path/to/dir1")

	// Verify dir1 is cleared but dir2 remains
	if _, exists := passwordCache["/path/to/dir1"]; exists {
		t.Error("Password cache for dir1 should be cleared")
	}
	if _, exists := passwordCache["/path/to/dir2"]; !exists {
		t.Error("Password cache for dir2 should still exist")
	}

	// Clean up
	delete(passwordCache, "/path/to/dir2")
}

// TestCopyFile tests the copyFile function
func TestCopyFile(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "pdfunlock_copy_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create source file
	srcPath := filepath.Join(tempDir, "source.txt")
	testContent := "This is test content for copy function"
	if err = os.WriteFile(srcPath, []byte(testContent), 0644); err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}

	// Test copying
	dstPath := filepath.Join(tempDir, "destination.txt")
	if err = copyFile(srcPath, dstPath); err != nil {
		t.Fatalf("copyFile failed: %v", err)
	}

	// Verify destination file exists and has correct content
	copiedContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("Failed to read copied file: %v", err)
	}

	if string(copiedContent) != testContent {
		t.Errorf("Copied content doesn't match. Expected %q, got %q", testContent, string(copiedContent))
	}

	// Test copying non-existent file
	nonExistentSrc := filepath.Join(tempDir, "nonexistent.txt")
	if err := copyFile(nonExistentSrc, dstPath); err == nil {
		t.Error("Expected error when copying non-existent file, but got nil")
	}
}

// Benchmark tests for performance
func BenchmarkFindPDFFiles(b *testing.B) {
	// Create a temporary directory with many files
	tempDir, err := os.MkdirTemp("", "pdfunlock_benchmark")
	if err != nil {
		b.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create 100 PDF files and 100 non-PDF files
	for i := 0; i < 100; i++ {
		pdfPath := filepath.Join(tempDir, fmt.Sprintf("test%d.pdf", i))
		txtPath := filepath.Join(tempDir, fmt.Sprintf("test%d.txt", i))
		os.WriteFile(pdfPath, []byte("dummy pdf"), 0644)
		os.WriteFile(txtPath, []byte("dummy txt"), 0644)
	}

	for b.Loop() {
		_, err := findPDFFiles(tempDir)
		if err != nil {
			b.Fatalf("findPDFFiles failed: %v", err)
		}
	}
}

func BenchmarkGenerateOutputPath(b *testing.B) {
	inputPath := "/very/long/path/to/some/deeply/nested/directory/structure/document.pdf"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = generateOutputPath(inputPath)
	}
}
