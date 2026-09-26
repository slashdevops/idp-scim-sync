FROM alpine:3.24@sha256:294b683cb724975bec92580e1e685676bd4b50bda910ddb8c51d4cabeaec77e6

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
