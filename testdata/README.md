# Test Data Directory

This directory contains test files used by the PDFUnlock test suite.

## Files

- `sample_encrypted.pdf` - A sample encrypted PDF for testing decryption
- `sample_unencrypted.pdf` - A sample unencrypted PDF for testing copy functionality
- `passwords.txt` - Contains test passwords used in tests

## Usage

These files are used by the unit and integration tests to verify:
- PDF file detection and processing
- Encryption/decryption functionality
- Error handling with various file types
- Performance with different file sizes

## Creating Test PDFs

To create new test PDFs for testing:

1. Use any PDF creation tool to create a simple PDF
2. For encrypted PDFs, use a tool like `qpdf` or Adobe Acrobat to add password protection
3. Place the files in this directory
4. Update the test files to reference the new PDFs

## Security Note

Test passwords and files in this directory are for testing purposes only and should not contain any sensitive information.