//go:build integration
// +build integration

package main

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// TestIntegrationCLIVersion tests the --version flag
func TestIntegrationCLIVersion(t *testing.T) {
	// Build the binary first
	binaryPath := buildTestBinary(t)
	defer os.Remove(binaryPath)

	cmd := exec.Command(binaryPath, "--version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Command failed: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "pdfunlock") {
		t.Errorf("Version output should contain 'pdfunlock', got: %s", outputStr)
	}
}

// TestIntegrationCLIHelp tests the help output when no arguments are provided
func TestIntegrationCLIHelp(t *testing.T) {
	binaryPath := buildTestBinary(t)
	defer os.Remove(binaryPath)

	cmd := exec.Command(binaryPath)
	output, err := cmd.CombinedOutput()

	// Help should exit with code 2
	if exitError, ok := err.(*exec.ExitError); ok {
		if exitError.ExitCode() != 2 {
			t.Errorf("Expected exit code 2, got %d", exitError.ExitCode())
		}
	} else {
		t.Errorf("Expected exit error, got: %v", err)
	}

	outputStr := string(output)
	expectedStrings := []string{"Usage:", "pdfunlock", "input.pdf", "directory"}
	for _, expected := range expectedStrings {
		if !strings.Contains(outputStr, expected) {
			t.Errorf("Help output should contain '%s', got: %s", expected, outputStr)
		}
	}
}

// TestIntegrationPDFUnlockSuccess tests successful PDF unlocking with correct password
func TestIntegrationPDFUnlockSuccess(t *testing.T) {
	binaryPath := buildTestBinary(t)
	defer os.Remove(binaryPath)

	// Create temporary directory for test output
	tempDir, err := os.MkdirTemp("", "pdfunlock_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy test PDF to temp directory
	testPDFPath := filepath.Join(tempDir, "test_protected.pdf")
	err = copyTestPDF("testdata/report_protected.pdf", testPDFPath)
	if err != nil {
		t.Fatalf("Failed to copy test PDF: %v", err)
	}

	// Set correct password as environment variable
	os.Setenv("PDF_PASSWORD", "I57x1IAX1|=F4ng8")
	defer os.Unsetenv("PDF_PASSWORD")

	// Run pdfunlock with correct password
	cmd := exec.Command(binaryPath, testPDFPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected PDF unlock to succeed, got error: %v, output: %s", err, string(output))
	}

	// Check that unlocked file was created
	unlockedPath := strings.TrimSuffix(testPDFPath, ".pdf") + "_unlocked.pdf"
	if _, err := os.Stat(unlockedPath); os.IsNotExist(err) {
		t.Errorf("Expected unlocked PDF file to be created at %s", unlockedPath)
	}

	// Verify the unlocked file is readable and not empty
	info, err := os.Stat(unlockedPath)
	if err != nil {
		t.Errorf("Failed to stat unlocked file: %v", err)
	} else if info.Size() == 0 {
		t.Errorf("Unlocked PDF file is empty")
	}

	// Check output contains success message
	if !strings.Contains(string(output), "Successfully processed") {
		t.Errorf("Expected success message in output, got: %s", string(output))
	}
}

// TestIntegrationPDFUnlockWrongPassword tests PDF unlocking with wrong password
func TestIntegrationPDFUnlockWrongPassword(t *testing.T) {
	binaryPath := buildTestBinary(t)
	defer os.Remove(binaryPath)

	// Create temporary directory for test output
	tempDir, err := os.MkdirTemp("", "pdfunlock_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy test PDF to temp directory
	testPDFPath := filepath.Join(tempDir, "test_protected.pdf")
	err = copyTestPDF("testdata/report_protected.pdf", testPDFPath)
	if err != nil {
		t.Fatalf("Failed to copy test PDF: %v", err)
	}

	// Set wrong password as environment variable
	os.Setenv("PDF_PASSWORD", "wrongpassword")
	defer os.Unsetenv("PDF_PASSWORD")

	// Run pdfunlock with wrong password
	cmd := exec.Command(binaryPath, testPDFPath)
	output, err := cmd.CombinedOutput()

	// Should fail with authentication error
	if err == nil {
		t.Errorf("Expected PDF unlock to fail with wrong password, but it succeeded")
	}

	// Check that unlocked file was NOT created
	unlockedPath := strings.TrimSuffix(testPDFPath, ".pdf") + "_unlocked.pdf"
	if _, err := os.Stat(unlockedPath); !os.IsNotExist(err) {
		t.Errorf("Expected unlocked PDF file NOT to be created with wrong password")
	}

	// Check output contains authentication error
	outputStr := string(output)
	if !strings.Contains(outputStr, "authentication") && !strings.Contains(outputStr, "password") {
		t.Errorf("Expected authentication error in output, got: %s", outputStr)
	}
}

// TestIntegrationPDFUnlockDirectory tests unlocking PDFs in a directory
func TestIntegrationPDFUnlockDirectory(t *testing.T) {
	binaryPath := buildTestBinary(t)
	defer os.Remove(binaryPath)

	// Create temporary directory structure
	tempDir, err := os.MkdirTemp("", "pdfunlock_test")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Copy test PDF to temp directory
	testPDFPath := filepath.Join(tempDir, "test_protected.pdf")
	err = copyTestPDF("testdata/report_protected.pdf", testPDFPath)
	if err != nil {
		t.Fatalf("Failed to copy test PDF: %v", err)
	}

	// Set correct password as environment variable
	os.Setenv("PDF_PASSWORD", "I57x1IAX1|=F4ng8")
	defer os.Unsetenv("PDF_PASSWORD")

	// Run pdfunlock on directory
	cmd := exec.Command(binaryPath, tempDir)
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Expected PDF unlock on directory to succeed, got error: %v, output: %s", err, string(output))
	}

	// Check that unlocked file was created
	unlockedPath := strings.TrimSuffix(testPDFPath, ".pdf") + "_unlocked.pdf"
	if _, err := os.Stat(unlockedPath); os.IsNotExist(err) {
		t.Errorf("Expected unlocked PDF file to be created at %s", unlockedPath)
	}

	// Verify the unlocked file is readable and not empty
	info, err := os.Stat(unlockedPath)
	if err != nil {
		t.Errorf("Failed to stat unlocked file: %v", err)
	} else if info.Size() == 0 {
		t.Errorf("Unlocked PDF file is empty")
	}
}

// TestIntegrationCLIWithNonExistentFile tests behavior with non-existent input file
func TestIntegrationCLIWithNonExistentFile(t *testing.T) {
	binaryPath := buildTestBinary(t)
	defer os.Remove(binaryPath)

	cmd := exec.Command(binaryPath, "nonexistent.pdf")
	output, err := cmd.CombinedOutput()

	// Should exit with error code 1
	if exitError, ok := err.(*exec.ExitError); ok {
		if exitError.ExitCode() != 1 {
			t.Errorf("Expected exit code 1, got %d", exitError.ExitCode())
		}
	} else {
		t.Errorf("Expected exit error, got: %v", err)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "Error accessing path") {
		t.Errorf("Output should contain error message, got: %s", outputStr)
	}
}

// TestIntegrationCLIWithEmptyDirectory tests behavior with empty directory
func TestIntegrationCLIWithEmptyDirectory(t *testing.T) {
	binaryPath := buildTestBinary(t)
	defer os.Remove(binaryPath)

	// Create empty temp directory
	tempDir, err := os.MkdirTemp("", "pdfunlock_empty_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	cmd := exec.Command(binaryPath, tempDir)
	output, err := cmd.CombinedOutput()

	// Should complete successfully even with empty directory
	if err != nil {
		t.Errorf("Command should succeed with empty directory, got error: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "No PDF files found") {
		t.Errorf("Output should indicate no PDF files found, got: %s", outputStr)
	}
}

// TestIntegrationCLIWithDirectoryContainingNonPDFs tests directory with non-PDF files
func TestIntegrationCLIWithDirectoryContainingNonPDFs(t *testing.T) {
	binaryPath := buildTestBinary(t)
	defer os.Remove(binaryPath)

	// Create temp directory with non-PDF files
	tempDir, err := os.MkdirTemp("", "pdfunlock_nonpdf_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create some non-PDF files
	testFiles := []string{"test.txt", "document.docx", "image.jpg", "readme.md"}
	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		if err := os.WriteFile(filePath, []byte("dummy content"), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", filename, err)
		}
	}

	cmd := exec.Command(binaryPath, tempDir)
	output, err := cmd.CombinedOutput()

	// Should complete successfully
	if err != nil {
		t.Errorf("Command should succeed, got error: %v\nOutput: %s", err, output)
	}

	outputStr := string(output)
	if !strings.Contains(outputStr, "No PDF files found") {
		t.Errorf("Output should indicate no PDF files found, got: %s", outputStr)
	}
}

// TestIntegrationCLIPerformance tests basic performance characteristics
func TestIntegrationCLIPerformance(t *testing.T) {
	binaryPath := buildTestBinary(t)
	defer os.Remove(binaryPath)

	// Create temp directory with many dummy PDF files
	tempDir, err := os.MkdirTemp("", "pdfunlock_perf_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create 50 dummy PDF files (they won't be real PDFs, but that's ok for this test)
	for i := 0; i < 50; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("test%d.pdf", i))
		if err := os.WriteFile(filename, []byte("dummy pdf content"), 0644); err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
	}

	// Measure time to scan directory
	start := time.Now()
	cmd := exec.Command(binaryPath, tempDir)
	_, err = cmd.CombinedOutput()
	duration := time.Since(start)

	// Should complete within reasonable time (10 seconds is very generous)
	if duration > 10*time.Second {
		t.Errorf("Directory scanning took too long: %v", duration)
	}

	// The command will likely fail because these aren't real PDFs,
	// but we're just testing the scanning performance
	t.Logf("Directory scan of 50 files took: %v", duration)
}

// buildTestBinary builds the pdfunlock binary for testing
func buildTestBinary(t *testing.T) string {
	tempDir, err := os.MkdirTemp("", "pdfunlock_build")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	binaryPath := filepath.Join(tempDir, "pdfunlock_test")
	cmd := exec.Command("go", "build", "-o", binaryPath, ".")
	cmd.Stdout = &bytes.Buffer{}
	cmd.Stderr = &bytes.Buffer{}

	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to build binary: %v\nStdout: %s\nStderr: %s",
			err, cmd.Stdout, cmd.Stderr)
	}

	return binaryPath
}

// Helper function to copy test PDF file
func copyTestPDF(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, sourceFile)
	return err
}
