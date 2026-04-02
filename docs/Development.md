# Development references

## Container Registry

Container images are published to [GitHub Container Registry](https://github.com/slashdevops/idp-scim-sync/pkgs/container/idp-scim-sync) (ghcr.io) using [podman](https://podman.io/).

### Build locally

```bash
# Build cross-platform binaries
make build-dist

# Build container images (requires podman)
GIT_VERSION=test make container-build

# Verify
podman images | grep idp-scim-sync
```

### Publish to ghcr.io

```bash
# Login
REPOSITORY_REGISTRY_TOKEN=<your-token> REPOSITORY_REGISTRY_USERNAME=<your-username> make container-login

# Publish
GIT_VERSION=<version> make container-publish
```
