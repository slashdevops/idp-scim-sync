FROM alpine:3.21@sha256:48b0309ca019d89d40f670aa1bc06e426dc0931948452e8491e3d65087abc07d

ARG SERVICE_NAME="idpscim"
ARG GOOS="linux"
ARG GOARCH="amd64"
ARG BUILD_DATE=""
ARG BUILD_VERSION=""
ARG DESCRIPTION="Container image for idp-scim-sync"
ARG REPO_URL="https://github.com/slashdevops/idp-scim-sync"

ENV HOME="/app"

LABEL name="${SERVICE_NAME}" \
  org.opencontainers.image.created="${BUILD_DATE}" \
  org.opencontainers.image.version="${BUILD_VERSION}" \
  org.opencontainers.image.description="${DESCRIPTION}" \
  org.opencontainers.image.url="${REPO_URL}" \
  org.opencontainers.image.source="${REPO_URL}"

RUN apk add --no-cache --update \
  ca-certificates \
  && rm -rf /tmp/* /var/tmp/* /var/cache/apk/*

RUN mkdir -p $HOME && \
  chown -R nobody:nobody $HOME

COPY dist/$SERVICE_NAME-$GOOS-$GOARCH $HOME/$SERVICE_NAME

ENV PATH="${PATH}:${HOME}"

USER nobody:nobody
WORKDIR $HOME

CMD ["/app/idpscim", "--help"]
