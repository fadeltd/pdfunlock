package main

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// writeMinimalPDF writes a one page, unencrypted PDF.
func writeMinimalPDF(t *testing.T, path string) {
	t.Helper()

	objs := []string{
		"<</Type/Catalog/Pages 2 0 R>>",
		"<</Type/Pages/Kids[3 0 R]/Count 1>>",
		"<</Type/Page/Parent 2 0 R/MediaBox[0 0 200 200]>>",
	}

	var b bytes.Buffer
	b.WriteString("%PDF-1.4\n")
	offsets := make([]int, len(objs))
	for i, obj := range objs {
		offsets[i] = b.Len()
		fmt.Fprintf(&b, "%d 0 obj\n%s\nendobj\n", i+1, obj)
	}
	xref := b.Len()
	fmt.Fprintf(&b, "xref\n0 %d\n0000000000 65535 f \n", len(objs)+1)
	for _, off := range offsets {
		fmt.Fprintf(&b, "%010d 00000 n \n", off)
	}
	fmt.Fprintf(&b, "trailer\n<</Size %d/Root 1 0 R>>\nstartxref\n%d\n%%%%EOF\n", len(objs)+1, xref)

	if err := os.WriteFile(path, b.Bytes(), 0o644); err != nil {
		t.Fatal(err)
	}
}

// writeV4PDFWithoutEncryptLength writes an AES-128 (V4) PDF whose encrypt dict has no top level /Length,
// which is valid since for V4 the key length comes from the crypt filter (issue #3).
func writeV4PDFWithoutEncryptLength(t *testing.T, path, userPW, ownerPW string) {
	t.Helper()

	plain := filepath.Join(t.TempDir(), "plain.pdf")
	writeMinimalPDF(t, plain)

	conf := model.NewAESConfiguration(userPW, ownerPW, 128)
	if err := api.EncryptFile(context.Background(), plain, path, conf); err != nil {
		t.Fatal(err)
	}

	// Rename the key keeping the byte count so the xref table stays valid.
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	old := []byte("/Filter/Standard/Length 128")
	if n := bytes.Count(b, old); n != 1 {
		t.Fatalf("found %d encrypt dict /Length entries, want 1", n)
	}
	b = bytes.Replace(b, old, []byte("/Filter/Standard/Lxngth 128"), 1)
	if err := os.WriteFile(path, b, 0o644); err != nil {
		t.Fatal(err)
	}
}

func TestProcessFileV4WithoutEncryptLength(t *testing.T) {
	dir := t.TempDir()
	in := filepath.Join(dir, "in.pdf")
	writeV4PDFWithoutEncryptLength(t, in, "secret", "owner")

	tests := []struct {
		name          string
		password      string
		owner         bool
		wantOK        bool
		wantAuthError bool
	}{
		{"user password", "secret", false, true, false},
		{"owner password", "owner", true, true, false},
		{"wrong password", "wrong", false, false, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out := filepath.Join(dir, "out.pdf")
			_ = os.Remove(out)

			ok, isAuthError := processFile(in, out, tc.password, tc.owner)
			if ok != tc.wantOK || isAuthError != tc.wantAuthError {
				t.Fatalf("processFile() = (%t, %t), want (%t, %t)", ok, isAuthError, tc.wantOK, tc.wantAuthError)
			}
			if !tc.wantOK {
				return
			}

			// The output must open without a password.
			if err := api.ValidateFile(context.Background(), out, model.NewDefaultConfiguration(), nil); err != nil {
				t.Fatalf("validate %s: %v", out, err)
			}
		})
	}
}
