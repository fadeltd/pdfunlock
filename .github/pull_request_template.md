Fixes #

## Why
<!-- The problem or request this solves. For a bug, the cause. -->

## Changes
<!-- What changed, file by file or area by area. -->

## Verification
<!-- What you ran and what you saw. -->
- [ ] `gofmt -l .` is empty, `go vet ./...` and `go test ./...` pass
- [ ] Tried the built binary on a real case (`make build`, then run the binary it prints)

## Checklist
- [ ] `CHANGELOG.md` has an entry under `## [Unreleased]`, or this PR only touches docs/CI/scripts. Merging releases from that section: `### Fixed`/`### Changed` → patch, `### Added` → minor, `### Breaking` → major. Label `skip-changelog` if it should not release.
- [ ] README and `docs/index.html` updated if flags or install steps changed
- [ ] No real passwords or personal PDFs in tests, `testdata/` or this description
