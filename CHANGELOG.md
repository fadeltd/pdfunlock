# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

Nothing yet.

## [1.2.0] - 2026-10-27

### Added
- Install on Windows with winget: `winget install fadeltd.pdfunlock` (#6)

### Changed
- The Windows release archive is now `pdfunlock_Windows_x86_64.zip` (was `.tar.gz`)

## [1.1.1] - 2026-10-25

### Fixed
- Correct password rejected with "please provide the correct password" for AES-128/RC4-128 PDFs whose encryption dictionary leaves out the optional key length (#3)

### Changed
- Pin Go 1.27.2 for local builds (was 1.26.8)
- Update pdfcpu to v0.16.3, built from the `fadeltd/pdfcpu` fork until the key length fix lands upstream

## [1.1.0] - 2026-09-27

### Added
- Install on macOS with Homebrew: `brew install --cask fadeltd/tap/pdfunlock`
- Project website at https://fadeltd.github.io/pdfunlock/

### Changed
- Pin Go 1.26.8 for local builds (was 1.26.4)

### Fixed
- Release notes now list the Windows download as `pdfunlock_Windows_x86_64.tar.gz`

## [1.0.1] - 2026-07-03

### Changed
- Windows release archive is `.tar.gz` instead of `.zip`
- Release archives no longer wrap files in a directory; the binary is at the top level
- Update GitHub Actions to `actions/checkout@v4.2.0` and `actions/setup-go@v5`

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
