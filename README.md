# neochess

A single-binary, locally-hosted chess interface. Play against UCI engines,
import and review PGN games, and keep a searchable library — all served from
`127.0.0.1` with no external services required.

## Requirements

- [Go 1.27+](https://go.dev/dl/)
- [just](https://github.com/casey/just) — task runner
- [goreleaser](https://goreleaser.com/install/) — required only for local release dry-runs

## Release process

Releases are fully automated via [GoReleaser](https://goreleaser.com) and GitHub Actions. Pushing a semver tag to `main` is the only manual step required.

### What is published

Each release produces the following and attaches them to the GitHub release:

| Artifact | Description |
|---|---|
| `neochess_<version>_linux_amd64.tar.gz` | Linux x86-64 binary |
| `neochess_<version>_linux_arm64.tar.gz` | Linux ARM64 binary |
| `neochess_<version>_darwin_amd64.tar.gz` | macOS Intel binary |
| `neochess_<version>_darwin_arm64.tar.gz` | macOS Apple Silicon binary |
| `checksums.txt` | SHA-256 checksums for all archives |

### Cutting a release

1. **Ensure `main` is green.** The CI pipeline (lint, race tests, vuln scan) must be passing before you tag.
2. **Choose a version.** This project follows [Semantic Versioning](https://semver.org):
- `v<MAJOR>.<MINOR>.<PATCH>` e.g. `v1.2.0`
- Increment PATCH for bug fixes, `MINOR` for new features, `MAJOR` for breaking changes.
3. **Tag and push** using the `just recipe`:
```shell
just release v1.2.0
```
This creates a signed, annotated git tag and pushes it to `origin`. That push triggers the release workflow automatically - no further action is needed.
4. **Monitor the workflow.** Open the [Actions Tab](../../actions/workflows/release.yml) in GitHub to watch GoReleaser build and publish the release. The job typically completes in under two mins.
5. **Verify the release.** Once the workflow succeeds, confirm that the [Releases Page](../../releases) shows the correct tag, changelog, and all expected artifacts with their checksums.

### Local dry-run

To validate the GoReleaser configuration without publishing anything:
```shell
just release-snapshot
```

Aritfacts are written to `dist/` and no tag or GitHub Release is created. Run this whenever you change `.goreleaser.yaml` to catch configuration errors before tagging.

### Changelog

The changelog included in each GitHub Release is generated automatically from commit messages between the previous and new tag. Use conventional commit prefixes to keep it readable.

| Prefix | Section in changelog |
|---|---|
| `feat:` | New Features |
| `fix:` | Bug Fixes |
| `chore:`, `docs:`, `test:`, `ci:` | Excluded (noise) |

### Recovering from a bad tag

If you tag the wrong commit or need to redo a release:

```shell
git tag -d v1.2.0 git push origin :refs/tags/v1.2.0
```

Then fix whatever needs fixing, re-run `just check` and re-tag. If the GitHub Release was already created, delete it from the Releases page before re-pushing the tag, or GoReleaser will error on the duplicate.
