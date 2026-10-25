// PDFUnlock - Free Open Source PDF Password Remover & Unlocker
// A powerful CLI tool to unlock password-protected PDF files instantly
// Keywords: unlock pdf free, pdf unlocker, batch pdf unlock, open source
module github.com/fadeltd/pdfunlock

go 1.27

require (
	github.com/pdfcpu/pdfcpu v0.16.3
	golang.org/x/term v0.47.0
)

require (
	github.com/clipperhouse/uax29/v2 v2.7.0 // indirect
	github.com/hhrutter/tiff v1.0.7 // indirect
	github.com/mattn/go-runewidth v0.0.31 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	go.yaml.in/yaml/v3 v3.0.5 // indirect
	golang.org/x/crypto v0.58.0 // indirect
	golang.org/x/image v0.47.0 // indirect
	golang.org/x/sys v0.49.0 // indirect
	golang.org/x/text v0.43.0 // indirect
)

// pdfcpu v0.16.3 plus a fix for V4 encryption without a top level /Length (issue #3),
// from https://github.com/fadeltd/pdfcpu/tree/fix-v4-key-length. Drop once fixed upstream.
replace github.com/pdfcpu/pdfcpu => github.com/fadeltd/pdfcpu v0.16.4-0.20261025023603-3d85045996bb
