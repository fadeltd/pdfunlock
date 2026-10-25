package main

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/pdfcpu/pdfcpu/pkg/api"
	"github.com/pdfcpu/pdfcpu/pkg/pdfcpu/model"
)

// Fixtures in testdata, user password "secret", owner password "owner":
//   - test.pdf: unencrypted source
//   - test_aes128.pdf: pdfcpu encrypt --mode aes --key 128 --upw secret --opw owner test.pdf test_aes128.pdf
//   - test_aes128_nolength.pdf: test_aes128.pdf with the encrypt dict's top-level /Length renamed to /Lxngth,
//     which is optional for V4 since the key length comes from the crypt filter (issue #3, pdfcpu/pdfcpu#1522)
func TestProcessFile(t *testing.T) {
	tests := []struct {
		name          string
		file          string
		password      string
		owner         bool
		wantOK        bool
		wantAuthError bool
	}{
		{"not encrypted", "test.pdf", "secret", false, true, false},
		{"user password", "test_aes128.pdf", "secret", false, true, false},
		{"owner password", "test_aes128.pdf", "owner", true, true, false},
		{"wrong password", "test_aes128.pdf", "wrong", false, false, true},
		{"no /Length, user password", "test_aes128_nolength.pdf", "secret", false, true, false},
		{"no /Length, owner password", "test_aes128_nolength.pdf", "owner", true, true, false},
		{"no /Length, wrong password", "test_aes128_nolength.pdf", "wrong", false, false, true},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			in := filepath.Join("testdata", tc.file)
			out := filepath.Join(t.TempDir(), "out.pdf")

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
