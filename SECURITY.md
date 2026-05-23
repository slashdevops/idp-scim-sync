# Security Policy

This project participates in the [OpenSSF Scorecard](https://github.com/ossf/scorecard) program and ships with several security controls enabled by default:

* [![OpenSSF Scorecard](https://api.securityscorecards.dev/projects/github.com/slashdevops/idp-scim-sync/badge)](https://securityscorecards.dev/viewer/?uri=github.com/slashdevops/idp-scim-sync) — continuous Scorecard analysis (`.github/workflows/scorecard.yml`).
* [![CodeQL](https://github.com/slashdevops/idp-scim-sync/actions/workflows/codeql.yml/badge.svg)](https://github.com/slashdevops/idp-scim-sync/actions/workflows/codeql.yml) — CodeQL SAST on every push and PR (`.github/workflows/codeql.yml`).
* `govulncheck` runs in the `Build` workflow and blocks PRs that introduce known Go vulnerabilities.
* All GitHub Actions are pinned by full commit SHA; container base images are pinned by `@sha256` digest.
* Release binaries are signed and shipped with SLSA Level 3 provenance (see [Verifying release artifacts](#verifying-release-artifacts) below).
* Container images published to `ghcr.io/slashdevops/idp-scim-sync` are signed with Cosign keyless (Sigstore).

## Supported Versions

The project follows the latest Go release line and updates its dependencies on a continuous basis. Only the most recent minor version receives security fixes.

| Version | Supported          |
| ------- | ------------------ |
| 0.44.x  | :white_check_mark: |
| 0.43.x  | :x:                |
| 0.42.x  | :x:                |
| 0.32.x  | :x:                |
| 0.31.x  | :x:                |
| 0.30.x  | :x:                |
| 0.2.x   | :x:                |
| 0.1.x   | :x:                |
| 0.0.x   | :x:                |

## Reporting a Vulnerability

Use the [Project Issues → Vulnerability template](https://github.com/slashdevops/idp-scim-sync/issues/new/choose) to report a security issue. For sensitive reports, please use [GitHub's private vulnerability reporting](https://github.com/slashdevops/idp-scim-sync/security/advisories/new) instead of a public issue.

## Verifying release artifacts

Starting with versions released after [PR-3 of the OpenSSF hardening effort](https://github.com/slashdevops/idp-scim-sync/issues?q=label%3Aopenssf), every release ships with a SLSA Level 3 provenance attestation (`multiple.intoto.jsonl`) and container images are signed with Cosign keyless (Sigstore).

### Binary release zips (SLSA provenance)

Download the release zip(s) plus the `multiple.intoto.jsonl` attestation from the same release page, then verify with [`slsa-verifier`](https://github.com/slsa-framework/slsa-verifier):

```shell
# Install slsa-verifier (one-time)
go install github.com/slsa-framework/slsa-verifier/v2/cli/slsa-verifier@latest

# From the directory where you downloaded the assets
slsa-verifier verify-artifact \
  --provenance-path multiple.intoto.jsonl \
  --source-uri github.com/slashdevops/idp-scim-sync \
  --source-tag v0.45.0 \
  idpscim-linux-amd64.zip
```

A successful run prints `PASSED: SLSA verification passed`.

### Container images (Cosign keyless)

Verify the signature on the published multi-arch manifest:

```shell
# Install cosign (one-time)
go install github.com/sigstore/cosign/v2/cmd/cosign@latest

cosign verify \
  --certificate-identity-regexp '^https://github\.com/slashdevops/idp-scim-sync/\.github/workflows/' \
  --certificate-oidc-issuer 'https://token.actions.githubusercontent.com' \
  ghcr.io/slashdevops/idp-scim-sync:v0.45.0
```

Cosign prints the verified signature, certificate, and transparency-log entry.
