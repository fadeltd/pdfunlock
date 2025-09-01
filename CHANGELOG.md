# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [1.0.0] - 2025-09-01

### Added
- Initial release of pdfunlock tool
- Password-protected PDF decryption functionality
- Interactive password prompting with secure input
- Directory scanning for batch PDF processing
- Automatic password retry with cache invalidation (up to 3 attempts)
- Support for both user and owner passwords
- Positional arguments with automatic file/directory detection
- Cross-platform binary releases (Linux, Windows, macOS)
- GitHub Actions CI/CD pipeline with GoReleaser
- Configurable UPX compression per platform
- Command-line flags and version information

### Features
- Single PDF file processing
- Batch directory processing
- Smart password caching and validation
- Flexible input methods
- Automated release pipeline

### Supported Platforms
- Linux (x86_64, ARM64)
- Windows (x86_64)
- macOS (Intel x86_64, Apple Silicon ARM64)
