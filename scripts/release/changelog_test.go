package main

import (
	"strings"
	"testing"
)

const header = "# Changelog\n\nIntro.\n\n"

func TestParseUnreleasedIgnoresFiller(t *testing.T) {
	u := ParseUnreleased(header + "## [Unreleased]\n\nNothing yet.\n\n## [1.0.0] - 2026-09-27\n\n### Added\n- First\n")
	if u.HasContent {
		t.Fatalf("filler counted as content: %+v", u)
	}
}

func TestParseUnreleasedHeadingWithoutBullets(t *testing.T) {
	u := ParseUnreleased(header + "## [Unreleased]\n\n### Fixed\n\n## [1.0.0] - 2026-09-27\n")
	if u.HasContent {
		t.Fatal("a heading with no bullets is not content")
	}
}

func TestParseUnreleasedStopsAtNextRelease(t *testing.T) {
	u := ParseUnreleased(header + "## [Unreleased]\n\n### Fixed\n- New\n\n## [1.0.0] - 2026-09-27\n\n### Added\n- Old\n")
	if !u.HasContent || strings.Contains(u.Body, "Old") || len(u.Sections) != 1 {
		t.Fatalf("unexpected parse: %+v", u)
	}
}

func TestParseUnreleasedMissing(t *testing.T) {
	if u := ParseUnreleased(header); u.HasContent {
		t.Fatal("no Unreleased heading should mean no content")
	}
}

func TestCurrentVersion(t *testing.T) {
	cl := header + "## [Unreleased]\n\n### Fixed\n- a\n\n## [1.0.1] - 2026-07-03\n\n## [1.0.0] - 2025-09-01\n"
	if v, err := CurrentVersion(cl); err != nil || v != "1.0.1" {
		t.Fatalf("CurrentVersion = %q, %v; want 1.0.1", v, err)
	}
	if _, err := CurrentVersion(header + "## [Unreleased]\n"); err == nil {
		t.Error("expected an error without a released version")
	}
}

func TestDetermineBump(t *testing.T) {
	cases := map[string]Bump{
		"### Fixed\n- a":                      BumpPatch,
		"### Changed\n- a":                    BumpPatch,
		"### Removed\n- a":                    BumpPatch,
		"### Fixed\n- a\n\n### Added\n- b":    BumpMinor,
		"### Added\n- a\n\n### Breaking\n- b": BumpMajor,
		"### Breaking Changes\n- a":           BumpMajor,
		"Nothing yet.":                        BumpNone,
	}
	for body, want := range cases {
		got := DetermineBump(ParseUnreleased("## [Unreleased]\n\n" + body + "\n"))
		if got != want {
			t.Errorf("%q: got %s, want %s", body, got, want)
		}
	}
}

func TestNextVersion(t *testing.T) {
	cases := []struct {
		current string
		bump    Bump
		want    string
	}{
		{"1.0.0", BumpPatch, "1.0.1"},
		{"1.0.3", BumpMinor, "1.1.0"},
		{"1.4.2", BumpMajor, "2.0.0"},
		{"0.3.1", BumpMajor, "0.4.0"},
		{"1.2.3", BumpNone, "1.2.3"},
	}
	for _, c := range cases {
		got, err := NextVersion(c.current, c.bump)
		if err != nil || got != c.want {
			t.Errorf("NextVersion(%s, %s) = %s, %v; want %s", c.current, c.bump, got, err, c.want)
		}
	}
	if _, err := NextVersion("v1.0", BumpPatch); err == nil {
		t.Error("expected an error for a non-semver version")
	}
}

func TestApplyRelease(t *testing.T) {
	in := header + "## [Unreleased]\n\n### Added\n- Homebrew cask\n\n## [1.0.1] - 2026-07-03\n\n### Changed\n- Tarballs\n"
	rel, err := ApplyRelease(in, "2026-10-01")
	if err != nil || rel == nil {
		t.Fatalf("ApplyRelease: %v, %v", rel, err)
	}
	want := header + "## [Unreleased]\n\nNothing yet.\n\n## [1.1.0] - 2026-10-01\n\n### Added\n- Homebrew cask\n\n## [1.0.1] - 2026-07-03\n\n### Changed\n- Tarballs\n"
	if rel.Changelog != want {
		t.Fatalf("changelog:\n%s\nwant:\n%s", rel.Changelog, want)
	}
	if rel.Version != "1.1.0" || rel.Notes != "### Added\n- Homebrew cask" {
		t.Fatalf("unexpected release: %+v", rel)
	}

	// Releasing the result again is a no-op, which is what makes re-runs safe.
	again, err := ApplyRelease(rel.Changelog, "2026-10-01")
	if err != nil || again != nil {
		t.Fatalf("second release should be nil, got %+v, %v", again, err)
	}
}

func TestNotesFor(t *testing.T) {
	cl := header + "## [Unreleased]\n\nNothing yet.\n\n## [1.1.0] - 2026-10-01\n\n### Added\n- B\n\n## [1.0.0] - 2026-09-27\n\n### Added\n- A\n"
	if got := NotesFor(cl, "1.1.0"); got != "### Added\n- B" {
		t.Errorf("1.1.0 notes = %q", got)
	}
	if got := NotesFor(cl, "1.0.0"); got != "### Added\n- A" {
		t.Errorf("1.0.0 notes = %q", got)
	}
	if got := NotesFor(cl, "9.9.9"); got != "" {
		t.Errorf("missing version notes = %q", got)
	}
}
