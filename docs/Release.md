# Release

This document describes the maintainer release flow for this repository.

The release automation is driven by the GitHub Actions workflow in [../.github/workflows/release.yml](../.github/workflows/release.yml), which runs when a tag matching `v<major>.<minor>.<patch>` is pushed.

## What The Release Workflow Does

When you push a semantic version tag such as `v0.44.1`, the release workflow:

1. Runs the test job.
2. Builds cross-platform distribution artifacts.
3. Creates zip assets for the release.
4. Publishes a GitHub Release.
5. Publishes container images.
6. Publishes the AWS SAM application through the reusable SAM workflow.

## Before Tagging

Before creating the tag, make sure at least these items are done:

* The branch contains the final code you want to release.
* [Whats-New.md](Whats-New.md) includes the user-facing release notes.
* Relevant user documentation such as [README.md](../README.md), [AWS-SAM.md](AWS-SAM.md), and [Configuration.md](Configuration.md) reflects the release.
* Local verification has been performed.

Recommended checks:

```bash
make go-fmt
make build
make test
```

## Tag Format

Use semantic version tags in the form:

```text
v<major>.<minor>.<patch>
```

Examples:

* `v0.44.1`
* `v0.45.0`
* `v1.0.0`

## Create And Push The Tag

Annotated tag example:

```bash
git tag -a v0.44.1 -m "release v0.44.1"
git push origin v0.44.1
```

If your team requires signed tags, create the tag with `-s` instead of `-a`.

## After Pushing

After the tag is pushed, monitor these outcomes:

* GitHub Release creation
* Distribution assets uploaded from `dist/assets`
* Container image publish job completion
* AWS SAM publish job completion

If the release includes a public serverless update, verify the published application page in the AWS Serverless Application Repository after the workflow finishes.

## Related Documentation

* [Whats-New.md](Whats-New.md)
* [Development.md](Development.md)
* [AWS-SAM.md](AWS-SAM.md)
