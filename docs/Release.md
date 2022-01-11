# Release

This document explain the basic of how to release a new version of the project.

## Process

The release process is divided into three steps:

1. Create a new git tag
2. This new tag is pushed to the remote repository
3. The Release Pipeline is triggered

## Commands

```bash
git tag -a v0.0.1 -m "testing the release process" -s
git push origin v0.0.1
```
