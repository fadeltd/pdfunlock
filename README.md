# PDFUnlock - Free PDF Password Remover & Unlocker

🔓 **Free Open Source PDF Unlocker** - A powerful CLI command-line tool to unlock password-protected PDF files instantly. Perfect for batch PDF unlock operations, this free PDF password remover supports both single file processing and bulk directory processing.

**Keywords:** unlock pdf free, pdf unlocker, open source pdf unlock, batch pdf unlock, free pdf password remover, decrypt pdf files, pdf security removal

## Features - Why Choose This Free PDF Unlocker?

- 🔓 **Free PDF Unlock**: Completely free and open source PDF password remover
- 📁 **Batch PDF Unlock**: Process entire directories of password-protected PDFs at once
- ⚡ **Fast PDF Decryption**: Lightning-fast PDF security removal built with Go
- 🖥️ **Cross-Platform PDF Unlocker**: Works on Linux, macOS, and Windows
- 🎯 **Simple PDF Unlock Tool**: Easy-to-use command-line interface, no complex setup
- 💯 **Open Source PDF Unlocker**: Transparent, secure, and community-driven
- 🚀 **Efficient PDF Password Remover**: Optimized for performance and reliability
- 📋 **Single & Batch Processing**: Unlock one PDF or thousands - your choice
- **Flexible input methods**: Positional arguments with auto-detection or traditional flags
- **Smart path detection**: Automatically detects files vs directories based on extension and path type
- Single PDF file processing with automatic output naming
- Smart password caching per directory (reuses password for files in same directory)
- Automatic password retry with cache invalidation on authentication failures
- Support for both user and owner passwords
- Automatic handling of non-encrypted PDFs (copies them unchanged)

> **Note**: This is a desktop CLI tool, not a web-based service. Download and run locally for maximum security and privacy of your PDF files.

## About This Free PDF Unlocker

Looking for a **free PDF unlocker** or **open source PDF password remover**? PDFUnlock is the perfect solution for anyone who needs to:

- **Unlock PDF files for free** without online services
- **Remove PDF passwords** from multiple files at once
- **Decrypt password-protected PDFs** securely on your own computer
- **Batch unlock PDF files** in entire directories
- Use a **free alternative to paid PDF unlockers**

Unlike web-based PDF unlock services, this tool runs entirely on your computer, ensuring your sensitive documents never leave your device. Perfect for businesses, students, and anyone who values privacy and security.

## Installation

### Download Pre-built Binaries

Download the latest release for your platform from the [Releases page](https://github.com/fadeltd/pdfunlock/releases).

Supported platforms:
- Linux (x86_64, ARM64)
- Windows (x86_64)
- macOS (Intel x86_64, Apple Silicon ARM64)

### Build from Source

```bash
git clone https://github.com/fadeltd/pdfunlock.git
cd pdfunlock
make build
```

### Install Globally

```bash
make install
```

## Usage - How to Unlock PDF Files for Free

### Unlock Single PDF File (Free PDF Password Removal)

```bash
# Unlock a password-protected PDF file instantly - Creates: document_unlocked.pdf
./pdfunlock document.pdf

# Specify both input and output files for PDF decryption
./pdfunlock document.pdf custom_output.pdf
./pdfunlock -in input.pdf -out output.pdf

# Auto-generate output filename for PDF unlock (creates input_unlocked.pdf)
./pdfunlock -in input.pdf

# Use owner password for PDF security removal
./pdfunlock document.pdf --owner
```

### Batch PDF Unlock (Process Multiple PDFs)

```bash
# Unlock all password-protected PDFs in a directory
./pdfunlock /path/to/pdf/directory

# Batch PDF password removal using flags
./pdfunlock -dir /path/to/pdf/directory

# Use owner password for directory processing
./pdfunlock /path/to/directory --owner
```

**Perfect for:**
- Students unlocking academic PDFs
- Businesses processing encrypted documents
- Anyone needing free PDF password removal
- Bulk PDF decryption tasks

When processing a directory, the tool will:
- Scan for all PDF files in the directory
- Prompt for password once per directory (cached for subsequent operations)
- Automatically retry with new password if authentication fails (up to 3 attempts)
- Create unlocked versions with "_unlocked" suffix
- Skip files that are already unlocked
- Reuse cached passwords when processing files from the same directory

### Options

- `-in`: Input PDF file path
- `-out`: Output PDF file path (required when using -in)
- `-dir`: Directory containing PDF files to process
- `--owner`: Treat password as owner password (default: user password)

## Build

### Local Development

```bash
# Build the binary
make build

# Clean build artifacts
make clean

# Run tests
make test

# Install to /usr/local/bin/
sudo make install

# Uninstall from /usr/local/bin/
sudo make uninstall

# Cross-compile for multiple platforms with UPX compression
make build-all
```

**Build Output:**
- Single platform builds: `./pdfunlock` (current directory)
- Multi-platform builds: `dist/pdfunlock_{platform}/pdfunlock[.exe]`
- Supports: Linux (amd64, arm64), Windows (amd64), macOS (amd64, arm64)
- Automatic UPX compression reduces binary size by ~37%

### Environment Configuration

The build system supports environment-based configuration for UPX compression:

```bash
# Copy the sample environment file
cp env.sample .env

# Edit .env to customize UPX settings
# Available options:
# ENABLE_UPX=true                    # Global UPX enable/disable
# ENABLE_UPX_LINUX=true             # Linux-specific UPX control
# ENABLE_UPX_WINDOWS=true           # Windows-specific UPX control  
# ENABLE_UPX_DARWIN_INTEL=true      # macOS Intel UPX control
# ENABLE_UPX_DARWIN_ARM=false       # macOS ARM UPX control (disabled by default)
# GITHUB_REPOSITORY_OWNER=fadeltd   # Repository owner for GoReleaser
```

**Note**: The `.env` file is ignored by git, so your local settings won't be committed to the repository.

### Automated Releases

This project uses GitHub Actions with GoReleaser for automated cross-platform releases:

**Supported Platforms:**
- Linux (x86_64, ARM64)
- Windows (x86_64)
- macOS (Intel x86_64, Apple Silicon ARM64)

**Release Process:**

*Using Makefile (Recommended):*
```bash
# Create and push tag in one command
make tag-release VERSION=v1.0.0

# Or step by step:
make tag VERSION=v1.0.0
make release VERSION=v1.0.0
```

*Manual process:*
```bash
# Create and push a git tag
git tag v1.0.0 && git push origin v1.0.0
```

2. GitHub Actions automatically builds and releases binaries
3. Binaries are available on the [Releases page](https://github.com/fadeltd/pdfunlock/releases)

**UPX Compression:**
UPX compression is automatically applied during cross-platform builds and can be configured via environment variables:

- **Local builds**: Configure via `.env` file (see Environment Configuration section)
- **GitHub Actions**: Set repository variables for automated releases:
  - Go to: Repository Settings → Secrets and variables → Actions → Variables
  - Available variables: `ENABLE_UPX_LINUX`, `ENABLE_UPX_WINDOWS`, `ENABLE_UPX_DARWIN_INTEL`, `ENABLE_UPX_DARWIN_ARM`
- **Compression results**: Typically reduces binary size by ~37% (from ~16MB to ~6MB)
- **Note**: UPX compression on macOS ARM may cause compatibility issues, so it's disabled by default

**Version Information:**
```bash
# Check version of installed binary
./pdfunlock --version
```

## Dependencies

- [pdfcpu](https://github.com/pdfcpu/pdfcpu) - PDF processing library

## Contributing

We welcome contributions to PDFUnlock! Here's how you can help:

### How to Contribute

1. **Report Issues**: Found a bug or have a feature request? [Open an issue](https://github.com/fadeltd/pdfunlock/issues) on GitHub
2. **Submit Pull Requests**: 
   - Fork the repository
   - Create a feature branch (`git checkout -b feature/amazing-feature`)
   - Make your changes
   - Add tests if applicable
   - Commit your changes (`git commit -m 'Add amazing feature'`)
   - Push to the branch (`git push origin feature/amazing-feature`)
   - Open a Pull Request
3. **Wait for Review**: Maintainers will review your PR and provide feedback
4. **Address Feedback**: Make any requested changes
5. **Merge**: Once approved, your contribution will be merged!

### Development Guidelines

- Follow Go best practices and conventions
- Add tests for new functionality
- Update documentation as needed
- Ensure all CI checks pass

## Support the Project

If you find PDFUnlock useful, consider supporting the development:

[![Ko-fi](https://ko-fi.com/img/githubbutton_sm.svg)](https://ko-fi.com/fadeltd)

Your support helps maintain and improve this free, open-source tool!

## License

MIT License

---

## SEO Tags & Keywords

**Primary Keywords:** unlock pdf free, pdf unlocker, open source pdf unlock, batch pdf unlock, free pdf password remover, decrypt pdf files, pdf security removal

**Secondary Keywords:** pdf unlock tool, free pdf decryption, remove pdf password, unlock password protected pdf, pdf password cracker, bulk pdf unlock, command line pdf unlocker, offline pdf unlock, secure pdf unlock, desktop pdf unlocker

**Use Cases:** academic pdf unlock, business document decryption, bulk pdf processing, password recovery, document management, file conversion preparation, archive processing

**Alternatives to:** online pdf unlocker, web pdf unlock, paid pdf tools, Adobe Acrobat password removal, SmallPDF unlock, iLovePDF unlock, PDF24 unlock

**Technical:** golang pdf tool, cross-platform pdf utility, CLI pdf processor, open source document tools, github pdf unlocker, MIT license pdf tool
