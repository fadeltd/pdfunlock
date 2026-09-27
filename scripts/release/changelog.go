package main

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Release logic, derived from CHANGELOG.md.
//
// Pure and tested, because a mistake here silently ships the wrong version or
// loses release notes, and neither is obvious from a green workflow run.

type Bump string

const (
	BumpMajor Bump = "major"
	BumpMinor Bump = "minor"
	BumpPatch Bump = "patch"
	BumpNone  Bump = "none"
)

var (
	unreleasedHeading = regexp.MustCompile(`(?m)^## \[Unreleased\][ \t]*$`)
	releaseHeading    = regexp.MustCompile(`(?m)^## \[(\d+\.\d+\.\d+)\]`)
	sectionHeading    = regexp.MustCompile(`(?m)^### (.+?)[ \t]*$`)
	bullet            = regexp.MustCompile(`(?m)^[-*] `)
	semver            = regexp.MustCompile(`^(\d+)\.(\d+)\.(\d+)$`)
)

// Unreleased is the block under "## [Unreleased]".
type Unreleased struct {
	Body       string   // section bodies, verbatim
	Sections   []string // heading names found, e.g. [Added Fixed]
	HasContent bool
}

// ParseUnreleased extracts the Unreleased block, ignoring filler like "Nothing yet."
func ParseUnreleased(changelog string) Unreleased {
	loc := unreleasedHeading.FindStringIndex(changelog)
	if loc == nil {
		return Unreleased{}
	}
	rest := changelog[loc[1]:]
	if next := releaseHeading.FindStringIndex(rest); next != nil {
		rest = rest[:next[0]]
	}
	body := strings.TrimSpace(rest)

	var sections []string
	for _, m := range sectionHeading.FindAllStringSubmatch(body, -1) {
		sections = append(sections, strings.TrimSpace(m[1]))
	}

	// A section heading with no bullets under it is not content.
	return Unreleased{
		Body:       body,
		Sections:   sections,
		HasContent: len(sections) > 0 && bullet.MatchString(body),
	}
}

// CurrentVersion is the newest released version: the first "## [x.y.z]"
// heading, since Keep a Changelog lists releases newest first.
func CurrentVersion(changelog string) (string, error) {
	m := releaseHeading.FindStringSubmatch(changelog)
	if m == nil {
		return "", fmt.Errorf("no released version in the changelog")
	}
	return m[1], nil
}

// DetermineBump maps Keep a Changelog sections to a semver bump.
//
// Breaking must be explicit: inferring a major from Removed would let a tidy-up
// silently become a 2.0, which is exactly the surprise this should avoid.
func DetermineBump(u Unreleased) Bump {
	if !u.HasContent {
		return BumpNone
	}
	bump := BumpPatch
	for _, s := range u.Sections {
		switch strings.ToLower(s) {
		case "breaking", "breaking changes":
			return BumpMajor
		case "added":
			bump = BumpMinor
		}
	}
	return bump
}

func parseVersion(version string) ([3]int, error) {
	m := semver.FindStringSubmatch(strings.TrimSpace(version))
	if m == nil {
		return [3]int{}, fmt.Errorf("not a semver version: %q", version)
	}
	var v [3]int
	for i := range v {
		v[i], _ = strconv.Atoi(m[i+1])
	}
	return v, nil
}

// NextVersion applies bump to current.
func NextVersion(current string, bump Bump) (string, error) {
	v, err := parseVersion(current)
	if err != nil {
		return "", err
	}
	switch bump {
	case BumpMajor:
		// A breaking change pre-1.0 is a minor bump by convention.
		if v[0] == 0 {
			return fmt.Sprintf("0.%d.0", v[1]+1), nil
		}
		return fmt.Sprintf("%d.0.0", v[0]+1), nil
	case BumpMinor:
		return fmt.Sprintf("%d.%d.0", v[0], v[1]+1), nil
	case BumpPatch:
		return fmt.Sprintf("%d.%d.%d", v[0], v[1], v[2]+1), nil
	}
	return current, nil
}

// Release is the outcome of ApplyRelease.
type Release struct {
	Changelog string
	Version   string
	Bump      Bump
	Notes     string
}

var excessBlankLines = regexp.MustCompile(`\n{4,}`)

// ApplyRelease moves the Unreleased block into a dated release section and
// leaves a fresh empty Unreleased behind. It returns nil when there is nothing
// to release, and writes nothing.
func ApplyRelease(changelog, today string) (*Release, error) {
	u := ParseUnreleased(changelog)
	bump := DetermineBump(u)
	if bump == BumpNone {
		return nil, nil
	}
	current, err := CurrentVersion(changelog)
	if err != nil {
		return nil, err
	}
	version, err := NextVersion(current, bump)
	if err != nil {
		return nil, err
	}

	loc := unreleasedHeading.FindStringIndex(changelog)
	head, rest := changelog[:loc[1]], changelog[loc[1]:]
	tail := ""
	if next := releaseHeading.FindStringIndex(rest); next != nil {
		tail = rest[next[0]:]
	}

	rebuilt := head + "\n\nNothing yet.\n\n" +
		fmt.Sprintf("## [%s] - %s\n\n", version, today) +
		u.Body + "\n\n" + tail
	rebuilt = strings.TrimRight(excessBlankLines.ReplaceAllString(rebuilt, "\n\n\n"), "\n") + "\n"

	return &Release{Changelog: rebuilt, Version: version, Bump: bump, Notes: u.Body}, nil
}

// NotesFor returns the body of the "## [version]" section, or "" if absent.
func NotesFor(changelog, version string) string {
	heading := regexp.MustCompile(`(?m)^## \[` + regexp.QuoteMeta(version) + `\].*$`)
	loc := heading.FindStringIndex(changelog)
	if loc == nil {
		return ""
	}
	rest := changelog[loc[1]:]
	if next := releaseHeading.FindStringIndex(rest); next != nil {
		rest = rest[:next[0]]
	}
	return strings.TrimSpace(rest)
}
