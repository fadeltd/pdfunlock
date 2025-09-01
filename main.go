// PDFUnlock - Free Open Source PDF Password Remover & Unlocker
// A powerful CLI command-line tool to unlock password-protected PDF files instantly.
// Perfect for batch PDF unlock operations, this free PDF password remover supports
// both single file processing and bulk directory processing.
//
// Features:
// - Free PDF unlock without online services
// - Batch PDF password removal for multiple files
// - Cross-platform PDF unlocker (Linux, macOS, Windows)
// - Open source alternative to paid PDF tools
// - Secure offline PDF decryption
//
// Keywords: unlock pdf free, pdf unlocker, batch pdf unlock, open source pdf unlock,
// free pdf password remover, decrypt pdf files, pdf security removal
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
	"golang.org/x/term"
)

// Version information (set by GoReleaser)
var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

// Global password cache to store passwords by directory
var passwordCache = make(map[string]string)

func main() {
	// Support both flag-based and positional arguments
	inPath := flag.String("in", "", "Path to input PDF")
	outPath := flag.String("out", "", "Path to output PDF (decrypted)")
	dirPath := flag.String("dir", "", "Directory containing PDF files to process")
	owner := flag.Bool("owner", false, "Treat password as owner password (default: user password)")
	versionFlag := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Handle version flag
	if *versionFlag {
		fmt.Printf("pdfunlock %s\n", version)
		fmt.Printf("Commit: %s\n", commit)
		fmt.Printf("Built: %s\n", date)
		os.Exit(0)
	}

	// Get positional arguments
	args := flag.Args()

	// Determine input path from flags or positional arguments
	var inputPath string
	var outputPath string

	if *dirPath != "" {
		// Flag-based directory processing
		processDirectory(*dirPath, *owner)
		return
	} else if *inPath != "" {
		// Flag-based file processing
		inputPath = *inPath
		outputPath = *outPath
	} else if len(args) >= 1 {
		// Positional argument processing
		inputPath = args[0]
		if len(args) >= 2 {
			outputPath = args[1]
		}
	} else {
		// No input provided
		fmt.Println("Usage:")
		fmt.Println("  Positional:  pdfunlock <input.pdf|directory> [output.pdf] [--owner]")
		fmt.Println("  Single file: pdfunlock -in input.pdf [-out output.pdf] [--owner]")
		fmt.Println("  Auto naming: pdfunlock -in input.pdf [--owner] (creates input_unlocked.pdf)")
		fmt.Println("  Directory:   pdfunlock -dir /path/to/pdfs [--owner]")
		fmt.Println("  Version:     pdfunlock --version")
		os.Exit(2)
	}

	// Auto-detect if input is a file or directory
	info, err := os.Stat(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error accessing path '%s': %v\n", inputPath, err)
		os.Exit(1)
	}

	if info.IsDir() {
		// Process as directory
		processDirectory(inputPath, *owner)
	} else {
		// Process as single file
		// Check if it's a PDF file
		if !strings.HasSuffix(strings.ToLower(inputPath), ".pdf") {
			fmt.Fprintf(os.Stderr, "Error: '%s' is not a PDF file (must have .pdf extension)\n", inputPath)
			os.Exit(1)
		}

		// If no output path specified, generate one with _unlocked suffix
		if outputPath == "" {
			outputPath = generateOutputPath(inputPath)
		}
		processSingleFile(inputPath, outputPath, *owner)
	}
}

func processDirectory(dirPath string, isOwnerPassword bool) {
	// Find all PDF files in the directory
	pdfFiles, err := findPDFFiles(dirPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning directory: %v\n", err)
		os.Exit(1)
	}

	if len(pdfFiles) == 0 {
		fmt.Println("No PDF files found in the directory.")
		return
	}

	fmt.Printf("Found %d PDF files. Processing...\n", len(pdfFiles))

	successCount := 0
	for _, pdfFile := range pdfFiles {
		outputFile := generateOutputPath(pdfFile)
		if processFileWithRetry(pdfFile, outputFile, dirPath, isOwnerPassword) {
			successCount++
		}
	}

	fmt.Printf("\nProcessed %d out of %d files successfully.\n", successCount, len(pdfFiles))
}

func processSingleFile(inPath, outPath string, isOwnerPassword bool) {
	dirPath := filepath.Dir(inPath)
	if processFileWithRetry(inPath, outPath, dirPath, isOwnerPassword) {
		fmt.Printf("Successfully processed: %s -> %s\n", inPath, outPath)
	} else {
		os.Exit(1)
	}
}

func getPasswordForDirectory(dirPath string) string {
	// Normalize the directory path
	absDir, err := filepath.Abs(dirPath)
	if err != nil {
		absDir = dirPath
	}

	// Check if we already have a cached password for this directory
	if cachedPassword, exists := passwordCache[absDir]; exists {
		fmt.Printf("Using cached password for directory: %s\n", absDir)
		return cachedPassword
	}

	// Get password and cache it
	password := getPassword()
	if password != "" {
		passwordCache[absDir] = password
		fmt.Printf("Password cached for directory: %s\n", absDir)
	}
	return password
}

func getPassword() string {
	// First check for environment variable
	if envPassword := os.Getenv("PDF_PASSWORD"); envPassword != "" {
		return envPassword
	}

	// Prompt for password if no environment variable is set
	fmt.Print("Enter PDF password: ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		fmt.Fprintf(os.Stderr, "\nError reading password: %v\n", err)
		return ""
	}
	fmt.Println() // Add newline after password input
	return string(bytePassword)
}

func findPDFFiles(dirPath string) ([]string, error) {
	var pdfFiles []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.ToLower(filepath.Ext(path)) == ".pdf" {
			pdfFiles = append(pdfFiles, path)
		}
		return nil
	})

	return pdfFiles, err
}

func generateOutputPath(inputPath string) string {
	dir := filepath.Dir(inputPath)
	filename := filepath.Base(inputPath)
	ext := filepath.Ext(filename)
	nameWithoutExt := strings.TrimSuffix(filename, ext)
	return filepath.Join(dir, nameWithoutExt+"_unlocked"+ext)
}

func processFileWithRetry(inPath, outPath, dirPath string, isOwnerPassword bool) bool {
	maxRetries := 3
	for attempt := 1; attempt <= maxRetries; attempt++ {
		password := getPasswordForDirectory(dirPath)
		if password == "" {
			fmt.Println("No password provided. Exiting.")
			return false
		}

		result, isAuthError := processFile(inPath, outPath, password, isOwnerPassword)
		if result {
			return true
		}

		// If it's an authentication error and we have retries left
		if isAuthError && attempt < maxRetries {
			fmt.Printf("Authentication failed for %s. Please try again (attempt %d/%d)\n", filepath.Base(inPath), attempt, maxRetries)
			// Clear the cached password for this directory
			clearPasswordCache(dirPath)
			continue
		}

		// If it's not an auth error or we're out of retries, fail
		return false
	}
	return false
}

func processFile(inPath, outPath, password string, isOwnerPassword bool) (bool, bool) {
	var userPW, ownerPW string
	if isOwnerPassword {
		ownerPW = password
	} else {
		userPW = password
	}

	// Build a configuration that provides the password to pdfcpu.
	conf := model.NewAESConfiguration(userPW, ownerPW, 0)

	// Try to decrypt.
	if err := api.DecryptFile(inPath, outPath, conf); err != nil {
		// If the file is not encrypted, just copy it through.
		if isNotEncrypted(err) {
			if err = copyFile(inPath, outPath); err != nil {
				fmt.Fprintf(os.Stderr, "Copy failed for %s: %v\n", inPath, err)
				return false, false
			}
			fmt.Printf("Not encrypted (copied): %s\n", filepath.Base(inPath))
			return true, false
		}

		// Check if it's an authentication error
		isAuthError := isAuthenticationError(err)
		if isAuthError {
			fmt.Fprintf(os.Stderr, "Authentication failed for %s: %v\n", inPath, err)
		} else {
			fmt.Fprintf(os.Stderr, "Decrypt failed for %s: %v\n", inPath, err)
		}
		return false, isAuthError
	}

	fmt.Printf("Decrypted: %s\n", filepath.Base(inPath))
	return true, false
}

func clearPasswordCache(dirPath string) {
	absDir, err := filepath.Abs(dirPath)
	if err != nil {
		absDir = dirPath
	}
	delete(passwordCache, absDir)
	fmt.Printf("Password cache cleared for directory: %s\n", absDir)
}

func isAuthenticationError(err error) bool {
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "authentication") ||
		strings.Contains(msg, "password") ||
		strings.Contains(msg, "invalid password") ||
		strings.Contains(msg, "wrong password") ||
		strings.Contains(msg, "bad password") ||
		strings.Contains(msg, "incorrect password") ||
		strings.Contains(msg, "unauthorized") ||
		strings.Contains(msg, "access denied")
}

func isNotEncrypted(err error) bool {
	// pdfcpu error messages commonly contain hints like these when the file isn't encrypted.
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "not encrypted") ||
		strings.Contains(msg, "no encryption") ||
		strings.Contains(msg, "this file is not encrypted")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	return out.Sync()
}
