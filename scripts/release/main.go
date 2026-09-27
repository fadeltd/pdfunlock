// Command release cuts a release from CHANGELOG.md.
//
// It reads the Unreleased section, decides the semver bump from its headings,
// rewrites CHANGELOG.md, and reports the outcome for the workflow to consume.
// The current version is the newest released section of the changelog. It
// writes the changelog but never touches git: the workflow owns committing,
// tagging and pushing, so this stays runnable locally as a dry run.
//
//	go run ./scripts/release -dry-run     # preview the next release
//	go run ./scripts/release              # apply it
//	go run ./scripts/release -notes 1.1.0 # print the notes of a released version
package main

import (
	"flag"
	"fmt"
	"os"
	"time"
)

const changelogPath = "CHANGELOG.md"

func main() {
	dryRun := flag.Bool("dry-run", false, "print what would be released without writing files")
	notes := flag.String("notes", "", "print the changelog notes of this version and exit")
	flag.Parse()

	if err := run(*dryRun, *notes); err != nil {
		fmt.Fprintln(os.Stderr, "release:", err)
		os.Exit(1)
	}
}

func run(dryRun bool, notesVersion string) error {
	changelog, err := os.ReadFile(changelogPath)
	if err != nil {
		return err
	}

	if notesVersion != "" {
		n := NotesFor(string(changelog), notesVersion)
		if n == "" {
			return fmt.Errorf("no [%s] section in %s", notesVersion, changelogPath)
		}
		fmt.Println(n)
		return nil
	}

	current, err := CurrentVersion(string(changelog))
	if err != nil {
		return err
	}
	today := time.Now().UTC().Format("2006-01-02")

	rel, err := ApplyRelease(string(changelog), today)
	if err != nil {
		return err
	}
	if rel == nil {
		fmt.Println("Nothing under [Unreleased]: no release.")
		return emit("released", "false")
	}

	fmt.Printf("bump    : %s\nversion : %s -> %s\ntag     : v%s\n", rel.Bump, current, rel.Version, rel.Version)
	if dryRun {
		fmt.Printf("\n--- notes ---\n%s\n\n(dry run: no files written)\n", rel.Notes)
		return nil
	}

	if err := os.WriteFile(changelogPath, []byte(rel.Changelog), 0o644); err != nil {
		return err
	}

	for k, v := range map[string]string{"released": "true", "version": rel.Version, "tag": "v" + rel.Version} {
		if err := emit(k, v); err != nil {
			return err
		}
	}

	fmt.Printf("\n%s updated.\n", changelogPath)
	return nil
}

// emit sets a single-line step output when running under GitHub Actions.
func emit(key, value string) error {
	out := os.Getenv("GITHUB_OUTPUT")
	if out == "" {
		return nil
	}
	f, err := os.OpenFile(out, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = fmt.Fprintln(f, key+"="+value)
	return err
}
