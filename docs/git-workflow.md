# Git Workflow: GitHub Flow with RC/Stable Tags

## Branch model

- **main** — single long-lived branch, always deployable
- Feature/fix branches created from main, merged back via PR

## Versioning

Semantic Versioning 2.0.0: `MAJOR.MINOR.PATCH`

## Commit format

Conventional Commits: `<type>(<scope>): <description>`

- Types: `feat`, `fix`, `refactor`, `docs`, `test`, `ci`, `build`, `chore`, `perf`
- Scope (optional): `api`, `cache`, `format`, `input`, `ci`
- Lowercase, imperative mood, no trailing period
- No footer (no sign-offs, no issue references)

## Tag strategy

1. **RC tag**: `v0.1.0-rc.1` — created on main when ready for testing
2. **Testing**: verify RC in staging/production
3. **Stable tag**: `v0.1.0` — created on same commit after successful testing
4. If RC fails: fix on main, tag next RC (`v0.1.0-rc.2`)

## Release process

1. Ensure all tests pass: `make check`
2. Update CHANGELOG.md — move Unreleased items to new version section
3. Commit: `docs: update changelog for v0.1.0`
4. Tag RC: `git tag v0.1.0-rc.1 && git push origin v0.1.0-rc.1`
5. CI builds and creates pre-release with binaries
6. After verification, tag stable: `git tag v0.1.0 && git push origin v0.1.0`
7. CI creates full release with binaries
